# CopyObjectRequest

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**XAmzMeta** | Pointer to **map[string]string** | A map of metadata to store with the object in S3. | [optional] |

## Methods

### NewCopyObjectRequest

`func NewCopyObjectRequest() *CopyObjectRequest`

NewCopyObjectRequest instantiates a new CopyObjectRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCopyObjectRequestWithDefaults

`func NewCopyObjectRequestWithDefaults() *CopyObjectRequest`

NewCopyObjectRequestWithDefaults instantiates a new CopyObjectRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetXAmzMeta

`func (o *CopyObjectRequest) GetXAmzMeta() map[string]string`

GetXAmzMeta returns the XAmzMeta field if non-nil, zero value otherwise.

### GetXAmzMetaOk

`func (o *CopyObjectRequest) GetXAmzMetaOk() (*map[string]string, bool)`

GetXAmzMetaOk returns a tuple with the XAmzMeta field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetXAmzMeta

`func (o *CopyObjectRequest) SetXAmzMeta(v map[string]string)`

SetXAmzMeta sets XAmzMeta field to given value.

### HasXAmzMeta

`func (o *CopyObjectRequest) HasXAmzMeta() bool`

HasXAmzMeta returns a boolean if a field has been set.


