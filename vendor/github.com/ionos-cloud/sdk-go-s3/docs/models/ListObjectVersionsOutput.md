# ListObjectVersionsOutput

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**IsTruncated** | Pointer to **bool** | A flag that indicates whether IONOS S3 Object Storage returned all of the results that satisfied the search criteria. If your results were truncated, you can make a follow-up paginated request using the NextKeyMarker and NextVersionIdMarker response parameters as a starting place in another request to return the rest of the results. | [optional] |
|**KeyMarker** | Pointer to **string** | Marks the last key returned in a truncated response. | [optional] |
|**VersionIdMarker** | Pointer to **string** | Marks the last version of the key returned in a truncated response. | [optional] |
|**NextKeyMarker** | Pointer to **string** | When the number of responses exceeds the value of &#x60;MaxKeys&#x60;, &#x60;NextKeyMarker&#x60; specifies the first key not returned that satisfies the search criteria. Use this value for the key-marker request parameter in a subsequent request. | [optional] |
|**NextVersionIdMarker** | Pointer to **string** | When the number of responses exceeds the value of &#x60;MaxKeys&#x60;, &#x60;NextVersionIdMarker&#x60; specifies the first object version not returned that satisfies the search criteria. Use this value for the version-id-marker request parameter in a subsequent request. | [optional] |
|**Versions** | Pointer to [**[]ObjectVersion**](ObjectVersion.md) | Container for version information. | [optional] |
|**Name** | Pointer to **string** | The bucket name. | [optional] |
|**Prefix** | Pointer to **string** | Selects objects that start with the value supplied by this parameter. | [optional] |
|**Delimiter** | Pointer to **string** |  | [optional] |
|**MaxKeys** | Pointer to **int32** | The maximum number of keys returned in the response. By default the operation returns up to 1000 key names. The response might contain fewer keys but will never contain more. | [optional] |
|**CommonPrefixes** | Pointer to [**[]CommonPrefix**](CommonPrefix.md) | All of the keys rolled up into a common prefix count as a single return when calculating the number of returns. | [optional] |
|**EncodingType** | Pointer to [**EncodingType**](EncodingType.md) |  | [optional] |

## Methods

### NewListObjectVersionsOutput

`func NewListObjectVersionsOutput() *ListObjectVersionsOutput`

NewListObjectVersionsOutput instantiates a new ListObjectVersionsOutput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewListObjectVersionsOutputWithDefaults

`func NewListObjectVersionsOutputWithDefaults() *ListObjectVersionsOutput`

NewListObjectVersionsOutputWithDefaults instantiates a new ListObjectVersionsOutput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetIsTruncated

`func (o *ListObjectVersionsOutput) GetIsTruncated() bool`

GetIsTruncated returns the IsTruncated field if non-nil, zero value otherwise.

### GetIsTruncatedOk

`func (o *ListObjectVersionsOutput) GetIsTruncatedOk() (*bool, bool)`

GetIsTruncatedOk returns a tuple with the IsTruncated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsTruncated

`func (o *ListObjectVersionsOutput) SetIsTruncated(v bool)`

SetIsTruncated sets IsTruncated field to given value.

### HasIsTruncated

`func (o *ListObjectVersionsOutput) HasIsTruncated() bool`

HasIsTruncated returns a boolean if a field has been set.

### GetKeyMarker

`func (o *ListObjectVersionsOutput) GetKeyMarker() string`

GetKeyMarker returns the KeyMarker field if non-nil, zero value otherwise.

### GetKeyMarkerOk

`func (o *ListObjectVersionsOutput) GetKeyMarkerOk() (*string, bool)`

GetKeyMarkerOk returns a tuple with the KeyMarker field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeyMarker

`func (o *ListObjectVersionsOutput) SetKeyMarker(v string)`

SetKeyMarker sets KeyMarker field to given value.

### HasKeyMarker

`func (o *ListObjectVersionsOutput) HasKeyMarker() bool`

HasKeyMarker returns a boolean if a field has been set.

### GetVersionIdMarker

`func (o *ListObjectVersionsOutput) GetVersionIdMarker() string`

GetVersionIdMarker returns the VersionIdMarker field if non-nil, zero value otherwise.

### GetVersionIdMarkerOk

`func (o *ListObjectVersionsOutput) GetVersionIdMarkerOk() (*string, bool)`

GetVersionIdMarkerOk returns a tuple with the VersionIdMarker field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersionIdMarker

`func (o *ListObjectVersionsOutput) SetVersionIdMarker(v string)`

SetVersionIdMarker sets VersionIdMarker field to given value.

### HasVersionIdMarker

`func (o *ListObjectVersionsOutput) HasVersionIdMarker() bool`

HasVersionIdMarker returns a boolean if a field has been set.

### GetNextKeyMarker

`func (o *ListObjectVersionsOutput) GetNextKeyMarker() string`

GetNextKeyMarker returns the NextKeyMarker field if non-nil, zero value otherwise.

### GetNextKeyMarkerOk

`func (o *ListObjectVersionsOutput) GetNextKeyMarkerOk() (*string, bool)`

GetNextKeyMarkerOk returns a tuple with the NextKeyMarker field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextKeyMarker

`func (o *ListObjectVersionsOutput) SetNextKeyMarker(v string)`

SetNextKeyMarker sets NextKeyMarker field to given value.

