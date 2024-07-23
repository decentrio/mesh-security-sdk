package keeper

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/mesh-security-sdk/x/meshsecurity/types"
)

// abstract keeper
type maxCapSource interface {
	HasMaxCapLimit(ctx sdk.Context, actor sdk.AccAddress) bool
}

// NewIntegrityHandler prevents any contract with max cap set to use staking
// or stargate messages. This ensures that staked "virtual" tokens are not bypassing
// the instant undelegate and burn mechanism provided by mesh-security.
//
// This handler should be chained before any other.
func NewIntegrityHandler(k maxCapSource) wasmkeeper.MessageHandlerFunc {
	return func(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) (
		events []sdk.Event,
		data [][]byte,
		err error,
	) {
		if msg.Stargate == nil && msg.Staking == nil ||
			!k.HasMaxCapLimit(ctx, contractAddr) {
			return nil, nil, wasmtypes.ErrUnknownMsg // pass down the chain
		}
		// reject
		return nil, nil, types.ErrUnsupported.Wrap("message type for contracts with max cap set")
	}
}
