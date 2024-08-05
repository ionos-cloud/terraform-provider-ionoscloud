# ListMultipartUploadsOutput

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Bucket** | Pointer to **string** | The bucket name. | [optional] |
|**KeyMarker** | Pointer to **string** | The key at or after which the listing began. | [optional] |
|**UploadIdMarker** | Pointer to **string** | Upload ID after which listing began. | [optional] |
|**NextKeyMarker** | Pointer to **string** | When a list is truncated, this element specifies the value that should be used for the key-marker request parameter in a subsequent request. | [optional] |
|**Prefix** | Pointer to **string** | When a prefix is provided in the request, this field contains the specified prefix. The result contains only keys starting with the specified prefix. | [optional] |
|**Delimiter** | Pointer to **string** | Contains the delimiter you specified in the request. If you don&#39;t specify a delimiter in your request, this element is absent from the response. | [optional] |
|**NextUploadIdMarker** | Pointer to **string** | When a list is truncated, this element specifies the value that should be used for the &#x60;upload-id-marker&#x60; request parameter in a subsequent request. | [optional] |
|**MaxUploads** | Pointer to **int32** | Maximum number of multipart uploads that could have been included in the response. | [optional] |
|**IsTruncated** | Pointer to **bool** | A flag that indicates whether IONOS S3 Object Storage returned all of the results that satisfied the search criteria. If your results were truncated, you can make a follow-up paginated request using the NextKeyMarker and NextVersionIdMarker response parameters as a starting place in another request to return the rest of the results. | [optional] |
|**Uploads** | Pointer to [**[]MultipartUpload**](MultipartUpload.md) | Container for elements related to a particular multipart upload. A response can contain zero or more &#x60;Upload&#x60; elements. | [optional] |
|**CommonPrefixes** | Pointer to [**[]CommonPrefix**](CommonPrefix.md) | All of the keys rolled up into a common prefix count as a single return when calculating the number of returns. | [optional] |
|**EncodingType** | Pointer to [**EncodingType**](EncodingType.md) |  | [optional] |

## Methods

### NewListMultipartUploadsOutput

`func NewListMultipartUploadsOutput() *ListMultipartUploadsOutput`

NewListMultipartUploadsOutput instantiates a new ListMultipartUploadsOutput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewListMultipartUploadsOutputWithDefaults

`func NewListMultipartUploadsOutputWithDefaults() *ListMultipartUploadsOutput`

NewListMultipartUploadsOutputWithDefaults instantiates a new ListMultipartUploadsOutput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBucket

`func (o *ListMultipartUploadsOutput) GetBucket() string`

GetBucket returns the Bucket field if non-nil, zero value otherwise.

### GetBucketOk

`func (o *ListMultipartUploadsOutput) GetBucketOk() (*string, bool)`

GetBucketOk returns a tuple with the Bucket field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBucket

`func (o *ListMultipartUploadsOutput) SetBucket(v string)`

SetBucket sets Bucket field to given value.

### HasBucket

`func (o *ListMultipartUploadsOutput) HasBucket() bool`

HasBucket returns a boolean if a field has been set.

### GetKeyMarker

`func (o *ListMultipartUploadsOutput) GetKeyMarker() string`

GetKeyMarker returns the KeyMarker field if non-nil, zero value otherwise.

### GetKeyMarkerOk

`func (o *ListMultipartUploadsOutput) GetKeyMarkerOk() (*string, bool)`

GetKeyMarkerOk returns a tuple with the KeyMarker field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeyMarker

`func (o *ListMultipartUploadsOutput) SetKeyMarker(v string)`

SetKeyMarker sets KeyMarker field to given value.

### HasKeyMarker

`func (o *ListMultipartUploadsOutput) HasKeyMarker() bool`

HasKeyMarker returns a boolean if a field has been set.

### GetUploadIdMarker

`func (o *ListMultipartUploadsOutput) GetUploadIdMarker() string`

GetUploadIdMarker returns the UploadIdMarker field if non-nil, zero value otherwise.

### GetUploadIdMarkerOk

`func (o *ListMultipartUploadsOutput) GetUploadIdMarkerOk() (*string, bool)`

GetUploadIdMarkerOk returns a tuple with the UploadIdMarker field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUploadIdMarker

`func (o *ListMultipartUploadsOutput) SetUploadIdMarker(v string)`

SetUploadIdMarker sets UploadIdMarker field to given value.

### HasUploadIdMarker

`func (o *ListMultipartUploadsOutput) HasUploadIdMarker() bool`

HasUploadIdMarker returns a boolean if a field has been set.

### GetNextKeyMarker

`func (o *ListMultipartUploadsOutput) GetNextKeyMarker() string`

GetNextKeyMarker returns the NextKeyMarker field if non-nil, zero value otherwise.

### GetNextKeyMarkerOk

`func (o *ListMultipartUploadsOutput) GetNextKeyMarkerOk() (*string, bool)`

GetNextKeyMarkerOk returns a tuple with the NextKeyMarker field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextKeyMarker

`func (o *ListMultipartUploadsOutput) SetNextKeyMarker(v string)`

SetNextKeyMarker sets NextKeyMarker field to given value.

### HasNextKeyMarker

`func (o *ListMultipartUploadsOutput) HasNextKeyMarker() bool`

HasNextKeyMarker returns a boolean if a field has been set.

### GetPrefix

`func (o *ListMultipartUploadsOutput) GetPrefix() string`

GetPrefix returns the Prefix field if non-nil, zero value otherwise.

### GetPrefixOk

