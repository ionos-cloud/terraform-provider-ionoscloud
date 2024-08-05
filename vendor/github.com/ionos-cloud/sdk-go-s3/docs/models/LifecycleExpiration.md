# LifecycleExpiration

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Date** | Pointer to [**time.Time**](time.Time.md) | Indicates at what date the object is to be moved or deleted. Should be in GMT ISO 8601 Format. | [optional] |
|**Days** | Pointer to **int32** | Indicates the lifetime, in days, of the objects that are subject to the rule. The value must be a non-zero positive integer. | [optional] |
|**ExpiredObjectDeleteMarker** | Pointer to **bool** | Indicates whether IONOS S3 Object Storage will remove a delete marker with no noncurrent versions. If set to true, the delete marker will be expired; if set to false the policy takes no operation. This cannot be specified with Days or Date in a Lifecycle Expiration Policy. | [optional] |

## Methods

### NewLifecycleExpiration

`func NewLifecycleExpiration() *LifecycleExpiration`

NewLifecycleExpiration instantiates a new LifecycleExpiration object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLifecycleExpirationWithDefaults

`func NewLifecycleExpirationWithDefaults() *LifecycleExpiration`

NewLifecycleExpirationWithDefaults instantiates a new LifecycleExpiration object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDate

`func (o *LifecycleExpiration) GetDate() time.Time`

GetDate returns the Date field if non-nil, zero value otherwise.

### GetDateOk

`func (o *LifecycleExpiration) GetDateOk() (*time.Time, bool)`

GetDateOk returns a tuple with the Date field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDate

`func (o *LifecycleExpiration) SetDate(v time.Time)`

SetDate sets Date field to given value.

### HasDate

`func (o *LifecycleExpiration) HasDate() bool`

HasDate returns a boolean if a field has been set.

### GetDays

`func (o *LifecycleExpiration) GetDays() int32`

GetDays returns the Days field if non-nil, zero value otherwise.

### GetDaysOk

`func (o *LifecycleExpiration) GetDaysOk() (*int32, bool)`

GetDaysOk returns a tuple with the Days field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDays

`func (o *LifecycleExpiration) SetDays(v int32)`

SetDays sets Days field to given value.

### HasDays

`func (o *LifecycleExpiration) HasDays() bool`

HasDays returns a boolean if a field has been set.

### GetExpiredObjectDeleteMarker

`func (o *LifecycleExpiration) GetExpiredObjectDeleteMarker() bool`

GetExpiredObjectDeleteMarker returns the ExpiredObjectDeleteMarker field if non-nil, zero value otherwise.

### GetExpiredObjectDeleteMarkerOk

`func (o *LifecycleExpiration) GetExpiredObjectDeleteMarkerOk() (*bool, bool)`

GetExpiredObjectDeleteMarkerOk returns a tuple with the ExpiredObjectDeleteMarker field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiredObjectDeleteMarker

`func (o *LifecycleExpiration) SetExpiredObjectDeleteMarker(v bool)`

SetExpiredObjectDeleteMarker sets ExpiredObjectDeleteMarker field to given value.

### HasExpiredObjectDeleteMarker

`func (o *LifecycleExpiration) HasExpiredObjectDeleteMarker() bool`

HasExpiredObjectDeleteMarker returns a boolean if a field has been set.


