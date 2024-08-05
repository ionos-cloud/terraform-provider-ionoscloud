# Destination

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Bucket** | **string** | Use the same \&quot;Bucket\&quot; value formatting as in the S3 API specification, that is, &#x60;arn:aws:s3:::{Bucket}&#x60;.  | |
|**StorageClass** | Pointer to [**StorageClass**](StorageClass.md) |  | [optional] |

## Methods

### NewDestination

`func NewDestination(bucket string, ) *Destination`

NewDestination instantiates a new Destination object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDestinationWithDefaults

`func NewDestinationWithDefaults() *Destination`

NewDestinationWithDefaults instantiates a new Destination object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBucket

`func (o *Destination) GetBucket() string`

GetBucket returns the Bucket field if non-nil, zero value otherwise.

### GetBucketOk

`func (o *Destination) GetBucketOk() (*string, bool)`

GetBucketOk returns a tuple with the Bucket field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBucket

`func (o *Destination) SetBucket(v string)`

SetBucket sets Bucket field to given value.


### GetStorageClass

`func (o *Destination) GetStorageClass() StorageClass`

GetStorageClass returns the StorageClass field if non-nil, zero value otherwise.

### GetStorageClassOk

`func (o *Destination) GetStorageClassOk() (*StorageClass, bool)`

GetStorageClassOk returns a tuple with the StorageClass field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStorageClass

`func (o *Destination) SetStorageClass(v StorageClass)`

SetStorageClass sets StorageClass field to given value.

### HasStorageClass

`func (o *Destination) HasStorageClass() bool`

HasStorageClass returns a boolean if a field has been set.


