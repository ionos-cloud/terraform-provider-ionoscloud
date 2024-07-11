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

// checks if the OutputSerialization type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &OutputSerialization{}

// OutputSerialization Describes how results of the Select job are serialized.
type OutputSerialization struct {
	XMLName xml.Name    `xml:"OutputSerialization"`
	CSV     *CSVOutput  `json:"CSV,omitempty" xml:"CSV"`
	JSON    *JSONOutput `json:"JSON,omitempty" xml:"JSON"`
}

// NewOutputSerialization instantiates a new OutputSerialization object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewOutputSerialization() *OutputSerialization {
	this := OutputSerialization{}

	return &this
}

// NewOutputSerializationWithDefaults instantiates a new OutputSerialization object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewOutputSerializationWithDefaults() *OutputSerialization {
	this := OutputSerialization{}
	return &this
}

// GetCSV returns the CSV field value if set, zero value otherwise.
func (o *OutputSerialization) GetCSV() CSVOutput {
	if o == nil || IsNil(o.CSV) {
		var ret CSVOutput
		return ret
	}
	return *o.CSV
}

// GetCSVOk returns a tuple with the CSV field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OutputSerialization) GetCSVOk() (*CSVOutput, bool) {
	if o == nil || IsNil(o.CSV) {
		return nil, false
	}
	return o.CSV, true
}

// HasCSV returns a boolean if a field has been set.
func (o *OutputSerialization) HasCSV() bool {
	if o != nil && !IsNil(o.CSV) {
		return true
	}

	return false
}

// SetCSV gets a reference to the given CSVOutput and assigns it to the CSV field.
func (o *OutputSerialization) SetCSV(v CSVOutput) {
	o.CSV = &v
}

// GetJSON returns the JSON field value if set, zero value otherwise.
func (o *OutputSerialization) GetJSON() JSONOutput {
	if o == nil || IsNil(o.JSON) {
		var ret JSONOutput
		return ret
	}
	return *o.JSON
}

// GetJSONOk returns a tuple with the JSON field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OutputSerialization) GetJSONOk() (*JSONOutput, bool) {
	if o == nil || IsNil(o.JSON) {
		return nil, false
	}
	return o.JSON, true
}

// HasJSON returns a boolean if a field has been set.
func (o *OutputSerialization) HasJSON() bool {
	if o != nil && !IsNil(o.JSON) {
		return true
	}

	return false
}

// SetJSON gets a reference to the given JSONOutput and assigns it to the JSON field.
func (o *OutputSerialization) SetJSON(v JSONOutput) {
	o.JSON = &v
}

func (o OutputSerialization) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o OutputSerialization) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.CSV) {
		toSerialize["CSV"] = o.CSV
	}
	if !IsNil(o.JSON) {
		toSerialize["JSON"] = o.JSON
	}
	return toSerialize, nil
}

type NullableOutputSerialization struct {
	value *OutputSerialization
	isSet bool
}

func (v NullableOutputSerialization) Get() *OutputSerialization {
	return v.value
}

func (v *NullableOutputSerialization) Set(val *OutputSerialization) {
	v.value = val
	v.isSet = true
}

func (v NullableOutputSerialization) IsSet() bool {
	return v.isSet
}

func (v *NullableOutputSerialization) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableOutputSerialization(val *OutputSerialization) *NullableOutputSerialization {
	return &NullableOutputSerialization{value: val, isSet: true}
}

func (v NullableOutputSerialization) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableOutputSerialization) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
