# InputSerialization

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**CSV** | Pointer to [**CSVInput**](CSVInput.md) |  | [optional] |
|**CompressionType** | Pointer to **string** | Specifies object&#39;s compression format. Valid values: NONE, GZIP, BZIP2. Default Value: NONE. | [optional] |
|**JSON** | Pointer to [**InputSerializationJSON**](InputSerializationJSON.md) |  | [optional] |
|**Parquet** | Pointer to **map[string]interface{}** | Specifies Parquet as object&#39;s input serialization format. | [optional] |

## Methods

### NewInputSerialization

`func NewInputSerialization() *InputSerialization`

NewInputSerialization instantiates a new InputSerialization object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInputSerializationWithDefaults

`func NewInputSerializationWithDefaults() *InputSerialization`

NewInputSerializationWithDefaults instantiates a new InputSerialization object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCSV

`func (o *InputSerialization) GetCSV() CSVInput`

GetCSV returns the CSV field if non-nil, zero value otherwise.

### GetCSVOk

`func (o *InputSerialization) GetCSVOk() (*CSVInput, bool)`

GetCSVOk returns a tuple with the CSV field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCSV

`func (o *InputSerialization) SetCSV(v CSVInput)`

SetCSV sets CSV field to given value.

### HasCSV

`func (o *InputSerialization) HasCSV() bool`

HasCSV returns a boolean if a field has been set.

### GetCompressionType

`func (o *InputSerialization) GetCompressionType() string`

GetCompressionType returns the CompressionType field if non-nil, zero value otherwise.

### GetCompressionTypeOk

`func (o *InputSerialization) GetCompressionTypeOk() (*string, bool)`

GetCompressionTypeOk returns a tuple with the CompressionType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompressionType

`func (o *InputSerialization) SetCompressionType(v string)`

SetCompressionType sets CompressionType field to given value.

### HasCompressionType

`func (o *InputSerialization) HasCompressionType() bool`

HasCompressionType returns a boolean if a field has been set.

### GetJSON

`func (o *InputSerialization) GetJSON() InputSerializationJSON`

GetJSON returns the JSON field if non-nil, zero value otherwise.

### GetJSONOk

`func (o *InputSerialization) GetJSONOk() (*InputSerializationJSON, bool)`

GetJSONOk returns a tuple with the JSON field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJSON

`func (o *InputSerialization) SetJSON(v InputSerializationJSON)`

SetJSON sets JSON field to given value.

### HasJSON

`func (o *InputSerialization) HasJSON() bool`

HasJSON returns a boolean if a field has been set.

### GetParquet

`func (o *InputSerialization) GetParquet() map[string]interface{}`

GetParquet returns the Parquet field if non-nil, zero value otherwise.

### GetParquetOk

`func (o *InputSerialization) GetParquetOk() (*map[string]interface{}, bool)`

GetParquetOk returns a tuple with the Parquet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetParquet

`func (o *InputSerialization) SetParquet(v map[string]interface{})`

SetParquet sets Parquet field to given value.

### HasParquet

`func (o *InputSerialization) HasParquet() bool`

HasParquet returns a boolean if a field has been set.


