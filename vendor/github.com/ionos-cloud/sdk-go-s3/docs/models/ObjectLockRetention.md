# ObjectLockRetention

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Mode** | Pointer to **string** | Indicates the Retention mode for the specified object. | [optional] |
|**RetainUntilDate** | Pointer to **string** | The date on which this Object Lock Retention will expire. | [optional] |

## Methods

### NewObjectLockRetention

`func NewObjectLockRetention() *ObjectLockRetention`

NewObjectLockRetention instantiates a new ObjectLockRetention object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewObjectLockRetentionWithDefaults

`func NewObjectLockRetentionWithDefaults() *ObjectLockRetention`

NewObjectLockRetentionWithDefaults instantiates a new ObjectLockRetention object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMode

`func (o *ObjectLockRetention) GetMode() string`

GetMode returns the Mode field if non-nil, zero value otherwise.

### GetModeOk

`func (o *ObjectLockRetention) GetModeOk() (*string, bool)`

GetModeOk returns a tuple with the Mode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMode

`func (o *ObjectLockRetention) SetMode(v string)`

SetMode sets Mode field to given value.

### HasMode

`func (o *ObjectLockRetention) HasMode() bool`

HasMode returns a boolean if a field has been set.

### GetRetainUntilDate

`func (o *ObjectLockRetention) GetRetainUntilDate() string`

GetRetainUntilDate returns the RetainUntilDate field if non-nil, zero value otherwise.

### GetRetainUntilDateOk

`func (o *ObjectLockRetention) GetRetainUntilDateOk() (*string, bool)`

GetRetainUntilDateOk returns a tuple with the RetainUntilDate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRetainUntilDate

`func (o *ObjectLockRetention) SetRetainUntilDate(v string)`

SetRetainUntilDate sets RetainUntilDate field to given value.

### HasRetainUntilDate

`func (o *ObjectLockRetention) HasRetainUntilDate() bool`

HasRetainUntilDate returns a boolean if a field has been set.


