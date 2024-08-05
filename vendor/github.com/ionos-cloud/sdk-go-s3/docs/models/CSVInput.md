# CSVInput

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**FileHeaderInfo** | Pointer to **string** | &lt;p&gt;Describes the first line of input. Valid values are:&lt;/p&gt; &lt;ul&gt; &lt;li&gt; &lt;p&gt; &#x60;NONE&#x60;: First line is not a header.&lt;/p&gt; &lt;/li&gt; &lt;li&gt; &lt;p&gt; &#x60;IGNORE&#x60;: First line is a header, but you can&#39;t use the header values to indicate the column in an expression. You can use column position (such as _1, _2, …) to indicate the column (&#x60;SELECT s._1 FROM OBJECT s&#x60;).&lt;/p&gt; &lt;/li&gt; &lt;li&gt; &lt;p&gt; &#x60;Use&#x60;: First line is a header, and you can use the header value to identify a column in an expression (&#x60;SELECT \&quot;name\&quot; FROM OBJECT&#x60;). &lt;/p&gt; &lt;/li&gt; &lt;/ul&gt; | [optional] |
|**Comments** | Pointer to **string** | A single character used to indicate that a row should be ignored when the character is present at the start of that row. You can specify any character to indicate a comment line. | [optional] |
|**QuoteEscapeCharacter** | Pointer to **string** | A single character used for escaping the quotation mark character inside an already escaped value. For example, the value \&quot;\&quot;\&quot; a , b \&quot;\&quot;\&quot; is parsed as \&quot; a , b \&quot;. | [optional] |
|**RecordDelimiter** | Pointer to **string** | A single character used to separate individual records in the input. Instead of the default value, you can specify an arbitrary delimiter. | [optional] |
|**FieldDelimiter** | Pointer to **string** | A single character used to separate individual fields in a record. You can specify an arbitrary delimiter. | [optional] |
|**QuoteCharacter** | Pointer to **string** | &lt;p&gt;A single character used for escaping when the field delimiter is part of the value. For example, if the value is &#x60;a, b&#x60;, IONOS S3 Object Storage wraps this field value in quotation marks, as follows: &#x60;\&quot; a , b \&quot;&#x60;.&lt;/p&gt; &lt;p&gt;Type: String&lt;/p&gt; &lt;p&gt;Default: &#x60;\&quot;&#x60; &lt;/p&gt; &lt;p&gt;Ancestors: &#x60;CSV&#x60; &lt;/p&gt; | [optional] |
|**AllowQuotedRecordDelimiter** | Pointer to **bool** | Specifies that CSV field values may contain quoted record delimiters and such records should be allowed. Default value is FALSE. Setting this value to TRUE may lower performance. | [optional] |

## Methods

### NewCSVInput

`func NewCSVInput() *CSVInput`

NewCSVInput instantiates a new CSVInput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCSVInputWithDefaults

`func NewCSVInputWithDefaults() *CSVInput`

NewCSVInputWithDefaults instantiates a new CSVInput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFileHeaderInfo

`func (o *CSVInput) GetFileHeaderInfo() string`

GetFileHeaderInfo returns the FileHeaderInfo field if non-nil, zero value otherwise.

### GetFileHeaderInfoOk

`func (o *CSVInput) GetFileHeaderInfoOk() (*string, bool)`

GetFileHeaderInfoOk returns a tuple with the FileHeaderInfo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFileHeaderInfo

`func (o *CSVInput) SetFileHeaderInfo(v string)`

SetFileHeaderInfo sets FileHeaderInfo field to given value.

### HasFileHeaderInfo

`func (o *CSVInput) HasFileHeaderInfo() bool`

HasFileHeaderInfo returns a boolean if a field has been set.

### GetComments

`func (o *CSVInput) GetComments() string`

GetComments returns the Comments field if non-nil, zero value otherwise.

### GetCommentsOk

`func (o *CSVInput) GetCommentsOk() (*string, bool)`

GetCommentsOk returns a tuple with the Comments field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetComments

`func (o *CSVInput) SetComments(v string)`

