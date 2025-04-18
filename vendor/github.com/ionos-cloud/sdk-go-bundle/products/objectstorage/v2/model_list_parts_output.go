/*
 * IONOS Object Storage API for contract-owned buckets
 *
 * ## Overview The IONOS Object Storage API for contract-owned buckets is a REST-based API that allows developers and applications to interact directly with IONOS' scalable storage solution, leveraging the S3 protocol for object storage operations. Its design ensures seamless compatibility with existing tools and libraries tailored for S3 systems.  ### API References - [S3 API Reference for contract-owned buckets](https://api.ionos.com/docs/s3-contract-owned-buckets/v2/) ### User documentation [IONOS Object Storage User Guide](https://docs.ionos.com/cloud/managed-services/s3-object-storage) * [Documentation on user-owned and contract-owned buckets](https://docs.ionos.com/cloud/managed-services/s3-object-storage/concepts/buckets) * [Documentation on S3 API Compatibility](https://docs.ionos.com/cloud/managed-services/s3-object-storage/concepts/s3-api-compatibility) * [S3 Tools](https://docs.ionos.com/cloud/managed-services/s3-object-storage/s3-tools)  ## Endpoints for contract-owned buckets | Location | Region Name | Bucket Type | Endpoint | | --- | --- | --- | --- | | **Berlin, Germany** | **eu-central-3** | Contract-owned | `https://s3.eu-central-3.ionoscloud.com` |  ## Changelog - 30.05.2024 Initial version
 *
 * API version: 2.0.2
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package objectstorage

import (
	"encoding/json"
)

import "encoding/xml"

// checks if the ListPartsOutput type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ListPartsOutput{}

// ListPartsOutput struct for ListPartsOutput
type ListPartsOutput struct {
	XMLName xml.Name `xml:"ListPartsOutput"`
	// The bucket name.
	Bucket *string `json:"Bucket,omitempty" xml:"Name"`
	// The object key.
	Key *string `json:"Key,omitempty" xml:"Key"`
	// ID of the multipart upload.
	UploadId *string `json:"UploadId,omitempty" xml:"UploadId"`
	// When a list is truncated, this element specifies the last part in the list, as well as the value to use for the part-number-marker request parameter in a subsequent request.
	PartNumberMarker *int32 `json:"PartNumberMarker,omitempty" xml:"PartNumberMarker"`
	// When a list is truncated, this element specifies the last part in the list, as well as the value to use for the part-number-marker request parameter in a subsequent request.
	NextPartNumberMarker *string `json:"NextPartNumberMarker,omitempty" xml:"NextPartNumberMarker"`
	// Maximum number of parts that were allowed in the response.
	MaxParts *string `json:"MaxParts,omitempty" xml:"MaxParts"`
	// A flag that indicates whether IONOS Object Storage returned all of the results that satisfied the search criteria. If your results were truncated, you can make a follow-up paginated request using the NextKeyMarker and NextVersionIdMarker response parameters as a starting place in another request to return the rest of the results.
	IsTruncated *bool `json:"IsTruncated,omitempty" xml:"IsTruncated"`
	//  Container for elements related to a particular part. A response can contain zero or more `Part` elements.
	Parts        []Part        `json:"Parts,omitempty" xml:"Parts"`
	Initiator    *Initiator    `json:"Initiator,omitempty" xml:"Initiator"`
	Owner        *Owner        `json:"Owner,omitempty" xml:"Owner"`
	StorageClass *StorageClass `json:"StorageClass,omitempty" xml:"StorageClass"`
}

// NewListPartsOutput instantiates a new ListPartsOutput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewListPartsOutput() *ListPartsOutput {
	this := ListPartsOutput{}

	return &this
}

// NewListPartsOutputWithDefaults instantiates a new ListPartsOutput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewListPartsOutputWithDefaults() *ListPartsOutput {
	this := ListPartsOutput{}
	return &this
}

// GetBucket returns the Bucket field value if set, zero value otherwise.
func (o *ListPartsOutput) GetBucket() string {
	if o == nil || IsNil(o.Bucket) {
		var ret string
		return ret
	}
	return *o.Bucket
}

// GetBucketOk returns a tuple with the Bucket field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListPartsOutput) GetBucketOk() (*string, bool) {
	if o == nil || IsNil(o.Bucket) {
		return nil, false
	}
	return o.Bucket, true
}

// HasBucket returns a boolean if a field has been set.
func (o *ListPartsOutput) HasBucket() bool {
	if o != nil && !IsNil(o.Bucket) {
		return true
	}

	return false
}

// SetBucket gets a reference to the given string and assigns it to the Bucket field.
func (o *ListPartsOutput) SetBucket(v string) {
	o.Bucket = &v
}

// GetKey returns the Key field value if set, zero value otherwise.
func (o *ListPartsOutput) GetKey() string {
	if o == nil || IsNil(o.Key) {
		var ret string
		return ret
	}
	return *o.Key
}

// GetKeyOk returns a tuple with the Key field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListPartsOutput) GetKeyOk() (*string, bool) {
	if o == nil || IsNil(o.Key) {
		return nil, false
	}
	return o.Key, true
}

// HasKey returns a boolean if a field has been set.
func (o *ListPartsOutput) HasKey() bool {
	if o != nil && !IsNil(o.Key) {
		return true
	}

	return false
}

// SetKey gets a reference to the given string and assigns it to the Key field.
func (o *ListPartsOutput) SetKey(v string) {
	o.Key = &v
}

// GetUploadId returns the UploadId field value if set, zero value otherwise.
func (o *ListPartsOutput) GetUploadId() string {
	if o == nil || IsNil(o.UploadId) {
		var ret string
		return ret
	}
	return *o.UploadId
}

// GetUploadIdOk returns a tuple with the UploadId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListPartsOutput) GetUploadIdOk() (*string, bool) {
	if o == nil || IsNil(o.UploadId) {
		return nil, false
	}
	return o.UploadId, true
}

// HasUploadId returns a boolean if a field has been set.
func (o *ListPartsOutput) HasUploadId() bool {
	if o != nil && !IsNil(o.UploadId) {
		return true
	}

	return false
}

// SetUploadId gets a reference to the given string and assigns it to the UploadId field.
func (o *ListPartsOutput) SetUploadId(v string) {
	o.UploadId = &v
}

// GetPartNumberMarker returns the PartNumberMarker field value if set, zero value otherwise.
func (o *ListPartsOutput) GetPartNumberMarker() int32 {
	if o == nil || IsNil(o.PartNumberMarker) {
		var ret int32
		return ret
	}
	return *o.PartNumberMarker
}

// GetPartNumberMarkerOk returns a tuple with the PartNumberMarker field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListPartsOutput) GetPartNumberMarkerOk() (*int32, bool) {
	if o == nil || IsNil(o.PartNumberMarker) {
		return nil, false
	}
	return o.PartNumberMarker, true
}

// HasPartNumberMarker returns a boolean if a field has been set.
func (o *ListPartsOutput) HasPartNumberMarker() bool {
	if o != nil && !IsNil(o.PartNumberMarker) {
		return true
	}

	return false
}

// SetPartNumberMarker gets a reference to the given int32 and assigns it to the PartNumberMarker field.
func (o *ListPartsOutput) SetPartNumberMarker(v int32) {
	o.PartNumberMarker = &v
}

// GetNextPartNumberMarker returns the NextPartNumberMarker field value if set, zero value otherwise.
func (o *ListPartsOutput) GetNextPartNumberMarker() string {
	if o == nil || IsNil(o.NextPartNumberMarker) {
		var ret string
		return ret
	}
	return *o.NextPartNumberMarker
}

// GetNextPartNumberMarkerOk returns a tuple with the NextPartNumberMarker field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListPartsOutput) GetNextPartNumberMarkerOk() (*string, bool) {
	if o == nil || IsNil(o.NextPartNumberMarker) {
		return nil, false
	}
	return o.NextPartNumberMarker, true
}

// HasNextPartNumberMarker returns a boolean if a field has been set.
func (o *ListPartsOutput) HasNextPartNumberMarker() bool {
	if o != nil && !IsNil(o.NextPartNumberMarker) {
		return true
	}

	return false
}

// SetNextPartNumberMarker gets a reference to the given string and assigns it to the NextPartNumberMarker field.
func (o *ListPartsOutput) SetNextPartNumberMarker(v string) {
	o.NextPartNumberMarker = &v
}

// GetMaxParts returns the MaxParts field value if set, zero value otherwise.
func (o *ListPartsOutput) GetMaxParts() string {
	if o == nil || IsNil(o.MaxParts) {
		var ret string
		return ret
	}
	return *o.MaxParts
}

// GetMaxPartsOk returns a tuple with the MaxParts field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListPartsOutput) GetMaxPartsOk() (*string, bool) {
	if o == nil || IsNil(o.MaxParts) {
		return nil, false
	}
	return o.MaxParts, true
}

// HasMaxParts returns a boolean if a field has been set.
func (o *ListPartsOutput) HasMaxParts() bool {
	if o != nil && !IsNil(o.MaxParts) {
		return true
	}

	return false
}

// SetMaxParts gets a reference to the given string and assigns it to the MaxParts field.
func (o *ListPartsOutput) SetMaxParts(v string) {
	o.MaxParts = &v
}

// GetIsTruncated returns the IsTruncated field value if set, zero value otherwise.
func (o *ListPartsOutput) GetIsTruncated() bool {
	if o == nil || IsNil(o.IsTruncated) {
		var ret bool
		return ret
	}
	return *o.IsTruncated
}

// GetIsTruncatedOk returns a tuple with the IsTruncated field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListPartsOutput) GetIsTruncatedOk() (*bool, bool) {
	if o == nil || IsNil(o.IsTruncated) {
		return nil, false
	}
	return o.IsTruncated, true
}

// HasIsTruncated returns a boolean if a field has been set.
func (o *ListPartsOutput) HasIsTruncated() bool {
	if o != nil && !IsNil(o.IsTruncated) {
		return true
	}

	return false
}

// SetIsTruncated gets a reference to the given bool and assigns it to the IsTruncated field.
func (o *ListPartsOutput) SetIsTruncated(v bool) {
	o.IsTruncated = &v
}

// GetParts returns the Parts field value if set, zero value otherwise.
func (o *ListPartsOutput) GetParts() []Part {
	if o == nil || IsNil(o.Parts) {
		var ret []Part
		return ret
	}
	return o.Parts
}

// GetPartsOk returns a tuple with the Parts field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListPartsOutput) GetPartsOk() ([]Part, bool) {
	if o == nil || IsNil(o.Parts) {
		return nil, false
	}
	return o.Parts, true
}

// HasParts returns a boolean if a field has been set.
func (o *ListPartsOutput) HasParts() bool {
	if o != nil && !IsNil(o.Parts) {
		return true
	}

	return false
}

// SetParts gets a reference to the given []Part and assigns it to the Parts field.
func (o *ListPartsOutput) SetParts(v []Part) {
	o.Parts = v
}

// GetInitiator returns the Initiator field value if set, zero value otherwise.
func (o *ListPartsOutput) GetInitiator() Initiator {
	if o == nil || IsNil(o.Initiator) {
		var ret Initiator
		return ret
	}
	return *o.Initiator
}

// GetInitiatorOk returns a tuple with the Initiator field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListPartsOutput) GetInitiatorOk() (*Initiator, bool) {
	if o == nil || IsNil(o.Initiator) {
		return nil, false
	}
	return o.Initiator, true
}

// HasInitiator returns a boolean if a field has been set.
func (o *ListPartsOutput) HasInitiator() bool {
	if o != nil && !IsNil(o.Initiator) {
		return true
	}

	return false
}

// SetInitiator gets a reference to the given Initiator and assigns it to the Initiator field.
func (o *ListPartsOutput) SetInitiator(v Initiator) {
	o.Initiator = &v
}

// GetOwner returns the Owner field value if set, zero value otherwise.
func (o *ListPartsOutput) GetOwner() Owner {
	if o == nil || IsNil(o.Owner) {
		var ret Owner
		return ret
	}
	return *o.Owner
}

// GetOwnerOk returns a tuple with the Owner field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListPartsOutput) GetOwnerOk() (*Owner, bool) {
	if o == nil || IsNil(o.Owner) {
		return nil, false
	}
	return o.Owner, true
}

// HasOwner returns a boolean if a field has been set.
func (o *ListPartsOutput) HasOwner() bool {
	if o != nil && !IsNil(o.Owner) {
		return true
	}

	return false
}

// SetOwner gets a reference to the given Owner and assigns it to the Owner field.
func (o *ListPartsOutput) SetOwner(v Owner) {
	o.Owner = &v
}

// GetStorageClass returns the StorageClass field value if set, zero value otherwise.
func (o *ListPartsOutput) GetStorageClass() StorageClass {
	if o == nil || IsNil(o.StorageClass) {
		var ret StorageClass
		return ret
	}
	return *o.StorageClass
}

// GetStorageClassOk returns a tuple with the StorageClass field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListPartsOutput) GetStorageClassOk() (*StorageClass, bool) {
	if o == nil || IsNil(o.StorageClass) {
		return nil, false
	}
	return o.StorageClass, true
}

// HasStorageClass returns a boolean if a field has been set.
func (o *ListPartsOutput) HasStorageClass() bool {
	if o != nil && !IsNil(o.StorageClass) {
		return true
	}

	return false
}

// SetStorageClass gets a reference to the given StorageClass and assigns it to the StorageClass field.
func (o *ListPartsOutput) SetStorageClass(v StorageClass) {
	o.StorageClass = &v
}

func (o ListPartsOutput) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ListPartsOutput) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Bucket) {
		toSerialize["Bucket"] = o.Bucket
	}
	if !IsNil(o.Key) {
		toSerialize["Key"] = o.Key
	}
	if !IsNil(o.UploadId) {
		toSerialize["UploadId"] = o.UploadId
	}
	if !IsNil(o.PartNumberMarker) {
		toSerialize["PartNumberMarker"] = o.PartNumberMarker
	}
	if !IsNil(o.NextPartNumberMarker) {
		toSerialize["NextPartNumberMarker"] = o.NextPartNumberMarker
	}
	if !IsNil(o.MaxParts) {
		toSerialize["MaxParts"] = o.MaxParts
	}
	if !IsNil(o.IsTruncated) {
		toSerialize["IsTruncated"] = o.IsTruncated
	}
	if !IsNil(o.Parts) {
		toSerialize["Parts"] = o.Parts
	}
	if !IsNil(o.Initiator) {
		toSerialize["Initiator"] = o.Initiator
	}
	if !IsNil(o.Owner) {
		toSerialize["Owner"] = o.Owner
	}
	if !IsNil(o.StorageClass) {
		toSerialize["StorageClass"] = o.StorageClass
	}
	return toSerialize, nil
}

type NullableListPartsOutput struct {
	value *ListPartsOutput
	isSet bool
}

func (v NullableListPartsOutput) Get() *ListPartsOutput {
	return v.value
}

func (v *NullableListPartsOutput) Set(val *ListPartsOutput) {
	v.value = val
	v.isSet = true
}

func (v NullableListPartsOutput) IsSet() bool {
	return v.isSet
}

func (v *NullableListPartsOutput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableListPartsOutput(val *ListPartsOutput) *NullableListPartsOutput {
	return &NullableListPartsOutput{value: val, isSet: true}
}

func (v NullableListPartsOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableListPartsOutput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
