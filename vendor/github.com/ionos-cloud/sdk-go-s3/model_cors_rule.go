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

// checks if the CORSRule type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CORSRule{}

// CORSRule Specifies a cross-origin access rule for an IONOS S3 Object Storage bucket.
type CORSRule struct {
	XMLName xml.Name `xml:"CORSRule"`
	// Container for the Contract Number of the owner.
	ID *int32 `json:"ID,omitempty" xml:"ID"`
	// Headers that are specified in the `Access-Control-Request-Headers` header. These headers are allowed in a preflight OPTIONS request. In response to any preflight OPTIONS request, IONOS S3 Object Storage returns any requested headers that are allowed.
	AllowedHeaders []string `json:"AllowedHeaders,omitempty" xml:"AllowedHeaders"`
	// An HTTP method that you allow the origin to execute. Valid values are `GET`, `PUT`, `HEAD`, `POST`, and `DELETE`.
	AllowedMethods []string `json:"AllowedMethods" xml:"AllowedMethods"`
	// One or more origins you want customers to be able to access the bucket from.
	AllowedOrigins []string `json:"AllowedOrigins" xml:"AllowedOrigins"`
	// One or more headers in the response that you want customers to be able to access from their applications (for example, from a JavaScript `XMLHttpRequest` object).
	ExposeHeaders []string `json:"ExposeHeaders,omitempty" xml:"ExposeHeaders"`
	// The time in seconds that your browser is to cache the preflight response for the specified resource.
	MaxAgeSeconds *int32 `json:"MaxAgeSeconds,omitempty" xml:"MaxAgeSeconds"`
}

// NewCORSRule instantiates a new CORSRule object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCORSRule(allowedMethods []string, allowedOrigins []string) *CORSRule {
	this := CORSRule{}

	this.AllowedMethods = allowedMethods
	this.AllowedOrigins = allowedOrigins

	return &this
}

// NewCORSRuleWithDefaults instantiates a new CORSRule object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCORSRuleWithDefaults() *CORSRule {
	this := CORSRule{}
	return &this
}

// GetID returns the ID field value if set, zero value otherwise.
func (o *CORSRule) GetID() int32 {
	if o == nil || IsNil(o.ID) {
		var ret int32
		return ret
	}
	return *o.ID
}

// GetIDOk returns a tuple with the ID field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CORSRule) GetIDOk() (*int32, bool) {
	if o == nil || IsNil(o.ID) {
		return nil, false
	}
	return o.ID, true
}

// HasID returns a boolean if a field has been set.
func (o *CORSRule) HasID() bool {
	if o != nil && !IsNil(o.ID) {
		return true
	}

	return false
}

// SetID gets a reference to the given int32 and assigns it to the ID field.
func (o *CORSRule) SetID(v int32) {
	o.ID = &v
}

// GetAllowedHeaders returns the AllowedHeaders field value if set, zero value otherwise.
func (o *CORSRule) GetAllowedHeaders() []string {
	if o == nil || IsNil(o.AllowedHeaders) {
		var ret []string
		return ret
	}
	return o.AllowedHeaders
}

// GetAllowedHeadersOk returns a tuple with the AllowedHeaders field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CORSRule) GetAllowedHeadersOk() ([]string, bool) {
	if o == nil || IsNil(o.AllowedHeaders) {
		return nil, false
	}
	return o.AllowedHeaders, true
}

// HasAllowedHeaders returns a boolean if a field has been set.
func (o *CORSRule) HasAllowedHeaders() bool {
	if o != nil && !IsNil(o.AllowedHeaders) {
		return true
	}

	return false
}

// SetAllowedHeaders gets a reference to the given []string and assigns it to the AllowedHeaders field.
func (o *CORSRule) SetAllowedHeaders(v []string) {
	o.AllowedHeaders = v
}

// GetAllowedMethods returns the AllowedMethods field value
func (o *CORSRule) GetAllowedMethods() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.AllowedMethods
}

// GetAllowedMethodsOk returns a tuple with the AllowedMethods field value
// and a boolean to check if the value has been set.
func (o *CORSRule) GetAllowedMethodsOk() ([]string, bool) {
	if o == nil {
		return nil, false
	}
	return o.AllowedMethods, true
}

