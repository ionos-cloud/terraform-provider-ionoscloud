# BucketPolicyConditionIpAddress

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**AwsSourceIp** | Pointer to **[]string** |  | [optional] |

## Methods

### NewBucketPolicyConditionIpAddress

`func NewBucketPolicyConditionIpAddress() *BucketPolicyConditionIpAddress`

NewBucketPolicyConditionIpAddress instantiates a new BucketPolicyConditionIpAddress object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBucketPolicyConditionIpAddressWithDefaults

`func NewBucketPolicyConditionIpAddressWithDefaults() *BucketPolicyConditionIpAddress`

NewBucketPolicyConditionIpAddressWithDefaults instantiates a new BucketPolicyConditionIpAddress object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAwsSourceIp

`func (o *BucketPolicyConditionIpAddress) GetAwsSourceIp() []string`

GetAwsSourceIp returns the AwsSourceIp field if non-nil, zero value otherwise.

### GetAwsSourceIpOk

`func (o *BucketPolicyConditionIpAddress) GetAwsSourceIpOk() (*[]string, bool)`

GetAwsSourceIpOk returns a tuple with the AwsSourceIp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAwsSourceIp

`func (o *BucketPolicyConditionIpAddress) SetAwsSourceIp(v []string)`

SetAwsSourceIp sets AwsSourceIp field to given value.

### HasAwsSourceIp

`func (o *BucketPolicyConditionIpAddress) HasAwsSourceIp() bool`

HasAwsSourceIp returns a boolean if a field has been set.


