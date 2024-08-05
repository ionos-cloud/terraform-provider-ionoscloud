# RedirectAllRequestsTo

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**HostName** | **string** | Name of the host where requests are redirected. | |
|**Protocol** | Pointer to **string** | Protocol to use when redirecting requests. The default is the protocol that is used in the original request. | [optional] |

## Methods

### NewRedirectAllRequestsTo

`func NewRedirectAllRequestsTo(hostName string, ) *RedirectAllRequestsTo`

NewRedirectAllRequestsTo instantiates a new RedirectAllRequestsTo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRedirectAllRequestsToWithDefaults

`func NewRedirectAllRequestsToWithDefaults() *RedirectAllRequestsTo`

NewRedirectAllRequestsToWithDefaults instantiates a new RedirectAllRequestsTo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetHostName

`func (o *RedirectAllRequestsTo) GetHostName() string`

GetHostName returns the HostName field if non-nil, zero value otherwise.

### GetHostNameOk

`func (o *RedirectAllRequestsTo) GetHostNameOk() (*string, bool)`

GetHostNameOk returns a tuple with the HostName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHostName

`func (o *RedirectAllRequestsTo) SetHostName(v string)`

SetHostName sets HostName field to given value.


### GetProtocol

`func (o *RedirectAllRequestsTo) GetProtocol() string`

GetProtocol returns the Protocol field if non-nil, zero value otherwise.

### GetProtocolOk

`func (o *RedirectAllRequestsTo) GetProtocolOk() (*string, bool)`

GetProtocolOk returns a tuple with the Protocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocol

`func (o *RedirectAllRequestsTo) SetProtocol(v string)`

SetProtocol sets Protocol field to given value.

### HasProtocol

`func (o *RedirectAllRequestsTo) HasProtocol() bool`

HasProtocol returns a boolean if a field has been set.


