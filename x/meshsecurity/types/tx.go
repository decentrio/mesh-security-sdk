package types

import (
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgSetVirtualStakingMaxCap) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners returns the expected signers for MsgSoftwareUpgrade.
func (msg MsgSetVirtualStakingMaxCap) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic validate basic constraints
func (msg MsgSetVirtualStakingMaxCap) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", err)
	}
	if _, err := sdk.AccAddressFromBech32(msg.Contract); err != nil {
		return errorsmod.Wrap(err, "contract")
	}
	if err := msg.MaxCap.Validate(); err != nil {
		return errorsmod.Wrap(err, "max cap")
	}
	return nil
}

func NewMsgDelegate(delAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin) *MsgDelegate {
	return &MsgDelegate{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: valAddr.String(),
		Amount:           amount,
	}
}

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgDelegate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners returns the expected signers for MsgSoftwareUpgrade.
func (msg MsgDelegate) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	return []sdk.AccAddress{addr}
}

// ValidateBasic validate basic constraints
func (msg MsgDelegate) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.DelegatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorAddress); err != nil {
		return errorsmod.Wrap(err, "contract")
	}
	if err := msg.Amount.Validate(); err != nil {
		return errorsmod.Wrap(err, "max cap")
	}
	return nil
}

func NewMsgUndelegate(delAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin) *MsgUndelegate {
	return &MsgUndelegate{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: valAddr.String(),
		Amount:           amount,
	}
}

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgUndelegate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners returns the expected signers for MsgSoftwareUpgrade.
func (msg MsgUndelegate) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	return []sdk.AccAddress{addr}
}

// ValidateBasic validate basic constraints
func (msg MsgUndelegate) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.DelegatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", err)
	}
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorAddress); err != nil {
		return errorsmod.Wrap(err, "contract")
	}
	if err := msg.Amount.Validate(); err != nil {
		return errorsmod.Wrap(err, "max cap")
	}
	return nil
}
