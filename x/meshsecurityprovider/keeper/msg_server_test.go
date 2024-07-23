package keeper

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/cometbft/cometbft/libs/rand"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/mesh-security-sdk/x/meshsecurityprovider/types"
)

func TestDelegateSuccessfully(t *testing.T) {
	pCtx, keepers := CreateDefaultTestInput(t)
	k := keepers.MeshKeeperProvider

	localTokenDenom := keepers.StakingKeeper.BondDenom(pCtx)
	rewardTokenDenom := "juno"
	delAddr := sdk.AccAddress(rand.Bytes(32))
	valAddr := sdk.ValAddress(rand.Bytes(32))
	balance, err := sdk.ParseCoinsNormalized("10000000" + rewardTokenDenom + ",10000000" + localTokenDenom)
	require.NoError(t, err)
	err = k.bank.MintCoins(pCtx, types.ModuleName, balance)
	require.NoError(t, err)
	k.bank.SendCoinsFromModuleToAccount(pCtx, types.ModuleName, delAddr, balance)
	localSlashRatioDoubleSign := "0.20"
	localSlashRatioOffline := "0.10"
	extSlashRatioDoubleSign := "0.20"
	extSlashRatioOffline := "0.10"
	connId := ""
	portID := ""
	unbondingPeriod := 21 * 24 * 60 * 60 // 21 days - make configurable?

	// BootstrapContracts
	vaultCodeID, err := StoreContractCode(pCtx, keepers, "/Users/donglieu/mesh-security-sdk/tests/testdata/mesh_vault.wasm.gz")
	require.NoError(t, err)
	proxyCodeID, err := StoreContractCode(pCtx, keepers, "/Users/donglieu/mesh-security-sdk/tests/testdata/mesh_native_staking_proxy.wasm.gz")
	require.NoError(t, err)
	nativeStakingCodeID, err := StoreContractCode(pCtx, keepers, "/Users/donglieu/mesh-security-sdk/tests/testdata/mesh_native_staking.wasm.gz")
	require.NoError(t, err)

	nativeInitMsg := []byte(fmt.Sprintf(`{"denom": %q, "proxy_code_id": %d, "slash_ratio_dsign": %q, "slash_ratio_offline": %q }`, localTokenDenom, proxyCodeID, localSlashRatioDoubleSign, localSlashRatioOffline))
	initMsg := []byte(fmt.Sprintf(`{"denom": %q, "local_staking": {"code_id": %d, "msg": %q}}`, localTokenDenom, nativeStakingCodeID, base64.StdEncoding.EncodeToString(nativeInitMsg)))
	vaultContract, err := InstantiateContract(pCtx, keepers, initMsg, vaultCodeID)
	require.NoError(t, err)

	extStakingCodeID, err := StoreContractCode(pCtx, keepers, "/Users/donglieu/mesh-security-sdk/tests/testdata/mesh_external_staking.wasm.gz")
	require.NoError(t, err)

	initMsg = []byte(fmt.Sprintf(
		`{"remote_contact": {"connection_id":%q, "port_id":%q}, "denom": %q, "vault": %q, "unbonding_period": %d, "rewards_denom": %q, "slash_ratio": { "double_sign": %q, "offline": %q }  }`,
		connId, portID, localTokenDenom, vaultContract.String(), unbondingPeriod, rewardTokenDenom, extSlashRatioDoubleSign, extSlashRatioOffline))
	externalStakingContract, err := InstantiateContract(pCtx, keepers, initMsg, extStakingCodeID)

	k.SetContractWithNativeDenom(pCtx, rewardTokenDenom, externalStakingContract)
	params := k.GetParams(pCtx)
	params.VaultContractAddress = vaultContract.String()
	k.SetParams(pCtx, params)

	testCase := []struct {
		name string
		rq   types.MsgDelegate
	}{
		{
			name: "local stake",
			rq: types.MsgDelegate{
				DelegatorAddress: delAddr.String(),
				ValidatorAddress: valAddr.String(),
				Amount:           sdk.NewInt64Coin(localTokenDenom, 1234),
			},
		},
		{
			name: "stake external",
			rq: types.MsgDelegate{
				DelegatorAddress: delAddr.String(),
				ValidatorAddress: valAddr.String(),
				Amount:           sdk.NewInt64Coin(rewardTokenDenom, 1234),
			},
		},
	}
	for _, j := range testCase {
		t.Run(j.name, func(t *testing.T) {
			m := NewMsgServer(k)
			ctx, _ := pCtx.CacheContext()
			_, err = m.Delegate(ctx, &j.rq)
			require.Error(t, err)
		})
	}
}
