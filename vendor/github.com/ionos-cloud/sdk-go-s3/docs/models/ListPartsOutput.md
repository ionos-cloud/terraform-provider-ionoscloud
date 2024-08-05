# ListPartsOutput

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Bucket** | Pointer to **string** | The bucket name. | [optional] |
|**Key** | Pointer to **string** | The object key. | [optional] |
|**UploadId** | Pointer to **string** | ID of the multipart upload. | [optional] |
|**PartNumberMarker** | Pointer to **int32** | When a list is truncated, this element specifies the last part in the list, as well as the value to use for the part-number-marker request parameter in a subsequent request. | [optional] |
|**NextPartNumberMarker** | Pointer to **string** | When a list is truncated, this element specifies the last part in the list, as well as the value to use for the part-number-marker request parameter in a subsequent request. | [optional] |
|**MaxParts** | Pointer to **string** | Maximum number of parts that were allowed in the response. | [optional] |
|**IsTruncated** | Pointer to **bool** | A flag that indicates whether IONOS S3 Object Storage returned all of the results that satisfied the search criteria. If your results were truncated, you can make a follow-up paginated request using the NextKeyMarker and NextVersionIdMarker response parameters as a starting place in another request to return the rest of the results. | [optional] |
|**Parts** | Pointer to [**[]Part**](Part.md) |  Container for elements related to a particular part. A response can contain zero or more &#x60;Part&#x60; elements. | [optional] |
|**Initiator** | Pointer to [**Initiator**](Initiator.md) |  | [optional] |
|**Owner** | Pointer to [**Owner**](Owner.md) |  | [optional] |
|**StorageClass** | Pointer to [**StorageClass**](StorageClass.md) |  | [optional] |

## Methods

### NewListPartsOutput

`func NewListPartsOutput() *ListPartsOutput`

NewListPartsOutput instantiates a new ListPartsOutput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewListPartsOutputWithDefaults

`func NewListPartsOutputWithDefaults() *ListPartsOutput`

NewListPartsOutputWithDefaults instantiates a new ListPartsOutput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBucket

`func (o *ListPartsOutput) GetBucket() string`

GetBucket returns the Bucket field if non-nil, zero value otherwise.

### GetBucketOk

`func (o *ListPartsOutput) GetBucketOk() (*string, bool)`

GetBucketOk returns a tuple with the Bucket field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBucket

`func (o *ListPartsOutput) SetBucket(v string)`

SetBucket sets Bucket field to given value.

### HasBucket

`func (o *ListPartsOutput) HasBucket() bool`

HasBucket returns a boolean if a field has been set.

### GetKey

`func (o *ListPartsOutput) GetKey() string`

GetKey returns the Key field if non-nil, zero value otherwise.

### GetKeyOk

`func (o *ListPartsOutput) GetKeyOk() (*string, bool)`

GetKeyOk returns a tuple with the Key field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKey

`func (o *ListPartsOutput) SetKey(v string)`

SetKey sets Key field to given value.

### HasKey

`func (o *ListPartsOutput) HasKey() bool`

HasKey returns a boolean if a field has been set.

### GetUploadId

`func (o *ListPartsOutput) GetUploadId() string`

GetUploadId returns the UploadId field if non-nil, zero value otherwise.

### GetUploadIdOk

`func (o *ListPartsOutput) GetUploadIdOk() (*string, bool)`

GetUploadIdOk returns a tuple with the UploadId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUploadId

`func (o *ListPartsOutput) SetUploadId(v string)`

SetUploadId sets UploadId field to given value.

### HasUploadId

`func (o *ListPartsOutput) HasUploadId() bool`

HasUploadId returns a boolean if a field has been set.

### GetPartNumberMarker

`func (o *ListPartsOutput) GetPartNumberMarker() int32`

GetPartNumberMarker returns the PartNumberMarker field if non-nil, zero value otherwise.

### GetPartNumberMarkerOk

`func (o *ListPartsOutput) GetPartNumberMarkerOk() (*int32, bool)`

GetPartNumberMarkerOk returns a tuple with the PartNumberMarker field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartNumberMarker

`func (o *ListPartsOutput) SetPartNumberMarker(v int32)`

SetPartNumberMarker sets PartNumberMarker field to given value.

### HasPartNumberMarker

`func (o *ListPartsOutput) HasPartNumberMarker() bool`

HasPartNumberMarker returns a boolean if a field has been set.

