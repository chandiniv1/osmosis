package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	gammmigration "github.com/osmosis-labs/osmosis/v24/x/gamm/types/migration"
)

// DefaultGenesis creates a default GenesisState object.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Pools:            []*codectypes.Any{},
		NextPoolNumber:   1,
		Params:           DefaultParams(),
		MigrationRecords: &gammmigration.MigrationRecords{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	return nil
}
