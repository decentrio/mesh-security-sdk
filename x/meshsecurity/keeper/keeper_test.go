package keeper

import (
	"testing"

	"github.com/cometbft/cometbft/libs/rand"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/mesh-security-sdk/x/meshsecurity/types"
)

func TestHasMaxCapLimit(t *testing.T) {
	pCtx, keepers := CreateDefaultTestInput(t)
	k := keepers.MeshKeeper
	myContractAddr := sdk.AccAddress(rand.Bytes(32))

	specs := map[string]struct {
		setup     func(ctx sdk.Context)
		expResult bool
	}{
		"limit set": {
			setup: func(ctx sdk.Context) {
				err := k.SetMaxCapLimit(ctx, myContractAddr, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1))
				require.NoError(t, err)
			},
			expResult: true,
		},
		"limit with empty amount set": {
			setup: func(ctx sdk.Context) {
				err := k.SetMaxCapLimit(ctx, myContractAddr, sdk.NewInt64Coin(sdk.DefaultBondDenom, 0))
				require.NoError(t, err)
			},
			expResult: true,
		},
		"limit not set": {
			setup:     func(ctx sdk.Context) {}, // noop
			expResult: false,
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			ctx, _ := pCtx.CacheContext()
			spec.setup(ctx)
			got := k.HasMaxCapLimit(ctx, myContractAddr)
			assert.Equal(t, spec.expResult, got)
		})
	}
}

func TestSetMaxCapLimit(t *testing.T) {
	pCtx, keepers := CreateDefaultTestInput(t)
	k := keepers.MeshKeeper
	var (
		myContractAddr = sdk.AccAddress(rand.Bytes(32))
		oneStakeCoin   = sdk.NewInt64Coin(sdk.DefaultBondDenom, 1)
		zeroStakeCoin  = sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)
	)

	specs := map[string]struct {
		setup  func(ctx sdk.Context) sdk.Coin
		expErr bool
	}{
		"all good": {
			setup: func(_ sdk.Context) sdk.Coin {
				return oneStakeCoin
			},
		},
		"zero amount allowed": {
			setup: func(_ sdk.Context) sdk.Coin {
				return zeroStakeCoin
			},
		},
		"overwrite existing value": {
			setup: func(ctx sdk.Context) sdk.Coin {
				err := k.SetMaxCapLimit(ctx, myContractAddr, oneStakeCoin)
				require.NoError(t, err)
				return oneStakeCoin.AddAmount(math.NewInt(1))
			},
		},
		"within total contracts max limit": {
			setup: func(ctx sdk.Context) sdk.Coin {
				p := k.GetParams(ctx)
				p.TotalContractsMaxCap = oneStakeCoin
				require.NoError(t, k.SetParams(ctx, p))
				return oneStakeCoin
			},
		},
		"non staking denom rejected": {
			setup: func(_ sdk.Context) sdk.Coin {
				return sdk.NewInt64Coin("NON", 1)
			},
			expErr: true,
		},
		"total contracts max exceeded - with other contract": {
			setup: func(ctx sdk.Context) sdk.Coin {
				p := k.GetParams(ctx)
				p.TotalContractsMaxCap = oneStakeCoin
				require.NoError(t, k.SetParams(ctx, p))
				require.NoError(t, k.SetMaxCapLimit(ctx, rand.Bytes(32), oneStakeCoin))
				return oneStakeCoin
			},
			expErr: true,
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			ctx, _ := pCtx.CacheContext()
			limitAmount := spec.setup(ctx)

			// when
			gotErr := k.SetMaxCapLimit(ctx, myContractAddr, limitAmount)
			// then
			if spec.expErr {
				require.Error(t, gotErr)
				limitAmount = zeroStakeCoin
			} else {
				require.NoError(t, gotErr)
			}
			assert.Equal(t, limitAmount, k.GetMaxCapLimit(ctx, myContractAddr))
		})
	}
}

func TestStoreDelegator(t *testing.T) {
	pCtx, keepers := CreateDefaultTestInput(t)
	k := keepers.MeshKeeper

	d1 := types.Depositors{
		Address: sdk.AccAddress(rand.Bytes(32)).String(),
		Tokens:  sdk.NewCoins([]sdk.Coin{sdk.NewCoin("osmo", sdk.NewInt(123))}...),
	}
	err := k.SetDepositors(pCtx, d1)
	require.NoError(t, err)

	d2, f := k.GetDepositors(pCtx, d1.Address)
	require.True(t, f)
	require.Equal(t, d1, d2)

	d2.Tokens = d2.Tokens.Add(sdk.NewCoin("juno", sdk.NewInt(10000)))
	d2.Tokens = d2.Tokens.Add(sdk.NewCoin("juno", sdk.NewInt(10000)))
	err = k.SetDepositors(pCtx, d2)
	require.NoError(t, err)

	d3, f := k.GetDepositors(pCtx, d1.Address)
	require.True(t, f)
	require.Equal(t, sdk.NewInt(20000), d3.Tokens.AmountOf("juno"))
}

func TestStoreIntermediary(t *testing.T) {
	pCtx, keepers := CreateDefaultTestInput(t)
	k := keepers.MeshKeeper

	coin := sdk.NewCoin("osmo", sdk.NewInt(123))
	e1 := types.Intermediary{
		ConsumerValidator: sdk.ValAddress(rand.Bytes(32)).String(),
		ChainId:           "osmo-1",
		ContractAddress:   sdk.AccAddress(rand.Bytes(32)).String(),
		Jailed:            false,
		Tombstoned:        false,
		Status:            types.Bonded,
		Token:             &coin,
	}
	err := k.SetIntermediary(pCtx, e1)
	require.NoError(t, err)

	e2, f := k.GetIntermediary(pCtx, e1.Token.Denom)
	require.True(t, f)
	require.Equal(t, e1, e2)
}
