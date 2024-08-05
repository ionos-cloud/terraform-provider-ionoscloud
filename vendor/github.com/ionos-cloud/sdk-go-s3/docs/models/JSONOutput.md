# JSONOutput

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**RecordDelimiter** | Pointer to **string** | The value used to separate individual records in the output. If no value is specified, IONOS S3 Object Storage uses a newline character (&#39;\\n&#39;). | [optional] |

## Methods

### NewJSONOutput

`func NewJSONOutput() *JSONOutput`

NewJSONOutput instantiates a new JSONOutput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewJSONOutputWithDefaults

`func NewJSONOutputWithDefaults() *JSONOutput`

NewJSONOutputWithDefaults instantiates a new JSONOutput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRecordDelimiter

`func (o *JSONOutput) GetRecordDelimiter() string`

GetRecordDelimiter returns the RecordDelimiter field if non-nil, zero value otherwise.

### GetRecordDelimiterOk

`func (o *JSONOutput) GetRecordDelimiterOk() (*string, bool)`

GetRecordDelimiterOk returns a tuple with the RecordDelimiter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRecordDelimiter

`func (o *JSONOutput) SetRecordDelimiter(v string)`

SetRecordDelimiter sets RecordDelimiter field to given value.

### HasRecordDelimiter

`func (o *JSONOutput) HasRecordDelimiter() bool`

HasRecordDelimiter returns a boolean if a field has been set.


