# InputSerializationJSON

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Type** | Pointer to **string** | Specifies JSON as object&#39;s input serialization format. | [optional] |

## Methods

### NewInputSerializationJSON

`func NewInputSerializationJSON() *InputSerializationJSON`

NewInputSerializationJSON instantiates a new InputSerializationJSON object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInputSerializationJSONWithDefaults

`func NewInputSerializationJSONWithDefaults() *InputSerializationJSON`

NewInputSerializationJSONWithDefaults instantiates a new InputSerializationJSON object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetType

`func (o *InputSerializationJSON) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *InputSerializationJSON) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *InputSerializationJSON) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *InputSerializationJSON) HasType() bool`

HasType returns a boolean if a field has been set.


