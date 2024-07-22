package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/osmosis-labs/mesh-security-sdk/x/meshsecurity/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	k *Keeper
}

// NewMsgServer constructor
func NewMsgServer(k *Keeper) *msgServer {
	return &msgServer{k: k}
}

// SetVirtualStakingMaxCap sets a new max cap limit for virtual staking
func (m msgServer) SetVirtualStakingMaxCap(goCtx context.Context, req *types.MsgSetVirtualStakingMaxCap) (*types.MsgSetVirtualStakingMaxCapResponse, error) {
	if err := req.ValidateBasic(); err != nil {
		return nil, err
	}

	if authority := m.k.GetAuthority(); authority != req.Authority {
		return nil, govtypes.ErrInvalidSigner.Wrapf("invalid authority; expected %s, got %s", authority, req.Authority)
	}

	acc, err := sdk.AccAddressFromBech32(req.Contract)
	if err != nil {
		return nil, errorsmod.Wrap(err, "contract")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := m.k.SetMaxCapLimit(ctx, acc, req.MaxCap); err != nil {
		return nil, err
	}
	if !m.k.HasScheduledTask(ctx, types.SchedulerTaskHandleEpoch, acc, true) {
		if err := m.k.ScheduleRegularRebalanceTask(ctx, acc); err != nil {
			return nil, errorsmod.Wrap(err, "schedule regular rebalance task")
		}
		return &types.MsgSetVirtualStakingMaxCapResponse{}, nil
	}
	if req.MaxCap.IsZero() {
		// no need to run regular rebalances with a new limit of 0
		if err := m.k.DeleteAllScheduledTasks(ctx, types.SchedulerTaskHandleEpoch, acc); err != nil {
			return nil, err
		}
	}

	// schedule last rebalance callback to let the contract do undelegates and housekeeping
	if err := m.k.ScheduleOneShotTask(ctx, types.SchedulerTaskHandleEpoch, acc, uint64(ctx.BlockHeight())); err != nil {
		return nil, errorsmod.Wrap(err, "schedule one shot rebalance task")
	}
	return &types.MsgSetVirtualStakingMaxCapResponse{}, nil
}

func (m msgServer) Delegate(goCtx context.Context, msg *types.MsgDelegate) (*types.MsgDelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := msg.Amount.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to delegate; Validate fail")
	}
	denomDelegate := msg.Amount.Denom

	vaultAdress := sdk.AccAddress(m.k.GetParams(ctx).GetVaultContractAddress())
	delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	balance := m.k.bank.GetBalance(ctx, delegatorAddress, denomDelegate)

	if balance.IsLT(msg.Amount) {
		return nil, fmt.Errorf("failed to delegate; %s is smaller than %s", balance, msg.Amount)
	}

	if err := m.k.bank.SendCoins(ctx, delegatorAddress, vaultAdress, sdk.NewCoins([]sdk.Coin{msg.Amount}...)); err != nil {
		return nil, err
	}

	bondDenom := m.k.Staking.BondDenom(ctx)
	if denomDelegate == bondDenom {
		// local stake
		err = m.k.LocalStake(ctx, msg.Amount, msg.ValidatorAddress)
		if err != nil {
			// if fail refunds
			m.k.bank.SendCoins(ctx, vaultAdress, delegatorAddress, sdk.NewCoins([]sdk.Coin{msg.Amount}...))
			return nil, err
		}
	} else {
		// remote stake
		err := m.k.RemoteStake(ctx, denomDelegate, msg.Amount)
		if err != nil {
			// if fail refunds
			m.k.bank.SendCoins(ctx, vaultAdress, delegatorAddress, sdk.NewCoins([]sdk.Coin{msg.Amount}...))
			return nil, err
		}
	}

	depositors, found := m.k.GetDepositors(ctx, msg.DelegatorAddress)
	if !found {
		depositors = types.NewDepositors(msg.DelegatorAddress, []sdk.Coin{msg.Amount})
	} else {
		depositors.Tokens = depositors.Tokens.Add(msg.Amount)
	}
	err = m.k.SetDepositors(ctx, depositors)
	if err != nil {
		return nil, err
	}

	return &types.MsgDelegateResponse{}, nil
}

func (m msgServer) Undelegate(goCtx context.Context, msg *types.MsgUndelegate) (*types.MsgUndelegateResponse, error) {
	return &types.MsgUndelegateResponse{}, nil
}
