# PutBucketTaggingRequest

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Tagging** | [**PutBucketTaggingRequestTagging**](PutBucketTaggingRequestTagging.md) |  | |

## Methods

### NewPutBucketTaggingRequest

`func NewPutBucketTaggingRequest(tagging PutBucketTaggingRequestTagging, ) *PutBucketTaggingRequest`

NewPutBucketTaggingRequest instantiates a new PutBucketTaggingRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPutBucketTaggingRequestWithDefaults

`func NewPutBucketTaggingRequestWithDefaults() *PutBucketTaggingRequest`

NewPutBucketTaggingRequestWithDefaults instantiates a new PutBucketTaggingRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTagging

`func (o *PutBucketTaggingRequest) GetTagging() PutBucketTaggingRequestTagging`

GetTagging returns the Tagging field if non-nil, zero value otherwise.

### GetTaggingOk

`func (o *PutBucketTaggingRequest) GetTaggingOk() (*PutBucketTaggingRequestTagging, bool)`

GetTaggingOk returns a tuple with the Tagging field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTagging

`func (o *PutBucketTaggingRequest) SetTagging(v PutBucketTaggingRequestTagging)`

SetTagging sets Tagging field to given value.



