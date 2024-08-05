# DeletedObject

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Key** | Pointer to **string** | The object key. | [optional] |
|**VersionId** | Pointer to **string** | Version ID of the deleted object | [optional] |
|**DeleteMarker** | Pointer to **bool** | Specifies whether the versioned object that was permanently deleted was (true) or was not (false) a delete marker. In a simple DELETE, this header indicates whether (true) or not (false) a delete marker was created. | [optional] |
|**DeleteMarkerVersionId** | Pointer to **string** | The version ID of the delete marker created as a result of the DELETE operation. If you delete a specific object version, the value returned by this header is the version ID of the object version deleted. | [optional] |

## Methods

### NewDeletedObject

`func NewDeletedObject() *DeletedObject`

NewDeletedObject instantiates a new DeletedObject object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeletedObjectWithDefaults

`func NewDeletedObjectWithDefaults() *DeletedObject`

NewDeletedObjectWithDefaults instantiates a new DeletedObject object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKey

`func (o *DeletedObject) GetKey() string`

GetKey returns the Key field if non-nil, zero value otherwise.

### GetKeyOk

`func (o *DeletedObject) GetKeyOk() (*string, bool)`

GetKeyOk returns a tuple with the Key field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKey

`func (o *DeletedObject) SetKey(v string)`

SetKey sets Key field to given value.

### HasKey

`func (o *DeletedObject) HasKey() bool`

HasKey returns a boolean if a field has been set.

### GetVersionId

`func (o *DeletedObject) GetVersionId() string`

GetVersionId returns the VersionId field if non-nil, zero value otherwise.

### GetVersionIdOk

`func (o *DeletedObject) GetVersionIdOk() (*string, bool)`

GetVersionIdOk returns a tuple with the VersionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersionId

`func (o *DeletedObject) SetVersionId(v string)`

SetVersionId sets VersionId field to given value.

### HasVersionId

`func (o *DeletedObject) HasVersionId() bool`

HasVersionId returns a boolean if a field has been set.

### GetDeleteMarker

`func (o *DeletedObject) GetDeleteMarker() bool`

GetDeleteMarker returns the DeleteMarker field if non-nil, zero value otherwise.

### GetDeleteMarkerOk

`func (o *DeletedObject) GetDeleteMarkerOk() (*bool, bool)`

GetDeleteMarkerOk returns a tuple with the DeleteMarker field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeleteMarker

`func (o *DeletedObject) SetDeleteMarker(v bool)`

SetDeleteMarker sets DeleteMarker field to given value.

### HasDeleteMarker

`func (o *DeletedObject) HasDeleteMarker() bool`

HasDeleteMarker returns a boolean if a field has been set.

### GetDeleteMarkerVersionId

`func (o *DeletedObject) GetDeleteMarkerVersionId() string`

GetDeleteMarkerVersionId returns the DeleteMarkerVersionId field if non-nil, zero value otherwise.

### GetDeleteMarkerVersionIdOk

`func (o *DeletedObject) GetDeleteMarkerVersionIdOk() (*string, bool)`

GetDeleteMarkerVersionIdOk returns a tuple with the DeleteMarkerVersionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeleteMarkerVersionId

`func (o *DeletedObject) SetDeleteMarkerVersionId(v string)`

SetDeleteMarkerVersionId sets DeleteMarkerVersionId field to given value.

### HasDeleteMarkerVersionId

`func (o *DeletedObject) HasDeleteMarkerVersionId() bool`

HasDeleteMarkerVersionId returns a boolean if a field has been set.


