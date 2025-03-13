/*
 * IONOS Cloud - DNS API
 *
 * Cloud DNS service helps IONOS Cloud customers to automate DNS Zone and Record management.
 *
 * API version: 1.17.0
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dns

import (
	"encoding/json"
)

// checks if the Zone type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Zone{}

// Zone struct for Zone
type Zone struct {
	// The zone name
	ZoneName string `json:"zoneName"`
	// The hosted zone is used for...
	Description *string `json:"description,omitempty"`
	// Users can activate and deactivate zones.
	Enabled *bool `json:"enabled,omitempty"`
}

// NewZone instantiates a new Zone object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewZone(zoneName string) *Zone {
	this := Zone{}

	this.ZoneName = zoneName
	var enabled bool = true
	this.Enabled = &enabled

	return &this
}

// NewZoneWithDefaults instantiates a new Zone object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewZoneWithDefaults() *Zone {
	this := Zone{}
	var enabled bool = true
	this.Enabled = &enabled
	return &this
}

// GetZoneName returns the ZoneName field value
func (o *Zone) GetZoneName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ZoneName
}

// GetZoneNameOk returns a tuple with the ZoneName field value
// and a boolean to check if the value has been set.
func (o *Zone) GetZoneNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ZoneName, true
}

// SetZoneName sets field value
func (o *Zone) SetZoneName(v string) {
	o.ZoneName = v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *Zone) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Zone) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *Zone) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *Zone) SetDescription(v string) {
	o.Description = &v
}

// GetEnabled returns the Enabled field value if set, zero value otherwise.
func (o *Zone) GetEnabled() bool {
	if o == nil || IsNil(o.Enabled) {
		var ret bool
		return ret
	}
	return *o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Zone) GetEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.Enabled) {
		return nil, false
	}
	return o.Enabled, true
}

// HasEnabled returns a boolean if a field has been set.
func (o *Zone) HasEnabled() bool {
	if o != nil && !IsNil(o.Enabled) {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given bool and assigns it to the Enabled field.
func (o *Zone) SetEnabled(v bool) {
	o.Enabled = &v
}

func (o Zone) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Zone) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["zoneName"] = o.ZoneName
	if !IsNil(o.Description) {
		toSerialize["description"] = o.Description
	}
	if !IsNil(o.Enabled) {
		toSerialize["enabled"] = o.Enabled
	}
	return toSerialize, nil
}

type NullableZone struct {
	value *Zone
	isSet bool
}

func (v NullableZone) Get() *Zone {
	return v.value
}

func (v *NullableZone) Set(val *Zone) {
	v.value = val
	v.isSet = true
}

func (v NullableZone) IsSet() bool {
	return v.isSet
}

func (v *NullableZone) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableZone(val *Zone) *NullableZone {
	return &NullableZone{value: val, isSet: true}
}

func (v NullableZone) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableZone) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
