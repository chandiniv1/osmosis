package poolmanager

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/osmosis-labs/osmosis/osmoutils"
	cosmwasmpooltypes "github.com/osmosis-labs/osmosis/v21/x/cosmwasmpool/types"
	gammtypes "github.com/osmosis-labs/osmosis/v21/x/gamm/types"
	"github.com/osmosis-labs/osmosis/v21/x/poolmanager/types"
)

var (
	OSMO                 = "uosmo"
	superfluidMultiplier = sdk.MustNewDecFromStr("1.5")
	directRouteCache     map[string]uint64
	spotPriceCache       map[string]osmomath.BigDec
)

func init() {
	directRouteCache = make(map[string]uint64)
	spotPriceCache = make(map[string]osmomath.BigDec)
}

// findRoutes finds all routes between two tokens that match the given hop count
func findRoutes(g types.RoutingGraphMap, start, end string, hops int) [][]*types.Route {
	if hops < 1 {
		return nil
	}

	var routeRoutes [][]*types.Route

	startRoutes, startExists := g.Graph[start]
	if !startExists {
		return routeRoutes
	}

	for token, routes := range startRoutes.InnerMap {
		if hops == 1 {
			if token == end {
				for _, route := range routes.Routes {
					route.Token = end
					routeRoutes = append(routeRoutes, []*types.Route{route})
				}
			}
		} else {
			subRoutes := findRoutes(g, token, end, hops-1)
			for _, subRoute := range subRoutes {
				for _, route := range routes.Routes {
					route.Token = token
					fullRoute := append([]*types.Route{route}, subRoute...)
					routeRoutes = append(routeRoutes, fullRoute)
				}
			}
		}
	}

	return routeRoutes
}

// SetDenomPairRoutes sets the route map to be used for route calculations
func (k Keeper) SetDenomPairRoutes(ctx sdk.Context) (types.RoutingGraph, error) {
	// Get all the pools
	pools, err := k.AllPools(ctx)
	if err != nil {
		return types.RoutingGraph{}, err
	}

	// Create a routingGraph to represent possible routes between tokens
	var routingGraph types.RoutingGraph

	// Iterate through the pools
	for _, pool := range pools {
		tokens := pool.GetPoolDenoms(ctx)
		poolID := pool.GetId()
		// Create edges for all possible combinations of tokens
		for i := 0; i < len(tokens); i++ {
			for j := i + 1; j < len(tokens); j++ {
				// Add edges with the associated token
				routingGraph.AddEdge(tokens[i], tokens[j], tokens[i], poolID)
				routingGraph.AddEdge(tokens[j], tokens[i], tokens[j], poolID)
			}
		}
	}

	// Set the route map in state
	// NOTE: This is done with the non map version of the route graph
	// If we used maps here, the serialization would be non-deterministic
	osmoutils.MustSet(ctx.KVStore(k.storeKey), types.KeyRouteMap, &routingGraph)
	return routingGraph, nil
}

