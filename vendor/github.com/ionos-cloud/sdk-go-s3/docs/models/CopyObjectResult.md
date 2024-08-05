# CopyObjectResult

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**ETag** | Pointer to **string** | Entity tag that identifies the object&#39;s data. Objects with different object data will have different entity tags. The entity tag is an opaque string. The entity tag may or may not be an MD5 digest of the object data. If the entity tag is not an MD5 digest of the object data, it will contain one or more nonhexadecimal characters and/or will consist of less than 32 or more than 32 hexadecimal digits.  | [optional] |
|**LastModified** | Pointer to [**time.Time**](time.Time.md) | Creation date of the object. | [optional] |

## Methods

### NewCopyObjectResult

`func NewCopyObjectResult() *CopyObjectResult`

NewCopyObjectResult instantiates a new CopyObjectResult object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCopyObjectResultWithDefaults

`func NewCopyObjectResultWithDefaults() *CopyObjectResult`

NewCopyObjectResultWithDefaults instantiates a new CopyObjectResult object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetETag

`func (o *CopyObjectResult) GetETag() string`

GetETag returns the ETag field if non-nil, zero value otherwise.

### GetETagOk

`func (o *CopyObjectResult) GetETagOk() (*string, bool)`

GetETagOk returns a tuple with the ETag field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetETag

`func (o *CopyObjectResult) SetETag(v string)`

SetETag sets ETag field to given value.

### HasETag

`func (o *CopyObjectResult) HasETag() bool`

HasETag returns a boolean if a field has been set.

### GetLastModified

`func (o *CopyObjectResult) GetLastModified() time.Time`

GetLastModified returns the LastModified field if non-nil, zero value otherwise.

### GetLastModifiedOk

`func (o *CopyObjectResult) GetLastModifiedOk() (*time.Time, bool)`

GetLastModifiedOk returns a tuple with the LastModified field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastModified

`func (o *CopyObjectResult) SetLastModified(v time.Time)`

SetLastModified sets LastModified field to given value.

### HasLastModified

`func (o *CopyObjectResult) HasLastModified() bool`

HasLastModified returns a boolean if a field has been set.


