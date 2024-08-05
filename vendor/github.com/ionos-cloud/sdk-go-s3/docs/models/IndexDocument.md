# IndexDocument

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Suffix** | **string** | A suffix that is appended to a request that is for a directory on the website endpoint (for example, if the suffix is index.html and you make a request to &#x60;samplebucket/images/&#x60; the data that is returned will be for the object with the key name &#x60;images/index.html&#x60;) The suffix must not be empty and must not include a slash character. Replacement must be made for object keys containing special characters (such as carriage returns) when using XML requests.  | |

## Methods

### NewIndexDocument

`func NewIndexDocument(suffix string, ) *IndexDocument`

NewIndexDocument instantiates a new IndexDocument object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewIndexDocumentWithDefaults

`func NewIndexDocumentWithDefaults() *IndexDocument`

NewIndexDocumentWithDefaults instantiates a new IndexDocument object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSuffix

`func (o *IndexDocument) GetSuffix() string`

GetSuffix returns the Suffix field if non-nil, zero value otherwise.

### GetSuffixOk

`func (o *IndexDocument) GetSuffixOk() (*string, bool)`

GetSuffixOk returns a tuple with the Suffix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuffix

`func (o *IndexDocument) SetSuffix(v string)`

SetSuffix sets Suffix field to given value.



