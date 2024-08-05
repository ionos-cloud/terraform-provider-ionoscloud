# Rule

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**ID** | Pointer to **string** | Unique identifier for the rule. The value can&#39;t be longer than 255 characters. | [optional] |
|**Prefix** | **string** | Object key prefix that identifies one or more objects to which this rule applies. Replacement must be made for object keys containing special characters (such as carriage returns) when using XML requests. | |
|**Status** | [**ExpirationStatus**](ExpirationStatus.md) |  | |
|**Expiration** | Pointer to [**LifecycleExpiration**](LifecycleExpiration.md) |  | [optional] |
|**NoncurrentVersionExpiration** | Pointer to [**NoncurrentVersionExpiration**](NoncurrentVersionExpiration.md) |  | [optional] |
|**AbortIncompleteMultipartUpload** | Pointer to [**AbortIncompleteMultipartUpload**](AbortIncompleteMultipartUpload.md) |  | [optional] |

## Methods

### NewRule

`func NewRule(prefix string, status ExpirationStatus, ) *Rule`

NewRule instantiates a new Rule object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRuleWithDefaults

`func NewRuleWithDefaults() *Rule`

NewRuleWithDefaults instantiates a new Rule object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetID

`func (o *Rule) GetID() string`

GetID returns the ID field if non-nil, zero value otherwise.

### GetIDOk

`func (o *Rule) GetIDOk() (*string, bool)`

GetIDOk returns a tuple with the ID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetID

`func (o *Rule) SetID(v string)`

SetID sets ID field to given value.

### HasID

`func (o *Rule) HasID() bool`

HasID returns a boolean if a field has been set.

### GetPrefix

`func (o *Rule) GetPrefix() string`

GetPrefix returns the Prefix field if non-nil, zero value otherwise.

### GetPrefixOk

`func (o *Rule) GetPrefixOk() (*string, bool)`

GetPrefixOk returns a tuple with the Prefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrefix

`func (o *Rule) SetPrefix(v string)`

SetPrefix sets Prefix field to given value.


### GetStatus

`func (o *Rule) GetStatus() ExpirationStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *Rule) GetStatusOk() (*ExpirationStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *Rule) SetStatus(v ExpirationStatus)`

SetStatus sets Status field to given value.


### GetExpiration

`func (o *Rule) GetExpiration() LifecycleExpiration`

GetExpiration returns the Expiration field if non-nil, zero value otherwise.

### GetExpirationOk

`func (o *Rule) GetExpirationOk() (*LifecycleExpiration, bool)`

GetExpirationOk returns a tuple with the Expiration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiration

`func (o *Rule) SetExpiration(v LifecycleExpiration)`

SetExpiration sets Expiration field to given value.

### HasExpiration

`func (o *Rule) HasExpiration() bool`

HasExpiration returns a boolean if a field has been set.

### GetNoncurrentVersionExpiration

`func (o *Rule) GetNoncurrentVersionExpiration() NoncurrentVersionExpiration`

GetNoncurrentVersionExpiration returns the NoncurrentVersionExpiration field if non-nil, zero value otherwise.

### GetNoncurrentVersionExpirationOk

`func (o *Rule) GetNoncurrentVersionExpirationOk() (*NoncurrentVersionExpiration, bool)`

GetNoncurrentVersionExpirationOk returns a tuple with the NoncurrentVersionExpiration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNoncurrentVersionExpiration

`func (o *Rule) SetNoncurrentVersionExpiration(v NoncurrentVersionExpiration)`

SetNoncurrentVersionExpiration sets NoncurrentVersionExpiration field to given value.

### HasNoncurrentVersionExpiration

`func (o *Rule) HasNoncurrentVersionExpiration() bool`

HasNoncurrentVersionExpiration returns a boolean if a field has been set.

### GetAbortIncompleteMultipartUpload

`func (o *Rule) GetAbortIncompleteMultipartUpload() AbortIncompleteMultipartUpload`

GetAbortIncompleteMultipartUpload returns the AbortIncompleteMultipartUpload field if non-nil, zero value otherwise.

### GetAbortIncompleteMultipartUploadOk

`func (o *Rule) GetAbortIncompleteMultipartUploadOk() (*AbortIncompleteMultipartUpload, bool)`

GetAbortIncompleteMultipartUploadOk returns a tuple with the AbortIncompleteMultipartUpload field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAbortIncompleteMultipartUpload

`func (o *Rule) SetAbortIncompleteMultipartUpload(v AbortIncompleteMultipartUpload)`

SetAbortIncompleteMultipartUpload sets AbortIncompleteMultipartUpload field to given value.

### HasAbortIncompleteMultipartUpload

`func (o *Rule) HasAbortIncompleteMultipartUpload() bool`

HasAbortIncompleteMultipartUpload returns a boolean if a field has been set.


