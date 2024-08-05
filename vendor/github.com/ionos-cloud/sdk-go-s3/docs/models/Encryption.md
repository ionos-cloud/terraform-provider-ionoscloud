# Encryption

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**EncryptionType** | **string** | The server-side encryption algorithm used when storing job results in IONOS S3 Object Storage (AES256). | |

## Methods

### NewEncryption

`func NewEncryption(encryptionType string, ) *Encryption`

NewEncryption instantiates a new Encryption object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEncryptionWithDefaults

`func NewEncryptionWithDefaults() *Encryption`

NewEncryptionWithDefaults instantiates a new Encryption object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEncryptionType

`func (o *Encryption) GetEncryptionType() string`

GetEncryptionType returns the EncryptionType field if non-nil, zero value otherwise.

### GetEncryptionTypeOk

`func (o *Encryption) GetEncryptionTypeOk() (*string, bool)`

GetEncryptionTypeOk returns a tuple with the EncryptionType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEncryptionType

`func (o *Encryption) SetEncryptionType(v string)`

SetEncryptionType sets EncryptionType field to given value.



