/*
 * IONOS DBaaS MariaDB REST API
 *
 * An enterprise-grade Database is provided as a Service (DBaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.  The API allows you to create additional MariaDB database clusters or modify existing ones. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 0.1.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package mariadb

import (
	"encoding/json"
)

// checks if the ClustersGet403Response type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ClustersGet403Response{}

// ClustersGet403Response struct for ClustersGet403Response
type ClustersGet403Response struct {
	// The HTTP status code of the operation.
	HttpStatus int32          `json:"httpStatus"`
	Messages   []ErrorMessage `json:"messages"`
}

// NewClustersGet403Response instantiates a new ClustersGet403Response object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClustersGet403Response(httpStatus int32, messages []ErrorMessage) *ClustersGet403Response {
	this := ClustersGet403Response{}

	this.HttpStatus = httpStatus
	this.Messages = messages

	return &this
}

// NewClustersGet403ResponseWithDefaults instantiates a new ClustersGet403Response object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClustersGet403ResponseWithDefaults() *ClustersGet403Response {
	this := ClustersGet403Response{}
	return &this
}

// GetHttpStatus returns the HttpStatus field value
func (o *ClustersGet403Response) GetHttpStatus() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.HttpStatus
}

// GetHttpStatusOk returns a tuple with the HttpStatus field value
// and a boolean to check if the value has been set.
func (o *ClustersGet403Response) GetHttpStatusOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.HttpStatus, true
}

// SetHttpStatus sets field value
func (o *ClustersGet403Response) SetHttpStatus(v int32) {
	o.HttpStatus = v
}

// GetMessages returns the Messages field value
func (o *ClustersGet403Response) GetMessages() []ErrorMessage {
	if o == nil {
		var ret []ErrorMessage
		return ret
	}

	return o.Messages
}

// GetMessagesOk returns a tuple with the Messages field value
// and a boolean to check if the value has been set.
func (o *ClustersGet403Response) GetMessagesOk() ([]ErrorMessage, bool) {
	if o == nil {
		return nil, false
	}
	return o.Messages, true
}

// SetMessages sets field value
func (o *ClustersGet403Response) SetMessages(v []ErrorMessage) {
	o.Messages = v
}

func (o ClustersGet403Response) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ClustersGet403Response) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["httpStatus"] = o.HttpStatus
	toSerialize["messages"] = o.Messages
	return toSerialize, nil
}

type NullableClustersGet403Response struct {
	value *ClustersGet403Response
	isSet bool
}

func (v NullableClustersGet403Response) Get() *ClustersGet403Response {
	return v.value
}

func (v *NullableClustersGet403Response) Set(val *ClustersGet403Response) {
	v.value = val
	v.isSet = true
}

func (v NullableClustersGet403Response) IsSet() bool {
	return v.isSet
}

func (v *NullableClustersGet403Response) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableClustersGet403Response(val *ClustersGet403Response) *NullableClustersGet403Response {
	return &NullableClustersGet403Response{value: val, isSet: true}
}

func (v NullableClustersGet403Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableClustersGet403Response) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
