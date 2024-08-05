# ReplicationConfiguration

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Role** | **string** | The Resource Name of the Identity and Access Management (IAM) role that IONOS S3 Object Storage assumes when replicating objects. | |
|**Rules** | [**[]ReplicationRule**](ReplicationRule.md) |  | |

## Methods

### NewReplicationConfiguration

`func NewReplicationConfiguration(role string, rules []ReplicationRule, ) *ReplicationConfiguration`

NewReplicationConfiguration instantiates a new ReplicationConfiguration object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewReplicationConfigurationWithDefaults

`func NewReplicationConfigurationWithDefaults() *ReplicationConfiguration`

NewReplicationConfigurationWithDefaults instantiates a new ReplicationConfiguration object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRole

`func (o *ReplicationConfiguration) GetRole() string`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *ReplicationConfiguration) GetRoleOk() (*string, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *ReplicationConfiguration) SetRole(v string)`

SetRole sets Role field to given value.


### GetRules

`func (o *ReplicationConfiguration) GetRules() []ReplicationRule`

GetRules returns the Rules field if non-nil, zero value otherwise.

### GetRulesOk

`func (o *ReplicationConfiguration) GetRulesOk() (*[]ReplicationRule, bool)`

GetRulesOk returns a tuple with the Rules field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRules

`func (o *ReplicationConfiguration) SetRules(v []ReplicationRule)`

SetRules sets Rules field to given value.