### HasNextKeyMarker

`func (o *ListObjectVersionsOutput) HasNextKeyMarker() bool`

HasNextKeyMarker returns a boolean if a field has been set.

### GetNextVersionIdMarker

`func (o *ListObjectVersionsOutput) GetNextVersionIdMarker() string`

GetNextVersionIdMarker returns the NextVersionIdMarker field if non-nil, zero value otherwise.

### GetNextVersionIdMarkerOk

`func (o *ListObjectVersionsOutput) GetNextVersionIdMarkerOk() (*string, bool)`

GetNextVersionIdMarkerOk returns a tuple with the NextVersionIdMarker field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextVersionIdMarker

`func (o *ListObjectVersionsOutput) SetNextVersionIdMarker(v string)`

SetNextVersionIdMarker sets NextVersionIdMarker field to given value.

### HasNextVersionIdMarker

`func (o *ListObjectVersionsOutput) HasNextVersionIdMarker() bool`

HasNextVersionIdMarker returns a boolean if a field has been set.

### GetVersions

`func (o *ListObjectVersionsOutput) GetVersions() []ObjectVersion`

GetVersions returns the Versions field if non-nil, zero value otherwise.

### GetVersionsOk

`func (o *ListObjectVersionsOutput) GetVersionsOk() (*[]ObjectVersion, bool)`

GetVersionsOk returns a tuple with the Versions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersions

`func (o *ListObjectVersionsOutput) SetVersions(v []ObjectVersion)`

SetVersions sets Versions field to given value.

### HasVersions

`func (o *ListObjectVersionsOutput) HasVersions() bool`

HasVersions returns a boolean if a field has been set.

### GetName

`func (o *ListObjectVersionsOutput) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ListObjectVersionsOutput) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ListObjectVersionsOutput) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ListObjectVersionsOutput) HasName() bool`

HasName returns a boolean if a field has been set.

### GetPrefix

`func (o *ListObjectVersionsOutput) GetPrefix() string`

GetPrefix returns the Prefix field if non-nil, zero value otherwise.

### GetPrefixOk

`func (o *ListObjectVersionsOutput) GetPrefixOk() (*string, bool)`

GetPrefixOk returns a tuple with the Prefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrefix

`func (o *ListObjectVersionsOutput) SetPrefix(v string)`

SetPrefix sets Prefix field to given value.

### HasPrefix

`func (o *ListObjectVersionsOutput) HasPrefix() bool`

HasPrefix returns a boolean if a field has been set.

### GetDelimiter

`func (o *ListObjectVersionsOutput) GetDelimiter() string`

GetDelimiter returns the Delimiter field if non-nil, zero value otherwise.

### GetDelimiterOk

`func (o *ListObjectVersionsOutput) GetDelimiterOk() (*string, bool)`

GetDelimiterOk returns a tuple with the Delimiter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDelimiter

`func (o *ListObjectVersionsOutput) SetDelimiter(v string)`

SetDelimiter sets Delimiter field to given value.

### HasDelimiter

`func (o *ListObjectVersionsOutput) HasDelimiter() bool`

HasDelimiter returns a boolean if a field has been set.

### GetMaxKeys

`func (o *ListObjectVersionsOutput) GetMaxKeys() int32`

GetMaxKeys returns the MaxKeys field if non-nil, zero value otherwise.

### GetMaxKeysOk

`func (o *ListObjectVersionsOutput) GetMaxKeysOk() (*int32, bool)`

GetMaxKeysOk returns a tuple with the MaxKeys field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxKeys

`func (o *ListObjectVersionsOutput) SetMaxKeys(v int32)`

SetMaxKeys sets MaxKeys field to given value.

### HasMaxKeys

`func (o *ListObjectVersionsOutput) HasMaxKeys() bool`

HasMaxKeys returns a boolean if a field has been set.

### GetCommonPrefixes

`func (o *ListObjectVersionsOutput) GetCommonPrefixes() []CommonPrefix`

GetCommonPrefixes returns the CommonPrefixes field if non-nil, zero value otherwise.

### GetCommonPrefixesOk

`func (o *ListObjectVersionsOutput) GetCommonPrefixesOk() (*[]CommonPrefix, bool)`

GetCommonPrefixesOk returns a tuple with the CommonPrefixes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommonPrefixes

`func (o *ListObjectVersionsOutput) SetCommonPrefixes(v []CommonPrefix)`

SetCommonPrefixes sets CommonPrefixes field to given value.

### HasCommonPrefixes

`func (o *ListObjectVersionsOutput) HasCommonPrefixes() bool`

HasCommonPrefixes returns a boolean if a field has been set.

### GetEncodingType

`func (o *ListObjectVersionsOutput) GetEncodingType() EncodingType`

GetEncodingType returns the EncodingType field if non-nil, zero value otherwise.

### GetEncodingTypeOk

`func (o *ListObjectVersionsOutput) GetEncodingTypeOk() (*EncodingType, bool)`

GetEncodingTypeOk returns a tuple with the EncodingType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEncodingType

`func (o *ListObjectVersionsOutput) SetEncodingType(v EncodingType)`

SetEncodingType sets EncodingType field to given value.

### HasEncodingType

`func (o *ListObjectVersionsOutput) HasEncodingType() bool`

HasEncodingType returns a boolean if a field has been set.


