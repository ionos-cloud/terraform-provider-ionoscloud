# CommonPrefix

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Prefix** | Pointer to **string** | Object key prefix that identifies one or more objects to which this rule applies. Replacement must be made for object keys containing special characters (such as carriage returns) when using XML requests. | [optional] |

## Methods

### NewCommonPrefix

`func NewCommonPrefix() *CommonPrefix`

NewCommonPrefix instantiates a new CommonPrefix object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonPrefixWithDefaults

`func NewCommonPrefixWithDefaults() *CommonPrefix`

NewCommonPrefixWithDefaults instantiates a new CommonPrefix object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPrefix

`func (o *CommonPrefix) GetPrefix() string`

GetPrefix returns the Prefix field if non-nil, zero value otherwise.

### GetPrefixOk

`func (o *CommonPrefix) GetPrefixOk() (*string, bool)`

GetPrefixOk returns a tuple with the Prefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrefix

`func (o *CommonPrefix) SetPrefix(v string)`

SetPrefix sets Prefix field to given value.

### HasPrefix

`func (o *CommonPrefix) HasPrefix() bool`

HasPrefix returns a boolean if a field has been set.


