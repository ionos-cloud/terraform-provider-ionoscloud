# OutputSerialization

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**CSV** | Pointer to [**CSVOutput**](CSVOutput.md) |  | [optional] |
|**JSON** | Pointer to [**JSONOutput**](JSONOutput.md) |  | [optional] |

## Methods

### NewOutputSerialization

`func NewOutputSerialization() *OutputSerialization`

NewOutputSerialization instantiates a new OutputSerialization object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOutputSerializationWithDefaults

`func NewOutputSerializationWithDefaults() *OutputSerialization`

NewOutputSerializationWithDefaults instantiates a new OutputSerialization object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCSV

`func (o *OutputSerialization) GetCSV() CSVOutput`

GetCSV returns the CSV field if non-nil, zero value otherwise.

### GetCSVOk

`func (o *OutputSerialization) GetCSVOk() (*CSVOutput, bool)`

GetCSVOk returns a tuple with the CSV field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCSV

`func (o *OutputSerialization) SetCSV(v CSVOutput)`

SetCSV sets CSV field to given value.

### HasCSV

`func (o *OutputSerialization) HasCSV() bool`

HasCSV returns a boolean if a field has been set.

### GetJSON

`func (o *OutputSerialization) GetJSON() JSONOutput`

GetJSON returns the JSON field if non-nil, zero value otherwise.

### GetJSONOk

`func (o *OutputSerialization) GetJSONOk() (*JSONOutput, bool)`

GetJSONOk returns a tuple with the JSON field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJSON

`func (o *OutputSerialization) SetJSON(v JSONOutput)`

SetJSON sets JSON field to given value.

### HasJSON

`func (o *OutputSerialization) HasJSON() bool`

HasJSON returns a boolean if a field has been set.


