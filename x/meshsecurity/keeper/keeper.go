package keeper

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/cometbft/cometbft/libs/log"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"

	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	conntypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	ibchost "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctmtypes "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/osmosis-labs/mesh-security-sdk/x/meshsecurity/types"
	cptypes "github.com/osmosis-labs/mesh-security-sdk/x/types"
)

// Option is an extension point to instantiate keeper with non default values
type Option interface {
	apply(*Keeper)
}

type Keeper struct {
	storeKey storetypes.StoreKey
	memKey   storetypes.StoreKey
	cdc      codec.Codec
	bank     types.XBankKeeper
	Staking  types.XStakingKeeper
	wasm     types.WasmKeeper

	scopedKeeper     types.ScopedKeeper
	channelKeeper    types.ChannelKeeper
	connectionKeeper types.ConnectionKeeper
	clientKeeper     types.ClientKeeper
	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

// NewKeeper constructor with vanilla sdk keepers
func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	memoryStoreKey storetypes.StoreKey,
	bank types.SDKBankKeeper,
	staking types.SDKStakingKeeper,
	wasm types.WasmKeeper,
	scopedKeeper types.ScopedKeeper,
	channelKeeper types.ChannelKeeper,
	connectionKeeper types.ConnectionKeeper,
	clientKeeper types.ClientKeeper,
	authority string,
	opts ...Option,
) *Keeper {
	return NewKeeperX(cdc, storeKey, memoryStoreKey, NewBankKeeperAdapter(bank), NewStakingKeeperAdapter(staking, bank), wasm, scopedKeeper, channelKeeper, connectionKeeper, clientKeeper, authority, opts...)
}

// NewKeeperX constructor with extended Osmosis SDK keepers
func NewKeeperX(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	memoryStoreKey storetypes.StoreKey,
	bank types.XBankKeeper,
	staking types.XStakingKeeper,
	wasm types.WasmKeeper,
	scopedKeeper types.ScopedKeeper,
	channelKeeper types.ChannelKeeper,
	connectionKeeper types.ConnectionKeeper,
	clientKeeper types.ClientKeeper,
	authority string,
	opts ...Option,
) *Keeper {
	k := &Keeper{
		storeKey:         storeKey,
		memKey:           memoryStoreKey,
		cdc:              cdc,
		bank:             bank,
		Staking:          staking,
		wasm:             wasm,
		scopedKeeper:     scopedKeeper,
		channelKeeper:    channelKeeper,
		connectionKeeper: connectionKeeper,
		clientKeeper:     clientKeeper,
		authority:        authority,
	}
	for _, o := range opts {
		o.apply(k)
	}

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// HasMaxCapLimit returns true when any max cap limit was set. The amount is not taken into account for the result.
// A 0 value would be true as well.
func (k Keeper) HasMaxCapLimit(ctx sdk.Context, actor sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.BuildMaxCapLimitKey(actor))
}

func (k Keeper) HasContractVault(ctx sdk.Context, contractAddr sdk.AccAddress) bool {
	return bytes.Equal(k.GetParams(ctx).GetVaultContractAddress(), contractAddr)
}

// GetMaxCapLimit the cap limit is set per consumer contract. Different providers can have different limits
// Returns zero amount when no limit is stored.
func (k Keeper) GetMaxCapLimit(ctx sdk.Context, actor sdk.AccAddress) sdk.Coin {
	return sdk.NewCoin(k.Staking.BondDenom(ctx), k.mustLoadInt(ctx, k.storeKey, types.BuildMaxCapLimitKey(actor)))
}

// SetMaxCapLimit stores the max cap limit for the given contract address.
// Any existing limit for this contract will be overwritten
func (k Keeper) SetMaxCapLimit(ctx sdk.Context, contract sdk.AccAddress, newAmount sdk.Coin) error {
	if k.Staking.BondDenom(ctx) != newAmount.Denom {
		return sdkerrors.ErrInvalidCoins
	}
	// ensure that the total max cap amount for all contracts is not exceeded
	total := math.ZeroInt()
	k.IterateMaxCapLimit(ctx, func(addr sdk.AccAddress, m math.Int) bool {
		if !addr.Equals(contract) {
			total = total.Add(m)
		}
		return false
	})
	totalMaxCap := k.GetTotalContractsMaxCap(ctx)
	if total.Add(newAmount.Amount).GT(totalMaxCap.Amount) {
		return types.ErrInvalid.Wrapf("amount exceeds total available max cap (used %s of %s)", total, totalMaxCap)
	}
	// persist
	store := ctx.KVStore(k.storeKey)
	bz, err := newAmount.Amount.Marshal()
	if err != nil { // always nil
		return errorsmod.Wrap(err, "marshal amount")
	}
	store.Set(types.BuildMaxCapLimitKey(contract), bz)

	types.EmitMaxCapLimitUpdatedEvent(ctx, contract, newAmount)
	return nil
}

