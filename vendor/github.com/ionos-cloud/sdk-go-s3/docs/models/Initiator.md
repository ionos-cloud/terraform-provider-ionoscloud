# Initiator

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**ID** | Pointer to **int32** | Container for the Contract Number of the owner. | [optional] |
|**DisplayName** | Pointer to **string** | Container for the display name of the owner. | [optional] |

## Methods

### NewInitiator

`func NewInitiator() *Initiator`

NewInitiator instantiates a new Initiator object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInitiatorWithDefaults

`func NewInitiatorWithDefaults() *Initiator`

NewInitiatorWithDefaults instantiates a new Initiator object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetID

`func (o *Initiator) GetID() int32`

GetID returns the ID field if non-nil, zero value otherwise.

### GetIDOk

`func (o *Initiator) GetIDOk() (*int32, bool)`

GetIDOk returns a tuple with the ID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetID

`func (o *Initiator) SetID(v int32)`

SetID sets ID field to given value.

### HasID

`func (o *Initiator) HasID() bool`

HasID returns a boolean if a field has been set.

### GetDisplayName

`func (o *Initiator) GetDisplayName() string`

GetDisplayName returns the DisplayName field if non-nil, zero value otherwise.

### GetDisplayNameOk

`func (o *Initiator) GetDisplayNameOk() (*string, bool)`

GetDisplayNameOk returns a tuple with the DisplayName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisplayName

`func (o *Initiator) SetDisplayName(v string)`

SetDisplayName sets DisplayName field to given value.

### HasDisplayName

`func (o *Initiator) HasDisplayName() bool`

HasDisplayName returns a boolean if a field has been set.


