package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	cptypes "github.com/osmosis-labs/mesh-security-sdk/x/types"
)

const (
	// ModuleName defines the module name.
	ModuleName = "meshsecurity-consummer"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "memory:meshsecurity-consummer"

	// RouterKey is the message route
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key.
	QuerierRoute = ModuleName
)

var (
	// ParamsKey is the prefix for the module parameters
	ParamsKey                     = []byte{0x10}
	MaxCapLimitKeyPrefix          = []byte{0x2}
	TotalDelegatedAmountKeyPrefix = []byte{0x3}
	SchedulerKeyPrefix            = []byte{0x4}

	PipedValsetPrefix             = []byte{0x5}
	ProviderChannelByteKey        = []byte{0x6}
	DepositorsKeyPrefix           = []byte{0x7}
	IntermediaryKeyPrefix         = []byte{0x8}
	ContractWithNativeDenomPrefix = []byte{0x9}
	ChainToChannelBytePrefix      = []byte{0x11}
	ChannelToChainBytePrefix      = []byte{0x12}
	ChainToClientBytePrefix       = []byte{0x13}
	InitChainHeightBytePrefix     = []byte{0x14}
)

// BuildMaxCapLimitKey build max cap limit store key
func BuildMaxCapLimitKey(contractAddr sdk.AccAddress) []byte {
	return append(MaxCapLimitKeyPrefix, contractAddr.Bytes()...)
}

// BuildTotalDelegatedAmountKey build delegated amount store key for given contract
func BuildTotalDelegatedAmountKey(contractAddr sdk.AccAddress) []byte {
	return append(TotalDelegatedAmountKeyPrefix, contractAddr.Bytes()...)
}

// BuildSchedulerTypeKeyPrefix internal scheduler store key
func BuildSchedulerTypeKeyPrefix(tp SchedulerTaskType) ([]byte, error) {
	if tp == SchedulerTaskUndefined {
		return nil, ErrInvalid.Wrapf("scheduler type: %x", tp)
	}
	return append(SchedulerKeyPrefix, byte(tp)), nil
}

// BuildSchedulerHeightKeyPrefix build store key prefix
func BuildSchedulerHeightKeyPrefix(tp SchedulerTaskType, blockHeight uint64) ([]byte, error) {
	prefix, err := BuildSchedulerTypeKeyPrefix(tp)
	if err != nil {
		return nil, err
	}
	return append(prefix, sdk.Uint64ToBigEndian(blockHeight)...), nil
}

// BuildSchedulerContractKey build store key
func BuildSchedulerContractKey(tp SchedulerTaskType, blockHeight uint64, contractAddr sdk.AccAddress) ([]byte, error) {
	prefix, err := BuildSchedulerHeightKeyPrefix(tp, blockHeight)
	if err != nil {
		return nil, err
	}
	return append(prefix, contractAddr.Bytes()...), nil
}

// BuildPipedValsetOpKey build store key for the temporary valset operation store
func BuildPipedValsetOpKey(op cptypes.PipedValsetOperation, val sdk.ValAddress) []byte {
	if op == cptypes.PipedValsetOperation_UNSPECIFIED {
		panic("empty operation")
	}
	k := append(append(PipedValsetPrefix, byte(op)), val...)
	// return
	return k
}

// ProviderChannelKey returns the key for storing channelID of the provider chain
func ProviderChannelKey() []byte {
	return ProviderChannelByteKey
}

func DepositorsKey(del string) []byte {
	return append(DepositorsKeyPrefix, []byte(del)...)
}

func IntermediaryKey(denom string) []byte {
	return append(IntermediaryKeyPrefix, []byte(denom)...)
}

func ContractWithNativeDenomKey(denom string) []byte {
	return append(ContractWithNativeDenomPrefix, []byte(denom)...)
}

func ChainToChannelKey(chainID string) []byte {
	return append(ChainToChannelBytePrefix, []byte(chainID)...)
}

func ChannelToChainKey(channelID string) []byte {
	return append(ChannelToChainBytePrefix, []byte(channelID)...)
}

func ChainToClientKey(chainID string) []byte {
	return append(ChainToClientBytePrefix, []byte(chainID)...)
}

func InitChainHeightKey(chainID string) []byte {
	return append(InitChainHeightBytePrefix, []byte(chainID)...)
}