// GetTotalDelegated returns the total amount delegated by the given consumer contract.
// This amount can be 0 is never negative.
func (k Keeper) GetTotalDelegated(ctx sdk.Context, actor sdk.AccAddress) sdk.Coin {
	v := k.mustLoadInt(ctx, k.storeKey, types.BuildTotalDelegatedAmountKey(actor))
	if v.IsNegative() {
		v = math.ZeroInt()
	}
	return sdk.NewCoin(k.Staking.BondDenom(ctx), v)
}

// internal setter. must only be used with bonding token denom or panics
func (k Keeper) setTotalDelegated(ctx sdk.Context, actor sdk.AccAddress, newAmount sdk.Coin) {
	if k.Staking.BondDenom(ctx) != newAmount.Denom {
		panic(sdkerrors.ErrInvalidCoins.Wrapf("not a staking denom: %s", newAmount.Denom))
	}

	store := ctx.KVStore(k.storeKey)
	bz, err := newAmount.Amount.Marshal()
	if err != nil { // always nil
		panic(err)
	}
	store.Set(types.BuildTotalDelegatedAmountKey(actor), bz)
}

// helper to deserialize a math.Int from store. Returns zero when key does not exist.
// Panics when Unmarshal fails
func (k Keeper) mustLoadInt(ctx sdk.Context, storeKey storetypes.StoreKey, key []byte) math.Int {
	store := ctx.KVStore(storeKey)
	bz := store.Get(key)
	if bz == nil {
		return sdk.ZeroInt()
	}
	var r math.Int
	if err := r.Unmarshal(bz); err != nil {
		panic(err)
	}
	return r
}

// IterateMaxCapLimit iterate over contract addresses with max cap limit set
// Callback can return true to stop early
func (k Keeper) IterateMaxCapLimit(ctx sdk.Context, cb func(sdk.AccAddress, math.Int) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.MaxCapLimitKeyPrefix)
	iter := prefixStore.Iterator(nil, nil)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var r math.Int
		if err := r.Unmarshal(iter.Value()); err != nil {
			panic(err)
		}
		// cb returns true to stop early
		if cb(iter.Key(), r) {
			return
		}
	}
}

// ModuleLogger returns logger with module attribute
func ModuleLogger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetProviderChannel gets the channelID for the channel to the provider.
func (k Keeper) GetProviderChannel(ctx sdk.Context) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	channelIdBytes := store.Get(types.ProviderChannelKey())
	if len(channelIdBytes) == 0 {
		return "", false
	}
	return string(channelIdBytes), true
}

// SetProviderChannel sets the channelID for the channel to the provider.
func (k Keeper) SetProviderChannel(ctx sdk.Context, channelID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ProviderChannelKey(), []byte(channelID))
}

func (k Keeper) ChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	capName := host.ChannelCapabilityPath(portID, channelID)
	chanCap, ok := k.scopedKeeper.GetCapability(ctx, capName)
	if !ok {
		return errorsmod.Wrapf(channeltypes.ErrChannelCapabilityNotFound, "could not retrieve channel capability at: %s", capName)
	}
	return k.channelKeeper.ChanCloseInit(ctx, portID, channelID, chanCap)
}

func (k Keeper) SetDepositors(ctx sdk.Context, del types.Depositors) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := types.ModuleCdc.Marshal(&del)
	if err != nil {
		err = fmt.Errorf("external staker marshalling failed: %s", err)
		ModuleLogger(ctx).Error(err.Error())
		return err
	}
	store.Set(types.DepositorsKey(del.Address), bz)

	return nil
}

func (k Keeper) GetDepositors(ctx sdk.Context, del string) (types.Depositors, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DepositorsKey(del))
	if bz == nil {
		return types.Depositors{}, false
	}
	var r types.Depositors
	if err := r.Unmarshal(bz); err != nil {
		panic(err)
	}
	return r, true
}

func (k Keeper) iterateDepositors(ctx sdk.Context, cb func(depositors types.Depositors) bool) error {
	pStore := prefix.NewStore(ctx.KVStore(k.memKey), types.DepositorsKeyPrefix)
	iter := pStore.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		var r types.Depositors
		if err := r.Unmarshal(iter.Value()); err != nil {
			panic(err)
		}
		if cb(r) {
			break
		}
	}
	return iter.Close()
}

