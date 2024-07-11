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

// checks if the PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo{}

// PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo struct for PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo
type PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo struct {
	XMLName xml.Name `xml:"PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo"`
	// Name of the host where requests are redirected.
	HostName string `json:"HostName" xml:"HostName"`
	// Protocol to use when redirecting requests. The default is the protocol that is used in the original request.
	Protocol *string `json:"Protocol,omitempty" xml:"Protocol"`
}

// NewPutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo instantiates a new PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo(hostName string) *PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo {
	this := PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo{}

	this.HostName = hostName

	return &this
}

// NewPutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsToWithDefaults instantiates a new PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsToWithDefaults() *PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo {
	this := PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo{}
	return &this
}

// GetHostName returns the HostName field value
func (o *PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo) GetHostName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.HostName
}

// GetHostNameOk returns a tuple with the HostName field value
// and a boolean to check if the value has been set.
func (o *PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo) GetHostNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.HostName, true
}

// SetHostName sets field value
func (o *PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo) SetHostName(v string) {
	o.HostName = v
}

// GetProtocol returns the Protocol field value if set, zero value otherwise.
func (o *PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo) GetProtocol() string {
	if o == nil || IsNil(o.Protocol) {
		var ret string
		return ret
	}
	return *o.Protocol
}

// GetProtocolOk returns a tuple with the Protocol field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo) GetProtocolOk() (*string, bool) {
	if o == nil || IsNil(o.Protocol) {
		return nil, false
	}
	return o.Protocol, true
}

// HasProtocol returns a boolean if a field has been set.
func (o *PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo) HasProtocol() bool {
	if o != nil && !IsNil(o.Protocol) {
		return true
	}

	return false
}

// SetProtocol gets a reference to the given string and assigns it to the Protocol field.
func (o *PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo) SetProtocol(v string) {
	o.Protocol = &v
}

func (o PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsZero(o.HostName) {
		toSerialize["HostName"] = o.HostName
	}
	if !IsNil(o.Protocol) {
		toSerialize["Protocol"] = o.Protocol
	}
	return toSerialize, nil
}

type NullablePutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo struct {
	value *PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo
	isSet bool
}

func (v NullablePutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo) Get() *PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo {
	return v.value
}

func (v *NullablePutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo) Set(val *PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo) {
	v.value = val
	v.isSet = true
}

func (v NullablePutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo) IsSet() bool {
	return v.isSet
}

func (v *NullablePutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo(val *PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo) *NullablePutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo {
	return &NullablePutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo{value: val, isSet: true}
}

func (v NullablePutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
