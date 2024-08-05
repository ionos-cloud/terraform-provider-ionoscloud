# DeletionError

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Key** | Pointer to **string** | The object key. | [optional] |
|**VersionId** | Pointer to **string** | The version ID of the object. | [optional] |
|**Code** | Pointer to **string** |  | [optional] |
|**Message** | Pointer to **string** |  | [optional] |

## Methods

### NewDeletionError

`func NewDeletionError() *DeletionError`

NewDeletionError instantiates a new DeletionError object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeletionErrorWithDefaults

`func NewDeletionErrorWithDefaults() *DeletionError`

NewDeletionErrorWithDefaults instantiates a new DeletionError object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKey

`func (o *DeletionError) GetKey() string`

GetKey returns the Key field if non-nil, zero value otherwise.

### GetKeyOk

`func (o *DeletionError) GetKeyOk() (*string, bool)`

GetKeyOk returns a tuple with the Key field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKey

`func (o *DeletionError) SetKey(v string)`

SetKey sets Key field to given value.

### HasKey

`func (o *DeletionError) HasKey() bool`

HasKey returns a boolean if a field has been set.

### GetVersionId

`func (o *DeletionError) GetVersionId() string`

GetVersionId returns the VersionId field if non-nil, zero value otherwise.

### GetVersionIdOk

`func (o *DeletionError) GetVersionIdOk() (*string, bool)`

GetVersionIdOk returns a tuple with the VersionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersionId

`func (o *DeletionError) SetVersionId(v string)`

SetVersionId sets VersionId field to given value.

### HasVersionId

`func (o *DeletionError) HasVersionId() bool`

HasVersionId returns a boolean if a field has been set.

### GetCode

`func (o *DeletionError) GetCode() string`

GetCode returns the Code field if non-nil, zero value otherwise.

### GetCodeOk

`func (o *DeletionError) GetCodeOk() (*string, bool)`

GetCodeOk returns a tuple with the Code field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCode

`func (o *DeletionError) SetCode(v string)`

SetCode sets Code field to given value.

### HasCode

`func (o *DeletionError) HasCode() bool`

HasCode returns a boolean if a field has been set.

### GetMessage

`func (o *DeletionError) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *DeletionError) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *DeletionError) SetMessage(v string)`

SetMessage sets Message field to given value.

### HasMessage

`func (o *DeletionError) HasMessage() bool`

HasMessage returns a boolean if a field has been set.