### GetNextPartNumberMarker

`func (o *ListPartsOutput) GetNextPartNumberMarker() string`

GetNextPartNumberMarker returns the NextPartNumberMarker field if non-nil, zero value otherwise.

### GetNextPartNumberMarkerOk

`func (o *ListPartsOutput) GetNextPartNumberMarkerOk() (*string, bool)`

GetNextPartNumberMarkerOk returns a tuple with the NextPartNumberMarker field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextPartNumberMarker

`func (o *ListPartsOutput) SetNextPartNumberMarker(v string)`

SetNextPartNumberMarker sets NextPartNumberMarker field to given value.

### HasNextPartNumberMarker

`func (o *ListPartsOutput) HasNextPartNumberMarker() bool`

HasNextPartNumberMarker returns a boolean if a field has been set.

### GetMaxParts

`func (o *ListPartsOutput) GetMaxParts() string`

GetMaxParts returns the MaxParts field if non-nil, zero value otherwise.

### GetMaxPartsOk

`func (o *ListPartsOutput) GetMaxPartsOk() (*string, bool)`

GetMaxPartsOk returns a tuple with the MaxParts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxParts

`func (o *ListPartsOutput) SetMaxParts(v string)`

SetMaxParts sets MaxParts field to given value.

### HasMaxParts

`func (o *ListPartsOutput) HasMaxParts() bool`

HasMaxParts returns a boolean if a field has been set.

### GetIsTruncated

`func (o *ListPartsOutput) GetIsTruncated() bool`

GetIsTruncated returns the IsTruncated field if non-nil, zero value otherwise.

### GetIsTruncatedOk

`func (o *ListPartsOutput) GetIsTruncatedOk() (*bool, bool)`

GetIsTruncatedOk returns a tuple with the IsTruncated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsTruncated

`func (o *ListPartsOutput) SetIsTruncated(v bool)`

SetIsTruncated sets IsTruncated field to given value.

### HasIsTruncated

`func (o *ListPartsOutput) HasIsTruncated() bool`

HasIsTruncated returns a boolean if a field has been set.

### GetParts

`func (o *ListPartsOutput) GetParts() []Part`

GetParts returns the Parts field if non-nil, zero value otherwise.

### GetPartsOk

`func (o *ListPartsOutput) GetPartsOk() (*[]Part, bool)`

GetPartsOk returns a tuple with the Parts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetParts

`func (o *ListPartsOutput) SetParts(v []Part)`

SetParts sets Parts field to given value.

### HasParts

`func (o *ListPartsOutput) HasParts() bool`

HasParts returns a boolean if a field has been set.

### GetInitiator

`func (o *ListPartsOutput) GetInitiator() Initiator`

GetInitiator returns the Initiator field if non-nil, zero value otherwise.

### GetInitiatorOk

`func (o *ListPartsOutput) GetInitiatorOk() (*Initiator, bool)`

GetInitiatorOk returns a tuple with the Initiator field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInitiator

`func (o *ListPartsOutput) SetInitiator(v Initiator)`

SetInitiator sets Initiator field to given value.

### HasInitiator

`func (o *ListPartsOutput) HasInitiator() bool`

HasInitiator returns a boolean if a field has been set.

### GetOwner

`func (o *ListPartsOutput) GetOwner() Owner`

GetOwner returns the Owner field if non-nil, zero value otherwise.

### GetOwnerOk

`func (o *ListPartsOutput) GetOwnerOk() (*Owner, bool)`

GetOwnerOk returns a tuple with the Owner field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOwner

`func (o *ListPartsOutput) SetOwner(v Owner)`

SetOwner sets Owner field to given value.

### HasOwner

`func (o *ListPartsOutput) HasOwner() bool`

HasOwner returns a boolean if a field has been set.

### GetStorageClass

`func (o *ListPartsOutput) GetStorageClass() StorageClass`

GetStorageClass returns the StorageClass field if non-nil, zero value otherwise.

### GetStorageClassOk

`func (o *ListPartsOutput) GetStorageClassOk() (*StorageClass, bool)`

GetStorageClassOk returns a tuple with the StorageClass field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStorageClass

`func (o *ListPartsOutput) SetStorageClass(v StorageClass)`

SetStorageClass sets StorageClass field to given value.

### HasStorageClass

`func (o *ListPartsOutput) HasStorageClass() bool`

HasStorageClass returns a boolean if a field has been set.


