# Object

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Key** | Pointer to **string** | The object key. | [optional] |
|**LastModified** | Pointer to [**time.Time**](time.Time.md) | Creation date of the object. | [optional] |
|**StorageClass** | Pointer to [**ObjectStorageClass**](ObjectStorageClass.md) |  | [optional] |
|**Size** | Pointer to **int32** | Size in bytes of the object | [optional] |
|**ETag** | Pointer to **string** | Entity tag that identifies the object&#39;s data. Objects with different object data will have different entity tags. The entity tag is an opaque string. The entity tag may or may not be an MD5 digest of the object data. If the entity tag is not an MD5 digest of the object data, it will contain one or more nonhexadecimal characters and/or will consist of less than 32 or more than 32 hexadecimal digits.  | [optional] |
|**Owner** | Pointer to [**Owner**](Owner.md) |  | [optional] |

## Methods

### NewObject

`func NewObject() *Object`

NewObject instantiates a new Object object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewObjectWithDefaults

`func NewObjectWithDefaults() *Object`

NewObjectWithDefaults instantiates a new Object object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKey

`func (o *Object) GetKey() string`

GetKey returns the Key field if non-nil, zero value otherwise.

### GetKeyOk

`func (o *Object) GetKeyOk() (*string, bool)`

GetKeyOk returns a tuple with the Key field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKey

`func (o *Object) SetKey(v string)`

SetKey sets Key field to given value.

### HasKey

`func (o *Object) HasKey() bool`

HasKey returns a boolean if a field has been set.

### GetLastModified

`func (o *Object) GetLastModified() time.Time`

GetLastModified returns the LastModified field if non-nil, zero value otherwise.

### GetLastModifiedOk

`func (o *Object) GetLastModifiedOk() (*time.Time, bool)`

GetLastModifiedOk returns a tuple with the LastModified field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastModified

`func (o *Object) SetLastModified(v time.Time)`

SetLastModified sets LastModified field to given value.

### HasLastModified

`func (o *Object) HasLastModified() bool`

HasLastModified returns a boolean if a field has been set.

### GetStorageClass

`func (o *Object) GetStorageClass() ObjectStorageClass`

GetStorageClass returns the StorageClass field if non-nil, zero value otherwise.

### GetStorageClassOk

`func (o *Object) GetStorageClassOk() (*ObjectStorageClass, bool)`

GetStorageClassOk returns a tuple with the StorageClass field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStorageClass

`func (o *Object) SetStorageClass(v ObjectStorageClass)`

SetStorageClass sets StorageClass field to given value.

### HasStorageClass

`func (o *Object) HasStorageClass() bool`

HasStorageClass returns a boolean if a field has been set.

### GetSize

`func (o *Object) GetSize() int32`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *Object) GetSizeOk() (*int32, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *Object) SetSize(v int32)`

SetSize sets Size field to given value.

### HasSize

`func (o *Object) HasSize() bool`

HasSize returns a boolean if a field has been set.

### GetETag

`func (o *Object) GetETag() string`

GetETag returns the ETag field if non-nil, zero value otherwise.

### GetETagOk

`func (o *Object) GetETagOk() (*string, bool)`

GetETagOk returns a tuple with the ETag field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetETag

`func (o *Object) SetETag(v string)`

SetETag sets ETag field to given value.

### HasETag

`func (o *Object) HasETag() bool`

HasETag returns a boolean if a field has been set.

### GetOwner

`func (o *Object) GetOwner() Owner`

GetOwner returns the Owner field if non-nil, zero value otherwise.

### GetOwnerOk

`func (o *Object) GetOwnerOk() (*Owner, bool)`

GetOwnerOk returns a tuple with the Owner field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOwner

`func (o *Object) SetOwner(v Owner)`

SetOwner sets Owner field to given value.

### HasOwner

`func (o *Object) HasOwner() bool`

HasOwner returns a boolean if a field has been set.


