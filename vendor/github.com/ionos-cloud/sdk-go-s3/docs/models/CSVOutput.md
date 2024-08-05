# CSVOutput

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**QuoteFields** | Pointer to **string** | &lt;p&gt;Indicates whether to use quotation marks around output fields. &lt;/p&gt; &lt;ul&gt; &lt;li&gt; &lt;p&gt; &#x60;ALWAYS&#x60;: Always use quotation marks for output fields.&lt;/p&gt; &lt;/li&gt; &lt;li&gt; &lt;p&gt; &#x60;ASNEEDED&#x60;: Use quotation marks for output fields when needed.&lt;/p&gt; &lt;/li&gt; &lt;/ul&gt; | [optional] |
|**QuoteEscapeCharacter** | Pointer to **string** | The single character used for escaping the quote character inside an already escaped value. | [optional] |
|**RecordDelimiter** | Pointer to **string** | A single character used to separate individual records in the output. Instead of the default value, you can specify an arbitrary delimiter. | [optional] |
|**FieldDelimiter** | Pointer to **interface{}** | The value used to separate individual fields in a record. You can specify an arbitrary delimiter. | [optional] |
|**QuoteCharacter** | Pointer to **string** | A single character used for escaping when the field delimiter is part of the value. For example, if the value is &#x60;a, b&#x60;, IONOS S3 Object Storage wraps this field value in quotation marks, as follows: &#x60;\&quot; a , b \&quot;&#x60;. | [optional] |

## Methods

### NewCSVOutput

`func NewCSVOutput() *CSVOutput`

NewCSVOutput instantiates a new CSVOutput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCSVOutputWithDefaults

`func NewCSVOutputWithDefaults() *CSVOutput`

NewCSVOutputWithDefaults instantiates a new CSVOutput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetQuoteFields

`func (o *CSVOutput) GetQuoteFields() string`

GetQuoteFields returns the QuoteFields field if non-nil, zero value otherwise.

### GetQuoteFieldsOk

`func (o *CSVOutput) GetQuoteFieldsOk() (*string, bool)`

GetQuoteFieldsOk returns a tuple with the QuoteFields field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQuoteFields

`func (o *CSVOutput) SetQuoteFields(v string)`

SetQuoteFields sets QuoteFields field to given value.

### HasQuoteFields

`func (o *CSVOutput) HasQuoteFields() bool`

HasQuoteFields returns a boolean if a field has been set.

### GetQuoteEscapeCharacter

`func (o *CSVOutput) GetQuoteEscapeCharacter() string`

GetQuoteEscapeCharacter returns the QuoteEscapeCharacter field if non-nil, zero value otherwise.

### GetQuoteEscapeCharacterOk

`func (o *CSVOutput) GetQuoteEscapeCharacterOk() (*string, bool)`

GetQuoteEscapeCharacterOk returns a tuple with the QuoteEscapeCharacter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQuoteEscapeCharacter

`func (o *CSVOutput) SetQuoteEscapeCharacter(v string)`

SetQuoteEscapeCharacter sets QuoteEscapeCharacter field to given value.

### HasQuoteEscapeCharacter

`func (o *CSVOutput) HasQuoteEscapeCharacter() bool`

HasQuoteEscapeCharacter returns a boolean if a field has been set.

### GetRecordDelimiter

`func (o *CSVOutput) GetRecordDelimiter() string`

GetRecordDelimiter returns the RecordDelimiter field if non-nil, zero value otherwise.

### GetRecordDelimiterOk

`func (o *CSVOutput) GetRecordDelimiterOk() (*string, bool)`

GetRecordDelimiterOk returns a tuple with the RecordDelimiter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRecordDelimiter

`func (o *CSVOutput) SetRecordDelimiter(v string)`

SetRecordDelimiter sets RecordDelimiter field to given value.

### HasRecordDelimiter

`func (o *CSVOutput) HasRecordDelimiter() bool`

HasRecordDelimiter returns a boolean if a field has been set.

### GetFieldDelimiter

`func (o *CSVOutput) GetFieldDelimiter() interface{}`

GetFieldDelimiter returns the FieldDelimiter field if non-nil, zero value otherwise.

### GetFieldDelimiterOk

`func (o *CSVOutput) GetFieldDelimiterOk() (*interface{}, bool)`

GetFieldDelimiterOk returns a tuple with the FieldDelimiter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFieldDelimiter

`func (o *CSVOutput) SetFieldDelimiter(v interface{})`

SetFieldDelimiter sets FieldDelimiter field to given value.

### HasFieldDelimiter

`func (o *CSVOutput) HasFieldDelimiter() bool`

HasFieldDelimiter returns a boolean if a field has been set.

### SetFieldDelimiterNil

`func (o *CSVOutput) SetFieldDelimiterNil(b bool)`

 SetFieldDelimiterNil sets the value for FieldDelimiter to be an explicit nil

### UnsetFieldDelimiter
`func (o *CSVOutput) UnsetFieldDelimiter()`

UnsetFieldDelimiter ensures that no value is present for FieldDelimiter, not even an explicit nil
### GetQuoteCharacter

`func (o *CSVOutput) GetQuoteCharacter() string`

GetQuoteCharacter returns the QuoteCharacter field if non-nil, zero value otherwise.

### GetQuoteCharacterOk

`func (o *CSVOutput) GetQuoteCharacterOk() (*string, bool)`

GetQuoteCharacterOk returns a tuple with the QuoteCharacter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQuoteCharacter

`func (o *CSVOutput) SetQuoteCharacter(v string)`

SetQuoteCharacter sets QuoteCharacter field to given value.

### HasQuoteCharacter

`func (o *CSVOutput) HasQuoteCharacter() bool`

HasQuoteCharacter returns a boolean if a field has been set.


