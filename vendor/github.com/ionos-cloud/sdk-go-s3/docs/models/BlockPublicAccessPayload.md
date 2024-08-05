# BlockPublicAccessPayload

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**BlockPublicAcls** | Pointer to **bool** | Indicates that access to the bucket via Access Control Lists (ACLs) that grant public access is blocked. In other words, ACLs that allow public access are not permitted.  | [optional] [default to false]|
|**IgnorePublicAcls** | Pointer to **bool** | Instructs the system to ignore any ACLs that grant public access. Even if ACLs are set to allow public access, they will be disregarded.  | [optional] [default to false]|
|**BlockPublicPolicy** | Pointer to **bool** | Blocks public access to the bucket via bucket policies. Bucket policies that grant public access will not be allowed.  | [optional] [default to false]|
|**RestrictPublicBuckets** | Pointer to **bool** | Restricts access to buckets that have public policies. Buckets with policies that grant public access will have their access restricted.  | [optional] [default to false]|

## Methods

### NewBlockPublicAccessPayload

`func NewBlockPublicAccessPayload() *BlockPublicAccessPayload`

NewBlockPublicAccessPayload instantiates a new BlockPublicAccessPayload object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBlockPublicAccessPayloadWithDefaults

`func NewBlockPublicAccessPayloadWithDefaults() *BlockPublicAccessPayload`

NewBlockPublicAccessPayloadWithDefaults instantiates a new BlockPublicAccessPayload object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBlockPublicAcls

`func (o *BlockPublicAccessPayload) GetBlockPublicAcls() bool`

GetBlockPublicAcls returns the BlockPublicAcls field if non-nil, zero value otherwise.

### GetBlockPublicAclsOk

`func (o *BlockPublicAccessPayload) GetBlockPublicAclsOk() (*bool, bool)`

GetBlockPublicAclsOk returns a tuple with the BlockPublicAcls field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBlockPublicAcls

`func (o *BlockPublicAccessPayload) SetBlockPublicAcls(v bool)`

SetBlockPublicAcls sets BlockPublicAcls field to given value.

### HasBlockPublicAcls

`func (o *BlockPublicAccessPayload) HasBlockPublicAcls() bool`

HasBlockPublicAcls returns a boolean if a field has been set.

### GetIgnorePublicAcls

`func (o *BlockPublicAccessPayload) GetIgnorePublicAcls() bool`

GetIgnorePublicAcls returns the IgnorePublicAcls field if non-nil, zero value otherwise.

### GetIgnorePublicAclsOk

`func (o *BlockPublicAccessPayload) GetIgnorePublicAclsOk() (*bool, bool)`

GetIgnorePublicAclsOk returns a tuple with the IgnorePublicAcls field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIgnorePublicAcls

`func (o *BlockPublicAccessPayload) SetIgnorePublicAcls(v bool)`

SetIgnorePublicAcls sets IgnorePublicAcls field to given value.

### HasIgnorePublicAcls

`func (o *BlockPublicAccessPayload) HasIgnorePublicAcls() bool`

HasIgnorePublicAcls returns a boolean if a field has been set.

### GetBlockPublicPolicy

`func (o *BlockPublicAccessPayload) GetBlockPublicPolicy() bool`

GetBlockPublicPolicy returns the BlockPublicPolicy field if non-nil, zero value otherwise.

### GetBlockPublicPolicyOk

`func (o *BlockPublicAccessPayload) GetBlockPublicPolicyOk() (*bool, bool)`

GetBlockPublicPolicyOk returns a tuple with the BlockPublicPolicy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBlockPublicPolicy

`func (o *BlockPublicAccessPayload) SetBlockPublicPolicy(v bool)`

SetBlockPublicPolicy sets BlockPublicPolicy field to given value.

### HasBlockPublicPolicy

`func (o *BlockPublicAccessPayload) HasBlockPublicPolicy() bool`

HasBlockPublicPolicy returns a boolean if a field has been set.

### GetRestrictPublicBuckets

`func (o *BlockPublicAccessPayload) GetRestrictPublicBuckets() bool`

GetRestrictPublicBuckets returns the RestrictPublicBuckets field if non-nil, zero value otherwise.

### GetRestrictPublicBucketsOk

`func (o *BlockPublicAccessPayload) GetRestrictPublicBucketsOk() (*bool, bool)`

GetRestrictPublicBucketsOk returns a tuple with the RestrictPublicBuckets field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRestrictPublicBuckets

`func (o *BlockPublicAccessPayload) SetRestrictPublicBuckets(v bool)`

SetRestrictPublicBuckets sets RestrictPublicBuckets field to given value.

### HasRestrictPublicBuckets

`func (o *BlockPublicAccessPayload) HasRestrictPublicBuckets() bool`

HasRestrictPublicBuckets returns a boolean if a field has been set.


