package keeper

import (
	// "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	// evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	// slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	// stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	// "github.com/osmosis-labs/mesh-security-sdk/x/meshsecurity/types"
)

func (k Keeper) HandleSlash(ctx sdk.Context, valAddr sdk.ValAddress, slashRatio sdk.Dec) error {
	if valAddr == nil {
		ModuleLogger(ctx).
			Error("can not propagate slash: validator not found", "validator", valAddr.String())
	}
	if err := k.ScheduleSlashed(ctx, valAddr, slashRatio); err != nil {
		ModuleLogger(ctx).
			Error("can not propagate slash: schedule event",
				"cause", err,
				"validator", valAddr.String())
	}
	if err := k.ScheduleJailed(ctx, valAddr); err != nil {
		ModuleLogger(ctx).
			Error("can not propagate jail: schedule event",
				"cause", err,
				"validator", valAddr.String())
	}

	return nil
}