func (k Keeper) SetIntermediary(ctx sdk.Context, inter types.Intermediary) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := types.ModuleCdc.Marshal(&inter)
	if err != nil {
		err = fmt.Errorf("external staker marshalling failed: %s", err)
		ModuleLogger(ctx).Error(err.Error())
		return err
	}
	store.Set(types.IntermediaryKey(inter.Token.Denom), bz)

	return nil
}

func (k Keeper) GetIntermediary(ctx sdk.Context, denom string) (types.Intermediary, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.IntermediaryKey(denom))
	if bz == nil {
		return types.Intermediary{}, false
	}
	var r types.Intermediary
	if err := r.Unmarshal(bz); err != nil {
		panic(err)
	}
	return r, true
}

func (k Keeper) GetContractWithNativeDenom(ctx sdk.Context, denom string) sdk.AccAddress {
	var contractAddr sdk.AccAddress

	store := ctx.KVStore(k.storeKey)
	contractAddr = store.Get(types.ContractWithNativeDenomKey(denom))
	return contractAddr
}

func (k Keeper) SetContractWithNativeDenom(ctx sdk.Context, denom string, contractAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ContractWithNativeDenomKey(denom), contractAddr)
}

func (k Keeper) GetChainToChannel(ctx sdk.Context, chainID string) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ChainToChannelKey(chainID))
	if bz == nil {
		return "", false
	}
	return string(bz), true
}

func (k Keeper) SetChainToChannel(ctx sdk.Context, chainID, channelID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ChainToChannelKey(chainID), []byte(channelID))
}

func (k Keeper) SetChannelToChain(ctx sdk.Context, channelID, chainID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ChannelToChainKey(channelID), []byte(chainID))
}

func (k Keeper) GetChannelToChain(ctx sdk.Context, channelID string) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ChannelToChainKey(channelID))
	if bz == nil {
		return "", false
	}
	return string(bz), true
}

func (k Keeper) GetConsumerClientId(ctx sdk.Context, chainID string) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	clientIdBytes := store.Get(types.ChainToClientKey(chainID))
	if clientIdBytes == nil {
		return "", false
	}
	return string(clientIdBytes), true
}

func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}

func (k Keeper) VerifyConsumerChain(ctx sdk.Context, channelID string, connectionHops []string) error {
	if len(connectionHops) != 1 {
		return errorsmod.Wrap(channeltypes.ErrTooManyConnectionHops, "must have direct connection to provider chain")
	}
	connectionID := connectionHops[0]
	clientID, tmClient, err := k.getUnderlyingClient(ctx, connectionID)
	if err != nil {
		return err
	}
	ccvClientId, found := k.GetConsumerClientId(ctx, tmClient.ChainId)
	if !found {
		return errorsmod.Wrapf(cptypes.ErrClientNotFound, "cannot find client for consumer chain %s", tmClient.ChainId)
	}
	if ccvClientId != clientID {
		return errorsmod.Wrapf(cptypes.ErrInvalidConsumerClient, "CCV channel must be built on top of CCV client. expected %s, got %s", ccvClientId, clientID)
	}

	// Verify that there isn't already a CCV channel for the consumer chain
	if prevChannel, ok := k.GetChainToChannel(ctx, tmClient.ChainId); ok {
		return errorsmod.Wrapf(cptypes.ErrDuplicateChannel, "CCV channel with ID: %s already created for consumer chain %s", prevChannel, tmClient.ChainId)
	}
	return nil
}

// Retrieves the underlying client state corresponding to a connection ID.
func (k Keeper) getUnderlyingClient(ctx sdk.Context, connectionID string) (
	clientID string, tmClient *ibctmtypes.ClientState, err error,
) {
	conn, ok := k.connectionKeeper.GetConnection(ctx, connectionID)
	if !ok {
		return "", nil, errorsmod.Wrapf(conntypes.ErrConnectionNotFound,
			"connection not found for connection ID: %s", connectionID)
	}
	clientID = conn.ClientId
	clientState, ok := k.clientKeeper.GetClientState(ctx, clientID)
	if !ok {
		return "", nil, errorsmod.Wrapf(clienttypes.ErrClientNotFound,
			"client not found for client ID: %s", conn.ClientId)
	}
	tmClient, ok = clientState.(*ibctmtypes.ClientState)
	if !ok {
		return "", nil, errorsmod.Wrapf(clienttypes.ErrInvalidClientType,
			"invalid client type. expected %s, got %s", ibchost.Tendermint, clientState.ClientType())
	}
	return clientID, tmClient, nil
}

func (k Keeper) GetInitChainHeight(ctx sdk.Context, chainID string) (uint64, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.InitChainHeightKey(chainID))
	if bz == nil {
		return 0, false
	}

	return binary.BigEndian.Uint64(bz), true
}
