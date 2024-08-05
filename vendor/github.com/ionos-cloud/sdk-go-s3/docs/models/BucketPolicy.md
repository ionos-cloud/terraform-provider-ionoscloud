# BucketPolicy

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Id** | Pointer to **string** | Specifies an optional identifier for the policy. | [optional] |
|**Version** | Pointer to **string** | Policy version | [optional] |
|**Statement** | [**[]BucketPolicyStatement**](BucketPolicyStatement.md) |  | |

## Methods

### NewBucketPolicy

`func NewBucketPolicy(statement []BucketPolicyStatement, ) *BucketPolicy`

NewBucketPolicy instantiates a new BucketPolicy object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBucketPolicyWithDefaults

`func NewBucketPolicyWithDefaults() *BucketPolicy`

NewBucketPolicyWithDefaults instantiates a new BucketPolicy object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *BucketPolicy) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *BucketPolicy) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *BucketPolicy) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *BucketPolicy) HasId() bool`

HasId returns a boolean if a field has been set.

### GetVersion

`func (o *BucketPolicy) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *BucketPolicy) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *BucketPolicy) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *BucketPolicy) HasVersion() bool`

HasVersion returns a boolean if a field has been set.

### GetStatement

`func (o *BucketPolicy) GetStatement() []BucketPolicyStatement`

GetStatement returns the Statement field if non-nil, zero value otherwise.

### GetStatementOk

`func (o *BucketPolicy) GetStatementOk() (*[]BucketPolicyStatement, bool)`

GetStatementOk returns a tuple with the Statement field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatement

`func (o *BucketPolicy) SetStatement(v []BucketPolicyStatement)`

SetStatement sets Statement field to given value.



