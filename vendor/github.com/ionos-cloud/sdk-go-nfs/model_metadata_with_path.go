/*
 * IONOS Cloud - Network File Storage API
 *
 * The RESTful API for managing Network File Storage.
 *
 * API version: 0.1.1
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
	"time"
)

// MetadataWithPath struct for MetadataWithPath
type MetadataWithPath struct {
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
	// The status of the resource can be one of the following:  * `AVAILABLE` - The resource exists and is healthy. * `PROVISIONING` - The resource is being created or updated. * `DESTROYING` - A delete command was issued, and the resource is being deleted. * `FAILED` - The resource failed, with details provided in `statusMessage`.
	Status *string `json:"status"`
	// The message of the failure if the status is `FAILED`.
	StatusMessage *string `json:"statusMessage,omitempty"`
	// The path of the NFS export.
	NfsPath *string `json:"nfsPath"`
}

// NewMetadataWithPath instantiates a new MetadataWithPath object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMetadataWithPath(status string, nfsPath string) *MetadataWithPath {
	this := MetadataWithPath{}

	this.Status = &status
	this.NfsPath = &nfsPath

	return &this
}

// NewMetadataWithPathWithDefaults instantiates a new MetadataWithPath object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMetadataWithPathWithDefaults() *MetadataWithPath {
	this := MetadataWithPath{}
	return &this
}

// GetCreatedDate returns the CreatedDate field value
// If the value is explicit nil, the zero value for time.Time will be returned
func (o *MetadataWithPath) GetCreatedDate() *time.Time {
	if o == nil {
		return nil
	}

	if o.CreatedDate == nil {
		return nil
	}
	return &o.CreatedDate.Time

}

// GetCreatedDateOk returns a tuple with the CreatedDate field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *MetadataWithPath) GetCreatedDateOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}

	if o.CreatedDate == nil {
		return nil, false
	}
	return &o.CreatedDate.Time, true

}

// SetCreatedDate sets field value
func (o *MetadataWithPath) SetCreatedDate(v time.Time) {

	o.CreatedDate = &IonosTime{v}

}

// HasCreatedDate returns a boolean if a field has been set.
func (o *MetadataWithPath) HasCreatedDate() bool {
	if o != nil && o.CreatedDate != nil {
		return true
	}

	return false
}

// GetCreatedBy returns the CreatedBy field value
// If the value is explicit nil, the zero value for string will be returned
func (o *MetadataWithPath) GetCreatedBy() *string {
	if o == nil {
		return nil
	}

	return o.CreatedBy

}

// GetCreatedByOk returns a tuple with the CreatedBy field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *MetadataWithPath) GetCreatedByOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.CreatedBy, true
}

// SetCreatedBy sets field value
func (o *MetadataWithPath) SetCreatedBy(v string) {

	o.CreatedBy = &v

}

// HasCreatedBy returns a boolean if a field has been set.
func (o *MetadataWithPath) HasCreatedBy() bool {
	if o != nil && o.CreatedBy != nil {
		return true
	}

	return false
}

// GetCreatedByUserId returns the CreatedByUserId field value
// If the value is explicit nil, the zero value for string will be returned
func (o *MetadataWithPath) GetCreatedByUserId() *string {
	if o == nil {
		return nil
	}

	return o.CreatedByUserId

}

// GetCreatedByUserIdOk returns a tuple with the CreatedByUserId field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *MetadataWithPath) GetCreatedByUserIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.CreatedByUserId, true
}

// SetCreatedByUserId sets field value
func (o *MetadataWithPath) SetCreatedByUserId(v string) {

	o.CreatedByUserId = &v

}

// HasCreatedByUserId returns a boolean if a field has been set.
func (o *MetadataWithPath) HasCreatedByUserId() bool {
	if o != nil && o.CreatedByUserId != nil {
		return true
	}

	return false
}

// GetLastModifiedDate returns the LastModifiedDate field value
// If the value is explicit nil, the zero value for time.Time will be returned
func (o *MetadataWithPath) GetLastModifiedDate() *time.Time {
	if o == nil {
		return nil
	}

	if o.LastModifiedDate == nil {
		return nil
	}
	return &o.LastModifiedDate.Time

}

// GetLastModifiedDateOk returns a tuple with the LastModifiedDate field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *MetadataWithPath) GetLastModifiedDateOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}

	if o.LastModifiedDate == nil {
		return nil, false
	}
	return &o.LastModifiedDate.Time, true

}

// SetLastModifiedDate sets field value
func (o *MetadataWithPath) SetLastModifiedDate(v time.Time) {

	o.LastModifiedDate = &IonosTime{v}

}

// HasLastModifiedDate returns a boolean if a field has been set.
func (o *MetadataWithPath) HasLastModifiedDate() bool {
	if o != nil && o.LastModifiedDate != nil {
		return true
	}

	return false
}

// GetLastModifiedBy returns the LastModifiedBy field value
// If the value is explicit nil, the zero value for string will be returned
func (o *MetadataWithPath) GetLastModifiedBy() *string {
	if o == nil {
		return nil
	}

	return o.LastModifiedBy

}

// GetLastModifiedByOk returns a tuple with the LastModifiedBy field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *MetadataWithPath) GetLastModifiedByOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.LastModifiedBy, true
}

// SetLastModifiedBy sets field value
func (o *MetadataWithPath) SetLastModifiedBy(v string) {

	o.LastModifiedBy = &v

}

// HasLastModifiedBy returns a boolean if a field has been set.
func (o *MetadataWithPath) HasLastModifiedBy() bool {
	if o != nil && o.LastModifiedBy != nil {
		return true
	}

	return false
}

// GetLastModifiedByUserId returns the LastModifiedByUserId field value
// If the value is explicit nil, the zero value for string will be returned
func (o *MetadataWithPath) GetLastModifiedByUserId() *string {
	if o == nil {
		return nil
	}

	return o.LastModifiedByUserId

}

// GetLastModifiedByUserIdOk returns a tuple with the LastModifiedByUserId field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *MetadataWithPath) GetLastModifiedByUserIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.LastModifiedByUserId, true
}

// SetLastModifiedByUserId sets field value
func (o *MetadataWithPath) SetLastModifiedByUserId(v string) {

	o.LastModifiedByUserId = &v

}

// HasLastModifiedByUserId returns a boolean if a field has been set.
func (o *MetadataWithPath) HasLastModifiedByUserId() bool {
	if o != nil && o.LastModifiedByUserId != nil {
		return true
	}

	return false
}

// GetResourceURN returns the ResourceURN field value
// If the value is explicit nil, the zero value for string will be returned
func (o *MetadataWithPath) GetResourceURN() *string {
	if o == nil {
		return nil
	}

	return o.ResourceURN

}

// GetResourceURNOk returns a tuple with the ResourceURN field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *MetadataWithPath) GetResourceURNOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.ResourceURN, true
}

// SetResourceURN sets field value
func (o *MetadataWithPath) SetResourceURN(v string) {

	o.ResourceURN = &v

}

// HasResourceURN returns a boolean if a field has been set.
func (o *MetadataWithPath) HasResourceURN() bool {
	if o != nil && o.ResourceURN != nil {
		return true
	}

	return false
}

// GetStatus returns the Status field value
// If the value is explicit nil, the zero value for string will be returned
func (o *MetadataWithPath) GetStatus() *string {
	if o == nil {
		return nil
	}

	return o.Status

}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *MetadataWithPath) GetStatusOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Status, true
}

// SetStatus sets field value
func (o *MetadataWithPath) SetStatus(v string) {

	o.Status = &v

}

// HasStatus returns a boolean if a field has been set.
func (o *MetadataWithPath) HasStatus() bool {
	if o != nil && o.Status != nil {
		return true
	}

	return false
}

// GetStatusMessage returns the StatusMessage field value
// If the value is explicit nil, the zero value for string will be returned
func (o *MetadataWithPath) GetStatusMessage() *string {
	if o == nil {
		return nil
	}

	return o.StatusMessage

}

// GetStatusMessageOk returns a tuple with the StatusMessage field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *MetadataWithPath) GetStatusMessageOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.StatusMessage, true
}

// SetStatusMessage sets field value
func (o *MetadataWithPath) SetStatusMessage(v string) {

	o.StatusMessage = &v

}

// HasStatusMessage returns a boolean if a field has been set.
func (o *MetadataWithPath) HasStatusMessage() bool {
	if o != nil && o.StatusMessage != nil {
		return true
	}

	return false
}

// GetNfsPath returns the NfsPath field value
// If the value is explicit nil, the zero value for string will be returned
func (o *MetadataWithPath) GetNfsPath() *string {
	if o == nil {
		return nil
	}

	return o.NfsPath

}

// GetNfsPathOk returns a tuple with the NfsPath field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *MetadataWithPath) GetNfsPathOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.NfsPath, true
}

// SetNfsPath sets field value
func (o *MetadataWithPath) SetNfsPath(v string) {

	o.NfsPath = &v

}

// HasNfsPath returns a boolean if a field has been set.
func (o *MetadataWithPath) HasNfsPath() bool {
	if o != nil && o.NfsPath != nil {
		return true
	}

	return false
}

func (o MetadataWithPath) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.CreatedDate != nil {
		toSerialize["createdDate"] = o.CreatedDate
	}

	if o.CreatedBy != nil {
		toSerialize["createdBy"] = o.CreatedBy
	}

	if o.CreatedByUserId != nil {
		toSerialize["createdByUserId"] = o.CreatedByUserId
	}

	if o.LastModifiedDate != nil {
		toSerialize["lastModifiedDate"] = o.LastModifiedDate
	}

	if o.LastModifiedBy != nil {
		toSerialize["lastModifiedBy"] = o.LastModifiedBy
	}

	if o.LastModifiedByUserId != nil {
		toSerialize["lastModifiedByUserId"] = o.LastModifiedByUserId
	}

	if o.ResourceURN != nil {
		toSerialize["resourceURN"] = o.ResourceURN
	}

	if o.Status != nil {
		toSerialize["status"] = o.Status
	}

	if o.StatusMessage != nil {
		toSerialize["statusMessage"] = o.StatusMessage
	}

	if o.NfsPath != nil {
		toSerialize["nfsPath"] = o.NfsPath
	}

	return json.Marshal(toSerialize)
}

type NullableMetadataWithPath struct {
	value *MetadataWithPath
	isSet bool
}

func (v NullableMetadataWithPath) Get() *MetadataWithPath {
	return v.value
}

func (v *NullableMetadataWithPath) Set(val *MetadataWithPath) {
	v.value = val
	v.isSet = true
}

func (v NullableMetadataWithPath) IsSet() bool {
	return v.isSet
}

func (v *NullableMetadataWithPath) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMetadataWithPath(val *MetadataWithPath) *NullableMetadataWithPath {
	return &NullableMetadataWithPath{value: val, isSet: true}
}

func (v NullableMetadataWithPath) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMetadataWithPath) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
