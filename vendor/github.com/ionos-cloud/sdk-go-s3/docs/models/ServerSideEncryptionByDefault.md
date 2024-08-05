# ServerSideEncryptionByDefault

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**SSEAlgorithm** | [**ServerSideEncryption**](ServerSideEncryption.md) |  | |

## Methods

### NewServerSideEncryptionByDefault

`func NewServerSideEncryptionByDefault(sSEAlgorithm ServerSideEncryption, ) *ServerSideEncryptionByDefault`

NewServerSideEncryptionByDefault instantiates a new ServerSideEncryptionByDefault object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServerSideEncryptionByDefaultWithDefaults

`func NewServerSideEncryptionByDefaultWithDefaults() *ServerSideEncryptionByDefault`

NewServerSideEncryptionByDefaultWithDefaults instantiates a new ServerSideEncryptionByDefault object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSSEAlgorithm

`func (o *ServerSideEncryptionByDefault) GetSSEAlgorithm() ServerSideEncryption`

GetSSEAlgorithm returns the SSEAlgorithm field if non-nil, zero value otherwise.

### GetSSEAlgorithmOk

`func (o *ServerSideEncryptionByDefault) GetSSEAlgorithmOk() (*ServerSideEncryption, bool)`

GetSSEAlgorithmOk returns a tuple with the SSEAlgorithm field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSSEAlgorithm

`func (o *ServerSideEncryptionByDefault) SetSSEAlgorithm(v ServerSideEncryption)`

SetSSEAlgorithm sets SSEAlgorithm field to given value.



