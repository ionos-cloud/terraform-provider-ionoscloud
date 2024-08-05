# MultipartUpload

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**UploadId** | Pointer to **string** | ID of the multipart upload. | [optional] |
|**Key** | Pointer to **string** | The object key. | [optional] |
|**Initiated** | Pointer to [**time.Time**](time.Time.md) | Date and time at which the multipart upload was initiated. | [optional] |
|**StorageClass** | Pointer to [**StorageClass**](StorageClass.md) |  | [optional] |
|**Owner** | Pointer to [**Owner**](Owner.md) |  | [optional] |
|**Initiator** | Pointer to [**Initiator**](Initiator.md) |  | [optional] |

## Methods

### NewMultipartUpload

`func NewMultipartUpload() *MultipartUpload`

NewMultipartUpload instantiates a new MultipartUpload object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMultipartUploadWithDefaults

`func NewMultipartUploadWithDefaults() *MultipartUpload`

NewMultipartUploadWithDefaults instantiates a new MultipartUpload object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUploadId

`func (o *MultipartUpload) GetUploadId() string`

GetUploadId returns the UploadId field if non-nil, zero value otherwise.

### GetUploadIdOk

`func (o *MultipartUpload) GetUploadIdOk() (*string, bool)`

GetUploadIdOk returns a tuple with the UploadId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUploadId

`func (o *MultipartUpload) SetUploadId(v string)`

SetUploadId sets UploadId field to given value.

### HasUploadId

`func (o *MultipartUpload) HasUploadId() bool`

HasUploadId returns a boolean if a field has been set.

### GetKey

`func (o *MultipartUpload) GetKey() string`

GetKey returns the Key field if non-nil, zero value otherwise.

### GetKeyOk

`func (o *MultipartUpload) GetKeyOk() (*string, bool)`

GetKeyOk returns a tuple with the Key field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKey

`func (o *MultipartUpload) SetKey(v string)`

SetKey sets Key field to given value.

### HasKey

`func (o *MultipartUpload) HasKey() bool`

HasKey returns a boolean if a field has been set.

### GetInitiated

`func (o *MultipartUpload) GetInitiated() time.Time`

GetInitiated returns the Initiated field if non-nil, zero value otherwise.

### GetInitiatedOk

`func (o *MultipartUpload) GetInitiatedOk() (*time.Time, bool)`

GetInitiatedOk returns a tuple with the Initiated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInitiated

`func (o *MultipartUpload) SetInitiated(v time.Time)`

SetInitiated sets Initiated field to given value.

### HasInitiated

`func (o *MultipartUpload) HasInitiated() bool`

HasInitiated returns a boolean if a field has been set.

### GetStorageClass

`func (o *MultipartUpload) GetStorageClass() StorageClass`

GetStorageClass returns the StorageClass field if non-nil, zero value otherwise.

### GetStorageClassOk

`func (o *MultipartUpload) GetStorageClassOk() (*StorageClass, bool)`

GetStorageClassOk returns a tuple with the StorageClass field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStorageClass

`func (o *MultipartUpload) SetStorageClass(v StorageClass)`

SetStorageClass sets StorageClass field to given value.

### HasStorageClass

`func (o *MultipartUpload) HasStorageClass() bool`

HasStorageClass returns a boolean if a field has been set.

### GetOwner

`func (o *MultipartUpload) GetOwner() Owner`

GetOwner returns the Owner field if non-nil, zero value otherwise.

### GetOwnerOk

`func (o *MultipartUpload) GetOwnerOk() (*Owner, bool)`

GetOwnerOk returns a tuple with the Owner field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOwner

`func (o *MultipartUpload) SetOwner(v Owner)`

SetOwner sets Owner field to given value.

### HasOwner

`func (o *MultipartUpload) HasOwner() bool`

HasOwner returns a boolean if a field has been set.

### GetInitiator

`func (o *MultipartUpload) GetInitiator() Initiator`

GetInitiator returns the Initiator field if non-nil, zero value otherwise.

### GetInitiatorOk

`func (o *MultipartUpload) GetInitiatorOk() (*Initiator, bool)`

GetInitiatorOk returns a tuple with the Initiator field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInitiator

`func (o *MultipartUpload) SetInitiator(v Initiator)`

SetInitiator sets Initiator field to given value.

### HasInitiator

`func (o *MultipartUpload) HasInitiator() bool`

HasInitiator returns a boolean if a field has been set.