// SetAllowedMethods sets field value
func (o *CORSRule) SetAllowedMethods(v []string) {
	o.AllowedMethods = v
}

// GetAllowedOrigins returns the AllowedOrigins field value
func (o *CORSRule) GetAllowedOrigins() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.AllowedOrigins
}

// GetAllowedOriginsOk returns a tuple with the AllowedOrigins field value
// and a boolean to check if the value has been set.
func (o *CORSRule) GetAllowedOriginsOk() ([]string, bool) {
	if o == nil {
		return nil, false
	}
	return o.AllowedOrigins, true
}

// SetAllowedOrigins sets field value
func (o *CORSRule) SetAllowedOrigins(v []string) {
	o.AllowedOrigins = v
}

// GetExposeHeaders returns the ExposeHeaders field value if set, zero value otherwise.
func (o *CORSRule) GetExposeHeaders() []string {
	if o == nil || IsNil(o.ExposeHeaders) {
		var ret []string
		return ret
	}
	return o.ExposeHeaders
}

// GetExposeHeadersOk returns a tuple with the ExposeHeaders field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CORSRule) GetExposeHeadersOk() ([]string, bool) {
	if o == nil || IsNil(o.ExposeHeaders) {
		return nil, false
	}
	return o.ExposeHeaders, true
}

// HasExposeHeaders returns a boolean if a field has been set.
func (o *CORSRule) HasExposeHeaders() bool {
	if o != nil && !IsNil(o.ExposeHeaders) {
		return true
	}

	return false
}

// SetExposeHeaders gets a reference to the given []string and assigns it to the ExposeHeaders field.
func (o *CORSRule) SetExposeHeaders(v []string) {
	o.ExposeHeaders = v
}

// GetMaxAgeSeconds returns the MaxAgeSeconds field value if set, zero value otherwise.
func (o *CORSRule) GetMaxAgeSeconds() int32 {
	if o == nil || IsNil(o.MaxAgeSeconds) {
		var ret int32
		return ret
	}
	return *o.MaxAgeSeconds
}

// GetMaxAgeSecondsOk returns a tuple with the MaxAgeSeconds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CORSRule) GetMaxAgeSecondsOk() (*int32, bool) {
	if o == nil || IsNil(o.MaxAgeSeconds) {
		return nil, false
	}
	return o.MaxAgeSeconds, true
}

// HasMaxAgeSeconds returns a boolean if a field has been set.
func (o *CORSRule) HasMaxAgeSeconds() bool {
	if o != nil && !IsNil(o.MaxAgeSeconds) {
		return true
	}

	return false
}

// SetMaxAgeSeconds gets a reference to the given int32 and assigns it to the MaxAgeSeconds field.
func (o *CORSRule) SetMaxAgeSeconds(v int32) {
	o.MaxAgeSeconds = &v
}

func (o CORSRule) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CORSRule) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.ID) {
		toSerialize["ID"] = o.ID
	}
	if !IsNil(o.AllowedHeaders) {
		toSerialize["AllowedHeaders"] = o.AllowedHeaders
	}
	if !IsZero(o.AllowedMethods) {
		toSerialize["AllowedMethods"] = o.AllowedMethods
	}
	if !IsZero(o.AllowedOrigins) {
		toSerialize["AllowedOrigins"] = o.AllowedOrigins
	}
	if !IsNil(o.ExposeHeaders) {
		toSerialize["ExposeHeaders"] = o.ExposeHeaders
	}
	if !IsNil(o.MaxAgeSeconds) {
		toSerialize["MaxAgeSeconds"] = o.MaxAgeSeconds
	}
	return toSerialize, nil
}

type NullableCORSRule struct {
	value *CORSRule
	isSet bool
}

func (v NullableCORSRule) Get() *CORSRule {
	return v.value
}

func (v *NullableCORSRule) Set(val *CORSRule) {
	v.value = val
	v.isSet = true
}

func (v NullableCORSRule) IsSet() bool {
	return v.isSet
}

func (v *NullableCORSRule) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCORSRule(val *CORSRule) *NullableCORSRule {
	return &NullableCORSRule{value: val, isSet: true}
}

func (v NullableCORSRule) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCORSRule) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}