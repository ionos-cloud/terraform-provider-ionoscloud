# BucketPolicyCondition

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**IpAddress** | Pointer to [**BucketPolicyConditionIpAddress**](BucketPolicyConditionIpAddress.md) |  | [optional] |
|**NotIpAddress** | Pointer to [**BucketPolicyConditionIpAddress**](BucketPolicyConditionIpAddress.md) |  | [optional] |
|**DateGreaterThan** | Pointer to [**BucketPolicyConditionDate**](BucketPolicyConditionDate.md) |  | [optional] |
|**DateLessThan** | Pointer to [**BucketPolicyConditionDate**](BucketPolicyConditionDate.md) |  | [optional] |

## Methods

### NewBucketPolicyCondition

`func NewBucketPolicyCondition() *BucketPolicyCondition`

NewBucketPolicyCondition instantiates a new BucketPolicyCondition object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBucketPolicyConditionWithDefaults

`func NewBucketPolicyConditionWithDefaults() *BucketPolicyCondition`

NewBucketPolicyConditionWithDefaults instantiates a new BucketPolicyCondition object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetIpAddress

`func (o *BucketPolicyCondition) GetIpAddress() BucketPolicyConditionIpAddress`

GetIpAddress returns the IpAddress field if non-nil, zero value otherwise.

### GetIpAddressOk

`func (o *BucketPolicyCondition) GetIpAddressOk() (*BucketPolicyConditionIpAddress, bool)`

GetIpAddressOk returns a tuple with the IpAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpAddress

`func (o *BucketPolicyCondition) SetIpAddress(v BucketPolicyConditionIpAddress)`

SetIpAddress sets IpAddress field to given value.

### HasIpAddress

`func (o *BucketPolicyCondition) HasIpAddress() bool`

HasIpAddress returns a boolean if a field has been set.

### GetNotIpAddress

`func (o *BucketPolicyCondition) GetNotIpAddress() BucketPolicyConditionIpAddress`

GetNotIpAddress returns the NotIpAddress field if non-nil, zero value otherwise.

### GetNotIpAddressOk

`func (o *BucketPolicyCondition) GetNotIpAddressOk() (*BucketPolicyConditionIpAddress, bool)`

GetNotIpAddressOk returns a tuple with the NotIpAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNotIpAddress

`func (o *BucketPolicyCondition) SetNotIpAddress(v BucketPolicyConditionIpAddress)`

SetNotIpAddress sets NotIpAddress field to given value.

### HasNotIpAddress

`func (o *BucketPolicyCondition) HasNotIpAddress() bool`

HasNotIpAddress returns a boolean if a field has been set.

### GetDateGreaterThan

`func (o *BucketPolicyCondition) GetDateGreaterThan() BucketPolicyConditionDate`

GetDateGreaterThan returns the DateGreaterThan field if non-nil, zero value otherwise.

### GetDateGreaterThanOk

`func (o *BucketPolicyCondition) GetDateGreaterThanOk() (*BucketPolicyConditionDate, bool)`

GetDateGreaterThanOk returns a tuple with the DateGreaterThan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDateGreaterThan

`func (o *BucketPolicyCondition) SetDateGreaterThan(v BucketPolicyConditionDate)`

SetDateGreaterThan sets DateGreaterThan field to given value.

### HasDateGreaterThan

`func (o *BucketPolicyCondition) HasDateGreaterThan() bool`

HasDateGreaterThan returns a boolean if a field has been set.

### GetDateLessThan

`func (o *BucketPolicyCondition) GetDateLessThan() BucketPolicyConditionDate`

GetDateLessThan returns the DateLessThan field if non-nil, zero value otherwise.

### GetDateLessThanOk

`func (o *BucketPolicyCondition) GetDateLessThanOk() (*BucketPolicyConditionDate, bool)`

GetDateLessThanOk returns a tuple with the DateLessThan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDateLessThan

`func (o *BucketPolicyCondition) SetDateLessThan(v BucketPolicyConditionDate)`

SetDateLessThan sets DateLessThan field to given value.

### HasDateLessThan

`func (o *BucketPolicyCondition) HasDateLessThan() bool`

HasDateLessThan returns a boolean if a field has been set.


