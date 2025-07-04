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

// checks if the GroupProperties type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GroupProperties{}

// GroupProperties struct for GroupProperties
type GroupProperties struct {
	Datacenter GroupPropertiesDatacenter `json:"datacenter"`
	// The data center location.
	Location string `json:"location"`
	// The maximum value for the number of replicas. Must be >= 0 and <= 100. Will be enforced for both automatic and manual changes.
	MaxReplicaCount *int64 `json:"maxReplicaCount,omitempty"`
	// The minimum value for the number of replicas. Must be >= 0 and <= 100. Will be enforced for both automatic and manual changes
	MinReplicaCount *int64 `json:"minReplicaCount,omitempty"`
	// The name of the VM Auto Scaling Group. This field must not be null or blank.
	Name                 *string                `json:"name,omitempty"`
	Policy               *GroupPolicy           `json:"policy,omitempty"`
	ReplicaConfiguration *ReplicaPropertiesPost `json:"replicaConfiguration,omitempty"`
}

// NewGroupProperties instantiates a new GroupProperties object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGroupProperties(datacenter GroupPropertiesDatacenter, location string) *GroupProperties {
	this := GroupProperties{}

	this.Datacenter = datacenter
	this.Location = location

	return &this
}

// NewGroupPropertiesWithDefaults instantiates a new GroupProperties object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGroupPropertiesWithDefaults() *GroupProperties {
	this := GroupProperties{}
	return &this
}

// GetDatacenter returns the Datacenter field value
func (o *GroupProperties) GetDatacenter() GroupPropertiesDatacenter {
	if o == nil {
		var ret GroupPropertiesDatacenter
		return ret
	}

	return o.Datacenter
}

// GetDatacenterOk returns a tuple with the Datacenter field value
// and a boolean to check if the value has been set.
func (o *GroupProperties) GetDatacenterOk() (*GroupPropertiesDatacenter, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Datacenter, true
}

// SetDatacenter sets field value
func (o *GroupProperties) SetDatacenter(v GroupPropertiesDatacenter) {
	o.Datacenter = v
}

// GetLocation returns the Location field value
func (o *GroupProperties) GetLocation() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Location
}

// GetLocationOk returns a tuple with the Location field value
// and a boolean to check if the value has been set.
func (o *GroupProperties) GetLocationOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Location, true
}

// SetLocation sets field value
func (o *GroupProperties) SetLocation(v string) {
	o.Location = v
}

// GetMaxReplicaCount returns the MaxReplicaCount field value if set, zero value otherwise.
func (o *GroupProperties) GetMaxReplicaCount() int64 {
	if o == nil || IsNil(o.MaxReplicaCount) {
		var ret int64
		return ret
	}
	return *o.MaxReplicaCount
}

// GetMaxReplicaCountOk returns a tuple with the MaxReplicaCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupProperties) GetMaxReplicaCountOk() (*int64, bool) {
	if o == nil || IsNil(o.MaxReplicaCount) {
		return nil, false
	}
	return o.MaxReplicaCount, true
}

// HasMaxReplicaCount returns a boolean if a field has been set.
func (o *GroupProperties) HasMaxReplicaCount() bool {
	if o != nil && !IsNil(o.MaxReplicaCount) {
		return true
	}

	return false
}

// SetMaxReplicaCount gets a reference to the given int64 and assigns it to the MaxReplicaCount field.
func (o *GroupProperties) SetMaxReplicaCount(v int64) {
	o.MaxReplicaCount = &v
}

// GetMinReplicaCount returns the MinReplicaCount field value if set, zero value otherwise.
func (o *GroupProperties) GetMinReplicaCount() int64 {
	if o == nil || IsNil(o.MinReplicaCount) {
		var ret int64
		return ret
	}
	return *o.MinReplicaCount
}

// GetMinReplicaCountOk returns a tuple with the MinReplicaCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupProperties) GetMinReplicaCountOk() (*int64, bool) {
	if o == nil || IsNil(o.MinReplicaCount) {
		return nil, false
	}
	return o.MinReplicaCount, true
}

