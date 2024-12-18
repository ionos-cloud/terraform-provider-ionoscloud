/*
 * IONOS Cloud - Monitoring REST API
 *
 * The monitoring service offers a centralized platform to collect and store metrics.
 *
 * API version: 0.0.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package monitoring

import (
	"encoding/json"

	"time"
)

// checks if the MetadataWithEndpoint type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &MetadataWithEndpoint{}

// MetadataWithEndpoint struct for MetadataWithEndpoint
type MetadataWithEndpoint struct {
	// The ISO 8601 creation timestamp.
	CreatedDate *IonosTime `json:"createdDate,omitempty"`
	// Unique name of the identity that created the resource.
	CreatedBy *string `json:"createdBy,omitempty"`
	// Unique id of the identity that created the resource.
	CreatedByUserId *string `json:"createdByUserId,omitempty"`
	// The ISO 8601 modified timestamp.
	LastModifiedDate *IonosTime `json:"lastModifiedDate,omitempty"`
	// Unique name of the identity that last modified the resource.
	LastModifiedBy *string `json:"lastModifiedBy,omitempty"`
	// Unique id of the identity that last modified the resource.
	LastModifiedByUserId *string `json:"lastModifiedByUserId,omitempty"`
	// Unique name of the resource.
	ResourceURN *string `json:"resourceURN,omitempty"`
	// The status of the object. The status can be: * `AVAILABLE` - resource exists and is healthy. * `PROVISIONING` - resource is being created or updated. * `DESTROYING` - delete command was issued, the resource is being deleted. * `FAILED` - resource failed, details in `failureMessage`.
	Status string `json:"status"`
	// The message of the failure if the status is `FAILED`.
	StatusMessage *string `json:"statusMessage,omitempty"`
	// The authentication key of the monitoring instance.
	Key string `json:"key"`
	// The endpoint of the Grafana instance.
	GrafanaEndpoint string `json:"grafanaEndpoint"`
	// The HTTP endpoint of the monitoring instance.
	HttpEndpoint string `json:"httpEndpoint"`
}

// NewMetadataWithEndpoint instantiates a new MetadataWithEndpoint object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMetadataWithEndpoint(status string, key string, grafanaEndpoint string, httpEndpoint string) *MetadataWithEndpoint {
	this := MetadataWithEndpoint{}

	this.Status = status
	this.Key = key
	this.GrafanaEndpoint = grafanaEndpoint
	this.HttpEndpoint = httpEndpoint

	return &this
}

// NewMetadataWithEndpointWithDefaults instantiates a new MetadataWithEndpoint object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMetadataWithEndpointWithDefaults() *MetadataWithEndpoint {
	this := MetadataWithEndpoint{}
	return &this
}

// GetCreatedDate returns the CreatedDate field value if set, zero value otherwise.
func (o *MetadataWithEndpoint) GetCreatedDate() time.Time {
	if o == nil || IsNil(o.CreatedDate) {
		var ret time.Time
		return ret
	}
	return o.CreatedDate.Time
}

// GetCreatedDateOk returns a tuple with the CreatedDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MetadataWithEndpoint) GetCreatedDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreatedDate) {
		return nil, false
	}
	return &o.CreatedDate.Time, true
}

// HasCreatedDate returns a boolean if a field has been set.
func (o *MetadataWithEndpoint) HasCreatedDate() bool {
	if o != nil && !IsNil(o.CreatedDate) {
		return true
	}

	return false
}

// SetCreatedDate gets a reference to the given time.Time and assigns it to the CreatedDate field.
func (o *MetadataWithEndpoint) SetCreatedDate(v time.Time) {
	o.CreatedDate = &IonosTime{v}
}

// GetCreatedBy returns the CreatedBy field value if set, zero value otherwise.
func (o *MetadataWithEndpoint) GetCreatedBy() string {
	if o == nil || IsNil(o.CreatedBy) {
		var ret string
		return ret
	}
	return *o.CreatedBy
}

// GetCreatedByOk returns a tuple with the CreatedBy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MetadataWithEndpoint) GetCreatedByOk() (*string, bool) {
	if o == nil || IsNil(o.CreatedBy) {
		return nil, false
	}
	return o.CreatedBy, true
}

// HasCreatedBy returns a boolean if a field has been set.
func (o *MetadataWithEndpoint) HasCreatedBy() bool {
	if o != nil && !IsNil(o.CreatedBy) {
		return true
	}

	return false
}

// SetCreatedBy gets a reference to the given string and assigns it to the CreatedBy field.
func (o *MetadataWithEndpoint) SetCreatedBy(v string) {
	o.CreatedBy = &v
}

// GetCreatedByUserId returns the CreatedByUserId field value if set, zero value otherwise.
func (o *MetadataWithEndpoint) GetCreatedByUserId() string {
	if o == nil || IsNil(o.CreatedByUserId) {
		var ret string
		return ret
	}
	return *o.CreatedByUserId
}

// GetCreatedByUserIdOk returns a tuple with the CreatedByUserId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MetadataWithEndpoint) GetCreatedByUserIdOk() (*string, bool) {
	if o == nil || IsNil(o.CreatedByUserId) {
		return nil, false
	}
	return o.CreatedByUserId, true
}

// HasCreatedByUserId returns a boolean if a field has been set.
func (o *MetadataWithEndpoint) HasCreatedByUserId() bool {
	if o != nil && !IsNil(o.CreatedByUserId) {
		return true
	}

	return false
}

// SetCreatedByUserId gets a reference to the given string and assigns it to the CreatedByUserId field.
func (o *MetadataWithEndpoint) SetCreatedByUserId(v string) {
	o.CreatedByUserId = &v
}

// GetLastModifiedDate returns the LastModifiedDate field value if set, zero value otherwise.
func (o *MetadataWithEndpoint) GetLastModifiedDate() time.Time {
	if o == nil || IsNil(o.LastModifiedDate) {
		var ret time.Time
		return ret
	}
	return o.LastModifiedDate.Time
}

// GetLastModifiedDateOk returns a tuple with the LastModifiedDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MetadataWithEndpoint) GetLastModifiedDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.LastModifiedDate) {
		return nil, false
	}
	return &o.LastModifiedDate.Time, true
}

// HasLastModifiedDate returns a boolean if a field has been set.
func (o *MetadataWithEndpoint) HasLastModifiedDate() bool {
	if o != nil && !IsNil(o.LastModifiedDate) {
		return true
	}

	return false
}

// SetLastModifiedDate gets a reference to the given time.Time and assigns it to the LastModifiedDate field.
func (o *MetadataWithEndpoint) SetLastModifiedDate(v time.Time) {
	o.LastModifiedDate = &IonosTime{v}
}

// GetLastModifiedBy returns the LastModifiedBy field value if set, zero value otherwise.
func (o *MetadataWithEndpoint) GetLastModifiedBy() string {
	if o == nil || IsNil(o.LastModifiedBy) {
		var ret string
		return ret
	}
	return *o.LastModifiedBy
}

// GetLastModifiedByOk returns a tuple with the LastModifiedBy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MetadataWithEndpoint) GetLastModifiedByOk() (*string, bool) {
	if o == nil || IsNil(o.LastModifiedBy) {
		return nil, false
	}
	return o.LastModifiedBy, true
}

// HasLastModifiedBy returns a boolean if a field has been set.
func (o *MetadataWithEndpoint) HasLastModifiedBy() bool {
	if o != nil && !IsNil(o.LastModifiedBy) {
		return true
	}

	return false
}

// SetLastModifiedBy gets a reference to the given string and assigns it to the LastModifiedBy field.
func (o *MetadataWithEndpoint) SetLastModifiedBy(v string) {
	o.LastModifiedBy = &v
}

// GetLastModifiedByUserId returns the LastModifiedByUserId field value if set, zero value otherwise.
func (o *MetadataWithEndpoint) GetLastModifiedByUserId() string {
	if o == nil || IsNil(o.LastModifiedByUserId) {
		var ret string
		return ret
	}
	return *o.LastModifiedByUserId
}

// GetLastModifiedByUserIdOk returns a tuple with the LastModifiedByUserId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MetadataWithEndpoint) GetLastModifiedByUserIdOk() (*string, bool) {
	if o == nil || IsNil(o.LastModifiedByUserId) {
		return nil, false
	}
	return o.LastModifiedByUserId, true
}

// HasLastModifiedByUserId returns a boolean if a field has been set.
func (o *MetadataWithEndpoint) HasLastModifiedByUserId() bool {
	if o != nil && !IsNil(o.LastModifiedByUserId) {
		return true
	}

	return false
}

// SetLastModifiedByUserId gets a reference to the given string and assigns it to the LastModifiedByUserId field.
func (o *MetadataWithEndpoint) SetLastModifiedByUserId(v string) {
	o.LastModifiedByUserId = &v
}

// GetResourceURN returns the ResourceURN field value if set, zero value otherwise.
func (o *MetadataWithEndpoint) GetResourceURN() string {
	if o == nil || IsNil(o.ResourceURN) {
		var ret string
		return ret
	}
	return *o.ResourceURN
}

// GetResourceURNOk returns a tuple with the ResourceURN field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MetadataWithEndpoint) GetResourceURNOk() (*string, bool) {
	if o == nil || IsNil(o.ResourceURN) {
		return nil, false
	}
	return o.ResourceURN, true
}

// HasResourceURN returns a boolean if a field has been set.
func (o *MetadataWithEndpoint) HasResourceURN() bool {
	if o != nil && !IsNil(o.ResourceURN) {
		return true
	}

	return false
}

// SetResourceURN gets a reference to the given string and assigns it to the ResourceURN field.
func (o *MetadataWithEndpoint) SetResourceURN(v string) {
	o.ResourceURN = &v
}

// GetStatus returns the Status field value
func (o *MetadataWithEndpoint) GetStatus() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Status
}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
func (o *MetadataWithEndpoint) GetStatusOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Status, true
}

// SetStatus sets field value
func (o *MetadataWithEndpoint) SetStatus(v string) {
	o.Status = v
}

// GetStatusMessage returns the StatusMessage field value if set, zero value otherwise.
func (o *MetadataWithEndpoint) GetStatusMessage() string {
	if o == nil || IsNil(o.StatusMessage) {
		var ret string
		return ret
	}
	return *o.StatusMessage
}

// GetStatusMessageOk returns a tuple with the StatusMessage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MetadataWithEndpoint) GetStatusMessageOk() (*string, bool) {
	if o == nil || IsNil(o.StatusMessage) {
		return nil, false
	}
	return o.StatusMessage, true
}

// HasStatusMessage returns a boolean if a field has been set.
func (o *MetadataWithEndpoint) HasStatusMessage() bool {
	if o != nil && !IsNil(o.StatusMessage) {
		return true
	}

	return false
}

// SetStatusMessage gets a reference to the given string and assigns it to the StatusMessage field.
func (o *MetadataWithEndpoint) SetStatusMessage(v string) {
	o.StatusMessage = &v
}

// GetKey returns the Key field value
func (o *MetadataWithEndpoint) GetKey() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Key
}

// GetKeyOk returns a tuple with the Key field value
// and a boolean to check if the value has been set.
func (o *MetadataWithEndpoint) GetKeyOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Key, true
}

// SetKey sets field value
func (o *MetadataWithEndpoint) SetKey(v string) {
	o.Key = v
}

// GetGrafanaEndpoint returns the GrafanaEndpoint field value
func (o *MetadataWithEndpoint) GetGrafanaEndpoint() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.GrafanaEndpoint
}

// GetGrafanaEndpointOk returns a tuple with the GrafanaEndpoint field value
// and a boolean to check if the value has been set.
func (o *MetadataWithEndpoint) GetGrafanaEndpointOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.GrafanaEndpoint, true
}

// SetGrafanaEndpoint sets field value
func (o *MetadataWithEndpoint) SetGrafanaEndpoint(v string) {
	o.GrafanaEndpoint = v
}

// GetHttpEndpoint returns the HttpEndpoint field value
func (o *MetadataWithEndpoint) GetHttpEndpoint() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.HttpEndpoint
}

// GetHttpEndpointOk returns a tuple with the HttpEndpoint field value
// and a boolean to check if the value has been set.
func (o *MetadataWithEndpoint) GetHttpEndpointOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.HttpEndpoint, true
}

// SetHttpEndpoint sets field value
func (o *MetadataWithEndpoint) SetHttpEndpoint(v string) {
	o.HttpEndpoint = v
}

func (o MetadataWithEndpoint) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.CreatedDate) {
		toSerialize["createdDate"] = o.CreatedDate
	}
	if !IsNil(o.CreatedBy) {
		toSerialize["createdBy"] = o.CreatedBy
	}
	if !IsNil(o.CreatedByUserId) {
		toSerialize["createdByUserId"] = o.CreatedByUserId
	}
	if !IsNil(o.LastModifiedDate) {
		toSerialize["lastModifiedDate"] = o.LastModifiedDate
	}
	if !IsNil(o.LastModifiedBy) {
		toSerialize["lastModifiedBy"] = o.LastModifiedBy
	}
	if !IsNil(o.LastModifiedByUserId) {
		toSerialize["lastModifiedByUserId"] = o.LastModifiedByUserId
	}
	if !IsNil(o.ResourceURN) {
		toSerialize["resourceURN"] = o.ResourceURN
	}
	toSerialize["status"] = o.Status
	if !IsNil(o.StatusMessage) {
		toSerialize["statusMessage"] = o.StatusMessage
	}
	toSerialize["key"] = o.Key
	toSerialize["grafanaEndpoint"] = o.GrafanaEndpoint
	toSerialize["httpEndpoint"] = o.HttpEndpoint
	return toSerialize, nil
}

type NullableMetadataWithEndpoint struct {
	value *MetadataWithEndpoint
	isSet bool
}

func (v NullableMetadataWithEndpoint) Get() *MetadataWithEndpoint {
	return v.value
}

func (v *NullableMetadataWithEndpoint) Set(val *MetadataWithEndpoint) {
	v.value = val
	v.isSet = true
}

func (v NullableMetadataWithEndpoint) IsSet() bool {
	return v.isSet
}

func (v *NullableMetadataWithEndpoint) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMetadataWithEndpoint(val *MetadataWithEndpoint) *NullableMetadataWithEndpoint {
	return &NullableMetadataWithEndpoint{value: val, isSet: true}
}

func (v NullableMetadataWithEndpoint) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMetadataWithEndpoint) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
