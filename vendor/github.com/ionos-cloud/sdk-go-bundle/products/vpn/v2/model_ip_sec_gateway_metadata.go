/*
 * VPN Gateways
 *
 * POC Docs for VPN gateway as service
 *
 * API version: 0.0.1
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"

	"time"
)

// checks if the IPSecGatewayMetadata type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &IPSecGatewayMetadata{}

// IPSecGatewayMetadata IPSec Gateway Metadata
type IPSecGatewayMetadata struct {
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
	// The current status of the resource. The status can be:  * `AVAILABLE` - resource exists and is healthy. * `PROVISIONING` - resource is being created or updated. * `DESTROYING` - delete command was issued, the resource is being deleted. * `FAILED`: - resource failed, details in `statusMessage`.
	Status string `json:"status"`
	// The message of the failure if the status is `FAILED`.
	StatusMessage *string `json:"statusMessage,omitempty"`
}

// NewIPSecGatewayMetadata instantiates a new IPSecGatewayMetadata object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewIPSecGatewayMetadata(status string) *IPSecGatewayMetadata {
	this := IPSecGatewayMetadata{}

	this.Status = status

	return &this
}

// NewIPSecGatewayMetadataWithDefaults instantiates a new IPSecGatewayMetadata object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewIPSecGatewayMetadataWithDefaults() *IPSecGatewayMetadata {
	this := IPSecGatewayMetadata{}
	return &this
}

// GetCreatedDate returns the CreatedDate field value if set, zero value otherwise.
func (o *IPSecGatewayMetadata) GetCreatedDate() time.Time {
	if o == nil || IsNil(o.CreatedDate) {
		var ret time.Time
		return ret
	}
	return o.CreatedDate.Time
}

// GetCreatedDateOk returns a tuple with the CreatedDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *IPSecGatewayMetadata) GetCreatedDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreatedDate) {
		return nil, false
	}
	return &o.CreatedDate.Time, true
}

// HasCreatedDate returns a boolean if a field has been set.
func (o *IPSecGatewayMetadata) HasCreatedDate() bool {
	if o != nil && !IsNil(o.CreatedDate) {
		return true
	}

	return false
}

// SetCreatedDate gets a reference to the given time.Time and assigns it to the CreatedDate field.
func (o *IPSecGatewayMetadata) SetCreatedDate(v time.Time) {
	o.CreatedDate = &IonosTime{v}
}

// GetCreatedBy returns the CreatedBy field value if set, zero value otherwise.
func (o *IPSecGatewayMetadata) GetCreatedBy() string {
	if o == nil || IsNil(o.CreatedBy) {
		var ret string
		return ret
	}
	return *o.CreatedBy
}

// GetCreatedByOk returns a tuple with the CreatedBy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *IPSecGatewayMetadata) GetCreatedByOk() (*string, bool) {
	if o == nil || IsNil(o.CreatedBy) {
		return nil, false
	}
	return o.CreatedBy, true
}

// HasCreatedBy returns a boolean if a field has been set.
func (o *IPSecGatewayMetadata) HasCreatedBy() bool {
	if o != nil && !IsNil(o.CreatedBy) {
		return true
	}

	return false
}

// SetCreatedBy gets a reference to the given string and assigns it to the CreatedBy field.
func (o *IPSecGatewayMetadata) SetCreatedBy(v string) {
	o.CreatedBy = &v
}

// GetCreatedByUserId returns the CreatedByUserId field value if set, zero value otherwise.
func (o *IPSecGatewayMetadata) GetCreatedByUserId() string {
	if o == nil || IsNil(o.CreatedByUserId) {
		var ret string
		return ret
	}
	return *o.CreatedByUserId
}

// GetCreatedByUserIdOk returns a tuple with the CreatedByUserId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *IPSecGatewayMetadata) GetCreatedByUserIdOk() (*string, bool) {
	if o == nil || IsNil(o.CreatedByUserId) {
		return nil, false
	}
	return o.CreatedByUserId, true
}

// HasCreatedByUserId returns a boolean if a field has been set.
func (o *IPSecGatewayMetadata) HasCreatedByUserId() bool {
	if o != nil && !IsNil(o.CreatedByUserId) {
		return true
	}

	return false
}

// SetCreatedByUserId gets a reference to the given string and assigns it to the CreatedByUserId field.
func (o *IPSecGatewayMetadata) SetCreatedByUserId(v string) {
	o.CreatedByUserId = &v
}

// GetLastModifiedDate returns the LastModifiedDate field value if set, zero value otherwise.
func (o *IPSecGatewayMetadata) GetLastModifiedDate() time.Time {
	if o == nil || IsNil(o.LastModifiedDate) {
		var ret time.Time
		return ret
	}
	return o.LastModifiedDate.Time
}

// GetLastModifiedDateOk returns a tuple with the LastModifiedDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *IPSecGatewayMetadata) GetLastModifiedDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.LastModifiedDate) {
		return nil, false
	}
	return &o.LastModifiedDate.Time, true
}

// HasLastModifiedDate returns a boolean if a field has been set.
func (o *IPSecGatewayMetadata) HasLastModifiedDate() bool {
	if o != nil && !IsNil(o.LastModifiedDate) {
		return true
	}

	return false
}

// SetLastModifiedDate gets a reference to the given time.Time and assigns it to the LastModifiedDate field.
func (o *IPSecGatewayMetadata) SetLastModifiedDate(v time.Time) {
	o.LastModifiedDate = &IonosTime{v}
}

// GetLastModifiedBy returns the LastModifiedBy field value if set, zero value otherwise.
func (o *IPSecGatewayMetadata) GetLastModifiedBy() string {
	if o == nil || IsNil(o.LastModifiedBy) {
		var ret string
		return ret
	}
	return *o.LastModifiedBy
}

// GetLastModifiedByOk returns a tuple with the LastModifiedBy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *IPSecGatewayMetadata) GetLastModifiedByOk() (*string, bool) {
	if o == nil || IsNil(o.LastModifiedBy) {
		return nil, false
	}
	return o.LastModifiedBy, true
}

// HasLastModifiedBy returns a boolean if a field has been set.
func (o *IPSecGatewayMetadata) HasLastModifiedBy() bool {
	if o != nil && !IsNil(o.LastModifiedBy) {
		return true
	}

	return false
}

// SetLastModifiedBy gets a reference to the given string and assigns it to the LastModifiedBy field.
func (o *IPSecGatewayMetadata) SetLastModifiedBy(v string) {
	o.LastModifiedBy = &v
}

// GetLastModifiedByUserId returns the LastModifiedByUserId field value if set, zero value otherwise.
func (o *IPSecGatewayMetadata) GetLastModifiedByUserId() string {
	if o == nil || IsNil(o.LastModifiedByUserId) {
		var ret string
		return ret
	}
	return *o.LastModifiedByUserId
}

// GetLastModifiedByUserIdOk returns a tuple with the LastModifiedByUserId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *IPSecGatewayMetadata) GetLastModifiedByUserIdOk() (*string, bool) {
	if o == nil || IsNil(o.LastModifiedByUserId) {
		return nil, false
	}
	return o.LastModifiedByUserId, true
}

// HasLastModifiedByUserId returns a boolean if a field has been set.
func (o *IPSecGatewayMetadata) HasLastModifiedByUserId() bool {
	if o != nil && !IsNil(o.LastModifiedByUserId) {
		return true
	}

	return false
}

// SetLastModifiedByUserId gets a reference to the given string and assigns it to the LastModifiedByUserId field.
func (o *IPSecGatewayMetadata) SetLastModifiedByUserId(v string) {
	o.LastModifiedByUserId = &v
}

// GetResourceURN returns the ResourceURN field value if set, zero value otherwise.
func (o *IPSecGatewayMetadata) GetResourceURN() string {
	if o == nil || IsNil(o.ResourceURN) {
		var ret string
		return ret
	}
	return *o.ResourceURN
}

// GetResourceURNOk returns a tuple with the ResourceURN field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *IPSecGatewayMetadata) GetResourceURNOk() (*string, bool) {
	if o == nil || IsNil(o.ResourceURN) {
		return nil, false
	}
	return o.ResourceURN, true
}

// HasResourceURN returns a boolean if a field has been set.
func (o *IPSecGatewayMetadata) HasResourceURN() bool {
	if o != nil && !IsNil(o.ResourceURN) {
		return true
	}

	return false
}

// SetResourceURN gets a reference to the given string and assigns it to the ResourceURN field.
func (o *IPSecGatewayMetadata) SetResourceURN(v string) {
	o.ResourceURN = &v
}

// GetStatus returns the Status field value
func (o *IPSecGatewayMetadata) GetStatus() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Status
}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
func (o *IPSecGatewayMetadata) GetStatusOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Status, true
}

// SetStatus sets field value
func (o *IPSecGatewayMetadata) SetStatus(v string) {
	o.Status = v
}

// GetStatusMessage returns the StatusMessage field value if set, zero value otherwise.
func (o *IPSecGatewayMetadata) GetStatusMessage() string {
	if o == nil || IsNil(o.StatusMessage) {
		var ret string
		return ret
	}
	return *o.StatusMessage
}

// GetStatusMessageOk returns a tuple with the StatusMessage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *IPSecGatewayMetadata) GetStatusMessageOk() (*string, bool) {
	if o == nil || IsNil(o.StatusMessage) {
		return nil, false
	}
	return o.StatusMessage, true
}

// HasStatusMessage returns a boolean if a field has been set.
func (o *IPSecGatewayMetadata) HasStatusMessage() bool {
	if o != nil && !IsNil(o.StatusMessage) {
		return true
	}

	return false
}

// SetStatusMessage gets a reference to the given string and assigns it to the StatusMessage field.
func (o *IPSecGatewayMetadata) SetStatusMessage(v string) {
	o.StatusMessage = &v
}

func (o IPSecGatewayMetadata) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o IPSecGatewayMetadata) ToMap() (map[string]interface{}, error) {
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
	if !IsZero(o.Status) {
		toSerialize["status"] = o.Status
	}
	if !IsNil(o.StatusMessage) {
		toSerialize["statusMessage"] = o.StatusMessage
	}
	return toSerialize, nil
}

type NullableIPSecGatewayMetadata struct {
	value *IPSecGatewayMetadata
	isSet bool
}

func (v NullableIPSecGatewayMetadata) Get() *IPSecGatewayMetadata {
	return v.value
}

func (v *NullableIPSecGatewayMetadata) Set(val *IPSecGatewayMetadata) {
	v.value = val
	v.isSet = true
}

func (v NullableIPSecGatewayMetadata) IsSet() bool {
	return v.isSet
}

func (v *NullableIPSecGatewayMetadata) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableIPSecGatewayMetadata(val *IPSecGatewayMetadata) *NullableIPSecGatewayMetadata {
	return &NullableIPSecGatewayMetadata{value: val, isSet: true}
}

func (v NullableIPSecGatewayMetadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableIPSecGatewayMetadata) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}