// HasMinReplicaCount returns a boolean if a field has been set.
func (o *GroupProperties) HasMinReplicaCount() bool {
	if o != nil && !IsNil(o.MinReplicaCount) {
		return true
	}

	return false
}

// SetMinReplicaCount gets a reference to the given int64 and assigns it to the MinReplicaCount field.
func (o *GroupProperties) SetMinReplicaCount(v int64) {
	o.MinReplicaCount = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *GroupProperties) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupProperties) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *GroupProperties) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *GroupProperties) SetName(v string) {
	o.Name = &v
}

// GetPolicy returns the Policy field value if set, zero value otherwise.
func (o *GroupProperties) GetPolicy() GroupPolicy {
	if o == nil || IsNil(o.Policy) {
		var ret GroupPolicy
		return ret
	}
	return *o.Policy
}

// GetPolicyOk returns a tuple with the Policy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupProperties) GetPolicyOk() (*GroupPolicy, bool) {
	if o == nil || IsNil(o.Policy) {
		return nil, false
	}
	return o.Policy, true
}

// HasPolicy returns a boolean if a field has been set.
func (o *GroupProperties) HasPolicy() bool {
	if o != nil && !IsNil(o.Policy) {
		return true
	}

	return false
}

// SetPolicy gets a reference to the given GroupPolicy and assigns it to the Policy field.
func (o *GroupProperties) SetPolicy(v GroupPolicy) {
	o.Policy = &v
}

// GetReplicaConfiguration returns the ReplicaConfiguration field value if set, zero value otherwise.
func (o *GroupProperties) GetReplicaConfiguration() ReplicaPropertiesPost {
	if o == nil || IsNil(o.ReplicaConfiguration) {
		var ret ReplicaPropertiesPost
		return ret
	}
	return *o.ReplicaConfiguration
}

// GetReplicaConfigurationOk returns a tuple with the ReplicaConfiguration field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupProperties) GetReplicaConfigurationOk() (*ReplicaPropertiesPost, bool) {
	if o == nil || IsNil(o.ReplicaConfiguration) {
		return nil, false
	}
	return o.ReplicaConfiguration, true
}

// HasReplicaConfiguration returns a boolean if a field has been set.
func (o *GroupProperties) HasReplicaConfiguration() bool {
	if o != nil && !IsNil(o.ReplicaConfiguration) {
		return true
	}

	return false
}

// SetReplicaConfiguration gets a reference to the given ReplicaPropertiesPost and assigns it to the ReplicaConfiguration field.
func (o *GroupProperties) SetReplicaConfiguration(v ReplicaPropertiesPost) {
	o.ReplicaConfiguration = &v
}

func (o GroupProperties) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GroupProperties) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["datacenter"] = o.Datacenter
	toSerialize["location"] = o.Location
	if !IsNil(o.MaxReplicaCount) {
		toSerialize["maxReplicaCount"] = o.MaxReplicaCount
	}
	if !IsNil(o.MinReplicaCount) {
		toSerialize["minReplicaCount"] = o.MinReplicaCount
	}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Policy) {
		toSerialize["policy"] = o.Policy
	}
	if !IsNil(o.ReplicaConfiguration) {
		toSerialize["replicaConfiguration"] = o.ReplicaConfiguration
	}
	return toSerialize, nil
}

type NullableGroupProperties struct {
	value *GroupProperties
	isSet bool
}

func (v NullableGroupProperties) Get() *GroupProperties {
	return v.value
}

func (v *NullableGroupProperties) Set(val *GroupProperties) {
	v.value = val
	v.isSet = true
}

func (v NullableGroupProperties) IsSet() bool {
	return v.isSet
}

func (v *NullableGroupProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGroupProperties(val *GroupProperties) *NullableGroupProperties {
	return &NullableGroupProperties{value: val, isSet: true}
}

func (v NullableGroupProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGroupProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
