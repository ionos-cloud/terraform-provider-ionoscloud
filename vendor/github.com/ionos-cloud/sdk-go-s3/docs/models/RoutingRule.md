# RoutingRule

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Condition** | Pointer to [**RoutingRuleCondition**](RoutingRuleCondition.md) |  | [optional] |
|**Redirect** | [**Redirect**](Redirect.md) |  | |

## Methods

### NewRoutingRule

`func NewRoutingRule(redirect Redirect, ) *RoutingRule`

NewRoutingRule instantiates a new RoutingRule object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRoutingRuleWithDefaults

`func NewRoutingRuleWithDefaults() *RoutingRule`

NewRoutingRuleWithDefaults instantiates a new RoutingRule object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCondition

`func (o *RoutingRule) GetCondition() RoutingRuleCondition`

GetCondition returns the Condition field if non-nil, zero value otherwise.

### GetConditionOk

`func (o *RoutingRule) GetConditionOk() (*RoutingRuleCondition, bool)`

GetConditionOk returns a tuple with the Condition field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCondition

`func (o *RoutingRule) SetCondition(v RoutingRuleCondition)`

SetCondition sets Condition field to given value.

### HasCondition

`func (o *RoutingRule) HasCondition() bool`

HasCondition returns a boolean if a field has been set.

### GetRedirect

`func (o *RoutingRule) GetRedirect() Redirect`

GetRedirect returns the Redirect field if non-nil, zero value otherwise.

### GetRedirectOk

`func (o *RoutingRule) GetRedirectOk() (*Redirect, bool)`

GetRedirectOk returns a tuple with the Redirect field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirect

`func (o *RoutingRule) SetRedirect(v Redirect)`

SetRedirect sets Redirect field to given value.



