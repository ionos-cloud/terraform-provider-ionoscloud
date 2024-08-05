# ListObjectsOutput

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**IsTruncated** | Pointer to **bool** | A flag that indicates whether IONOS S3 Object Storage returned all of the results that satisfied the search criteria. If your results were truncated, you can make a follow-up paginated request using the NextKeyMarker and NextVersionIdMarker response parameters as a starting place in another request to return the rest of the results. | [optional] |
|**Marker** | Pointer to **string** | Indicates where in the bucket listing begins. Marker is included in the response if it was sent with the request. | [optional] |
|**NextMarker** | Pointer to **string** | When response is truncated (the IsTruncated element value in the response is true), you can use the key name in this field as marker in the subsequent request to get next set of objects. IONOS S3 Object Storage lists objects in alphabetical order Note: This element is returned only if you have delimiter request parameter specified. If response does not include the NextMarker and it is truncated, you can use the value of the last Key in the response as the marker in the subsequent request to get the next set of object keys. | [optional] |
|**Contents** | Pointer to [**[]Object**](Object.md) | Metadata about each object returned. | [optional] |
|**Name** | Pointer to **string** | The bucket name. | [optional] |
|**Prefix** | Pointer to **string** | Object key prefix that identifies one or more objects to which this rule applies. Replacement must be made for object keys containing special characters (such as carriage returns) when using XML requests. | [optional] |
|**Delimiter** | Pointer to **string** |  | [optional] |
|**MaxKeys** | Pointer to **int32** | The maximum number of keys returned in the response. By default the operation returns up to 1000 key names. The response might contain fewer keys but will never contain more. | [optional] |
|**CommonPrefixes** | Pointer to [**[]CommonPrefix**](CommonPrefix.md) | All of the keys rolled up into a common prefix count as a single return when calculating the number of returns. | [optional] |
|**EncodingType** | Pointer to [**EncodingType**](EncodingType.md) |  | [optional] |

## Methods

### NewListObjectsOutput

`func NewListObjectsOutput() *ListObjectsOutput`

NewListObjectsOutput instantiates a new ListObjectsOutput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewListObjectsOutputWithDefaults

`func NewListObjectsOutputWithDefaults() *ListObjectsOutput`

NewListObjectsOutputWithDefaults instantiates a new ListObjectsOutput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetIsTruncated

`func (o *ListObjectsOutput) GetIsTruncated() bool`

GetIsTruncated returns the IsTruncated field if non-nil, zero value otherwise.

### GetIsTruncatedOk

`func (o *ListObjectsOutput) GetIsTruncatedOk() (*bool, bool)`

GetIsTruncatedOk returns a tuple with the IsTruncated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsTruncated

`func (o *ListObjectsOutput) SetIsTruncated(v bool)`

SetIsTruncated sets IsTruncated field to given value.

### HasIsTruncated

`func (o *ListObjectsOutput) HasIsTruncated() bool`

HasIsTruncated returns a boolean if a field has been set.

### GetMarker

`func (o *ListObjectsOutput) GetMarker() string`

GetMarker returns the Marker field if non-nil, zero value otherwise.

### GetMarkerOk

`func (o *ListObjectsOutput) GetMarkerOk() (*string, bool)`

GetMarkerOk returns a tuple with the Marker field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMarker

`func (o *ListObjectsOutput) SetMarker(v string)`

SetMarker sets Marker field to given value.

### HasMarker

`func (o *ListObjectsOutput) HasMarker() bool`

HasMarker returns a boolean if a field has been set.

### GetNextMarker

`func (o *ListObjectsOutput) GetNextMarker() string`

GetNextMarker returns the NextMarker field if non-nil, zero value otherwise.

### GetNextMarkerOk

`func (o *ListObjectsOutput) GetNextMarkerOk() (*string, bool)`

GetNextMarkerOk returns a tuple with the NextMarker field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextMarker

`func (o *ListObjectsOutput) SetNextMarker(v string)`

SetNextMarker sets NextMarker field to given value.

### HasNextMarker

`func (o *ListObjectsOutput) HasNextMarker() bool`

