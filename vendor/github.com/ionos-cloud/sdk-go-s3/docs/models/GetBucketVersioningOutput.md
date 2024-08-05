# GetBucketVersioningOutput

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Status** | Pointer to [**BucketVersioningStatus**](BucketVersioningStatus.md) |  | [optional] |
|**MfaDelete** | Pointer to [**MfaDeleteStatus**](MfaDeleteStatus.md) |  | [optional] |

## Methods

### NewGetBucketVersioningOutput

`func NewGetBucketVersioningOutput() *GetBucketVersioningOutput`

NewGetBucketVersioningOutput instantiates a new GetBucketVersioningOutput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGetBucketVersioningOutputWithDefaults

`func NewGetBucketVersioningOutputWithDefaults() *GetBucketVersioningOutput`

NewGetBucketVersioningOutputWithDefaults instantiates a new GetBucketVersioningOutput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *GetBucketVersioningOutput) GetStatus() BucketVersioningStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *GetBucketVersioningOutput) GetStatusOk() (*BucketVersioningStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *GetBucketVersioningOutput) SetStatus(v BucketVersioningStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *GetBucketVersioningOutput) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetMfaDelete

`func (o *GetBucketVersioningOutput) GetMfaDelete() MfaDeleteStatus`

GetMfaDelete returns the MfaDelete field if non-nil, zero value otherwise.

### GetMfaDeleteOk

`func (o *GetBucketVersioningOutput) GetMfaDeleteOk() (*MfaDeleteStatus, bool)`

GetMfaDeleteOk returns a tuple with the MfaDelete field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMfaDelete

`func (o *GetBucketVersioningOutput) SetMfaDelete(v MfaDeleteStatus)`

SetMfaDelete sets MfaDelete field to given value.

### HasMfaDelete

`func (o *GetBucketVersioningOutput) HasMfaDelete() bool`

HasMfaDelete returns a boolean if a field has been set.


