/*
 * VM Auto Scaling API
 *
 * The VM Auto Scaling Service enables IONOS clients to horizontally scale the number of VM replicas based on configured rules. You can use VM Auto Scaling to ensure that you have a sufficient number of replicas to handle your application loads at all times.  For this purpose, create a VM Auto Scaling Group that contains the server replicas. The VM Auto Scaling Service ensures that the number of replicas in the group is always within the defined limits.   When scaling policies are set, VM Auto Scaling creates or deletes replicas according to the requirements of your applications. For each policy, specified 'scale-in' and 'scale-out' actions are performed when the corresponding thresholds are reached.
 *
 * API version: 1-SDK.1
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// ErrorMessage struct for ErrorMessage
type ErrorMessage struct {
	ErrorCode *string `json:"errorCode,omitempty"`
	Message   *string `json:"message,omitempty"`
}

// NewErrorMessage instantiates a new ErrorMessage object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewErrorMessage() *ErrorMessage {
	this := ErrorMessage{}

	return &this
}

// NewErrorMessageWithDefaults instantiates a new ErrorMessage object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewErrorMessageWithDefaults() *ErrorMessage {
	this := ErrorMessage{}
	return &this
}

// GetErrorCode returns the ErrorCode field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ErrorMessage) GetErrorCode() *string {
	if o == nil {
		return nil
	}

	return o.ErrorCode

}

// GetErrorCodeOk returns a tuple with the ErrorCode field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ErrorMessage) GetErrorCodeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.ErrorCode, true
}

// SetErrorCode sets field value
func (o *ErrorMessage) SetErrorCode(v string) {

	o.ErrorCode = &v

}

// HasErrorCode returns a boolean if a field has been set.
func (o *ErrorMessage) HasErrorCode() bool {
	if o != nil && o.ErrorCode != nil {
		return true
	}

	return false
}

// GetMessage returns the Message field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ErrorMessage) GetMessage() *string {
	if o == nil {
		return nil
	}

	return o.Message

}

// GetMessageOk returns a tuple with the Message field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ErrorMessage) GetMessageOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Message, true
}

// SetMessage sets field value
func (o *ErrorMessage) SetMessage(v string) {

	o.Message = &v

}

// HasMessage returns a boolean if a field has been set.
func (o *ErrorMessage) HasMessage() bool {
	if o != nil && o.Message != nil {
		return true
	}

	return false
}

func (o ErrorMessage) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.ErrorCode != nil {
		toSerialize["errorCode"] = o.ErrorCode
	}

	if o.Message != nil {
		toSerialize["message"] = o.Message
	}

	return json.Marshal(toSerialize)
}

type NullableErrorMessage struct {
	value *ErrorMessage
	isSet bool
}

func (v NullableErrorMessage) Get() *ErrorMessage {
	return v.value
}

func (v *NullableErrorMessage) Set(val *ErrorMessage) {
	v.value = val
	v.isSet = true
}

func (v NullableErrorMessage) IsSet() bool {
	return v.isSet
}

func (v *NullableErrorMessage) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableErrorMessage(val *ErrorMessage) *NullableErrorMessage {
	return &NullableErrorMessage{value: val, isSet: true}
}

func (v NullableErrorMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableErrorMessage) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
