# Principal

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**AWS** | **[]string** |  | |

## Methods

### NewPrincipal

`func NewPrincipal(aWS []string, ) *Principal`

NewPrincipal instantiates a new Principal object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPrincipalWithDefaults

`func NewPrincipalWithDefaults() *Principal`

NewPrincipalWithDefaults instantiates a new Principal object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAWS

`func (o *Principal) GetAWS() []string`

GetAWS returns the AWS field if non-nil, zero value otherwise.

### GetAWSOk

`func (o *Principal) GetAWSOk() (*[]string, bool)`

GetAWSOk returns a tuple with the AWS field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAWS

`func (o *Principal) SetAWS(v []string)`

SetAWS sets AWS field to given value.



