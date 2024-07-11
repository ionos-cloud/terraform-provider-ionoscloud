/*
 * IONOS S3 Object Storage API for contract-owned buckets
 *
 * ## Overview The IONOS S3 Object Storage API for contract-owned buckets is a REST-based API that allows developers and applications to interact directly with IONOS' scalable storage solution, leveraging the S3 protocol for object storage operations. Its design ensures seamless compatibility with existing tools and libraries tailored for S3 systems.  ### API References - [S3 Management API Reference](https://api.ionos.com/docs/s3-management/v1/) for managing Access Keys - S3 API Reference for contract-owned buckets - current document - [S3 API Reference for user-owned buckets](https://api.ionos.com/docs/s3-user-owned-buckets/v2/)  ### User documentation [IONOS S3 Object Storage User Guide](https://docs.ionos.com/cloud/managed-services/s3-object-storage) * [Documentation on user-owned and contract-owned buckets](https://docs.ionos.com/cloud/managed-services/s3-object-storage/concepts/buckets) * [Documentation on S3 API Compatibility](https://docs.ionos.com/cloud/managed-services/s3-object-storage/concepts/s3-api-compatibility) * [S3 Tools](https://docs.ionos.com/cloud/managed-services/s3-object-storage/s3-tools)  ## Endpoints for contract-owned buckets | Location | Region Name | Bucket Type | Endpoint | | --- | --- | --- | --- | | **Berlin, Germany** | **eu-central-3** | Contract-owned | `https://s3.eu-central-3.ionoscloud.com` |  ## Changelog - 30.05.2024 Initial version
 *
 * API version: 2.0.2
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

import "encoding/xml"

// checks if the ReplicationConfiguration type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ReplicationConfiguration{}

// ReplicationConfiguration A container for replication rules. You can add up to 1,000 rules. The maximum size of a replication configuration is 2 MB.
type ReplicationConfiguration struct {
	XMLName xml.Name `xml:"ReplicationConfiguration"`
	// The Resource Name of the Identity and Access Management (IAM) role that IONOS S3 Object Storage assumes when replicating objects.
	Role  string            `json:"Role" xml:"Role"`
	Rules []ReplicationRule `json:"Rules" xml:"Rules"`
}

// NewReplicationConfiguration instantiates a new ReplicationConfiguration object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewReplicationConfiguration(role string, rules []ReplicationRule) *ReplicationConfiguration {
	this := ReplicationConfiguration{}

	this.Role = role
	this.Rules = rules

	return &this
}

// NewReplicationConfigurationWithDefaults instantiates a new ReplicationConfiguration object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewReplicationConfigurationWithDefaults() *ReplicationConfiguration {
	this := ReplicationConfiguration{}
	return &this
}

// GetRole returns the Role field value
func (o *ReplicationConfiguration) GetRole() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Role
}

// GetRoleOk returns a tuple with the Role field value
// and a boolean to check if the value has been set.
func (o *ReplicationConfiguration) GetRoleOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Role, true
}

// SetRole sets field value
func (o *ReplicationConfiguration) SetRole(v string) {
	o.Role = v
}

// GetRules returns the Rules field value
func (o *ReplicationConfiguration) GetRules() []ReplicationRule {
	if o == nil {
		var ret []ReplicationRule
		return ret
	}

	return o.Rules
}

// GetRulesOk returns a tuple with the Rules field value
// and a boolean to check if the value has been set.
func (o *ReplicationConfiguration) GetRulesOk() ([]ReplicationRule, bool) {
	if o == nil {
		return nil, false
	}
	return o.Rules, true
}

// SetRules sets field value
func (o *ReplicationConfiguration) SetRules(v []ReplicationRule) {
	o.Rules = v
}

func (o ReplicationConfiguration) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ReplicationConfiguration) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsZero(o.Role) {
		toSerialize["Role"] = o.Role
	}
	if !IsZero(o.Rules) {
		toSerialize["Rules"] = o.Rules
	}
	return toSerialize, nil
}

type NullableReplicationConfiguration struct {
	value *ReplicationConfiguration
	isSet bool
}

func (v NullableReplicationConfiguration) Get() *ReplicationConfiguration {
	return v.value
}

func (v *NullableReplicationConfiguration) Set(val *ReplicationConfiguration) {
	v.value = val
	v.isSet = true
}

func (v NullableReplicationConfiguration) IsSet() bool {
	return v.isSet
}

func (v *NullableReplicationConfiguration) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableReplicationConfiguration(val *ReplicationConfiguration) *NullableReplicationConfiguration {
	return &NullableReplicationConfiguration{value: val, isSet: true}
}

func (v NullableReplicationConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableReplicationConfiguration) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