SetComments sets Comments field to given value.

### HasComments

`func (o *CSVInput) HasComments() bool`

HasComments returns a boolean if a field has been set.

### GetQuoteEscapeCharacter

`func (o *CSVInput) GetQuoteEscapeCharacter() string`

GetQuoteEscapeCharacter returns the QuoteEscapeCharacter field if non-nil, zero value otherwise.

### GetQuoteEscapeCharacterOk

`func (o *CSVInput) GetQuoteEscapeCharacterOk() (*string, bool)`

GetQuoteEscapeCharacterOk returns a tuple with the QuoteEscapeCharacter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQuoteEscapeCharacter

`func (o *CSVInput) SetQuoteEscapeCharacter(v string)`

SetQuoteEscapeCharacter sets QuoteEscapeCharacter field to given value.

### HasQuoteEscapeCharacter

`func (o *CSVInput) HasQuoteEscapeCharacter() bool`

HasQuoteEscapeCharacter returns a boolean if a field has been set.

### GetRecordDelimiter

`func (o *CSVInput) GetRecordDelimiter() string`

GetRecordDelimiter returns the RecordDelimiter field if non-nil, zero value otherwise.

### GetRecordDelimiterOk

`func (o *CSVInput) GetRecordDelimiterOk() (*string, bool)`

GetRecordDelimiterOk returns a tuple with the RecordDelimiter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRecordDelimiter

`func (o *CSVInput) SetRecordDelimiter(v string)`

SetRecordDelimiter sets RecordDelimiter field to given value.

### HasRecordDelimiter

`func (o *CSVInput) HasRecordDelimiter() bool`

HasRecordDelimiter returns a boolean if a field has been set.

### GetFieldDelimiter

`func (o *CSVInput) GetFieldDelimiter() string`

GetFieldDelimiter returns the FieldDelimiter field if non-nil, zero value otherwise.

### GetFieldDelimiterOk

`func (o *CSVInput) GetFieldDelimiterOk() (*string, bool)`

GetFieldDelimiterOk returns a tuple with the FieldDelimiter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFieldDelimiter

`func (o *CSVInput) SetFieldDelimiter(v string)`

SetFieldDelimiter sets FieldDelimiter field to given value.

### HasFieldDelimiter

`func (o *CSVInput) HasFieldDelimiter() bool`

HasFieldDelimiter returns a boolean if a field has been set.

### GetQuoteCharacter

`func (o *CSVInput) GetQuoteCharacter() string`

GetQuoteCharacter returns the QuoteCharacter field if non-nil, zero value otherwise.

### GetQuoteCharacterOk

`func (o *CSVInput) GetQuoteCharacterOk() (*string, bool)`

GetQuoteCharacterOk returns a tuple with the QuoteCharacter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQuoteCharacter

`func (o *CSVInput) SetQuoteCharacter(v string)`

SetQuoteCharacter sets QuoteCharacter field to given value.

### HasQuoteCharacter

`func (o *CSVInput) HasQuoteCharacter() bool`

HasQuoteCharacter returns a boolean if a field has been set.

### GetAllowQuotedRecordDelimiter

`func (o *CSVInput) GetAllowQuotedRecordDelimiter() bool`

GetAllowQuotedRecordDelimiter returns the AllowQuotedRecordDelimiter field if non-nil, zero value otherwise.

### GetAllowQuotedRecordDelimiterOk

`func (o *CSVInput) GetAllowQuotedRecordDelimiterOk() (*bool, bool)`

GetAllowQuotedRecordDelimiterOk returns a tuple with the AllowQuotedRecordDelimiter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowQuotedRecordDelimiter

`func (o *CSVInput) SetAllowQuotedRecordDelimiter(v bool)`

SetAllowQuotedRecordDelimiter sets AllowQuotedRecordDelimiter field to given value.

### HasAllowQuotedRecordDelimiter

`func (o *CSVInput) HasAllowQuotedRecordDelimiter() bool`

HasAllowQuotedRecordDelimiter returns a boolean if a field has been set.


