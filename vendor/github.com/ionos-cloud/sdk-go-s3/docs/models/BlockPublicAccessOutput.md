# BlockPublicAccessOutput

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**BlockPublicAcls** | Pointer to **bool** | Indicates that access to the bucket via Access Control Lists (ACLs) that grant public access is blocked. In other words, ACLs that allow public access are not permitted.  | [optional] |
|**IgnorePublicAcls** | Pointer to **bool** | Instructs the system to ignore any ACLs that grant public access. Even if ACLs are set to allow public access, they will be disregarded.  | [optional] |
|**BlockPublicPolicy** | Pointer to **bool** | Blocks public access to the bucket via bucket policies. Bucket policies that grant public access will not be allowed.  | [optional] |
|**RestrictPublicBuckets** | Pointer to **bool** | Restricts access to buckets that have public policies. Buckets with policies that grant public access will have their access restricted.  | [optional] |

## Methods

### NewBlockPublicAccessOutput

`func NewBlockPublicAccessOutput() *BlockPublicAccessOutput`

NewBlockPublicAccessOutput instantiates a new BlockPublicAccessOutput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBlockPublicAccessOutputWithDefaults

`func NewBlockPublicAccessOutputWithDefaults() *BlockPublicAccessOutput`

NewBlockPublicAccessOutputWithDefaults instantiates a new BlockPublicAccessOutput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBlockPublicAcls

`func (o *BlockPublicAccessOutput) GetBlockPublicAcls() bool`

GetBlockPublicAcls returns the BlockPublicAcls field if non-nil, zero value otherwise.

### GetBlockPublicAclsOk

`func (o *BlockPublicAccessOutput) GetBlockPublicAclsOk() (*bool, bool)`

GetBlockPublicAclsOk returns a tuple with the BlockPublicAcls field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBlockPublicAcls

`func (o *BlockPublicAccessOutput) SetBlockPublicAcls(v bool)`

SetBlockPublicAcls sets BlockPublicAcls field to given value.

### HasBlockPublicAcls

`func (o *BlockPublicAccessOutput) HasBlockPublicAcls() bool`

HasBlockPublicAcls returns a boolean if a field has been set.

### GetIgnorePublicAcls

`func (o *BlockPublicAccessOutput) GetIgnorePublicAcls() bool`

GetIgnorePublicAcls returns the IgnorePublicAcls field if non-nil, zero value otherwise.

### GetIgnorePublicAclsOk

`func (o *BlockPublicAccessOutput) GetIgnorePublicAclsOk() (*bool, bool)`

GetIgnorePublicAclsOk returns a tuple with the IgnorePublicAcls field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIgnorePublicAcls

`func (o *BlockPublicAccessOutput) SetIgnorePublicAcls(v bool)`

SetIgnorePublicAcls sets IgnorePublicAcls field to given value.

### HasIgnorePublicAcls

`func (o *BlockPublicAccessOutput) HasIgnorePublicAcls() bool`

HasIgnorePublicAcls returns a boolean if a field has been set.

### GetBlockPublicPolicy

`func (o *BlockPublicAccessOutput) GetBlockPublicPolicy() bool`

GetBlockPublicPolicy returns the BlockPublicPolicy field if non-nil, zero value otherwise.

### GetBlockPublicPolicyOk

`func (o *BlockPublicAccessOutput) GetBlockPublicPolicyOk() (*bool, bool)`

GetBlockPublicPolicyOk returns a tuple with the BlockPublicPolicy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBlockPublicPolicy

`func (o *BlockPublicAccessOutput) SetBlockPublicPolicy(v bool)`

SetBlockPublicPolicy sets BlockPublicPolicy field to given value.

### HasBlockPublicPolicy

`func (o *BlockPublicAccessOutput) HasBlockPublicPolicy() bool`

HasBlockPublicPolicy returns a boolean if a field has been set.

### GetRestrictPublicBuckets

`func (o *BlockPublicAccessOutput) GetRestrictPublicBuckets() bool`

GetRestrictPublicBuckets returns the RestrictPublicBuckets field if non-nil, zero value otherwise.

### GetRestrictPublicBucketsOk

`func (o *BlockPublicAccessOutput) GetRestrictPublicBucketsOk() (*bool, bool)`

GetRestrictPublicBucketsOk returns a tuple with the RestrictPublicBuckets field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRestrictPublicBuckets

`func (o *BlockPublicAccessOutput) SetRestrictPublicBuckets(v bool)`

SetRestrictPublicBuckets sets RestrictPublicBuckets field to given value.

### HasRestrictPublicBuckets

`func (o *BlockPublicAccessOutput) HasRestrictPublicBuckets() bool`

HasRestrictPublicBuckets returns a boolean if a field has been set.


