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

// checks if the GroupPostResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GroupPostResponse{}

// GroupPostResponse A group of virtual servers where the number of replicas can be scaled automatically.
type GroupPostResponse struct {
	// The unique resource identifier.
	Id string `json:"id"`
	// The resource type.
	Type *string `json:"type,omitempty"`
	// The absolute URL to the resource's representation.
	Href       *string            `json:"href,omitempty"`
	Metadata   *Metadata          `json:"metadata,omitempty"`
	Properties GroupProperties    `json:"properties"`
	Entities   *GroupPostEntities `json:"entities,omitempty"`
	// Any background activity caused by this request. You can use this to track the progress of such activities.
	StartedActions []ActionResource `json:"startedActions,omitempty"`
}

// NewGroupPostResponse instantiates a new GroupPostResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGroupPostResponse(id string, properties GroupProperties) *GroupPostResponse {
	this := GroupPostResponse{}

	this.Id = id
	this.Properties = properties

	return &this
}

// NewGroupPostResponseWithDefaults instantiates a new GroupPostResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGroupPostResponseWithDefaults() *GroupPostResponse {
	this := GroupPostResponse{}
	return &this
}

// GetId returns the Id field value
func (o *GroupPostResponse) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *GroupPostResponse) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *GroupPostResponse) SetId(v string) {
	o.Id = v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *GroupPostResponse) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupPostResponse) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *GroupPostResponse) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *GroupPostResponse) SetType(v string) {
	o.Type = &v
}

// GetHref returns the Href field value if set, zero value otherwise.
func (o *GroupPostResponse) GetHref() string {
	if o == nil || IsNil(o.Href) {
		var ret string
		return ret
	}
	return *o.Href
}

// GetHrefOk returns a tuple with the Href field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupPostResponse) GetHrefOk() (*string, bool) {
	if o == nil || IsNil(o.Href) {
		return nil, false
	}
	return o.Href, true
}

// HasHref returns a boolean if a field has been set.
func (o *GroupPostResponse) HasHref() bool {
	if o != nil && !IsNil(o.Href) {
		return true
	}

	return false
}

// SetHref gets a reference to the given string and assigns it to the Href field.
func (o *GroupPostResponse) SetHref(v string) {
	o.Href = &v
}

// GetMetadata returns the Metadata field value if set, zero value otherwise.
func (o *GroupPostResponse) GetMetadata() Metadata {
	if o == nil || IsNil(o.Metadata) {
		var ret Metadata
		return ret
	}
	return *o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupPostResponse) GetMetadataOk() (*Metadata, bool) {
	if o == nil || IsNil(o.Metadata) {
		return nil, false
	}
	return o.Metadata, true
}

// HasMetadata returns a boolean if a field has been set.
func (o *GroupPostResponse) HasMetadata() bool {
	if o != nil && !IsNil(o.Metadata) {
		return true
	}

	return false
}

// SetMetadata gets a reference to the given Metadata and assigns it to the Metadata field.
func (o *GroupPostResponse) SetMetadata(v Metadata) {
	o.Metadata = &v
}

// GetProperties returns the Properties field value
func (o *GroupPostResponse) GetProperties() GroupProperties {
	if o == nil {
		var ret GroupProperties
		return ret
	}

	return o.Properties
}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
func (o *GroupPostResponse) GetPropertiesOk() (*GroupProperties, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Properties, true
}

// SetProperties sets field value
func (o *GroupPostResponse) SetProperties(v GroupProperties) {
	o.Properties = v
}

// GetEntities returns the Entities field value if set, zero value otherwise.
func (o *GroupPostResponse) GetEntities() GroupPostEntities {
	if o == nil || IsNil(o.Entities) {
		var ret GroupPostEntities
		return ret
	}
	return *o.Entities
}

// GetEntitiesOk returns a tuple with the Entities field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupPostResponse) GetEntitiesOk() (*GroupPostEntities, bool) {
	if o == nil || IsNil(o.Entities) {
		return nil, false
	}
	return o.Entities, true
}

// HasEntities returns a boolean if a field has been set.
func (o *GroupPostResponse) HasEntities() bool {
	if o != nil && !IsNil(o.Entities) {
		return true
	}

	return false
}

// SetEntities gets a reference to the given GroupPostEntities and assigns it to the Entities field.
func (o *GroupPostResponse) SetEntities(v GroupPostEntities) {
	o.Entities = &v
}

// GetStartedActions returns the StartedActions field value if set, zero value otherwise.
func (o *GroupPostResponse) GetStartedActions() []ActionResource {
	if o == nil || IsNil(o.StartedActions) {
		var ret []ActionResource
		return ret
	}
	return o.StartedActions
}

// GetStartedActionsOk returns a tuple with the StartedActions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupPostResponse) GetStartedActionsOk() ([]ActionResource, bool) {
	if o == nil || IsNil(o.StartedActions) {
		return nil, false
	}
	return o.StartedActions, true
}

// HasStartedActions returns a boolean if a field has been set.
func (o *GroupPostResponse) HasStartedActions() bool {
	if o != nil && !IsNil(o.StartedActions) {
		return true
	}

	return false
}

// SetStartedActions gets a reference to the given []ActionResource and assigns it to the StartedActions field.
func (o *GroupPostResponse) SetStartedActions(v []ActionResource) {
	o.StartedActions = v
}

func (o GroupPostResponse) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GroupPostResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	if !IsNil(o.Type) {
		toSerialize["type"] = o.Type
	}
	if !IsNil(o.Href) {
		toSerialize["href"] = o.Href
	}
	if !IsNil(o.Metadata) {
		toSerialize["metadata"] = o.Metadata
	}
	toSerialize["properties"] = o.Properties
	if !IsNil(o.Entities) {
		toSerialize["entities"] = o.Entities
	}
	if !IsNil(o.StartedActions) {
		toSerialize["startedActions"] = o.StartedActions
	}
	return toSerialize, nil
}

type NullableGroupPostResponse struct {
	value *GroupPostResponse
	isSet bool
}

func (v NullableGroupPostResponse) Get() *GroupPostResponse {
	return v.value
}

func (v *NullableGroupPostResponse) Set(val *GroupPostResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableGroupPostResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableGroupPostResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGroupPostResponse(val *GroupPostResponse) *NullableGroupPostResponse {
	return &NullableGroupPostResponse{value: val, isSet: true}
}

func (v NullableGroupPostResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGroupPostResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
