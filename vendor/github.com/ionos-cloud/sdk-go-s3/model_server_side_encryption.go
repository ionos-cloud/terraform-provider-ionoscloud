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

	"fmt"
)

// ServerSideEncryption Server-side encryption algorithm for the default encryption. The valid value is `AES256`.
type ServerSideEncryption string

// List of ServerSideEncryption
const (
	SERVERSIDEENCRYPTION_AES256 ServerSideEncryption = "AES256"
)

func (v *ServerSideEncryption) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ServerSideEncryption(value)
	for _, existing := range []ServerSideEncryption{"AES256"} {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid ServerSideEncryption", value)
}

// Ptr returns reference to ServerSideEncryption value
func (v ServerSideEncryption) Ptr() *ServerSideEncryption {
	return &v
}

type NullableServerSideEncryption struct {
	value *ServerSideEncryption
	isSet bool
}

func (v NullableServerSideEncryption) Get() *ServerSideEncryption {
	return v.value
}

func (v *NullableServerSideEncryption) Set(val *ServerSideEncryption) {
	v.value = val
	v.isSet = true
}

func (v NullableServerSideEncryption) IsSet() bool {
	return v.isSet
}

func (v *NullableServerSideEncryption) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableServerSideEncryption(val *ServerSideEncryption) *NullableServerSideEncryption {
	return &NullableServerSideEncryption{value: val, isSet: true}
}

func (v NullableServerSideEncryption) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableServerSideEncryption) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
