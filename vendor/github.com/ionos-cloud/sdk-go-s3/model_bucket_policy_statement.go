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

// checks if the BucketPolicyStatement type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &BucketPolicyStatement{}

// BucketPolicyStatement struct for BucketPolicyStatement
type BucketPolicyStatement struct {
	XMLName xml.Name `xml:"BucketPolicyStatement"`
	// Custom string identifying the statement.
	Sid *string `json:"Sid,omitempty" xml:"Sid"`
	// The array of allowed or denied actions.   IONOS S3 Object Storage supports the use of a wildcard in your Action configuration (`\"Action\":[\"s3:*\"]`). When an Action wildcard is used together with an object-level Resource element (`\"arn:aws:s3:::<bucketName>/_*\"` or `\"arn:aws:s3:::<bucketName>/<objectName>\"`), the wildcard denotes all supported Object actions. When an Action wildcard is used together with bucket-level Resource element (`\"arn:aws:s3:::<bucketName>\"`), the wildcard denotes all the bucket actions and bucket subresource actions that IONOS S3 Object Storage supports.
	Action []string `json:"Action" xml:"Action"`
	// Specify the outcome when the user requests a particular action.
	Effect string `json:"Effect" xml:"Effect"`
	// The bucket or object that the policy applies to.   Must be one of the following: - `\"arn:aws:s3:::<bucketName>\"` - For bucket actions (such as `s3:ListBucket`) and bucket subresource actions (such as `s3:GetBucketAcl`). - `\"arn:aws:s3:::<bucketName>/_*\"` or `\"arn:aws:s3:::<bucketName>/<objectName>\"` - For object actions (such as `s3:PutObject`).
	Resource  []string                        `json:"Resource" xml:"Resource"`
	Condition *BucketPolicyStatementCondition `json:"Condition,omitempty" xml:"Condition"`
	Principal *BucketPolicyStatementPrincipal `json:"Principal,omitempty" xml:"Principal"`
}

// NewBucketPolicyStatement instantiates a new BucketPolicyStatement object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBucketPolicyStatement(action []string, effect string, resource []string) *BucketPolicyStatement {
	this := BucketPolicyStatement{}

	this.Action = action
	this.Effect = effect
	this.Resource = resource

	return &this
}

// NewBucketPolicyStatementWithDefaults instantiates a new BucketPolicyStatement object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBucketPolicyStatementWithDefaults() *BucketPolicyStatement {
	this := BucketPolicyStatement{}
	return &this
}

// GetSid returns the Sid field value if set, zero value otherwise.
func (o *BucketPolicyStatement) GetSid() string {
	if o == nil || IsNil(o.Sid) {
		var ret string
		return ret
	}
	return *o.Sid
}

// GetSidOk returns a tuple with the Sid field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BucketPolicyStatement) GetSidOk() (*string, bool) {
	if o == nil || IsNil(o.Sid) {
		return nil, false
	}
	return o.Sid, true
}

// HasSid returns a boolean if a field has been set.
func (o *BucketPolicyStatement) HasSid() bool {
	if o != nil && !IsNil(o.Sid) {
		return true
	}

	return false
}

// SetSid gets a reference to the given string and assigns it to the Sid field.
func (o *BucketPolicyStatement) SetSid(v string) {
	o.Sid = &v
}

// GetAction returns the Action field value
func (o *BucketPolicyStatement) GetAction() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.Action
}

// GetActionOk returns a tuple with the Action field value
// and a boolean to check if the value has been set.
func (o *BucketPolicyStatement) GetActionOk() ([]string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Action, true
}

// SetAction sets field value
func (o *BucketPolicyStatement) SetAction(v []string) {
	o.Action = v
}

// GetEffect returns the Effect field value
func (o *BucketPolicyStatement) GetEffect() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Effect
}

// GetEffectOk returns a tuple with the Effect field value
// and a boolean to check if the value has been set.
func (o *BucketPolicyStatement) GetEffectOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Effect, true
}

