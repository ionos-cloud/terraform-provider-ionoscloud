# ExampleCompleteMultipartUpload

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Parts** | Pointer to [**[]CompletedPart**](CompletedPart.md) | Array of CompletedPart data types. | [optional] |

## Methods

### NewExampleCompleteMultipartUpload

`func NewExampleCompleteMultipartUpload() *ExampleCompleteMultipartUpload`

NewExampleCompleteMultipartUpload instantiates a new ExampleCompleteMultipartUpload object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewExampleCompleteMultipartUploadWithDefaults

`func NewExampleCompleteMultipartUploadWithDefaults() *ExampleCompleteMultipartUpload`

NewExampleCompleteMultipartUploadWithDefaults instantiates a new ExampleCompleteMultipartUpload object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetParts

`func (o *ExampleCompleteMultipartUpload) GetParts() []CompletedPart`

GetParts returns the Parts field if non-nil, zero value otherwise.

### GetPartsOk

`func (o *ExampleCompleteMultipartUpload) GetPartsOk() (*[]CompletedPart, bool)`

GetPartsOk returns a tuple with the Parts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetParts

`func (o *ExampleCompleteMultipartUpload) SetParts(v []CompletedPart)`

SetParts sets Parts field to given value.

### HasParts

`func (o *ExampleCompleteMultipartUpload) HasParts() bool`

HasParts returns a boolean if a field has been set.


