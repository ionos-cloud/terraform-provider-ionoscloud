# PutObjectTaggingRequest

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Tagging** | [**PutObjectTaggingRequestTagging**](PutObjectTaggingRequestTagging.md) |  | |

## Methods

### NewPutObjectTaggingRequest

`func NewPutObjectTaggingRequest(tagging PutObjectTaggingRequestTagging, ) *PutObjectTaggingRequest`

NewPutObjectTaggingRequest instantiates a new PutObjectTaggingRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPutObjectTaggingRequestWithDefaults

`func NewPutObjectTaggingRequestWithDefaults() *PutObjectTaggingRequest`

NewPutObjectTaggingRequestWithDefaults instantiates a new PutObjectTaggingRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTagging

`func (o *PutObjectTaggingRequest) GetTagging() PutObjectTaggingRequestTagging`

GetTagging returns the Tagging field if non-nil, zero value otherwise.

### GetTaggingOk

`func (o *PutObjectTaggingRequest) GetTaggingOk() (*PutObjectTaggingRequestTagging, bool)`

GetTaggingOk returns a tuple with the Tagging field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTagging

`func (o *PutObjectTaggingRequest) SetTagging(v PutObjectTaggingRequestTagging)`

SetTagging sets Tagging field to given value.