`func (o *ListMultipartUploadsOutput) GetPrefixOk() (*string, bool)`

GetPrefixOk returns a tuple with the Prefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrefix

`func (o *ListMultipartUploadsOutput) SetPrefix(v string)`

SetPrefix sets Prefix field to given value.

### HasPrefix

`func (o *ListMultipartUploadsOutput) HasPrefix() bool`

HasPrefix returns a boolean if a field has been set.

### GetDelimiter

`func (o *ListMultipartUploadsOutput) GetDelimiter() string`

GetDelimiter returns the Delimiter field if non-nil, zero value otherwise.

### GetDelimiterOk

`func (o *ListMultipartUploadsOutput) GetDelimiterOk() (*string, bool)`

GetDelimiterOk returns a tuple with the Delimiter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDelimiter

`func (o *ListMultipartUploadsOutput) SetDelimiter(v string)`

SetDelimiter sets Delimiter field to given value.

### HasDelimiter

`func (o *ListMultipartUploadsOutput) HasDelimiter() bool`

HasDelimiter returns a boolean if a field has been set.

### GetNextUploadIdMarker

`func (o *ListMultipartUploadsOutput) GetNextUploadIdMarker() string`

GetNextUploadIdMarker returns the NextUploadIdMarker field if non-nil, zero value otherwise.

### GetNextUploadIdMarkerOk

`func (o *ListMultipartUploadsOutput) GetNextUploadIdMarkerOk() (*string, bool)`

GetNextUploadIdMarkerOk returns a tuple with the NextUploadIdMarker field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextUploadIdMarker

`func (o *ListMultipartUploadsOutput) SetNextUploadIdMarker(v string)`

SetNextUploadIdMarker sets NextUploadIdMarker field to given value.

### HasNextUploadIdMarker

`func (o *ListMultipartUploadsOutput) HasNextUploadIdMarker() bool`

HasNextUploadIdMarker returns a boolean if a field has been set.

### GetMaxUploads

`func (o *ListMultipartUploadsOutput) GetMaxUploads() int32`

GetMaxUploads returns the MaxUploads field if non-nil, zero value otherwise.

### GetMaxUploadsOk

`func (o *ListMultipartUploadsOutput) GetMaxUploadsOk() (*int32, bool)`

GetMaxUploadsOk returns a tuple with the MaxUploads field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxUploads

`func (o *ListMultipartUploadsOutput) SetMaxUploads(v int32)`

SetMaxUploads sets MaxUploads field to given value.

### HasMaxUploads

`func (o *ListMultipartUploadsOutput) HasMaxUploads() bool`

HasMaxUploads returns a boolean if a field has been set.

### GetIsTruncated

`func (o *ListMultipartUploadsOutput) GetIsTruncated() bool`

GetIsTruncated returns the IsTruncated field if non-nil, zero value otherwise.

### GetIsTruncatedOk

`func (o *ListMultipartUploadsOutput) GetIsTruncatedOk() (*bool, bool)`

GetIsTruncatedOk returns a tuple with the IsTruncated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsTruncated

`func (o *ListMultipartUploadsOutput) SetIsTruncated(v bool)`

SetIsTruncated sets IsTruncated field to given value.

### HasIsTruncated

`func (o *ListMultipartUploadsOutput) HasIsTruncated() bool`

HasIsTruncated returns a boolean if a field has been set.

### GetUploads

`func (o *ListMultipartUploadsOutput) GetUploads() []MultipartUpload`

GetUploads returns the Uploads field if non-nil, zero value otherwise.

### GetUploadsOk

`func (o *ListMultipartUploadsOutput) GetUploadsOk() (*[]MultipartUpload, bool)`

GetUploadsOk returns a tuple with the Uploads field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUploads

`func (o *ListMultipartUploadsOutput) SetUploads(v []MultipartUpload)`

SetUploads sets Uploads field to given value.

### HasUploads

`func (o *ListMultipartUploadsOutput) HasUploads() bool`

HasUploads returns a boolean if a field has been set.

### GetCommonPrefixes

`func (o *ListMultipartUploadsOutput) GetCommonPrefixes() []CommonPrefix`

GetCommonPrefixes returns the CommonPrefixes field if non-nil, zero value otherwise.

### GetCommonPrefixesOk

`func (o *ListMultipartUploadsOutput) GetCommonPrefixesOk() (*[]CommonPrefix, bool)`

GetCommonPrefixesOk returns a tuple with the CommonPrefixes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommonPrefixes

`func (o *ListMultipartUploadsOutput) SetCommonPrefixes(v []CommonPrefix)`

SetCommonPrefixes sets CommonPrefixes field to given value.

### HasCommonPrefixes

`func (o *ListMultipartUploadsOutput) HasCommonPrefixes() bool`

HasCommonPrefixes returns a boolean if a field has been set.

### GetEncodingType

`func (o *ListMultipartUploadsOutput) GetEncodingType() EncodingType`

GetEncodingType returns the EncodingType field if non-nil, zero value otherwise.

### GetEncodingTypeOk

`func (o *ListMultipartUploadsOutput) GetEncodingTypeOk() (*EncodingType, bool)`

GetEncodingTypeOk returns a tuple with the EncodingType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEncodingType

`func (o *ListMultipartUploadsOutput) SetEncodingType(v EncodingType)`

SetEncodingType sets EncodingType field to given value.

### HasEncodingType

`func (o *ListMultipartUploadsOutput) HasEncodingType() bool`

HasEncodingType returns a boolean if a field has been set.