// SetEffect sets field value
func (o *BucketPolicyStatement) SetEffect(v string) {
	o.Effect = v
}

// GetResource returns the Resource field value
func (o *BucketPolicyStatement) GetResource() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.Resource
}

// GetResourceOk returns a tuple with the Resource field value
// and a boolean to check if the value has been set.
func (o *BucketPolicyStatement) GetResourceOk() ([]string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Resource, true
}

// SetResource sets field value
func (o *BucketPolicyStatement) SetResource(v []string) {
	o.Resource = v
}

// GetCondition returns the Condition field value if set, zero value otherwise.
func (o *BucketPolicyStatement) GetCondition() BucketPolicyStatementCondition {
	if o == nil || IsNil(o.Condition) {
		var ret BucketPolicyStatementCondition
		return ret
	}
	return *o.Condition
}

// GetConditionOk returns a tuple with the Condition field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BucketPolicyStatement) GetConditionOk() (*BucketPolicyStatementCondition, bool) {
	if o == nil || IsNil(o.Condition) {
		return nil, false
	}
	return o.Condition, true
}

// HasCondition returns a boolean if a field has been set.
func (o *BucketPolicyStatement) HasCondition() bool {
	if o != nil && !IsNil(o.Condition) {
		return true
	}

	return false
}

// SetCondition gets a reference to the given BucketPolicyStatementCondition and assigns it to the Condition field.
func (o *BucketPolicyStatement) SetCondition(v BucketPolicyStatementCondition) {
	o.Condition = &v
}

// GetPrincipal returns the Principal field value if set, zero value otherwise.
func (o *BucketPolicyStatement) GetPrincipal() BucketPolicyStatementPrincipal {
	if o == nil || IsNil(o.Principal) {
		var ret BucketPolicyStatementPrincipal
		return ret
	}
	return *o.Principal
}

// GetPrincipalOk returns a tuple with the Principal field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BucketPolicyStatement) GetPrincipalOk() (*BucketPolicyStatementPrincipal, bool) {
	if o == nil || IsNil(o.Principal) {
		return nil, false
	}
	return o.Principal, true
}

// HasPrincipal returns a boolean if a field has been set.
func (o *BucketPolicyStatement) HasPrincipal() bool {
	if o != nil && !IsNil(o.Principal) {
		return true
	}

	return false
}

// SetPrincipal gets a reference to the given BucketPolicyStatementPrincipal and assigns it to the Principal field.
func (o *BucketPolicyStatement) SetPrincipal(v BucketPolicyStatementPrincipal) {
	o.Principal = &v
}

func (o BucketPolicyStatement) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o BucketPolicyStatement) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Sid) {
		toSerialize["Sid"] = o.Sid
	}
	if !IsZero(o.Action) {
		toSerialize["Action"] = o.Action
	}
	if !IsZero(o.Effect) {
		toSerialize["Effect"] = o.Effect
	}
	if !IsZero(o.Resource) {
		toSerialize["Resource"] = o.Resource
	}
	if !IsNil(o.Condition) {
		toSerialize["Condition"] = o.Condition
	}
	if !IsNil(o.Principal) {
		toSerialize["Principal"] = o.Principal
	}
	return toSerialize, nil
}

type NullableBucketPolicyStatement struct {
	value *BucketPolicyStatement
	isSet bool
}

func (v NullableBucketPolicyStatement) Get() *BucketPolicyStatement {
	return v.value
}

func (v *NullableBucketPolicyStatement) Set(val *BucketPolicyStatement) {
	v.value = val
	v.isSet = true
}

func (v NullableBucketPolicyStatement) IsSet() bool {
	return v.isSet
}

func (v *NullableBucketPolicyStatement) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBucketPolicyStatement(val *BucketPolicyStatement) *NullableBucketPolicyStatement {
	return &NullableBucketPolicyStatement{value: val, isSet: true}
}

func (v NullableBucketPolicyStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBucketPolicyStatement) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
