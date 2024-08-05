# BucketPolicyStatement

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Sid** | Pointer to **string** | Custom string identifying the statement. | [optional] |
|**Action** | **[]string** | The array of allowed or denied actions.   IONOS S3 Object Storage supports the use of a wildcard in your Action configuration (&#x60;\&quot;Action\&quot;:[\&quot;s3:*\&quot;]&#x60;). When an Action wildcard is used together with an object-level Resource element (&#x60;\&quot;arn:aws:s3:::&lt;bucketName&gt;/_*\&quot;&#x60; or &#x60;\&quot;arn:aws:s3:::&lt;bucketName&gt;/&lt;objectName&gt;\&quot;&#x60;), the wildcard denotes all supported Object actions. When an Action wildcard is used together with bucket-level Resource element (&#x60;\&quot;arn:aws:s3:::&lt;bucketName&gt;\&quot;&#x60;), the wildcard denotes all the bucket actions and bucket subresource actions that IONOS S3 Object Storage supports.  | |
|**Effect** | **string** | Specify the outcome when the user requests a particular action. | |
|**Resource** | **[]string** | The bucket or object that the policy applies to.   Must be one of the following: - &#x60;\&quot;arn:aws:s3:::&lt;bucketName&gt;\&quot;&#x60; - For bucket actions (such as &#x60;s3:ListBucket&#x60;) and bucket subresource actions (such as &#x60;s3:GetBucketAcl&#x60;). - &#x60;\&quot;arn:aws:s3:::&lt;bucketName&gt;/_*\&quot;&#x60; or &#x60;\&quot;arn:aws:s3:::&lt;bucketName&gt;/&lt;objectName&gt;\&quot;&#x60; - For object actions (such as &#x60;s3:PutObject&#x60;).  | |
|**Condition** | Pointer to [**BucketPolicyCondition**](BucketPolicyCondition.md) |  | [optional] |
|**Principal** | Pointer to [**Principal**](Principal.md) |  | [optional] |

## Methods

### NewBucketPolicyStatement

`func NewBucketPolicyStatement(action []string, effect string, resource []string, ) *BucketPolicyStatement`

NewBucketPolicyStatement instantiates a new BucketPolicyStatement object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBucketPolicyStatementWithDefaults

`func NewBucketPolicyStatementWithDefaults() *BucketPolicyStatement`

NewBucketPolicyStatementWithDefaults instantiates a new BucketPolicyStatement object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSid

`func (o *BucketPolicyStatement) GetSid() string`

GetSid returns the Sid field if non-nil, zero value otherwise.

### GetSidOk

`func (o *BucketPolicyStatement) GetSidOk() (*string, bool)`

GetSidOk returns a tuple with the Sid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSid

`func (o *BucketPolicyStatement) SetSid(v string)`

SetSid sets Sid field to given value.

### HasSid

`func (o *BucketPolicyStatement) HasSid() bool`

HasSid returns a boolean if a field has been set.

### GetAction

`func (o *BucketPolicyStatement) GetAction() []string`

GetAction returns the Action field if non-nil, zero value otherwise.

### GetActionOk

`func (o *BucketPolicyStatement) GetActionOk() (*[]string, bool)`

GetActionOk returns a tuple with the Action field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAction

`func (o *BucketPolicyStatement) SetAction(v []string)`

SetAction sets Action field to given value.


### GetEffect

`func (o *BucketPolicyStatement) GetEffect() string`

GetEffect returns the Effect field if non-nil, zero value otherwise.

### GetEffectOk

`func (o *BucketPolicyStatement) GetEffectOk() (*string, bool)`

GetEffectOk returns a tuple with the Effect field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEffect

`func (o *BucketPolicyStatement) SetEffect(v string)`

SetEffect sets Effect field to given value.


### GetResource

`func (o *BucketPolicyStatement) GetResource() []string`

GetResource returns the Resource field if non-nil, zero value otherwise.

### GetResourceOk

`func (o *BucketPolicyStatement) GetResourceOk() (*[]string, bool)`

GetResourceOk returns a tuple with the Resource field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResource

`func (o *BucketPolicyStatement) SetResource(v []string)`

SetResource sets Resource field to given value.


### GetCondition

`func (o *BucketPolicyStatement) GetCondition() BucketPolicyCondition`

GetCondition returns the Condition field if non-nil, zero value otherwise.

### GetConditionOk

`func (o *BucketPolicyStatement) GetConditionOk() (*BucketPolicyCondition, bool)`

GetConditionOk returns a tuple with the Condition field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCondition

`func (o *BucketPolicyStatement) SetCondition(v BucketPolicyCondition)`

SetCondition sets Condition field to given value.

### HasCondition

`func (o *BucketPolicyStatement) HasCondition() bool`

HasCondition returns a boolean if a field has been set.

### GetPrincipal

`func (o *BucketPolicyStatement) GetPrincipal() Principal`

GetPrincipal returns the Principal field if non-nil, zero value otherwise.

### GetPrincipalOk

`func (o *BucketPolicyStatement) GetPrincipalOk() (*Principal, bool)`

GetPrincipalOk returns a tuple with the Principal field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrincipal

`func (o *BucketPolicyStatement) SetPrincipal(v Principal)`

SetPrincipal sets Principal field to given value.

### HasPrincipal

`func (o *BucketPolicyStatement) HasPrincipal() bool`

HasPrincipal returns a boolean if a field has been set.