// GetDenomPairRoute returns the route with the lowest slippage between an input denom and an output denom
// The method starts by finding all direct routes, two hop, and three hop routes.
// It then calculates the liquidity of each route and selects the route with the highest liquidity for each hop count.
// In other words, at a max, there will be one direct route, one two hop route, and one three hop route.
// The method then simulates a swap against each route and selects the route with the highest swap amount out.
func (k Keeper) GetDenomPairRoute(ctx sdk.Context, inputCoin sdk.Coin, outputDenom string) ([]types.SwapAmountInRoute, error) {
	inputDenom := inputCoin.Denom

	// Ensure the caches are restored when the function exits
	defer func() {
		directRouteCache = make(map[string]uint64)
		spotPriceCache = make(map[string]osmomath.BigDec)
	}()

	// Retrieve the route map from state
	routeMap, err := k.GetRouteMap(ctx)
	if err != nil {
		return nil, err
	}

	// Get all direct routes, two hop, and three hop routes
	directPoolIDs := findRoutes(routeMap, inputDenom, outputDenom, 1)

	var twoHopPoolIDs [][]*types.Route
	if inputDenom != OSMO && outputDenom != OSMO {
		twoHopPoolIDs = findRoutes(routeMap, inputDenom, outputDenom, 2)
	}

	var threeHopPoolIDs [][]*types.Route
	if inputDenom != OSMO && outputDenom != OSMO {
		threeHopPoolIDs = findRoutes(routeMap, inputDenom, outputDenom, 3)
	}

	// Map the total liquidity of each route (using the route string as the key)
	routeLiquidity := make(map[string]osmomath.Int)

	// Check liquidity for all direct routes
	for _, route := range directPoolIDs {
		pool, err := k.GetPool(ctx, route[0].PoolId)
		if err != nil {
			return nil, err
		}
		poolDenoms := pool.GetPoolDenoms(ctx)
		liqInOsmo := osmomath.ZeroInt()
		for _, denom := range poolDenoms {
			liquidity, err := k.getPoolLiquidityOfDenom(ctx, route[0].PoolId, denom)
			if err != nil {
				return nil, err
			}
			liqInOsmoInternal, err := k.inputAmountToOSMO(ctx, denom, liquidity, routeMap)
			if err != nil {
				return nil, err
			}

			if pool.GetType() == types.Concentrated {
				liqInOsmoInternal = liqInOsmoInternal.ToLegacyDec().Mul(superfluidMultiplier).TruncateInt()
			}

			liqInOsmo = liqInOsmo.Add(liqInOsmoInternal)
		}
		routeKey := fmt.Sprintf("%v", route)
		routeLiquidity[routeKey] = liqInOsmo
	}

	// Check liquidity for all two-hop routes
	for _, routes := range twoHopPoolIDs {
		totalLiquidityInOsmo := osmomath.ZeroInt()
		routeKey := fmt.Sprintf("%v", routes)
		for _, route := range routes {
			pool, err := k.GetPool(ctx, route.PoolId)
			if err != nil {
				return nil, err
			}
			poolDenoms := pool.GetPoolDenoms(ctx)
			liqInOsmo := osmomath.ZeroInt()
			for _, denom := range poolDenoms {
				liquidity, err := k.getPoolLiquidityOfDenom(ctx, route.PoolId, denom)
				if err != nil {
					return nil, err
				}
				liqInOsmoInternal, err := k.inputAmountToOSMO(ctx, denom, liquidity, routeMap)
				if err != nil {
					return nil, err
				}

				if pool.GetType() == types.Concentrated {
					liqInOsmoInternal = liqInOsmoInternal.ToLegacyDec().Mul(superfluidMultiplier).TruncateInt()
				}
				liqInOsmo = liqInOsmo.Add(liqInOsmoInternal)
			}

			totalLiquidityInOsmo = totalLiquidityInOsmo.Add(liqInOsmo)
		}
		routeLiquidity[routeKey] = totalLiquidityInOsmo
	}

	// Check liquidity for all three-hop routes
	for _, routes := range threeHopPoolIDs {
		totalLiquidityInOsmo := osmomath.ZeroInt()
		routeKey := fmt.Sprintf("%v", routes)
		for _, route := range routes {
			pool, err := k.GetPool(ctx, route.PoolId)
			if err != nil {
				return nil, err
			}
			poolDenoms := pool.GetPoolDenoms(ctx)
			liqInOsmo := osmomath.ZeroInt()
			for _, denom := range poolDenoms {
				liquidity, err := k.getPoolLiquidityOfDenom(ctx, route.PoolId, denom)
				if err != nil {
					return nil, err
				}
				liqInOsmoInternal, err := k.inputAmountToOSMO(ctx, denom, liquidity, routeMap)
				if err != nil {
					return nil, err
				}

				if pool.GetType() == types.Concentrated {
					liqInOsmoInternal = liqInOsmoInternal.ToLegacyDec().Mul(superfluidMultiplier).TruncateInt()
				}
				liqInOsmo = liqInOsmo.Add(liqInOsmoInternal)
			}
			totalLiquidityInOsmo = totalLiquidityInOsmo.Add(liqInOsmo)
		}
		routeLiquidity[routeKey] = totalLiquidityInOsmo
	}

	// Extract and sort the keys from the routeLiquidity map
	// This ensures deterministic selection of the best route
	var keys []string
	for k := range routeLiquidity {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Variables to store the best route key for each hop count
	var bestSingleHopRouteKey, bestDoubleHopRouteKey, bestTripleHopRouteKey string
	maxSingleHopLiquidity, maxDoubleHopLiquidity, maxTripleHopLiquidity := osmomath.ZeroInt(), osmomath.ZeroInt(), osmomath.ZeroInt()

	// Iterate through all viable routes and select the route with the highest liquidity for each hop count
	for _, routeKey := range keys {
		liquidity := routeLiquidity[routeKey]
		hopCount := len(strings.Fields(routeKey)) / 2

		switch hopCount {
		case 1: // Single hop
			if liquidity.GT(maxSingleHopLiquidity) {
				maxSingleHopLiquidity = liquidity
				bestSingleHopRouteKey = routeKey
			}
		case 2: // Double hop
			if liquidity.GT(maxDoubleHopLiquidity) {
				maxDoubleHopLiquidity = liquidity
				bestDoubleHopRouteKey = routeKey
			}
		case 3: // Triple hop
			if liquidity.GT(maxTripleHopLiquidity) {
				maxTripleHopLiquidity = liquidity
				bestTripleHopRouteKey = routeKey
			}
		}
	}

	// Construct the result map
	result := make(map[string][]types.Route)

	// Parse the route keys and store the result in the result map
	singleHopRoute, err := parseRouteKey(bestSingleHopRouteKey)
	if err != nil {
		return nil, fmt.Errorf("error parsing single hop route key: %v", err)
	}
	result["singleHop"] = singleHopRoute

	doubleHopRoute, err := parseRouteKey(bestDoubleHopRouteKey)
	if err != nil {
		return nil, fmt.Errorf("error parsing double hop route key: %v", err)
	}
	result["doubleHop"] = doubleHopRoute

	tripleHopRoute, err := parseRouteKey(bestTripleHopRouteKey)
	if err != nil {
		return nil, fmt.Errorf("error parsing triple hop route key: %v", err)
	}
	result["tripleHop"] = tripleHopRoute

	maxAmtOut := sdk.ZeroInt()

	var resultAsString []string
	for k := range result {
		resultAsString = append(resultAsString, k)
	}
	sort.Strings(resultAsString)

	// Iterate through the result map and simulate a swap against each route
	// Select the route with the highest swap amount out
	var maxKey string
	for _, key := range resultAsString {
		value := result[key]
		swapRoute := []types.SwapAmountInRoute{}
		for _, route := range value {
			// Construct SwapAmountInRoute for each poolID
			swapRoute = append(swapRoute, types.SwapAmountInRoute{
				PoolId:        route.PoolId,
				TokenOutDenom: route.Token,
			})
		}

		// Call MultihopEstimateOutGivenExactAmountIn with swapRoute
		amtOut, err := k.MultihopEstimateOutGivenExactAmountIn(ctx, swapRoute, inputCoin)
		if err != nil {
			continue
		}

		// Update maxAmtOut and maxKey if the current amtOut is greater
		if amtOut.GT(maxAmtOut) {
			maxAmtOut = amtOut
			maxKey = key
		}
	}

	// Return the SwapAmountInRoute for the best route
	var swapRoutes []types.SwapAmountInRoute
	for _, route := range result[maxKey] {
		swapRoutes = append(swapRoutes, types.SwapAmountInRoute{
			PoolId:        route.PoolId,
			TokenOutDenom: route.Token,
		})
	}

	return swapRoutes, nil
}

// parseRouteKey is a helper function to parse route key into a slice of Route.
func parseRouteKey(routeKey string) ([]types.Route, error) {
	var route []types.Route
	if routeKey == "" {
		return route, nil
	}
	cleanedRouteKey := strings.Trim(routeKey, "[]")

	// Regular expression to match pool_id and token
	re := regexp.MustCompile(`pool_id:(\d+) token:"([^"]+)"`)

	matches := re.FindAllStringSubmatch(cleanedRouteKey, -1)
	for _, match := range matches {
		if len(match) != 3 {
			return nil, fmt.Errorf("invalid route key format: %v", match)
		}
		id, err := strconv.ParseUint(match[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing pool ID: %v", err)
		}
		token := match[2]
		route = append(route, types.Route{PoolId: id, Token: token})
	}
	return route, nil
}

// getDirectOSMORouteWithMostLiquidity returns the route with the highest liquidity between an input denom and uosmo
func (k Keeper) getDirectOSMORouteWithMostLiquidity(ctx sdk.Context, inputDenom string, routeMap types.RoutingGraphMap) (uint64, error) {
	// Get all direct routes from the input denom to uosmo
	directRoutes := findRoutes(routeMap, inputDenom, OSMO, 1)

	// Store liquidity for all direct routes found
	routeLiquidity := make(map[string]osmomath.Int)
	for _, route := range directRoutes {
		liquidity, err := k.getPoolLiquidityOfDenom(ctx, route[0].PoolId, OSMO)
		if err != nil {
			return 0, err
		}
		routeKey := fmt.Sprintf("%v", route[0].PoolId)
		routeLiquidity[routeKey] = liquidity
	}

	// Extract and sort the keys from the routeLiquidity map
	// This ensures deterministic selection of the best route
	var keys []string
	for k := range routeLiquidity {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Find the route (single or double hop) with the highest liquidity
	var bestRouteKey string
	maxLiquidity := osmomath.ZeroInt()
	for _, routeKey := range keys {
		liquidity := routeLiquidity[routeKey]
		// Update best route if a higher liquidity is found,
		// or if the liquidity is equal but the routeKey is encountered earlier in the sorted order
		if liquidity.GT(maxLiquidity) || (liquidity.Equal(maxLiquidity) && bestRouteKey == "") {
			bestRouteKey = routeKey
			maxLiquidity = liquidity
		}
	}
	if bestRouteKey == "" {
		return 0, fmt.Errorf("no route found with sufficient liquidity, likely no direct pairing with osmo")
	}

	// Convert the best route key back to []uint64
	var bestRoute []uint64
	cleanedRouteKey := strings.Trim(bestRouteKey, "[]")
	idStrs := strings.Split(cleanedRouteKey, " ")

	for _, idStr := range idStrs {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("error parsing pool ID: %v", err)
		}
		bestRoute = append(bestRoute, id)
	}

	// Return the route with the highest liquidity
	return bestRoute[0], nil
}

// inputAmountToOSMO transforms an input denom and its amount to uosmo
// If a route is not found, returns 0 with no error.
func (k Keeper) inputAmountToOSMO(ctx sdk.Context, inputDenom string, amount osmomath.Int, routeMap types.RoutingGraphMap) (osmomath.Int, error) {
	if inputDenom == OSMO {
		return amount, nil
	}

	var route uint64
	var err error

	// Check if the route is cached
	if cachedRoute, ok := directRouteCache[inputDenom]; ok {
		route = cachedRoute
	} else {
		// If not, get the route and cache it
		route, err = k.getDirectOSMORouteWithMostLiquidity(ctx, inputDenom, routeMap)
		if err != nil {
			return osmomath.ZeroInt(), nil
		}
		directRouteCache[inputDenom] = route
	}

	var osmoPerInputToken osmomath.BigDec

	// Check if the spot price is cached
	spotPriceKey := fmt.Sprintf("%d:%s", route, inputDenom)
	if cachedSpotPrice, ok := spotPriceCache[spotPriceKey]; ok {
		osmoPerInputToken = cachedSpotPrice
	} else {
		// If not, calculate the spot price and cache it
		osmoPerInputToken, err = k.RouteCalculateSpotPrice(ctx, route, OSMO, inputDenom)
		if err != nil {
			return osmomath.ZeroInt(), err
		}
		spotPriceCache[spotPriceKey] = osmoPerInputToken
	}

	// Convert the input denom to uosmo
	// Rounding is fine here
	uosmoAmount := amount.ToLegacyDec().Mul(osmoPerInputToken.Dec())
	return uosmoAmount.RoundInt(), nil
}

// getPoolLiquidityOfDenom returns the liquidity of a denom in a pool.
// This calls different methods depending on the pool type.
func (k Keeper) getPoolLiquidityOfDenom(ctx sdk.Context, poolId uint64, denom string) (osmomath.Int, error) {
	pool, err := k.GetPool(ctx, poolId)
	if err != nil {
		return osmomath.ZeroInt(), err
	}

	// Check the pool type, and check the pool liquidity based on the type
	switch pool.GetType() {
	case types.Balancer:
		pool, ok := pool.(gammtypes.CFMMPoolI)
		if !ok {
			return osmomath.ZeroInt(), fmt.Errorf("invalid pool type")
		}
		totalPoolLiquidity := pool.GetTotalPoolLiquidity(ctx)
		return totalPoolLiquidity.AmountOf(denom), nil
	case types.Stableswap:
		pool, ok := pool.(gammtypes.CFMMPoolI)
		if !ok {
			return osmomath.ZeroInt(), fmt.Errorf("invalid pool type")
		}
		totalPoolLiquidity := pool.GetTotalPoolLiquidity(ctx)
		return totalPoolLiquidity.AmountOf(denom), nil
	case types.Concentrated:
		poolAddress := pool.GetAddress()
		poolAddressBalances := k.bankKeeper.GetAllBalances(ctx, poolAddress)
		return poolAddressBalances.AmountOf(denom), nil
	case types.CosmWasm:
		pool, ok := pool.(cosmwasmpooltypes.CosmWasmExtension)
		if !ok {
			return osmomath.ZeroInt(), fmt.Errorf("invalid pool type")
		}
		totalPoolLiquidity := pool.GetTotalPoolLiquidity(ctx)
		return totalPoolLiquidity.AmountOf(denom), nil
	default:
		return osmomath.ZeroInt(), fmt.Errorf("invalid pool type")
	}
}

// GetRouteMap returns the route map that is stored in state.
// It converts the route graph stored in the KVStore to a map for easier access.
func (k Keeper) GetRouteMap(ctx sdk.Context) (types.RoutingGraphMap, error) {
	var routeGraph types.RoutingGraph

	found, err := osmoutils.Get(ctx.KVStore(k.storeKey), types.KeyRouteMap, &routeGraph)
	if err != nil {
		return types.RoutingGraphMap{}, err
	}
	if !found {
		return types.RoutingGraphMap{}, fmt.Errorf("route map not found")
	}

	routeMap := convertToMap(&routeGraph)

	return routeMap, nil
}

// convertToMap converts a RoutingGraph to a RoutingGraphMap
// This is done to take advantage of the map data structure for easier access.
func convertToMap(routingGraph *types.RoutingGraph) types.RoutingGraphMap {
	result := types.RoutingGraphMap{
		Graph: make(map[string]*types.InnerMap),
	}

	for _, graphEntry := range routingGraph.Entries {
		innerMap := &types.InnerMap{
			InnerMap: make(map[string]*types.Routes),
		}
		for _, innerMapEntry := range graphEntry.Value.Entries {
			routes := make([]*types.Route, len(innerMapEntry.Value.Routes))
			for i, route := range innerMapEntry.Value.Routes {
				routes[i] = &types.Route{PoolId: route.PoolId, Token: route.Token}
			}
			innerMap.InnerMap[innerMapEntry.Key] = &types.Routes{Routes: routes}
		}
		result.Graph[graphEntry.Key] = innerMap
	}

	return result
}