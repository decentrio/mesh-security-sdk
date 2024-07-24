package keeper

import (
	"testing"

	"github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/osmosis-labs/mesh-security-sdk/x/meshsecurityprovider/types"
)

func TestQueryParams(t *testing.T) {
	ctx, keepers := CreateDefaultTestInput(t)
	k := keepers.MeshKeeperProvider

	vaultContract := sdk.AccAddress(rand.Bytes(32))

	params := types.Params{
		TimeoutPeriod:        60,
		VaultContractAddress: vaultContract.String(),
	}
	err := k.SetParams(ctx, params)
	require.NoError(t, err)

	gotRsp, gotErr := NewQuerier(keepers.EncodingConfig.Marshaler, k).Params(sdk.WrapSDKContext(ctx), &types.QueryParamsRequest{})
	require.NoError(t, gotErr)
	assert.Equal(t, params.VaultContractAddress, gotRsp.Params.VaultContractAddress)
	assert.Equal(t, params.TimeoutPeriod, gotRsp.Params.TimeoutPeriod)
}
