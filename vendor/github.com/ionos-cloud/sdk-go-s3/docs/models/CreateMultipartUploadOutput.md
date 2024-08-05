# CreateMultipartUploadOutput

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Bucket** | Pointer to **string** | The bucket name. | [optional] |
|**Key** | Pointer to **string** | The object key. | [optional] |
|**UploadId** | Pointer to **string** | ID of the multipart upload. | [optional] |

## Methods

### NewCreateMultipartUploadOutput

`func NewCreateMultipartUploadOutput() *CreateMultipartUploadOutput`

NewCreateMultipartUploadOutput instantiates a new CreateMultipartUploadOutput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateMultipartUploadOutputWithDefaults

`func NewCreateMultipartUploadOutputWithDefaults() *CreateMultipartUploadOutput`

NewCreateMultipartUploadOutputWithDefaults instantiates a new CreateMultipartUploadOutput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBucket

`func (o *CreateMultipartUploadOutput) GetBucket() string`

GetBucket returns the Bucket field if non-nil, zero value otherwise.

### GetBucketOk

`func (o *CreateMultipartUploadOutput) GetBucketOk() (*string, bool)`

GetBucketOk returns a tuple with the Bucket field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBucket

`func (o *CreateMultipartUploadOutput) SetBucket(v string)`

SetBucket sets Bucket field to given value.

### HasBucket

`func (o *CreateMultipartUploadOutput) HasBucket() bool`

HasBucket returns a boolean if a field has been set.

### GetKey

`func (o *CreateMultipartUploadOutput) GetKey() string`

GetKey returns the Key field if non-nil, zero value otherwise.

### GetKeyOk

`func (o *CreateMultipartUploadOutput) GetKeyOk() (*string, bool)`

GetKeyOk returns a tuple with the Key field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKey

`func (o *CreateMultipartUploadOutput) SetKey(v string)`

SetKey sets Key field to given value.

### HasKey

`func (o *CreateMultipartUploadOutput) HasKey() bool`

HasKey returns a boolean if a field has been set.

### GetUploadId

`func (o *CreateMultipartUploadOutput) GetUploadId() string`

GetUploadId returns the UploadId field if non-nil, zero value otherwise.

### GetUploadIdOk

`func (o *CreateMultipartUploadOutput) GetUploadIdOk() (*string, bool)`

GetUploadIdOk returns a tuple with the UploadId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUploadId

`func (o *CreateMultipartUploadOutput) SetUploadId(v string)`

SetUploadId sets UploadId field to given value.

### HasUploadId

`func (o *CreateMultipartUploadOutput) HasUploadId() bool`

HasUploadId returns a boolean if a field has been set.


