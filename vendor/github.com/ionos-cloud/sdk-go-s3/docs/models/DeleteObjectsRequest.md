# DeleteObjectsRequest

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Objects** | Pointer to [**[]ObjectIdentifier**](ObjectIdentifier.md) | The objects to delete. | [optional] |
|**Quiet** | Pointer to **bool** |  | [optional] |

## Methods

### NewDeleteObjectsRequest

`func NewDeleteObjectsRequest() *DeleteObjectsRequest`

NewDeleteObjectsRequest instantiates a new DeleteObjectsRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeleteObjectsRequestWithDefaults

`func NewDeleteObjectsRequestWithDefaults() *DeleteObjectsRequest`

NewDeleteObjectsRequestWithDefaults instantiates a new DeleteObjectsRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetObjects

`func (o *DeleteObjectsRequest) GetObjects() []ObjectIdentifier`

GetObjects returns the Objects field if non-nil, zero value otherwise.

### GetObjectsOk

`func (o *DeleteObjectsRequest) GetObjectsOk() (*[]ObjectIdentifier, bool)`

GetObjectsOk returns a tuple with the Objects field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjects

`func (o *DeleteObjectsRequest) SetObjects(v []ObjectIdentifier)`

SetObjects sets Objects field to given value.

### HasObjects

`func (o *DeleteObjectsRequest) HasObjects() bool`

HasObjects returns a boolean if a field has been set.

### GetQuiet

`func (o *DeleteObjectsRequest) GetQuiet() bool`

GetQuiet returns the Quiet field if non-nil, zero value otherwise.

### GetQuietOk

`func (o *DeleteObjectsRequest) GetQuietOk() (*bool, bool)`

GetQuietOk returns a tuple with the Quiet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQuiet

`func (o *DeleteObjectsRequest) SetQuiet(v bool)`

SetQuiet sets Quiet field to given value.

### HasQuiet

`func (o *DeleteObjectsRequest) HasQuiet() bool`

HasQuiet returns a boolean if a field has been set.


