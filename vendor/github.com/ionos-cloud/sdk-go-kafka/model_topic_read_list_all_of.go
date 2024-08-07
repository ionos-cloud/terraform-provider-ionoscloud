/*
 * Kafka as a Service API
 *
 * An managed Apache Kafka cluster is designed to be highly fault-tolerant and scalable, allowing large volumes of data to be ingested, stored, and processed in real-time. By distributing data across multiple brokers, Kafka achieves high throughput and low latency, making it suitable for applications requiring real-time data processing and analytics.
 *
 * API version: 1.7.1
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// TopicReadListAllOf struct for TopicReadListAllOf
type TopicReadListAllOf struct {
	// ID of the list of Topic resources.
	Id *string `json:"id"`
	// The type of the resource.
	Type *string `json:"type"`
	// The URL of the list of Topic resources.
	Href *string `json:"href"`
	// The list of Topic resources.
	Items *[]TopicRead `json:"items,omitempty"`
}

// NewTopicReadListAllOf instantiates a new TopicReadListAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTopicReadListAllOf(id string, type_ string, href string) *TopicReadListAllOf {
	this := TopicReadListAllOf{}

	this.Id = &id
	this.Type = &type_
	this.Href = &href

	return &this
}

// NewTopicReadListAllOfWithDefaults instantiates a new TopicReadListAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTopicReadListAllOfWithDefaults() *TopicReadListAllOf {
	this := TopicReadListAllOf{}
	return &this
}

// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *TopicReadListAllOf) GetId() *string {
	if o == nil {
		return nil
	}

	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TopicReadListAllOf) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Id, true
}

// SetId sets field value
func (o *TopicReadListAllOf) SetId(v string) {

	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *TopicReadListAllOf) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// GetType returns the Type field value
// If the value is explicit nil, the zero value for string will be returned
func (o *TopicReadListAllOf) GetType() *string {
	if o == nil {
		return nil
	}

	return o.Type

}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TopicReadListAllOf) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Type, true
}

// SetType sets field value
func (o *TopicReadListAllOf) SetType(v string) {

	o.Type = &v

}

// HasType returns a boolean if a field has been set.
func (o *TopicReadListAllOf) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// GetHref returns the Href field value
// If the value is explicit nil, the zero value for string will be returned
func (o *TopicReadListAllOf) GetHref() *string {
	if o == nil {
		return nil
	}

	return o.Href

}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TopicReadListAllOf) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Href, true
}

// SetHref sets field value
func (o *TopicReadListAllOf) SetHref(v string) {

	o.Href = &v

}

// HasHref returns a boolean if a field has been set.
func (o *TopicReadListAllOf) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}

// GetItems returns the Items field value
// If the value is explicit nil, the zero value for []TopicRead will be returned
func (o *TopicReadListAllOf) GetItems() *[]TopicRead {
	if o == nil {
		return nil
	}

	return o.Items

}

// GetItemsOk returns a tuple with the Items field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TopicReadListAllOf) GetItemsOk() (*[]TopicRead, bool) {
	if o == nil {
		return nil, false
	}

	return o.Items, true
}

// SetItems sets field value
func (o *TopicReadListAllOf) SetItems(v []TopicRead) {

	o.Items = &v

}

// HasItems returns a boolean if a field has been set.
func (o *TopicReadListAllOf) HasItems() bool {
	if o != nil && o.Items != nil {
		return true
	}

	return false
}

func (o TopicReadListAllOf) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Id != nil {
		toSerialize["id"] = o.Id
	}

	if o.Type != nil {
		toSerialize["type"] = o.Type
	}

	if o.Href != nil {
		toSerialize["href"] = o.Href
	}

	if o.Items != nil {
		toSerialize["items"] = o.Items
	}

	return json.Marshal(toSerialize)
}

type NullableTopicReadListAllOf struct {
	value *TopicReadListAllOf
	isSet bool
}

func (v NullableTopicReadListAllOf) Get() *TopicReadListAllOf {
	return v.value
}

func (v *NullableTopicReadListAllOf) Set(val *TopicReadListAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableTopicReadListAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableTopicReadListAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTopicReadListAllOf(val *TopicReadListAllOf) *NullableTopicReadListAllOf {
	return &NullableTopicReadListAllOf{value: val, isSet: true}
}

func (v NullableTopicReadListAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTopicReadListAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
