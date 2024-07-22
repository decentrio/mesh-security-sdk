package types

import (
	errorsmod "cosmossdk.io/errors"
)

// meshsecurity sentinel errors
var (
	ErrInvalidVersion        = errorsmod.Register(ModuleName, 2, "invalid meshsecurity version")
	ErrInvalidChannelFlow    = errorsmod.Register(ModuleName, 3, "invalid message sent to channel end")
	ErrDuplicateChannel      = errorsmod.Register(ModuleName, 5, "meshsecurity channel already exists")
	ErrClientNotFound        = errorsmod.Register(ModuleName, 6, "client not found")
	ErrInvalidPacketData     = errorsmod.Register(ModuleName, 1, "invalid CCV packet data")
	ErrInvalidConsumerClient = errorsmod.Register(ModuleName, 7, "ccv channel is not built on correct client")
)