HasNextMarker returns a boolean if a field has been set.

### GetContents

`func (o *ListObjectsOutput) GetContents() []Object`

GetContents returns the Contents field if non-nil, zero value otherwise.

### GetContentsOk

`func (o *ListObjectsOutput) GetContentsOk() (*[]Object, bool)`

GetContentsOk returns a tuple with the Contents field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContents

`func (o *ListObjectsOutput) SetContents(v []Object)`

SetContents sets Contents field to given value.

### HasContents

`func (o *ListObjectsOutput) HasContents() bool`

HasContents returns a boolean if a field has been set.

### GetName

`func (o *ListObjectsOutput) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ListObjectsOutput) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ListObjectsOutput) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ListObjectsOutput) HasName() bool`

HasName returns a boolean if a field has been set.

### GetPrefix

`func (o *ListObjectsOutput) GetPrefix() string`

GetPrefix returns the Prefix field if non-nil, zero value otherwise.

### GetPrefixOk

`func (o *ListObjectsOutput) GetPrefixOk() (*string, bool)`

GetPrefixOk returns a tuple with the Prefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrefix

`func (o *ListObjectsOutput) SetPrefix(v string)`

SetPrefix sets Prefix field to given value.

### HasPrefix

`func (o *ListObjectsOutput) HasPrefix() bool`

HasPrefix returns a boolean if a field has been set.

### GetDelimiter

`func (o *ListObjectsOutput) GetDelimiter() string`

GetDelimiter returns the Delimiter field if non-nil, zero value otherwise.

### GetDelimiterOk

`func (o *ListObjectsOutput) GetDelimiterOk() (*string, bool)`

GetDelimiterOk returns a tuple with the Delimiter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDelimiter

`func (o *ListObjectsOutput) SetDelimiter(v string)`

SetDelimiter sets Delimiter field to given value.

### HasDelimiter

`func (o *ListObjectsOutput) HasDelimiter() bool`

HasDelimiter returns a boolean if a field has been set.

### GetMaxKeys

`func (o *ListObjectsOutput) GetMaxKeys() int32`

GetMaxKeys returns the MaxKeys field if non-nil, zero value otherwise.

### GetMaxKeysOk

`func (o *ListObjectsOutput) GetMaxKeysOk() (*int32, bool)`

GetMaxKeysOk returns a tuple with the MaxKeys field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxKeys

`func (o *ListObjectsOutput) SetMaxKeys(v int32)`

SetMaxKeys sets MaxKeys field to given value.

### HasMaxKeys

`func (o *ListObjectsOutput) HasMaxKeys() bool`

HasMaxKeys returns a boolean if a field has been set.

### GetCommonPrefixes

`func (o *ListObjectsOutput) GetCommonPrefixes() []CommonPrefix`

GetCommonPrefixes returns the CommonPrefixes field if non-nil, zero value otherwise.

### GetCommonPrefixesOk

`func (o *ListObjectsOutput) GetCommonPrefixesOk() (*[]CommonPrefix, bool)`

GetCommonPrefixesOk returns a tuple with the CommonPrefixes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommonPrefixes

`func (o *ListObjectsOutput) SetCommonPrefixes(v []CommonPrefix)`

SetCommonPrefixes sets CommonPrefixes field to given value.

### HasCommonPrefixes

`func (o *ListObjectsOutput) HasCommonPrefixes() bool`

HasCommonPrefixes returns a boolean if a field has been set.

### GetEncodingType

`func (o *ListObjectsOutput) GetEncodingType() EncodingType`

GetEncodingType returns the EncodingType field if non-nil, zero value otherwise.

### GetEncodingTypeOk

`func (o *ListObjectsOutput) GetEncodingTypeOk() (*EncodingType, bool)`

GetEncodingTypeOk returns a tuple with the EncodingType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEncodingType

`func (o *ListObjectsOutput) SetEncodingType(v EncodingType)`

SetEncodingType sets EncodingType field to given value.

### HasEncodingType

`func (o *ListObjectsOutput) HasEncodingType() bool`

HasEncodingType returns a boolean if a field has been set.


