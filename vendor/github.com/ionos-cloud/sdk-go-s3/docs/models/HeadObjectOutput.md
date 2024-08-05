# HeadObjectOutput

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Metadata** | Pointer to **map[string]string** | A map of metadata to store with the object. Each key must start with  &#x60;x-amz-meta-&#x60; prefix.  | [optional] |

## Methods

### NewHeadObjectOutput

`func NewHeadObjectOutput() *HeadObjectOutput`

NewHeadObjectOutput instantiates a new HeadObjectOutput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewHeadObjectOutputWithDefaults

`func NewHeadObjectOutputWithDefaults() *HeadObjectOutput`

NewHeadObjectOutputWithDefaults instantiates a new HeadObjectOutput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMetadata

`func (o *HeadObjectOutput) GetMetadata() map[string]string`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *HeadObjectOutput) GetMetadataOk() (*map[string]string, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *HeadObjectOutput) SetMetadata(v map[string]string)`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *HeadObjectOutput) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.


