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

// InputSerialization Describes the serialization format of the object.
type InputSerialization struct {
	XMLName xml.Name  `xml:"InputSerialization"`
	CSV     *CSVInput `json:"CSV,omitempty" xml:"CSV"`
	// Specifies object's compression format. Valid values: NONE, GZIP, BZIP2. Default Value: NONE.
	CompressionType *string                 `json:"CompressionType,omitempty" xml:"CompressionType"`
	JSON            *InputSerializationJSON `json:"JSON,omitempty" xml:"JSON"`
	// Specifies Parquet as object's input serialization format.
	Parquet *map[string]interface{} `json:"Parquet,omitempty" xml:"Parquet"`
}

// NewInputSerialization instantiates a new InputSerialization object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewInputSerialization() *InputSerialization {
	this := InputSerialization{}

	return &this
}

// NewInputSerializationWithDefaults instantiates a new InputSerialization object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewInputSerializationWithDefaults() *InputSerialization {
	this := InputSerialization{}
	return &this
}

// GetCSV returns the CSV field value
// If the value is explicit nil, the zero value for CSVInput will be returned
func (o *InputSerialization) GetCSV() *CSVInput {
	if o == nil {
		return nil
	}

	return o.CSV

}

// GetCSVOk returns a tuple with the CSV field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *InputSerialization) GetCSVOk() (*CSVInput, bool) {
	if o == nil {
		return nil, false
	}

	return o.CSV, true
}

// SetCSV sets field value
func (o *InputSerialization) SetCSV(v CSVInput) {

	o.CSV = &v

}

// HasCSV returns a boolean if a field has been set.
func (o *InputSerialization) HasCSV() bool {
	if o != nil && o.CSV != nil {
		return true
	}

	return false
}

// GetCompressionType returns the CompressionType field value
// If the value is explicit nil, the zero value for string will be returned
func (o *InputSerialization) GetCompressionType() *string {
	if o == nil {
		return nil
	}

	return o.CompressionType

}

// GetCompressionTypeOk returns a tuple with the CompressionType field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *InputSerialization) GetCompressionTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.CompressionType, true
}

// SetCompressionType sets field value
func (o *InputSerialization) SetCompressionType(v string) {

	o.CompressionType = &v

}

// HasCompressionType returns a boolean if a field has been set.
func (o *InputSerialization) HasCompressionType() bool {
	if o != nil && o.CompressionType != nil {
		return true
	}

	return false
}

// GetJSON returns the JSON field value
// If the value is explicit nil, the zero value for InputSerializationJSON will be returned
func (o *InputSerialization) GetJSON() *InputSerializationJSON {
	if o == nil {
		return nil
	}

	return o.JSON

}

// GetJSONOk returns a tuple with the JSON field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *InputSerialization) GetJSONOk() (*InputSerializationJSON, bool) {
	if o == nil {
		return nil, false
	}

	return o.JSON, true
}

// SetJSON sets field value
func (o *InputSerialization) SetJSON(v InputSerializationJSON) {

	o.JSON = &v

}

// HasJSON returns a boolean if a field has been set.
func (o *InputSerialization) HasJSON() bool {
	if o != nil && o.JSON != nil {
		return true
	}

	return false
}

// GetParquet returns the Parquet field value
// If the value is explicit nil, the zero value for map[string]interface{} will be returned
func (o *InputSerialization) GetParquet() *map[string]interface{} {
	if o == nil {
		return nil
	}

	return o.Parquet

}

// GetParquetOk returns a tuple with the Parquet field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *InputSerialization) GetParquetOk() (*map[string]interface{}, bool) {
	if o == nil {
		return nil, false
	}

	return o.Parquet, true
}

// SetParquet sets field value
func (o *InputSerialization) SetParquet(v map[string]interface{}) {

	o.Parquet = &v

}

// HasParquet returns a boolean if a field has been set.
func (o *InputSerialization) HasParquet() bool {
	if o != nil && o.Parquet != nil {
		return true
	}

	return false
}

func (o InputSerialization) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.CSV != nil {
		toSerialize["CSV"] = o.CSV
	}

	if o.CompressionType != nil {
		toSerialize["CompressionType"] = o.CompressionType
	}

	if o.JSON != nil {
		toSerialize["JSON"] = o.JSON
	}

	if o.Parquet != nil {
		toSerialize["Parquet"] = o.Parquet
	}

	return json.Marshal(toSerialize)
}

type NullableInputSerialization struct {
	value *InputSerialization
	isSet bool
}

func (v NullableInputSerialization) Get() *InputSerialization {
	return v.value
}

func (v *NullableInputSerialization) Set(val *InputSerialization) {
	v.value = val
	v.isSet = true
}

func (v NullableInputSerialization) IsSet() bool {
	return v.isSet
}

func (v *NullableInputSerialization) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableInputSerialization(val *InputSerialization) *NullableInputSerialization {
	return &NullableInputSerialization{value: val, isSet: true}
}

func (v NullableInputSerialization) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableInputSerialization) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
