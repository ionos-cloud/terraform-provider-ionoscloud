# PolicyStatus

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**IsPublic** | Pointer to **bool** | The policy status for this bucket: - &#x60;true&#x60; indicates that this bucket is public. - &#x60;false&#x60; indicates that this bucket is private.  | [optional] |

## Methods

### NewPolicyStatus

`func NewPolicyStatus() *PolicyStatus`

NewPolicyStatus instantiates a new PolicyStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPolicyStatusWithDefaults

`func NewPolicyStatusWithDefaults() *PolicyStatus`

NewPolicyStatusWithDefaults instantiates a new PolicyStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetIsPublic

`func (o *PolicyStatus) GetIsPublic() bool`

GetIsPublic returns the IsPublic field if non-nil, zero value otherwise.

### GetIsPublicOk

`func (o *PolicyStatus) GetIsPublicOk() (*bool, bool)`

GetIsPublicOk returns a tuple with the IsPublic field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsPublic

`func (o *PolicyStatus) SetIsPublic(v bool)`

SetIsPublic sets IsPublic field to given value.

### HasIsPublic

`func (o *PolicyStatus) HasIsPublic() bool`

HasIsPublic returns a boolean if a field has been set.


