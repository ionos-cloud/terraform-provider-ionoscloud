/*
 * VM Auto Scaling API
 *
 * The VM Auto Scaling Service enables IONOS clients to horizontally scale the number of VM replicas based on configured rules. You can use VM Auto Scaling to ensure that you have a sufficient number of replicas to handle your application loads at all times.  For this purpose, create a VM Auto Scaling Group that contains the server replicas. The VM Auto Scaling Service ensures that the number of replicas in the group is always within the defined limits.   When scaling policies are set, VM Auto Scaling creates or deletes replicas according to the requirements of your applications. For each policy, specified 'scale-in' and 'scale-out' actions are performed when the corresponding thresholds are reached.
 *
 * API version: 1.0.0
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package vmautoscaling

import (
	"encoding/json"
)

// checks if the GroupPolicy type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GroupPolicy{}

// GroupPolicy Defines the behavior of this VM Auto Scaling Group. A policy consists of triggers and actions, where an action is an automated behavior, and the trigger defines the circumstances under which the action is triggered. Currently, two separate actions are supported, namely scaling inward and outward, triggered by the thresholds defined for a particular metric.
type GroupPolicy struct {
	Metric Metric `json:"metric"`
	// Specifies the time range for which the samples are to be aggregated. Must be >= 2 minutes.
	Range         *string                  `json:"range,omitempty"`
	ScaleInAction GroupPolicyScaleInAction `json:"scaleInAction"`
	// The lower threshold for the value of the 'metric'. Used with the `less than` (<) operator. When this value is exceeded, a scale-in action is triggered, specified by the 'scaleInAction' property. The value must have a higher minimum delta to the 'scaleOutThreshold', depending on the 'metric', to avoid competing for actions at the same time.
	ScaleInThreshold float32                   `json:"scaleInThreshold"`
	ScaleOutAction   GroupPolicyScaleOutAction `json:"scaleOutAction"`
	// The upper threshold for the value of the 'metric'. Used with the 'greater than' (>) operator. A scale-out action is triggered when this value is exceeded, specified by the 'scaleOutAction' property. The value must have a lower minimum delta to the 'scaleInThreshold', depending on the metric, to avoid competing for actions simultaneously. If 'properties.policy.unit=TOTAL', a value >= 40 must be chosen.
	ScaleOutThreshold float32   `json:"scaleOutThreshold"`
	Unit              QueryUnit `json:"unit"`
}

// NewGroupPolicy instantiates a new GroupPolicy object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGroupPolicy(metric Metric, scaleInAction GroupPolicyScaleInAction, scaleInThreshold float32, scaleOutAction GroupPolicyScaleOutAction, scaleOutThreshold float32, unit QueryUnit) *GroupPolicy {
	this := GroupPolicy{}

	this.Metric = metric
	var range_ string = "120s"
	this.Range = &range_
	this.ScaleInAction = scaleInAction
	this.ScaleInThreshold = scaleInThreshold
	this.ScaleOutAction = scaleOutAction
	this.ScaleOutThreshold = scaleOutThreshold
	this.Unit = unit

	return &this
}

// NewGroupPolicyWithDefaults instantiates a new GroupPolicy object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGroupPolicyWithDefaults() *GroupPolicy {
	this := GroupPolicy{}
	var range_ string = "120s"
	this.Range = &range_
	var unit QueryUnit = QUERYUNIT_TOTAL
	this.Unit = unit
	return &this
}

// GetMetric returns the Metric field value
func (o *GroupPolicy) GetMetric() Metric {
	if o == nil {
		var ret Metric
		return ret
	}

	return o.Metric
}

// GetMetricOk returns a tuple with the Metric field value
// and a boolean to check if the value has been set.
func (o *GroupPolicy) GetMetricOk() (*Metric, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Metric, true
}

// SetMetric sets field value
func (o *GroupPolicy) SetMetric(v Metric) {
	o.Metric = v
}

// GetRange returns the Range field value if set, zero value otherwise.
func (o *GroupPolicy) GetRange() string {
	if o == nil || IsNil(o.Range) {
		var ret string
		return ret
	}
	return *o.Range
}

// GetRangeOk returns a tuple with the Range field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupPolicy) GetRangeOk() (*string, bool) {
	if o == nil || IsNil(o.Range) {
		return nil, false
	}
	return o.Range, true
}

// HasRange returns a boolean if a field has been set.
func (o *GroupPolicy) HasRange() bool {
	if o != nil && !IsNil(o.Range) {
		return true
	}

	return false
}

// SetRange gets a reference to the given string and assigns it to the Range field.
func (o *GroupPolicy) SetRange(v string) {
	o.Range = &v
}

// GetScaleInAction returns the ScaleInAction field value
func (o *GroupPolicy) GetScaleInAction() GroupPolicyScaleInAction {
	if o == nil {
		var ret GroupPolicyScaleInAction
		return ret
	}

	return o.ScaleInAction
}

// GetScaleInActionOk returns a tuple with the ScaleInAction field value
// and a boolean to check if the value has been set.
func (o *GroupPolicy) GetScaleInActionOk() (*GroupPolicyScaleInAction, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ScaleInAction, true
}

// SetScaleInAction sets field value
func (o *GroupPolicy) SetScaleInAction(v GroupPolicyScaleInAction) {
	o.ScaleInAction = v
}

// GetScaleInThreshold returns the ScaleInThreshold field value
func (o *GroupPolicy) GetScaleInThreshold() float32 {
	if o == nil {
		var ret float32
		return ret
	}

	return o.ScaleInThreshold
}

// GetScaleInThresholdOk returns a tuple with the ScaleInThreshold field value
// and a boolean to check if the value has been set.
func (o *GroupPolicy) GetScaleInThresholdOk() (*float32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ScaleInThreshold, true
}

// SetScaleInThreshold sets field value
func (o *GroupPolicy) SetScaleInThreshold(v float32) {
	o.ScaleInThreshold = v
}

// GetScaleOutAction returns the ScaleOutAction field value
func (o *GroupPolicy) GetScaleOutAction() GroupPolicyScaleOutAction {
	if o == nil {
		var ret GroupPolicyScaleOutAction
		return ret
	}

	return o.ScaleOutAction
}

// GetScaleOutActionOk returns a tuple with the ScaleOutAction field value
// and a boolean to check if the value has been set.
func (o *GroupPolicy) GetScaleOutActionOk() (*GroupPolicyScaleOutAction, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ScaleOutAction, true
}

// SetScaleOutAction sets field value
func (o *GroupPolicy) SetScaleOutAction(v GroupPolicyScaleOutAction) {
	o.ScaleOutAction = v
}

// GetScaleOutThreshold returns the ScaleOutThreshold field value
func (o *GroupPolicy) GetScaleOutThreshold() float32 {
	if o == nil {
		var ret float32
		return ret
	}

	return o.ScaleOutThreshold
}

// GetScaleOutThresholdOk returns a tuple with the ScaleOutThreshold field value
// and a boolean to check if the value has been set.
func (o *GroupPolicy) GetScaleOutThresholdOk() (*float32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ScaleOutThreshold, true
}

// SetScaleOutThreshold sets field value
func (o *GroupPolicy) SetScaleOutThreshold(v float32) {
	o.ScaleOutThreshold = v
}

// GetUnit returns the Unit field value
func (o *GroupPolicy) GetUnit() QueryUnit {
	if o == nil {
		var ret QueryUnit
		return ret
	}

	return o.Unit
}

// GetUnitOk returns a tuple with the Unit field value
// and a boolean to check if the value has been set.
func (o *GroupPolicy) GetUnitOk() (*QueryUnit, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Unit, true
}

// SetUnit sets field value
func (o *GroupPolicy) SetUnit(v QueryUnit) {
	o.Unit = v
}

func (o GroupPolicy) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GroupPolicy) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["metric"] = o.Metric
	if !IsNil(o.Range) {
		toSerialize["range"] = o.Range
	}
	toSerialize["scaleInAction"] = o.ScaleInAction
	toSerialize["scaleInThreshold"] = o.ScaleInThreshold
	toSerialize["scaleOutAction"] = o.ScaleOutAction
	toSerialize["scaleOutThreshold"] = o.ScaleOutThreshold
	toSerialize["unit"] = o.Unit
	return toSerialize, nil
}

type NullableGroupPolicy struct {
	value *GroupPolicy
	isSet bool
}

func (v NullableGroupPolicy) Get() *GroupPolicy {
	return v.value
}

func (v *NullableGroupPolicy) Set(val *GroupPolicy) {
	v.value = val
	v.isSet = true
}

func (v NullableGroupPolicy) IsSet() bool {
	return v.isSet
}

func (v *NullableGroupPolicy) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGroupPolicy(val *GroupPolicy) *NullableGroupPolicy {
	return &NullableGroupPolicy{value: val, isSet: true}
}

func (v NullableGroupPolicy) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGroupPolicy) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
