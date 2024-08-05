# RoutingRuleCondition

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**HttpErrorCodeReturnedEquals** | Pointer to **string** | The HTTP error code when the redirect is applied. In the event of an error, if the error code equals this value, then the specified redirect is applied. Required when parent element &#x60;Condition&#x60; is specified and sibling &#x60;KeyPrefixEquals&#x60; is not specified. If both are specified, then both must be true for the redirect to be applied. | [optional] |
|**KeyPrefixEquals** | Pointer to **string** | &lt;p&gt;The object key name prefix when the redirect is applied. For example, to redirect requests for &#x60;ExamplePage.html&#x60;, the key prefix will be &#x60;ExamplePage.html&#x60;. To redirect request for all pages with the prefix &#x60;docs/&#x60;, the key prefix will be &#x60;/docs&#x60;, which identifies all objects in the &#x60;docs/&#x60; folder. Required when the parent element &#x60;Condition&#x60; is specified and sibling &#x60;HttpErrorCodeReturnedEquals&#x60; is not specified. If both conditions are specified, both must be true for the redirect to be applied.&lt;/p&gt; &lt;important&gt; &lt;p&gt;Replacement must be made for object keys containing special characters (such as carriage returns) when using XML requests.&lt;/p&gt; &lt;/important&gt; | [optional] |

## Methods

### NewRoutingRuleCondition

`func NewRoutingRuleCondition() *RoutingRuleCondition`

NewRoutingRuleCondition instantiates a new RoutingRuleCondition object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRoutingRuleConditionWithDefaults

`func NewRoutingRuleConditionWithDefaults() *RoutingRuleCondition`

NewRoutingRuleConditionWithDefaults instantiates a new RoutingRuleCondition object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetHttpErrorCodeReturnedEquals

`func (o *RoutingRuleCondition) GetHttpErrorCodeReturnedEquals() string`

GetHttpErrorCodeReturnedEquals returns the HttpErrorCodeReturnedEquals field if non-nil, zero value otherwise.

### GetHttpErrorCodeReturnedEqualsOk

`func (o *RoutingRuleCondition) GetHttpErrorCodeReturnedEqualsOk() (*string, bool)`

GetHttpErrorCodeReturnedEqualsOk returns a tuple with the HttpErrorCodeReturnedEquals field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHttpErrorCodeReturnedEquals

`func (o *RoutingRuleCondition) SetHttpErrorCodeReturnedEquals(v string)`

SetHttpErrorCodeReturnedEquals sets HttpErrorCodeReturnedEquals field to given value.

### HasHttpErrorCodeReturnedEquals

`func (o *RoutingRuleCondition) HasHttpErrorCodeReturnedEquals() bool`

HasHttpErrorCodeReturnedEquals returns a boolean if a field has been set.

### GetKeyPrefixEquals

`func (o *RoutingRuleCondition) GetKeyPrefixEquals() string`

GetKeyPrefixEquals returns the KeyPrefixEquals field if non-nil, zero value otherwise.

### GetKeyPrefixEqualsOk

`func (o *RoutingRuleCondition) GetKeyPrefixEqualsOk() (*string, bool)`

GetKeyPrefixEqualsOk returns a tuple with the KeyPrefixEquals field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeyPrefixEquals

`func (o *RoutingRuleCondition) SetKeyPrefixEquals(v string)`

SetKeyPrefixEquals sets KeyPrefixEquals field to given value.

### HasKeyPrefixEquals

`func (o *RoutingRuleCondition) HasKeyPrefixEquals() bool`

HasKeyPrefixEquals returns a boolean if a field has been set.


