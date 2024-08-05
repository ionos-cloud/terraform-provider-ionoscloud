# GetBucketEncryptionOutput

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**ServerSideEncryptionConfiguration** | Pointer to [**ServerSideEncryptionConfiguration**](ServerSideEncryptionConfiguration.md) |  | [optional] |

## Methods

### NewGetBucketEncryptionOutput

`func NewGetBucketEncryptionOutput() *GetBucketEncryptionOutput`

NewGetBucketEncryptionOutput instantiates a new GetBucketEncryptionOutput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGetBucketEncryptionOutputWithDefaults

`func NewGetBucketEncryptionOutputWithDefaults() *GetBucketEncryptionOutput`

NewGetBucketEncryptionOutputWithDefaults instantiates a new GetBucketEncryptionOutput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetServerSideEncryptionConfiguration

`func (o *GetBucketEncryptionOutput) GetServerSideEncryptionConfiguration() ServerSideEncryptionConfiguration`

GetServerSideEncryptionConfiguration returns the ServerSideEncryptionConfiguration field if non-nil, zero value otherwise.

### GetServerSideEncryptionConfigurationOk

`func (o *GetBucketEncryptionOutput) GetServerSideEncryptionConfigurationOk() (*ServerSideEncryptionConfiguration, bool)`

GetServerSideEncryptionConfigurationOk returns a tuple with the ServerSideEncryptionConfiguration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServerSideEncryptionConfiguration

`func (o *GetBucketEncryptionOutput) SetServerSideEncryptionConfiguration(v ServerSideEncryptionConfiguration)`

SetServerSideEncryptionConfiguration sets ServerSideEncryptionConfiguration field to given value.

### HasServerSideEncryptionConfiguration

`func (o *GetBucketEncryptionOutput) HasServerSideEncryptionConfiguration() bool`

HasServerSideEncryptionConfiguration returns a boolean if a field has been set.


