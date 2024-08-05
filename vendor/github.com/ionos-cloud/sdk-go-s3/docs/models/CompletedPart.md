# CompletedPart

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**ETag** | Pointer to **string** | Entity tag that identifies the object&#39;s data. Objects with different object data will have different entity tags. The entity tag is an opaque string. The entity tag may or may not be an MD5 digest of the object data. If the entity tag is not an MD5 digest of the object data, it will contain one or more nonhexadecimal characters and/or will consist of less than 32 or more than 32 hexadecimal digits.  | [optional] |
|**PartNumber** | Pointer to **int32** | Part number that identifies the part. | [optional] |

## Methods

### NewCompletedPart

`func NewCompletedPart() *CompletedPart`

NewCompletedPart instantiates a new CompletedPart object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCompletedPartWithDefaults

`func NewCompletedPartWithDefaults() *CompletedPart`

NewCompletedPartWithDefaults instantiates a new CompletedPart object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetETag

`func (o *CompletedPart) GetETag() string`

GetETag returns the ETag field if non-nil, zero value otherwise.

### GetETagOk

`func (o *CompletedPart) GetETagOk() (*string, bool)`

GetETagOk returns a tuple with the ETag field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetETag

`func (o *CompletedPart) SetETag(v string)`

SetETag sets ETag field to given value.

### HasETag

`func (o *CompletedPart) HasETag() bool`

HasETag returns a boolean if a field has been set.

### GetPartNumber

`func (o *CompletedPart) GetPartNumber() int32`

GetPartNumber returns the PartNumber field if non-nil, zero value otherwise.

### GetPartNumberOk

`func (o *CompletedPart) GetPartNumberOk() (*int32, bool)`

GetPartNumberOk returns a tuple with the PartNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPartNumber

`func (o *CompletedPart) SetPartNumber(v int32)`

SetPartNumber sets PartNumber field to given value.

### HasPartNumber

`func (o *CompletedPart) HasPartNumber() bool`

HasPartNumber returns a boolean if a field has been set.


