# DefaultRetention

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Mode** | Pointer to **string** | The default Object Lock retention mode for new objects placed in the specified bucket. Must be used with either &#x60;Days&#x60; or &#x60;Years&#x60;.    | [optional] |
|**Days** | Pointer to **int32** | The number of days that you want to specify for the default retention period. Must be used with &#x60;Mode&#x60;. | [optional] |
|**Years** | Pointer to **int32** | The number of years that you want to specify for the default retention period. Must be used with &#x60;Mode&#x60;. | [optional] |

## Methods

### NewDefaultRetention

`func NewDefaultRetention() *DefaultRetention`

NewDefaultRetention instantiates a new DefaultRetention object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDefaultRetentionWithDefaults

`func NewDefaultRetentionWithDefaults() *DefaultRetention`

NewDefaultRetentionWithDefaults instantiates a new DefaultRetention object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMode

`func (o *DefaultRetention) GetMode() string`

GetMode returns the Mode field if non-nil, zero value otherwise.

### GetModeOk

`func (o *DefaultRetention) GetModeOk() (*string, bool)`

GetModeOk returns a tuple with the Mode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMode

`func (o *DefaultRetention) SetMode(v string)`

SetMode sets Mode field to given value.

### HasMode

`func (o *DefaultRetention) HasMode() bool`

HasMode returns a boolean if a field has been set.

### GetDays

`func (o *DefaultRetention) GetDays() int32`

GetDays returns the Days field if non-nil, zero value otherwise.

### GetDaysOk

`func (o *DefaultRetention) GetDaysOk() (*int32, bool)`

GetDaysOk returns a tuple with the Days field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDays

`func (o *DefaultRetention) SetDays(v int32)`

SetDays sets Days field to given value.

### HasDays

`func (o *DefaultRetention) HasDays() bool`

HasDays returns a boolean if a field has been set.

### GetYears

`func (o *DefaultRetention) GetYears() int32`

GetYears returns the Years field if non-nil, zero value otherwise.

### GetYearsOk

`func (o *DefaultRetention) GetYearsOk() (*int32, bool)`

GetYearsOk returns a tuple with the Years field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetYears

`func (o *DefaultRetention) SetYears(v int32)`

SetYears sets Years field to given value.

### HasYears

`func (o *DefaultRetention) HasYears() bool`

HasYears returns a boolean if a field has been set.


