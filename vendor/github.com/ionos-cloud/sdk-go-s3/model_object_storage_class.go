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

// ObjectStorageClass the model 'ObjectStorageClass'
type ObjectStorageClass string

// List of ObjectStorageClass
const (
	OBJECTSTORAGECLASS_STANDARD ObjectStorageClass = "STANDARD"
)

func (v *ObjectStorageClass) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ObjectStorageClass(value)
	for _, existing := range []ObjectStorageClass{"STANDARD"} {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid ObjectStorageClass", value)
}

// Ptr returns reference to ObjectStorageClass value
func (v ObjectStorageClass) Ptr() *ObjectStorageClass {
	return &v
}

type NullableObjectStorageClass struct {
	value *ObjectStorageClass
	isSet bool
}

func (v NullableObjectStorageClass) Get() *ObjectStorageClass {
	return v.value
}

func (v *NullableObjectStorageClass) Set(val *ObjectStorageClass) {
	v.value = val
	v.isSet = true
}

func (v NullableObjectStorageClass) IsSet() bool {
	return v.isSet
}

func (v *NullableObjectStorageClass) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableObjectStorageClass(val *ObjectStorageClass) *NullableObjectStorageClass {
	return &NullableObjectStorageClass{value: val, isSet: true}
}

func (v NullableObjectStorageClass) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableObjectStorageClass) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
