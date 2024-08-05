# CompleteMultipartUploadOutput

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Location** | Pointer to **string** | The URI that identifies the newly created object. | [optional] |
|**Bucket** | Pointer to **string** | The bucket name. | [optional] |
|**Key** | Pointer to **string** | The object key. | [optional] |
|**ETag** | Pointer to **string** | Entity tag that identifies the object&#39;s data. Objects with different object data will have different entity tags. The entity tag is an opaque string. The entity tag may or may not be an MD5 digest of the object data. If the entity tag is not an MD5 digest of the object data, it will contain one or more nonhexadecimal characters and/or will consist of less than 32 or more than 32 hexadecimal digits.  | [optional] |

## Methods

### NewCompleteMultipartUploadOutput

`func NewCompleteMultipartUploadOutput() *CompleteMultipartUploadOutput`

NewCompleteMultipartUploadOutput instantiates a new CompleteMultipartUploadOutput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCompleteMultipartUploadOutputWithDefaults

`func NewCompleteMultipartUploadOutputWithDefaults() *CompleteMultipartUploadOutput`

NewCompleteMultipartUploadOutputWithDefaults instantiates a new CompleteMultipartUploadOutput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLocation

`func (o *CompleteMultipartUploadOutput) GetLocation() string`

GetLocation returns the Location field if non-nil, zero value otherwise.

### GetLocationOk

`func (o *CompleteMultipartUploadOutput) GetLocationOk() (*string, bool)`

GetLocationOk returns a tuple with the Location field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocation

`func (o *CompleteMultipartUploadOutput) SetLocation(v string)`

SetLocation sets Location field to given value.

### HasLocation

`func (o *CompleteMultipartUploadOutput) HasLocation() bool`

HasLocation returns a boolean if a field has been set.

### GetBucket

`func (o *CompleteMultipartUploadOutput) GetBucket() string`

GetBucket returns the Bucket field if non-nil, zero value otherwise.

### GetBucketOk

`func (o *CompleteMultipartUploadOutput) GetBucketOk() (*string, bool)`

GetBucketOk returns a tuple with the Bucket field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBucket

`func (o *CompleteMultipartUploadOutput) SetBucket(v string)`

SetBucket sets Bucket field to given value.

### HasBucket

`func (o *CompleteMultipartUploadOutput) HasBucket() bool`

HasBucket returns a boolean if a field has been set.

### GetKey

`func (o *CompleteMultipartUploadOutput) GetKey() string`

GetKey returns the Key field if non-nil, zero value otherwise.

### GetKeyOk

`func (o *CompleteMultipartUploadOutput) GetKeyOk() (*string, bool)`

GetKeyOk returns a tuple with the Key field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKey

`func (o *CompleteMultipartUploadOutput) SetKey(v string)`

SetKey sets Key field to given value.

### HasKey

`func (o *CompleteMultipartUploadOutput) HasKey() bool`

HasKey returns a boolean if a field has been set.

### GetETag

`func (o *CompleteMultipartUploadOutput) GetETag() string`

GetETag returns the ETag field if non-nil, zero value otherwise.

### GetETagOk

`func (o *CompleteMultipartUploadOutput) GetETagOk() (*string, bool)`

GetETagOk returns a tuple with the ETag field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetETag

`func (o *CompleteMultipartUploadOutput) SetETag(v string)`

SetETag sets ETag field to given value.

### HasETag

`func (o *CompleteMultipartUploadOutput) HasETag() bool`

HasETag returns a boolean if a field has been set.


