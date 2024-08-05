# Part

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**PartNumber** | Pointer to **int32** | Part number that identifies the part. | [optional] |
|**LastModified** | Pointer to [**time.Time**](time.Time.md) | Creation date of the object. | [optional] |
|**ETag** | Pointer to **string** | Entity tag that identifies the object&#39;s data. Objects with different object data will have different entity tags. The entity tag is an opaque string. The entity tag may or may not be an MD5 digest of the object data. If the entity tag is not an MD5 digest of the object data, it will contain one or more nonhexadecimal characters and/or will consist of less than 32 or more than 32 hexadecimal digits.  | [optional] |
|**Size** | Pointer to **int32** | Size in bytes of the object | [optional] |

## Methods

### NewPart

`func NewPart() *Part`

NewPart instantiates a new Part object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPartWithDefaults

`func NewPartWithDefaults() *Part`

NewPartWithDefaults instantiates a new Part object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPartNumber

`func (o *Part) GetPartNumber() int32`

GetPartNumber returns the PartNumber field if non-nil, zero value otherwise.

### GetPartNumberOk

`func (o *Part) GetPartNumberOk() (*int32, bool)`

GetPartNumberOk returns a tuple with the PartNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartNumber

`func (o *Part) SetPartNumber(v int32)`

SetPartNumber sets PartNumber field to given value.

### HasPartNumber

`func (o *Part) HasPartNumber() bool`

HasPartNumber returns a boolean if a field has been set.

### GetLastModified

`func (o *Part) GetLastModified() time.Time`

GetLastModified returns the LastModified field if non-nil, zero value otherwise.

### GetLastModifiedOk

`func (o *Part) GetLastModifiedOk() (*time.Time, bool)`

GetLastModifiedOk returns a tuple with the LastModified field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastModified

`func (o *Part) SetLastModified(v time.Time)`

SetLastModified sets LastModified field to given value.

### HasLastModified

`func (o *Part) HasLastModified() bool`

HasLastModified returns a boolean if a field has been set.

### GetETag

`func (o *Part) GetETag() string`

GetETag returns the ETag field if non-nil, zero value otherwise.

### GetETagOk

`func (o *Part) GetETagOk() (*string, bool)`

GetETagOk returns a tuple with the ETag field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetETag

`func (o *Part) SetETag(v string)`

SetETag sets ETag field to given value.

### HasETag

`func (o *Part) HasETag() bool`

HasETag returns a boolean if a field has been set.

### GetSize

`func (o *Part) GetSize() int32`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *Part) GetSizeOk() (*int32, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *Part) SetSize(v int32)`

SetSize sets Size field to given value.

### HasSize

`func (o *Part) HasSize() bool`

HasSize returns a boolean if a field has been set.


