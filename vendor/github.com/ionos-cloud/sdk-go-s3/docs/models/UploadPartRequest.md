# UploadPartRequest

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Body** | Pointer to **string** | Object data. | [optional] |

## Methods

### NewUploadPartRequest

`func NewUploadPartRequest() *UploadPartRequest`

NewUploadPartRequest instantiates a new UploadPartRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUploadPartRequestWithDefaults

`func NewUploadPartRequestWithDefaults() *UploadPartRequest`

NewUploadPartRequestWithDefaults instantiates a new UploadPartRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBody

`func (o *UploadPartRequest) GetBody() string`

GetBody returns the Body field if non-nil, zero value otherwise.

### GetBodyOk

`func (o *UploadPartRequest) GetBodyOk() (*string, bool)`

GetBodyOk returns a tuple with the Body field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBody

`func (o *UploadPartRequest) SetBody(v string)`

SetBody sets Body field to given value.

### HasBody

`func (o *UploadPartRequest) HasBody() bool`

HasBody returns a boolean if a field has been set.


