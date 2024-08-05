# CORSRule

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**ID** | Pointer to **int32** | Container for the Contract Number of the owner. | [optional] |
|**AllowedHeaders** | Pointer to **[]string** | Headers that are specified in the &#x60;Access-Control-Request-Headers&#x60; header. These headers are allowed in a preflight OPTIONS request. In response to any preflight OPTIONS request, IONOS S3 Object Storage returns any requested headers that are allowed. | [optional] |
|**AllowedMethods** | **[]string** | An HTTP method that you allow the origin to execute. Valid values are &#x60;GET&#x60;, &#x60;PUT&#x60;, &#x60;HEAD&#x60;, &#x60;POST&#x60;, and &#x60;DELETE&#x60;. | |
|**AllowedOrigins** | **[]string** | One or more origins you want customers to be able to access the bucket from. | |
|**ExposeHeaders** | Pointer to **[]string** | One or more headers in the response that you want customers to be able to access from their applications (for example, from a JavaScript &#x60;XMLHttpRequest&#x60; object). | [optional] |
|**MaxAgeSeconds** | Pointer to **int32** | The time in seconds that your browser is to cache the preflight response for the specified resource. | [optional] |

## Methods

### NewCORSRule

`func NewCORSRule(allowedMethods []string, allowedOrigins []string, ) *CORSRule`

NewCORSRule instantiates a new CORSRule object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCORSRuleWithDefaults

`func NewCORSRuleWithDefaults() *CORSRule`

NewCORSRuleWithDefaults instantiates a new CORSRule object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetID

`func (o *CORSRule) GetID() int32`

GetID returns the ID field if non-nil, zero value otherwise.

### GetIDOk

`func (o *CORSRule) GetIDOk() (*int32, bool)`

GetIDOk returns a tuple with the ID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetID

`func (o *CORSRule) SetID(v int32)`

SetID sets ID field to given value.

### HasID

`func (o *CORSRule) HasID() bool`

HasID returns a boolean if a field has been set.

### GetAllowedHeaders

`func (o *CORSRule) GetAllowedHeaders() []string`

GetAllowedHeaders returns the AllowedHeaders field if non-nil, zero value otherwise.

### GetAllowedHeadersOk

`func (o *CORSRule) GetAllowedHeadersOk() (*[]string, bool)`

GetAllowedHeadersOk returns a tuple with the AllowedHeaders field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowedHeaders

`func (o *CORSRule) SetAllowedHeaders(v []string)`

SetAllowedHeaders sets AllowedHeaders field to given value.

### HasAllowedHeaders

`func (o *CORSRule) HasAllowedHeaders() bool`

HasAllowedHeaders returns a boolean if a field has been set.

### GetAllowedMethods

`func (o *CORSRule) GetAllowedMethods() []string`

GetAllowedMethods returns the AllowedMethods field if non-nil, zero value otherwise.

### GetAllowedMethodsOk

`func (o *CORSRule) GetAllowedMethodsOk() (*[]string, bool)`

GetAllowedMethodsOk returns a tuple with the AllowedMethods field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowedMethods

`func (o *CORSRule) SetAllowedMethods(v []string)`

SetAllowedMethods sets AllowedMethods field to given value.


### GetAllowedOrigins

`func (o *CORSRule) GetAllowedOrigins() []string`

GetAllowedOrigins returns the AllowedOrigins field if non-nil, zero value otherwise.

### GetAllowedOriginsOk

`func (o *CORSRule) GetAllowedOriginsOk() (*[]string, bool)`

GetAllowedOriginsOk returns a tuple with the AllowedOrigins field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowedOrigins

`func (o *CORSRule) SetAllowedOrigins(v []string)`

SetAllowedOrigins sets AllowedOrigins field to given value.


### GetExposeHeaders

`func (o *CORSRule) GetExposeHeaders() []string`

GetExposeHeaders returns the ExposeHeaders field if non-nil, zero value otherwise.

### GetExposeHeadersOk

`func (o *CORSRule) GetExposeHeadersOk() (*[]string, bool)`

GetExposeHeadersOk returns a tuple with the ExposeHeaders field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExposeHeaders

`func (o *CORSRule) SetExposeHeaders(v []string)`

SetExposeHeaders sets ExposeHeaders field to given value.

### HasExposeHeaders

`func (o *CORSRule) HasExposeHeaders() bool`

HasExposeHeaders returns a boolean if a field has been set.

### GetMaxAgeSeconds

`func (o *CORSRule) GetMaxAgeSeconds() int32`

GetMaxAgeSeconds returns the MaxAgeSeconds field if non-nil, zero value otherwise.

### GetMaxAgeSecondsOk

`func (o *CORSRule) GetMaxAgeSecondsOk() (*int32, bool)`

GetMaxAgeSecondsOk returns a tuple with the MaxAgeSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxAgeSeconds

`func (o *CORSRule) SetMaxAgeSeconds(v int32)`

SetMaxAgeSeconds sets MaxAgeSeconds field to given value.

### HasMaxAgeSeconds

`func (o *CORSRule) HasMaxAgeSeconds() bool`

HasMaxAgeSeconds returns a boolean if a field has been set.


