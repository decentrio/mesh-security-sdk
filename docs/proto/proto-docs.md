<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [osmosis/meshsecurityprovider/v1beta1/delegate.proto](#osmosis/meshsecurityprovider/v1beta1/delegate.proto)
    - [Depositors](#osmosis.meshsecurityprovider.v1beta1.Depositors)
    - [Intermediary](#osmosis.meshsecurityprovider.v1beta1.Intermediary)
  
    - [BondStatus](#osmosis.meshsecurityprovider.v1beta1.BondStatus)
  
- [osmosis/meshsecurityprovider/v1beta1/params.proto](#osmosis/meshsecurityprovider/v1beta1/params.proto)
    - [Params](#osmosis.meshsecurityprovider.v1beta1.Params)
  
- [osmosis/meshsecurityprovider/v1beta1/genesis.proto](#osmosis/meshsecurityprovider/v1beta1/genesis.proto)
    - [GenesisState](#osmosis.meshsecurityprovider.v1beta1.GenesisState)
  
- [osmosis/meshsecurityprovider/v1beta1/query.proto](#osmosis/meshsecurityprovider/v1beta1/query.proto)
    - [QueryParamsRequest](#osmosis.meshsecurityprovider.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#osmosis.meshsecurityprovider.v1beta1.QueryParamsResponse)
  
    - [Query](#osmosis.meshsecurityprovider.v1beta1.Query)
  
- [osmosis/meshsecurityprovider/v1beta1/tx.proto](#osmosis/meshsecurityprovider/v1beta1/tx.proto)
    - [MsgDelegate](#osmosis.meshsecurityprovider.v1beta1.MsgDelegate)
    - [MsgDelegateResponse](#osmosis.meshsecurityprovider.v1beta1.MsgDelegateResponse)
    - [MsgSetConsumerCommissionRate](#osmosis.meshsecurityprovider.v1beta1.MsgSetConsumerCommissionRate)
    - [MsgSetConsumerCommissionRateResponse](#osmosis.meshsecurityprovider.v1beta1.MsgSetConsumerCommissionRateResponse)
    - [MsgUndelegate](#osmosis.meshsecurityprovider.v1beta1.MsgUndelegate)
    - [MsgUndelegateResponse](#osmosis.meshsecurityprovider.v1beta1.MsgUndelegateResponse)
  
    - [Msg](#osmosis.meshsecurityprovider.v1beta1.Msg)
  
- [Scalar Value Types](#scalar-value-types)



<a name="osmosis/meshsecurityprovider/v1beta1/delegate.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## osmosis/meshsecurityprovider/v1beta1/delegate.proto



<a name="osmosis.meshsecurityprovider.v1beta1.Depositors"></a>

### Depositors
Params defines the parameters for the x/meshsecurity module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |
| `tokens` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="osmosis.meshsecurityprovider.v1beta1.Intermediary"></a>

### Intermediary
vault-staker


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `consumer_validator` | [string](#string) |  |  |
| `chain_id` | [string](#string) |  |  |
| `contract_address` | [string](#string) |  |  |
| `jailed` | [bool](#bool) |  |  |
| `tombstoned` | [bool](#bool) |  |  |
| `status` | [BondStatus](#osmosis.meshsecurityprovider.v1beta1.BondStatus) |  |  |
| `token` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |





 <!-- end messages -->


<a name="osmosis.meshsecurityprovider.v1beta1.BondStatus"></a>

### BondStatus
BondStatus is the status of a validator.

| Name | Number | Description |
| ---- | ------ | ----------- |
| BOND_STATUS_UNSPECIFIED | 0 | UNSPECIFIED defines an invalid validator status. |
| BOND_STATUS_UNBONDED | 1 | UNBONDED defines a validator that is not bonded. |
| BOND_STATUS_UNBONDING | 2 | UNBONDING defines a validator that is unbonding. |
| BOND_STATUS_BONDED | 3 | BONDED defines a validator that is bonded. |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="osmosis/meshsecurityprovider/v1beta1/params.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## osmosis/meshsecurityprovider/v1beta1/params.proto



<a name="osmosis.meshsecurityprovider.v1beta1.Params"></a>

### Params
Params defines the parameters for the x/meshsecurity module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `timeout_period` | [uint64](#uint64) |  | TimeoutPeriod has the unit time.Millisecond |
| `vault_contract_address` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="osmosis/meshsecurityprovider/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## osmosis/meshsecurityprovider/v1beta1/genesis.proto



<a name="osmosis.meshsecurityprovider.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines meshsecurity module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#osmosis.meshsecurityprovider.v1beta1.Params) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="osmosis/meshsecurityprovider/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## osmosis/meshsecurityprovider/v1beta1/query.proto



<a name="osmosis.meshsecurityprovider.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the
Query/Params RPC method






<a name="osmosis.meshsecurityprovider.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the
Query/Params RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#osmosis.meshsecurityprovider.v1beta1.Params) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="osmosis.meshsecurityprovider.v1beta1.Query"></a>

### Query
Query provides defines the gRPC querier service

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#osmosis.meshsecurityprovider.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#osmosis.meshsecurityprovider.v1beta1.QueryParamsResponse) | Params queries the parameters of x/meshsecurity module. | GET|/osmosis/meshsecurityprovider/v1beta1/params|

 <!-- end services -->



<a name="osmosis/meshsecurityprovider/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## osmosis/meshsecurityprovider/v1beta1/tx.proto



<a name="osmosis.meshsecurityprovider.v1beta1.MsgDelegate"></a>

### MsgDelegate



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `validator_address` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="osmosis.meshsecurityprovider.v1beta1.MsgDelegateResponse"></a>

### MsgDelegateResponse







<a name="osmosis.meshsecurityprovider.v1beta1.MsgSetConsumerCommissionRate"></a>

### MsgSetConsumerCommissionRate



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `provider_addr` | [string](#string) |  |  |
| `chain_id` | [string](#string) |  |  |
| `rate` | [string](#string) |  |  |






<a name="osmosis.meshsecurityprovider.v1beta1.MsgSetConsumerCommissionRateResponse"></a>

### MsgSetConsumerCommissionRateResponse







<a name="osmosis.meshsecurityprovider.v1beta1.MsgUndelegate"></a>

### MsgUndelegate



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `validator_address` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="osmosis.meshsecurityprovider.v1beta1.MsgUndelegateResponse"></a>

### MsgUndelegateResponse
MsgUndelegateResponse defines the Msg/Undelegate response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="osmosis.meshsecurityprovider.v1beta1.Msg"></a>

### Msg


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Delegate` | [MsgDelegate](#osmosis.meshsecurityprovider.v1beta1.MsgDelegate) | [MsgDelegateResponse](#osmosis.meshsecurityprovider.v1beta1.MsgDelegateResponse) |  | |
| `Undelegate` | [MsgUndelegate](#osmosis.meshsecurityprovider.v1beta1.MsgUndelegate) | [MsgUndelegateResponse](#osmosis.meshsecurityprovider.v1beta1.MsgUndelegateResponse) |  | |
| `SetConsumerCommissionRate` | [MsgSetConsumerCommissionRate](#osmosis.meshsecurityprovider.v1beta1.MsgSetConsumerCommissionRate) | [MsgSetConsumerCommissionRateResponse](#osmosis.meshsecurityprovider.v1beta1.MsgSetConsumerCommissionRateResponse) |  | |

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

