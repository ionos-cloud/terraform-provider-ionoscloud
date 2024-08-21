/*
 * IONOS Cloud - CDN Distribution API
 *
 * This API manages CDN distributions.
 *
 * API version: 0.1.7
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package cdn

import (
	"encoding/json"
)

// checks if the Distribution type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Distribution{}

// Distribution struct for Distribution
type Distribution struct {
	// The ID (UUID) of the Distribution.
	Id string `json:"id"`
	// The type of the resource.
	Type string `json:"type"`
	// The URL of the Distribution.
	Href       string                 `json:"href"`
	Metadata   DistributionMetadata   `json:"metadata"`
	Properties DistributionProperties `json:"properties"`
}

// NewDistribution instantiates a new Distribution object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDistribution(id string, type_ string, href string, metadata DistributionMetadata, properties DistributionProperties) *Distribution {
	this := Distribution{}

	this.Id = id
	this.Type = type_
	this.Href = href
	this.Metadata = metadata
	this.Properties = properties

	return &this
}

// NewDistributionWithDefaults instantiates a new Distribution object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDistributionWithDefaults() *Distribution {
	this := Distribution{}
	return &this
}

// GetId returns the Id field value
func (o *Distribution) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *Distribution) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *Distribution) SetId(v string) {
	o.Id = v
}

// GetType returns the Type field value
func (o *Distribution) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *Distribution) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *Distribution) SetType(v string) {
	o.Type = v
}

// GetHref returns the Href field value
func (o *Distribution) GetHref() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Href
}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
func (o *Distribution) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Href, true
}

// SetHref sets field value
func (o *Distribution) SetHref(v string) {
	o.Href = v
}

// GetMetadata returns the Metadata field value
func (o *Distribution) GetMetadata() DistributionMetadata {
	if o == nil {
		var ret DistributionMetadata
		return ret
	}

	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
func (o *Distribution) GetMetadataOk() (*DistributionMetadata, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Metadata, true
}

// SetMetadata sets field value
func (o *Distribution) SetMetadata(v DistributionMetadata) {
	o.Metadata = v
}

// GetProperties returns the Properties field value
func (o *Distribution) GetProperties() DistributionProperties {
	if o == nil {
		var ret DistributionProperties
		return ret
	}

	return o.Properties
}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
func (o *Distribution) GetPropertiesOk() (*DistributionProperties, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Properties, true
}

// SetProperties sets field value
func (o *Distribution) SetProperties(v DistributionProperties) {
	o.Properties = v
}

func (o Distribution) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["type"] = o.Type
	toSerialize["href"] = o.Href
	toSerialize["metadata"] = o.Metadata
	toSerialize["properties"] = o.Properties
	return toSerialize, nil
}

type NullableDistribution struct {
	value *Distribution
	isSet bool
}

func (v NullableDistribution) Get() *Distribution {
	return v.value
}

func (v *NullableDistribution) Set(val *Distribution) {
	v.value = val
	v.isSet = true
}

func (v NullableDistribution) IsSet() bool {
	return v.isSet
}

func (v *NullableDistribution) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDistribution(val *Distribution) *NullableDistribution {
	return &NullableDistribution{value: val, isSet: true}
}

func (v NullableDistribution) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDistribution) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
