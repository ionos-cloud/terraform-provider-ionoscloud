# ListAllMyBucketsResult

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Owner** | Pointer to [**Owner**](Owner.md) |  | [optional] |
|**Buckets** | Pointer to [**[]Bucket**](Bucket.md) |  | [optional] |

## Methods

### NewListAllMyBucketsResult

`func NewListAllMyBucketsResult() *ListAllMyBucketsResult`

NewListAllMyBucketsResult instantiates a new ListAllMyBucketsResult object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewListAllMyBucketsResultWithDefaults

`func NewListAllMyBucketsResultWithDefaults() *ListAllMyBucketsResult`

NewListAllMyBucketsResultWithDefaults instantiates a new ListAllMyBucketsResult object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetOwner

`func (o *ListAllMyBucketsResult) GetOwner() Owner`

GetOwner returns the Owner field if non-nil, zero value otherwise.

### GetOwnerOk

`func (o *ListAllMyBucketsResult) GetOwnerOk() (*Owner, bool)`

GetOwnerOk returns a tuple with the Owner field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOwner

`func (o *ListAllMyBucketsResult) SetOwner(v Owner)`

SetOwner sets Owner field to given value.

### HasOwner

`func (o *ListAllMyBucketsResult) HasOwner() bool`

HasOwner returns a boolean if a field has been set.

### GetBuckets

`func (o *ListAllMyBucketsResult) GetBuckets() []Bucket`

GetBuckets returns the Buckets field if non-nil, zero value otherwise.

### GetBucketsOk

`func (o *ListAllMyBucketsResult) GetBucketsOk() (*[]Bucket, bool)`

GetBucketsOk returns a tuple with the Buckets field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBuckets

`func (o *ListAllMyBucketsResult) SetBuckets(v []Bucket)`

SetBuckets sets Buckets field to given value.

### HasBuckets

`func (o *ListAllMyBucketsResult) HasBuckets() bool`

HasBuckets returns a boolean if a field has been set.


