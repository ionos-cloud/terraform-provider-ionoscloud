# ObjectIdentifier

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Key** | **string** | The object key. | |
|**VersionId** | Pointer to **string** | VersionId for the specific version of the object to delete. | [optional] |

## Methods

### NewObjectIdentifier

`func NewObjectIdentifier(key string, ) *ObjectIdentifier`

NewObjectIdentifier instantiates a new ObjectIdentifier object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewObjectIdentifierWithDefaults

`func NewObjectIdentifierWithDefaults() *ObjectIdentifier`

NewObjectIdentifierWithDefaults instantiates a new ObjectIdentifier object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKey

`func (o *ObjectIdentifier) GetKey() string`

GetKey returns the Key field if non-nil, zero value otherwise.

### GetKeyOk

`func (o *ObjectIdentifier) GetKeyOk() (*string, bool)`

GetKeyOk returns a tuple with the Key field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKey

`func (o *ObjectIdentifier) SetKey(v string)`

SetKey sets Key field to given value.


### GetVersionId

`func (o *ObjectIdentifier) GetVersionId() string`

GetVersionId returns the VersionId field if non-nil, zero value otherwise.

### GetVersionIdOk

`func (o *ObjectIdentifier) GetVersionIdOk() (*string, bool)`

GetVersionIdOk returns a tuple with the VersionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersionId

`func (o *ObjectIdentifier) SetVersionId(v string)`

SetVersionId sets VersionId field to given value.

### HasVersionId

`func (o *ObjectIdentifier) HasVersionId() bool`

HasVersionId returns a boolean if a field has been set.


