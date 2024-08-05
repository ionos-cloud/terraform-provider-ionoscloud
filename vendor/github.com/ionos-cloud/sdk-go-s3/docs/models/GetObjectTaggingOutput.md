# GetObjectTaggingOutput

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**TagSet** | [**[]Tag**](Tag.md) | Contains the tag set. | |

## Methods

### NewGetObjectTaggingOutput

`func NewGetObjectTaggingOutput(tagSet []Tag, ) *GetObjectTaggingOutput`

NewGetObjectTaggingOutput instantiates a new GetObjectTaggingOutput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGetObjectTaggingOutputWithDefaults

`func NewGetObjectTaggingOutputWithDefaults() *GetObjectTaggingOutput`

NewGetObjectTaggingOutputWithDefaults instantiates a new GetObjectTaggingOutput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTagSet

`func (o *GetObjectTaggingOutput) GetTagSet() []Tag`

GetTagSet returns the TagSet field if non-nil, zero value otherwise.

### GetTagSetOk

`func (o *GetObjectTaggingOutput) GetTagSetOk() (*[]Tag, bool)`

GetTagSetOk returns a tuple with the TagSet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTagSet

`func (o *GetObjectTaggingOutput) SetTagSet(v []Tag)`

SetTagSet sets TagSet field to given value.



