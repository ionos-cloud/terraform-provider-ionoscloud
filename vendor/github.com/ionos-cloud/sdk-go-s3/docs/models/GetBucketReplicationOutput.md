# GetBucketReplicationOutput

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**ReplicationConfiguration** | Pointer to [**ReplicationConfiguration**](ReplicationConfiguration.md) |  | [optional] |

## Methods

### NewGetBucketReplicationOutput

`func NewGetBucketReplicationOutput() *GetBucketReplicationOutput`

NewGetBucketReplicationOutput instantiates a new GetBucketReplicationOutput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGetBucketReplicationOutputWithDefaults

`func NewGetBucketReplicationOutputWithDefaults() *GetBucketReplicationOutput`

NewGetBucketReplicationOutputWithDefaults instantiates a new GetBucketReplicationOutput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetReplicationConfiguration

`func (o *GetBucketReplicationOutput) GetReplicationConfiguration() ReplicationConfiguration`

GetReplicationConfiguration returns the ReplicationConfiguration field if non-nil, zero value otherwise.

### GetReplicationConfigurationOk

`func (o *GetBucketReplicationOutput) GetReplicationConfigurationOk() (*ReplicationConfiguration, bool)`

GetReplicationConfigurationOk returns a tuple with the ReplicationConfiguration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReplicationConfiguration

`func (o *GetBucketReplicationOutput) SetReplicationConfiguration(v ReplicationConfiguration)`

SetReplicationConfiguration sets ReplicationConfiguration field to given value.

### HasReplicationConfiguration

`func (o *GetBucketReplicationOutput) HasReplicationConfiguration() bool`

HasReplicationConfiguration returns a boolean if a field has been set.


