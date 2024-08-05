# Redirect

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**HostName** | Pointer to **string** | The host name to use in the redirect request. | [optional] |
|**HttpRedirectCode** | Pointer to **string** | The HTTP redirect code to use on the response. Not required if one of the siblings is present. | [optional] |
|**Protocol** | Pointer to **string** | Protocol to use when redirecting requests. The default is the protocol that is used in the original request. | [optional] |
|**ReplaceKeyPrefixWith** | Pointer to **string** | &lt;p&gt;The object key prefix to use in the redirect request. For example, to redirect requests for all pages with prefix &#x60;docs/&#x60; (objects in the &#x60;docs/&#x60; folder) to &#x60;documents/&#x60;, you can set a condition block with &#x60;KeyPrefixEquals&#x60; set to &#x60;docs/&#x60; and in the Redirect set &#x60;ReplaceKeyPrefixWith&#x60; to &#x60;/documents&#x60;. Not required if one of the siblings is present. Can be present only if &#x60;ReplaceKeyWith&#x60; is not provided.&lt;/p&gt; &lt;p&gt;Replacement must be made for object keys containing special characters (such as carriage returns) when using XML requests. &lt;/p&gt; | [optional] |
|**ReplaceKeyWith** | Pointer to **string** | The specific object key to use in the redirect request. For example, redirect request to &#x60;error.html&#x60;. Not required if one of the siblings is present. Can be present only if &#x60;ReplaceKeyPrefixWith&#x60; is not provided. Replacement must be made for object keys containing special characters (such as carriage returns) when using XML requests. | [optional] |

## Methods

### NewRedirect

`func NewRedirect() *Redirect`

NewRedirect instantiates a new Redirect object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRedirectWithDefaults

`func NewRedirectWithDefaults() *Redirect`

NewRedirectWithDefaults instantiates a new Redirect object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetHostName

`func (o *Redirect) GetHostName() string`

GetHostName returns the HostName field if non-nil, zero value otherwise.

### GetHostNameOk

`func (o *Redirect) GetHostNameOk() (*string, bool)`

GetHostNameOk returns a tuple with the HostName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHostName

`func (o *Redirect) SetHostName(v string)`

SetHostName sets HostName field to given value.

### HasHostName

`func (o *Redirect) HasHostName() bool`

HasHostName returns a boolean if a field has been set.

### GetHttpRedirectCode

`func (o *Redirect) GetHttpRedirectCode() string`

GetHttpRedirectCode returns the HttpRedirectCode field if non-nil, zero value otherwise.

### GetHttpRedirectCodeOk

`func (o *Redirect) GetHttpRedirectCodeOk() (*string, bool)`

GetHttpRedirectCodeOk returns a tuple with the HttpRedirectCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHttpRedirectCode

`func (o *Redirect) SetHttpRedirectCode(v string)`

SetHttpRedirectCode sets HttpRedirectCode field to given value.

### HasHttpRedirectCode

`func (o *Redirect) HasHttpRedirectCode() bool`

HasHttpRedirectCode returns a boolean if a field has been set.

### GetProtocol

`func (o *Redirect) GetProtocol() string`

GetProtocol returns the Protocol field if non-nil, zero value otherwise.

### GetProtocolOk

`func (o *Redirect) GetProtocolOk() (*string, bool)`

GetProtocolOk returns a tuple with the Protocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocol

`func (o *Redirect) SetProtocol(v string)`

SetProtocol sets Protocol field to given value.

### HasProtocol

`func (o *Redirect) HasProtocol() bool`

HasProtocol returns a boolean if a field has been set.

### GetReplaceKeyPrefixWith

`func (o *Redirect) GetReplaceKeyPrefixWith() string`

GetReplaceKeyPrefixWith returns the ReplaceKeyPrefixWith field if non-nil, zero value otherwise.

### GetReplaceKeyPrefixWithOk

`func (o *Redirect) GetReplaceKeyPrefixWithOk() (*string, bool)`

GetReplaceKeyPrefixWithOk returns a tuple with the ReplaceKeyPrefixWith field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReplaceKeyPrefixWith

`func (o *Redirect) SetReplaceKeyPrefixWith(v string)`

SetReplaceKeyPrefixWith sets ReplaceKeyPrefixWith field to given value.

### HasReplaceKeyPrefixWith

`func (o *Redirect) HasReplaceKeyPrefixWith() bool`

HasReplaceKeyPrefixWith returns a boolean if a field has been set.

### GetReplaceKeyWith

`func (o *Redirect) GetReplaceKeyWith() string`

GetReplaceKeyWith returns the ReplaceKeyWith field if non-nil, zero value otherwise.

### GetReplaceKeyWithOk

`func (o *Redirect) GetReplaceKeyWithOk() (*string, bool)`

GetReplaceKeyWithOk returns a tuple with the ReplaceKeyWith field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReplaceKeyWith

`func (o *Redirect) SetReplaceKeyWith(v string)`

SetReplaceKeyWith sets ReplaceKeyWith field to given value.

### HasReplaceKeyWith

`func (o *Redirect) HasReplaceKeyWith() bool`

HasReplaceKeyWith returns a boolean if a field has been set.


