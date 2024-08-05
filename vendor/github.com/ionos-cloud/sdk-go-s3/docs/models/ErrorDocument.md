# ErrorDocument

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Key** | **string** | The object key. | |

## Methods

### NewErrorDocument

`func NewErrorDocument(key string, ) *ErrorDocument`

NewErrorDocument instantiates a new ErrorDocument object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewErrorDocumentWithDefaults

`func NewErrorDocumentWithDefaults() *ErrorDocument`

NewErrorDocumentWithDefaults instantiates a new ErrorDocument object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKey

`func (o *ErrorDocument) GetKey() string`

GetKey returns the Key field if non-nil, zero value otherwise.

### GetKeyOk

`func (o *ErrorDocument) GetKeyOk() (*string, bool)`

GetKeyOk returns a tuple with the Key field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKey

`func (o *ErrorDocument) SetKey(v string)`

SetKey sets Key field to given value.



