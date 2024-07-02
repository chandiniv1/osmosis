package service

import (
	"context"
	"sync"

	"github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	commondomain "github.com/osmosis-labs/osmosis/v25/ingest/common/domain"
	"github.com/osmosis-labs/osmosis/v25/ingest/indexer/domain"
	"github.com/osmosis-labs/osmosis/v25/ingest/indexer/service/blockprocessor"
)

var _ baseapp.StreamingService = (*indexerStreamingService)(nil)

// ind is a streaming service that processes block data and ingests it into the indexer
type indexerStreamingService struct {
	writeListeners map[storetypes.StoreKey][]storetypes.WriteListener

	// manages tracking of whether all the data should be processed or only the changed in the block
	pushStrategyManager commondomain.PushStrategyManager

	client domain.Publisher

	keepers domain.Keepers

	// extracts the pools from chain state
	poolExtracter commondomain.PoolExtracter
}

// New creates a new sqsStreamingService.
// writeListeners is a map of store keys to write listeners.
// sqsIngester is an ingester that ingests the block data into SQS.
// poolTracker is a tracker that tracks the pools that were changed in the block.
// nodeStatusChecker is a checker that checks if the node is syncing.
func New(writeListeners map[storetypes.StoreKey][]storetypes.WriteListener, coldStartManager commondomain.PushStrategyManager, client domain.Publisher, poolExtracter commondomain.PoolExtracter, keepers domain.Keepers) baseapp.StreamingService {
	return &indexerStreamingService{

		writeListeners: writeListeners,

		pushStrategyManager: coldStartManager,

		poolExtracter: poolExtracter,

		client: client,

		keepers: keepers,
	}
}

// Close implements baseapp.StreamingService.
func (s *indexerStreamingService) Close() error {
	return nil
}

// ListenBeginBlock implements baseapp.StreamingService.
func (s *indexerStreamingService) ListenBeginBlock(ctx context.Context, req types.RequestBeginBlock, res types.ResponseBeginBlock) error {
	return nil
}

// ListenCommit implements baseapp.StreamingService.
func (s *indexerStreamingService) ListenCommit(ctx context.Context, res types.ResponseCommit) error {
	return nil
}

// ListenDeliverTx implements baseapp.StreamingService.
func (s *indexerStreamingService) ListenDeliverTx(ctx context.Context, req types.RequestDeliverTx, res types.ResponseDeliverTx) error {
	return nil
}

// publishBlock publishes the block data to the indexer.
func (s *indexerStreamingService) publishBlock(ctx context.Context, req types.RequestEndBlock) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	height := (uint64)(req.GetHeight())
	timeEndBlock := sdkCtx.BlockTime().UTC()
	chainId := sdkCtx.ChainID()
	gasConsumed := sdkCtx.GasMeter().GasConsumed()
	block := domain.Block{
		ChainId:     chainId,
		Height:      height,
		BlockTime:   timeEndBlock,
		GasConsumed: gasConsumed,
	}
	return s.client.PublishBlock(sdkCtx, block)
}

// ListenEndBlock implements baseapp.StreamingService.
func (s *indexerStreamingService) ListenEndBlock(ctx context.Context, req types.RequestEndBlock, res types.ResponseEndBlock) error {
	// Publish the block data
	err := s.publishBlock(ctx, req)
	if err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Create block processor
	blockProcessor := blockprocessor.NewBlockProcessor(s.pushStrategyManager, s.client, s.poolExtracter, s.keepers)

	// Process block.
	if err := blockProcessor.ProcessBlock(sdkCtx); err != nil {
		return err
	}

	return nil
}

// Listeners implements baseapp.StreamingService.
func (s *indexerStreamingService) Listeners() map[storetypes.StoreKey][]storetypes.WriteListener {
	return s.writeListeners
}

// Stream implements baseapp.StreamingService.
func (s *indexerStreamingService) Stream(wg *sync.WaitGroup) error {
	return nil
}
