# NoncurrentVersionExpiration

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**NoncurrentDays** | Pointer to **int32** | Specifies the number of days an object is noncurrent before IONOS S3 Object Storage can perform the associated operation. | [optional] |

## Methods

### NewNoncurrentVersionExpiration

`func NewNoncurrentVersionExpiration() *NoncurrentVersionExpiration`

NewNoncurrentVersionExpiration instantiates a new NoncurrentVersionExpiration object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNoncurrentVersionExpirationWithDefaults

`func NewNoncurrentVersionExpirationWithDefaults() *NoncurrentVersionExpiration`

NewNoncurrentVersionExpirationWithDefaults instantiates a new NoncurrentVersionExpiration object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNoncurrentDays

`func (o *NoncurrentVersionExpiration) GetNoncurrentDays() int32`

GetNoncurrentDays returns the NoncurrentDays field if non-nil, zero value otherwise.

### GetNoncurrentDaysOk

`func (o *NoncurrentVersionExpiration) GetNoncurrentDaysOk() (*int32, bool)`

GetNoncurrentDaysOk returns a tuple with the NoncurrentDays field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNoncurrentDays

`func (o *NoncurrentVersionExpiration) SetNoncurrentDays(v int32)`

SetNoncurrentDays sets NoncurrentDays field to given value.

### HasNoncurrentDays

`func (o *NoncurrentVersionExpiration) HasNoncurrentDays() bool`

HasNoncurrentDays returns a boolean if a field has been set.


