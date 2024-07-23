package keeper

import (
	"testing"

	"github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/osmosis-labs/mesh-security-sdk/x/meshsecurityprovider/types"
)

func TestStoreParams(t *testing.T) {
	ctx, keepers := CreateDefaultTestInput(t)
	k := keepers.MeshKeeperProvider

	vaultContract := sdk.AccAddress(rand.Bytes(32))

	params := types.Params{
		TimeoutPeriod:        60,
		VaultContractAddress: vaultContract.String(),
	}
	err := k.SetParams(ctx, params)
	require.NoError(t, err)

	gotParams := k.GetParams(ctx)
	require.Equal(t, params.VaultContractAddress, gotParams.VaultContractAddress)
	require.Equal(t, params.TimeoutPeriod, gotParams.TimeoutPeriod)
}
