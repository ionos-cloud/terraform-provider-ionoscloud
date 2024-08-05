# DeleteMarkerEntry

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Owner** | Pointer to [**Owner**](Owner.md) |  | [optional] |
|**Key** | Pointer to **string** | The object key. | [optional] |
|**VersionId** | Pointer to **string** | Version ID of the Deletion Marker | [optional] |
|**IsLatest** | Pointer to **bool** | Specifies whether the object is (true) or is not (false) the latest version of an object. | [optional] |
|**LastModified** | Pointer to [**time.Time**](time.Time.md) | Creation date of the object. | [optional] |

## Methods

### NewDeleteMarkerEntry

`func NewDeleteMarkerEntry() *DeleteMarkerEntry`

NewDeleteMarkerEntry instantiates a new DeleteMarkerEntry object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeleteMarkerEntryWithDefaults

`func NewDeleteMarkerEntryWithDefaults() *DeleteMarkerEntry`

NewDeleteMarkerEntryWithDefaults instantiates a new DeleteMarkerEntry object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetOwner

`func (o *DeleteMarkerEntry) GetOwner() Owner`

GetOwner returns the Owner field if non-nil, zero value otherwise.

### GetOwnerOk

`func (o *DeleteMarkerEntry) GetOwnerOk() (*Owner, bool)`

GetOwnerOk returns a tuple with the Owner field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOwner

`func (o *DeleteMarkerEntry) SetOwner(v Owner)`

SetOwner sets Owner field to given value.

### HasOwner

`func (o *DeleteMarkerEntry) HasOwner() bool`

HasOwner returns a boolean if a field has been set.

### GetKey

`func (o *DeleteMarkerEntry) GetKey() string`

GetKey returns the Key field if non-nil, zero value otherwise.

### GetKeyOk

`func (o *DeleteMarkerEntry) GetKeyOk() (*string, bool)`

GetKeyOk returns a tuple with the Key field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKey

`func (o *DeleteMarkerEntry) SetKey(v string)`

SetKey sets Key field to given value.

### HasKey

`func (o *DeleteMarkerEntry) HasKey() bool`

HasKey returns a boolean if a field has been set.

### GetVersionId

`func (o *DeleteMarkerEntry) GetVersionId() string`

GetVersionId returns the VersionId field if non-nil, zero value otherwise.

### GetVersionIdOk

`func (o *DeleteMarkerEntry) GetVersionIdOk() (*string, bool)`

GetVersionIdOk returns a tuple with the VersionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersionId

`func (o *DeleteMarkerEntry) SetVersionId(v string)`

SetVersionId sets VersionId field to given value.

### HasVersionId

`func (o *DeleteMarkerEntry) HasVersionId() bool`

HasVersionId returns a boolean if a field has been set.

### GetIsLatest

`func (o *DeleteMarkerEntry) GetIsLatest() bool`

GetIsLatest returns the IsLatest field if non-nil, zero value otherwise.

### GetIsLatestOk

`func (o *DeleteMarkerEntry) GetIsLatestOk() (*bool, bool)`

GetIsLatestOk returns a tuple with the IsLatest field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsLatest

`func (o *DeleteMarkerEntry) SetIsLatest(v bool)`

SetIsLatest sets IsLatest field to given value.

### HasIsLatest

`func (o *DeleteMarkerEntry) HasIsLatest() bool`

HasIsLatest returns a boolean if a field has been set.

### GetLastModified

`func (o *DeleteMarkerEntry) GetLastModified() time.Time`

GetLastModified returns the LastModified field if non-nil, zero value otherwise.

### GetLastModifiedOk

`func (o *DeleteMarkerEntry) GetLastModifiedOk() (*time.Time, bool)`

GetLastModifiedOk returns a tuple with the LastModified field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastModified

`func (o *DeleteMarkerEntry) SetLastModified(v time.Time)`

SetLastModified sets LastModified field to given value.

### HasLastModified

`func (o *DeleteMarkerEntry) HasLastModified() bool`

HasLastModified returns a boolean if a field has been set.


