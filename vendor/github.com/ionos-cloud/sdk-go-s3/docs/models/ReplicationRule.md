# ReplicationRule

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**ID** | Pointer to **int32** | Container for the Contract Number of the owner. | [optional] |
|**Prefix** | Pointer to **string** | An object key name prefix that identifies the subset of objects to which the rule applies. Replace the Object keys containing special characters, such as carriage returns, when using XML requests.  | [optional] |
|**Status** | **string** | Specifies whether the rule is enabled. | |
|**Destination** | [**Destination**](Destination.md) |  | |

## Methods

### NewReplicationRule

`func NewReplicationRule(status string, destination Destination, ) *ReplicationRule`

NewReplicationRule instantiates a new ReplicationRule object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewReplicationRuleWithDefaults

`func NewReplicationRuleWithDefaults() *ReplicationRule`

NewReplicationRuleWithDefaults instantiates a new ReplicationRule object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetID

`func (o *ReplicationRule) GetID() int32`

GetID returns the ID field if non-nil, zero value otherwise.

### GetIDOk

`func (o *ReplicationRule) GetIDOk() (*int32, bool)`

GetIDOk returns a tuple with the ID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetID

`func (o *ReplicationRule) SetID(v int32)`

SetID sets ID field to given value.

### HasID

`func (o *ReplicationRule) HasID() bool`

HasID returns a boolean if a field has been set.

### GetPrefix

`func (o *ReplicationRule) GetPrefix() string`

GetPrefix returns the Prefix field if non-nil, zero value otherwise.

### GetPrefixOk

`func (o *ReplicationRule) GetPrefixOk() (*string, bool)`

GetPrefixOk returns a tuple with the Prefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrefix

`func (o *ReplicationRule) SetPrefix(v string)`

SetPrefix sets Prefix field to given value.

### HasPrefix

`func (o *ReplicationRule) HasPrefix() bool`

HasPrefix returns a boolean if a field has been set.

### GetStatus

`func (o *ReplicationRule) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ReplicationRule) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ReplicationRule) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetDestination

`func (o *ReplicationRule) GetDestination() Destination`

GetDestination returns the Destination field if non-nil, zero value otherwise.

### GetDestinationOk

`func (o *ReplicationRule) GetDestinationOk() (*Destination, bool)`

GetDestinationOk returns a tuple with the Destination field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDestination

`func (o *ReplicationRule) SetDestination(v Destination)`

SetDestination sets Destination field to given value.



