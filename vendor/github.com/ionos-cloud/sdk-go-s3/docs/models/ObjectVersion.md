# ObjectVersion

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**ETag** | Pointer to **string** | Entity tag that identifies the object&#39;s data. Objects with different object data will have different entity tags. The entity tag is an opaque string. The entity tag may or may not be an MD5 digest of the object data. If the entity tag is not an MD5 digest of the object data, it will contain one or more nonhexadecimal characters and/or will consist of less than 32 or more than 32 hexadecimal digits.  | [optional] |
|**Size** | Pointer to **int32** | Size in bytes of the object | [optional] |
|**StorageClass** | Pointer to [**ObjectVersionStorageClass**](ObjectVersionStorageClass.md) |  | [optional] |
|**Key** | Pointer to **string** | The object key. | [optional] |
|**VersionId** | Pointer to **string** | Version ID of an object. | [optional] |
|**IsLatest** | Pointer to **bool** | Specifies whether the object is (true) or is not (false) the latest version of an object. | [optional] |
|**LastModified** | Pointer to [**time.Time**](time.Time.md) | Creation date of the object. | [optional] |
|**Owner** | Pointer to [**Owner**](Owner.md) |  | [optional] |

## Methods

### NewObjectVersion

`func NewObjectVersion() *ObjectVersion`

NewObjectVersion instantiates a new ObjectVersion object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewObjectVersionWithDefaults

`func NewObjectVersionWithDefaults() *ObjectVersion`

NewObjectVersionWithDefaults instantiates a new ObjectVersion object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetETag

`func (o *ObjectVersion) GetETag() string`

GetETag returns the ETag field if non-nil, zero value otherwise.

### GetETagOk

`func (o *ObjectVersion) GetETagOk() (*string, bool)`

GetETagOk returns a tuple with the ETag field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetETag

`func (o *ObjectVersion) SetETag(v string)`

SetETag sets ETag field to given value.

### HasETag

`func (o *ObjectVersion) HasETag() bool`

HasETag returns a boolean if a field has been set.

### GetSize

`func (o *ObjectVersion) GetSize() int32`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *ObjectVersion) GetSizeOk() (*int32, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *ObjectVersion) SetSize(v int32)`

SetSize sets Size field to given value.

### HasSize

`func (o *ObjectVersion) HasSize() bool`

HasSize returns a boolean if a field has been set.

### GetStorageClass

`func (o *ObjectVersion) GetStorageClass() ObjectVersionStorageClass`

GetStorageClass returns the StorageClass field if non-nil, zero value otherwise.

### GetStorageClassOk

`func (o *ObjectVersion) GetStorageClassOk() (*ObjectVersionStorageClass, bool)`

GetStorageClassOk returns a tuple with the StorageClass field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStorageClass

`func (o *ObjectVersion) SetStorageClass(v ObjectVersionStorageClass)`

SetStorageClass sets StorageClass field to given value.

### HasStorageClass

`func (o *ObjectVersion) HasStorageClass() bool`

HasStorageClass returns a boolean if a field has been set.

### GetKey

`func (o *ObjectVersion) GetKey() string`

GetKey returns the Key field if non-nil, zero value otherwise.

### GetKeyOk

`func (o *ObjectVersion) GetKeyOk() (*string, bool)`

GetKeyOk returns a tuple with the Key field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKey

`func (o *ObjectVersion) SetKey(v string)`

SetKey sets Key field to given value.

### HasKey

`func (o *ObjectVersion) HasKey() bool`

HasKey returns a boolean if a field has been set.

### GetVersionId

`func (o *ObjectVersion) GetVersionId() string`

GetVersionId returns the VersionId field if non-nil, zero value otherwise.

### GetVersionIdOk

`func (o *ObjectVersion) GetVersionIdOk() (*string, bool)`

GetVersionIdOk returns a tuple with the VersionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersionId

`func (o *ObjectVersion) SetVersionId(v string)`

SetVersionId sets VersionId field to given value.

### HasVersionId

`func (o *ObjectVersion) HasVersionId() bool`

HasVersionId returns a boolean if a field has been set.

### GetIsLatest

`func (o *ObjectVersion) GetIsLatest() bool`

GetIsLatest returns the IsLatest field if non-nil, zero value otherwise.

### GetIsLatestOk

`func (o *ObjectVersion) GetIsLatestOk() (*bool, bool)`

GetIsLatestOk returns a tuple with the IsLatest field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsLatest

`func (o *ObjectVersion) SetIsLatest(v bool)`

SetIsLatest sets IsLatest field to given value.

### HasIsLatest

`func (o *ObjectVersion) HasIsLatest() bool`

HasIsLatest returns a boolean if a field has been set.

### GetLastModified

`func (o *ObjectVersion) GetLastModified() time.Time`

GetLastModified returns the LastModified field if non-nil, zero value otherwise.

### GetLastModifiedOk

`func (o *ObjectVersion) GetLastModifiedOk() (*time.Time, bool)`

GetLastModifiedOk returns a tuple with the LastModified field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastModified

`func (o *ObjectVersion) SetLastModified(v time.Time)`

SetLastModified sets LastModified field to given value.

### HasLastModified

`func (o *ObjectVersion) HasLastModified() bool`

HasLastModified returns a boolean if a field has been set.

### GetOwner

`func (o *ObjectVersion) GetOwner() Owner`

GetOwner returns the Owner field if non-nil, zero value otherwise.

### GetOwnerOk

`func (o *ObjectVersion) GetOwnerOk() (*Owner, bool)`

GetOwnerOk returns a tuple with the Owner field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOwner

`func (o *ObjectVersion) SetOwner(v Owner)`

SetOwner sets Owner field to given value.

### HasOwner

`func (o *ObjectVersion) HasOwner() bool`

HasOwner returns a boolean if a field has been set.


