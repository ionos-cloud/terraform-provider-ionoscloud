# BucketLocation

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**LocationConstraint** | Pointer to **string** | Specifies the Region where the bucket resides. | [optional] |

## Methods

### NewBucketLocation

`func NewBucketLocation() *BucketLocation`

NewBucketLocation instantiates a new BucketLocation object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBucketLocationWithDefaults

`func NewBucketLocationWithDefaults() *BucketLocation`

NewBucketLocationWithDefaults instantiates a new BucketLocation object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLocationConstraint

`func (o *BucketLocation) GetLocationConstraint() string`

GetLocationConstraint returns the LocationConstraint field if non-nil, zero value otherwise.

### GetLocationConstraintOk

`func (o *BucketLocation) GetLocationConstraintOk() (*string, bool)`

GetLocationConstraintOk returns a tuple with the LocationConstraint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocationConstraint

`func (o *BucketLocation) SetLocationConstraint(v string)`

SetLocationConstraint sets LocationConstraint field to given value.

### HasLocationConstraint

`func (o *BucketLocation) HasLocationConstraint() bool`

HasLocationConstraint returns a boolean if a field has been set.


