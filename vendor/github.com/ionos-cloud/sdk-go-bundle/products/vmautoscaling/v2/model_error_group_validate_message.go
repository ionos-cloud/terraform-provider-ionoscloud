/*
 * VM Auto Scaling API
 *
 * The VM Auto Scaling Service enables IONOS clients to horizontally scale the number of VM replicas based on configured rules. You can use VM Auto Scaling to ensure that you have a sufficient number of replicas to handle your application loads at all times.  For this purpose, create a VM Auto Scaling Group that contains the server replicas. The VM Auto Scaling Service ensures that the number of replicas in the group is always within the defined limits.   When scaling policies are set, VM Auto Scaling creates or deletes replicas according to the requirements of your applications. For each policy, specified 'scale-in' and 'scale-out' actions are performed when the corresponding thresholds are reached.
 *
 * API version: 1.0.0
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package vmautoscaling

import (
	"encoding/json"
)

// checks if the ErrorGroupValidateMessage type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ErrorGroupValidateMessage{}

// ErrorGroupValidateMessage struct for ErrorGroupValidateMessage
type ErrorGroupValidateMessage struct {
	ErrorCode *string `json:"errorCode,omitempty"`
	Message   *string `json:"message,omitempty"`
}

// NewErrorGroupValidateMessage instantiates a new ErrorGroupValidateMessage object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewErrorGroupValidateMessage() *ErrorGroupValidateMessage {
	this := ErrorGroupValidateMessage{}

	return &this
}

// NewErrorGroupValidateMessageWithDefaults instantiates a new ErrorGroupValidateMessage object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewErrorGroupValidateMessageWithDefaults() *ErrorGroupValidateMessage {
	this := ErrorGroupValidateMessage{}
	return &this
}

// GetErrorCode returns the ErrorCode field value if set, zero value otherwise.
func (o *ErrorGroupValidateMessage) GetErrorCode() string {
	if o == nil || IsNil(o.ErrorCode) {
		var ret string
		return ret
	}
	return *o.ErrorCode
}

// GetErrorCodeOk returns a tuple with the ErrorCode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ErrorGroupValidateMessage) GetErrorCodeOk() (*string, bool) {
	if o == nil || IsNil(o.ErrorCode) {
		return nil, false
	}
	return o.ErrorCode, true
}

// HasErrorCode returns a boolean if a field has been set.
func (o *ErrorGroupValidateMessage) HasErrorCode() bool {
	if o != nil && !IsNil(o.ErrorCode) {
		return true
	}

	return false
}

// SetErrorCode gets a reference to the given string and assigns it to the ErrorCode field.
func (o *ErrorGroupValidateMessage) SetErrorCode(v string) {
	o.ErrorCode = &v
}

// GetMessage returns the Message field value if set, zero value otherwise.
func (o *ErrorGroupValidateMessage) GetMessage() string {
	if o == nil || IsNil(o.Message) {
		var ret string
		return ret
	}
	return *o.Message
}

// GetMessageOk returns a tuple with the Message field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ErrorGroupValidateMessage) GetMessageOk() (*string, bool) {
	if o == nil || IsNil(o.Message) {
		return nil, false
	}
	return o.Message, true
}

// HasMessage returns a boolean if a field has been set.
func (o *ErrorGroupValidateMessage) HasMessage() bool {
	if o != nil && !IsNil(o.Message) {
		return true
	}

	return false
}

// SetMessage gets a reference to the given string and assigns it to the Message field.
func (o *ErrorGroupValidateMessage) SetMessage(v string) {
	o.Message = &v
}

func (o ErrorGroupValidateMessage) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ErrorGroupValidateMessage) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.ErrorCode) {
		toSerialize["errorCode"] = o.ErrorCode
	}
	if !IsNil(o.Message) {
		toSerialize["message"] = o.Message
	}
	return toSerialize, nil
}

type NullableErrorGroupValidateMessage struct {
	value *ErrorGroupValidateMessage
	isSet bool
}

func (v NullableErrorGroupValidateMessage) Get() *ErrorGroupValidateMessage {
	return v.value
}

func (v *NullableErrorGroupValidateMessage) Set(val *ErrorGroupValidateMessage) {
	v.value = val
	v.isSet = true
}

func (v NullableErrorGroupValidateMessage) IsSet() bool {
	return v.isSet
}

func (v *NullableErrorGroupValidateMessage) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableErrorGroupValidateMessage(val *ErrorGroupValidateMessage) *NullableErrorGroupValidateMessage {
	return &NullableErrorGroupValidateMessage{value: val, isSet: true}
}

func (v NullableErrorGroupValidateMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableErrorGroupValidateMessage) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
