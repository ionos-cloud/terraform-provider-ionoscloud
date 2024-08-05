# AbortIncompleteMultipartUpload

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**DaysAfterInitiation** | Pointer to **int32** | Specifies the number of days after which IONOS S3 Object Storage aborts an incomplete multipart upload. | [optional] |

## Methods

### NewAbortIncompleteMultipartUpload

`func NewAbortIncompleteMultipartUpload() *AbortIncompleteMultipartUpload`

NewAbortIncompleteMultipartUpload instantiates a new AbortIncompleteMultipartUpload object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAbortIncompleteMultipartUploadWithDefaults

`func NewAbortIncompleteMultipartUploadWithDefaults() *AbortIncompleteMultipartUpload`

NewAbortIncompleteMultipartUploadWithDefaults instantiates a new AbortIncompleteMultipartUpload object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDaysAfterInitiation

`func (o *AbortIncompleteMultipartUpload) GetDaysAfterInitiation() int32`

GetDaysAfterInitiation returns the DaysAfterInitiation field if non-nil, zero value otherwise.

### GetDaysAfterInitiationOk

`func (o *AbortIncompleteMultipartUpload) GetDaysAfterInitiationOk() (*int32, bool)`

GetDaysAfterInitiationOk returns a tuple with the DaysAfterInitiation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDaysAfterInitiation

`func (o *AbortIncompleteMultipartUpload) SetDaysAfterInitiation(v int32)`

SetDaysAfterInitiation sets DaysAfterInitiation field to given value.

### HasDaysAfterInitiation

`func (o *AbortIncompleteMultipartUpload) HasDaysAfterInitiation() bool`

HasDaysAfterInitiation returns a boolean if a field has been set.


