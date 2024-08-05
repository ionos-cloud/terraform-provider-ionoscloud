# MetadataEntry

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Name** | Pointer to **string** | Name of the Object. | [optional] |
|**Value** | Pointer to **string** | Value of the Object. | [optional] |

## Methods

### NewMetadataEntry

`func NewMetadataEntry() *MetadataEntry`

NewMetadataEntry instantiates a new MetadataEntry object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMetadataEntryWithDefaults

`func NewMetadataEntryWithDefaults() *MetadataEntry`

NewMetadataEntryWithDefaults instantiates a new MetadataEntry object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *MetadataEntry) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *MetadataEntry) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *MetadataEntry) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *MetadataEntry) HasName() bool`

HasName returns a boolean if a field has been set.

### GetValue

`func (o *MetadataEntry) GetValue() string`

GetValue returns the Value field if non-nil, zero value otherwise.

### GetValueOk

`func (o *MetadataEntry) GetValueOk() (*string, bool)`

GetValueOk returns a tuple with the Value field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValue

`func (o *MetadataEntry) SetValue(v string)`

SetValue sets Value field to given value.

### HasValue

`func (o *MetadataEntry) HasValue() bool`

HasValue returns a boolean if a field has been set.


