# Error

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Code** | Pointer to **string** | The error code is a string that uniquely identifies an error condition. It is meant to be read and understood by programs that detect and handle errors by type.  ## IONOS S3 Object Storage error codes - AccessDenied   - Description: Access Denied   - HTTPStatus Code: 403 Forbidden - AccountProblem   - Description: There is a problem with your IONOS S3 Object Storage account that prevents the operation from completing successfully. Contact IONOS for further assistance.   - HTTP Status Code: 403 Forbidden - AmbiguousGrantByEmailAddress   - Description: The email address you provided is associated with more than one account.   - HTTP Status Code: 400 Bad Request - BadDigest   - Description: The Content-MD5 you specified did not match what we received.   - HTTP Status Code: 400 Bad Request - BucketAlreadyExists   - Description: The requested bucket name is not available. The bucket namespace is shared by all users of the system. Please select a different name and try again.   - HTTP Status Code: 409 Conflict - BucketAlreadyOwnedByYou   - Description: The bucket you tried to create already exists, and you own it.   - HTTP Code: 409 Conflict - BucketNotEmpty   - Description: The bucket you tried to delete is not empty.   - HTTP Status Code: 409 Conflict - CrossLocationLoggingProhibited   - Description: Cross-location logging not allowed. Buckets in one geographic location cannot log information to a bucket in another location.   - HTTP Status Code: 403 Forbidden - EntityTooSmall   - Description: Your proposed upload is smaller than the minimum allowed object size.   - HTTP Status Code: 400 Bad Request - EntityTooLarge   - Description: Your proposed upload exceeds the maximum allowed object size.   - HTTP Status Code: 400 Bad Request - IllegalVersioningConfigurationException   - Description: Indicates that the versioning configuration specified in the request is invalid.   - HTTP Status Code: 400 Bad Request - IncorrectNumberOfFilesInPostRequest   - Description: POST requires exactly one file upload per request.   - HTTP Status Code: 400 Bad Request - InternalError   - Description: We encountered an internal error. Please try again.   - HTTP Status Code: 500 Internal Server Error  - InvalidAccessKeyId   - Description: The IONOS S3 Object Storage access key ID you provided does not exist in our records.   - HTTP Status Code: 403 Forbidden - InvalidArgument   - Description: Invalid Argument   - HTTP Status Code: 400 Bad Request - InvalidBucketName   - Description: The specified bucket is not valid.   - HTTP Status Code: 400 Bad Request - InvalidBucketState   - Description: The request is not valid with the current state of the bucket.   - HTTP Status Code: 409 Conflict - InvalidDigest   - Description: The Content-MD5 you specified is not valid.   - HTTP Status Code: 400 Bad Request - InvalidEncryptionAlgorithmError   - Description: The encryption request you specified is not valid. The valid value is AES256.   - HTTP Status Code: 400 Bad Request - InvalidLocationConstraint   - HTTP Status Code: 400 Bad Request - InvalidObjectState   - Description: The operation is not valid for the current state of the object.   - HTTP Status Code: 403 Forbidden - InvalidPart   - Description: One or more of the specified parts could not be found. The part might not have been uploaded, or the specified entity tag might not have matched the part&#39;s entity tag.   - HTTP Status Code: 400 Bad Request - InvalidPartOrder   - Description: The list of parts was not in ascending order. Parts list must be specified in order by part number.   - HTTP Status Code: 400 Bad Request - InvalidPolicyDocument   - Description: The content of the form does not meet the conditions specified in the policy document.   - HTTP Status Code: 400 Bad Request - InvalidRange   - Description: The requested range cannot be satisfied.   - HTTP Status Code: 416 Requested Range Not Satisfiable - InvalidRequest   - Description: Please use &#x60;AWS4-HMAC-SHA256&#x60;.   - HTTP Status Code: 400 Bad Request - InvalidSecurity   - Description: The provided security credentials are not valid.   - HTTP Status Code: 403 Forbidden - InvalidTargetBucketForLogging   - Description: The target bucket for logging does not exist, is not owned by you, or does not have the appropriate grants for the log-delivery group.   - Status Code: 400 Bad Request - InvalidURI   - Description: Couldn&#39;t parse the specified URI.   - HTTP Status Code: 400 Bad Request - KeyTooLong   - Description: Your key is too long.   - HTTP Status Code: 400 Bad Request - MalformedACLError   - Description: The XML you provided was not well-formed or did not validate against our published schema.   - HTTP Status Code: 400 Bad Request - MalformedPOSTRequest   - Description: The body of your POST request is not well-formed multipart/form-data.   - HTTP Status Code: 400 Bad Request - MalformedXML   - Description: This happens when the user sends malformed XML (XML that doesn&#39;t conform to the published XSD) for the configuration. The error message is, \&quot;The XML you provided was not well-formed or did not validate against our published schema.\&quot;   - HTTP Status Code: 400 Bad Request - MaxMessageLengthExceeded   - Description: Your request was too big.   - HTTP Status Code: 400 Bad Request - MaxPostPreDataLengthExceededError   - Description: Your POST request fields preceding the upload file were too large.   - HTTP Status Code: 400 Bad Request - MetadataTooLarge   - Description: Your metadata headers exceed the maximum allowed metadata size.   - HTTP Status Code: 400 Bad Request - MethodNotAllowed   - Description: The specified method is not allowed against this resource.   - HTTP Status Code: 405 Method Not Allowed - MissingContentLength   - Description: You must provide the Content-Length HTTP header.   - HTTP Status Code: 411 Length Required - MissingSecurityHeader   - Description: Your request is missing a required header.   - HTTP Status Code: 400 Bad Request - NoSuchBucket   - Description: The specified bucket does not exist.   - HTTP Status Code: 404 Not Found - NoSuchBucketPolicy   - Description: The specified bucket does not have a bucket policy.   - HTTP Status Code: 404 Not Found - NoSuchKey   - Description: The specified key does not exist.   - HTTP Status Code:404 Not Found - NoSuchLifecycleConfiguration   - Description: The lifecycle configuration does not exist.   - HTTP Status Code: 404 Not Found - NoSuchReplicationConfiguration   - Description: The replication configuration does not exist.   - HTTP Status Code: 404 Not Found - NoSuchUpload   - Description: The specified multipart upload does not exist. The upload ID might be invalid, or the multipart upload might have been aborted or completed.   - HTTP Status Code: 404 Not Found - NoSuchVersion   - Description: Indicates that the version ID specified in the request does not match an existing version.   - HTTP Status Code: 404 Not Found - NotImplemented   - Description: A header you provided implies functionality that is not implemented.   - HTTP Status Code: 501 Not Implemented - PermanentRedirect   - Description: The bucket you are attempting to access must be addressed using the specified endpoint. Send all future requests to this endpoint.   - HTTP Status Code: 301 Moved Permanently - PreconditionFailed   - Description: At least one of the preconditions you specified did not hold.   - HTTP Status Code: 412 Precondition Failed - Redirect   - Description: Temporary redirect.   - HTTP Status Code: 307 Moved Temporarily - RestoreAlreadyInProgress   - Description: Object restore is already in progress.   - HTTP Status Code: 409 Conflict - RequestIsNotMultiPartContent   - Description: Bucket POST must be of the enclosure-type multipart/form-data.   - HTTP Status Code: 400 Bad Request - RequestTimeout   - Description: Your socket connection to the server was not read from or written to within the timeout period.   - HTTP Status Code: 400 Bad Request - RequestTimeTooSkewed   - Description: The difference between the request time and the server&#39;s time is too large.   - HTTP Status Code: 403 Forbidden - SignatureDoesNotMatch   - HTTP Status Code: 403 Forbidden - ServiceUnavailable   - Description: Reduce your request rate.   - HTTP Status Code: 503 Service Unavailable - SlowDown   - Description: Reduce your request rate.   - HTTP Status Code: 503 Slow Down - TemporaryRedirect   - Description: You are being redirected to the bucket while DNS updates.   - HTTP Status Code: 307 Moved Temporarily - TooManyBuckets   - Description: You have attempted to create more buckets than allowed.   - HTTP Status Code: 400 Bad Request - UnexpectedContent   - Description: This request does not support content.   - HTTP Status Code: 400 Bad Request - UnresolvableGrantByEmailAddress   - Description: The email address you provided does not match any account on record.   - HTTP Status Code: 400 Bad Request - UserKeyMustBeSpecified   - Description: The bucket POST must contain the specified field name. If it is specified, check the order of the fields.   - HTTP Status Code: 400 Bad Request  | [optional] |
|**Message** | Pointer to **string** | Gives a brief English description of the issue. | [optional] |
|**RequestId** | Pointer to **string** |  | [optional] |
|**HostId** | Pointer to **string** |  | [optional] |

## Methods

### NewError

`func NewError() *Error`

NewError instantiates a new Error object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewErrorWithDefaults

`func NewErrorWithDefaults() *Error`

NewErrorWithDefaults instantiates a new Error object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCode

`func (o *Error) GetCode() string`

GetCode returns the Code field if non-nil, zero value otherwise.

### GetCodeOk

`func (o *Error) GetCodeOk() (*string, bool)`

GetCodeOk returns a tuple with the Code field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCode

`func (o *Error) SetCode(v string)`

SetCode sets Code field to given value.

### HasCode

`func (o *Error) HasCode() bool`

HasCode returns a boolean if a field has been set.

### GetMessage

`func (o *Error) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *Error) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *Error) SetMessage(v string)`

SetMessage sets Message field to given value.

### HasMessage

`func (o *Error) HasMessage() bool`

HasMessage returns a boolean if a field has been set.

### GetRequestId

`func (o *Error) GetRequestId() string`

GetRequestId returns the RequestId field if non-nil, zero value otherwise.

### GetRequestIdOk

`func (o *Error) GetRequestIdOk() (*string, bool)`

GetRequestIdOk returns a tuple with the RequestId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestId

`func (o *Error) SetRequestId(v string)`

SetRequestId sets RequestId field to given value.

### HasRequestId

`func (o *Error) HasRequestId() bool`

HasRequestId returns a boolean if a field has been set.

### GetHostId

`func (o *Error) GetHostId() string`

GetHostId returns the HostId field if non-nil, zero value otherwise.

### GetHostIdOk

`func (o *Error) GetHostIdOk() (*string, bool)`

GetHostIdOk returns a tuple with the HostId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHostId

`func (o *Error) SetHostId(v string)`

SetHostId sets HostId field to given value.

### HasHostId

`func (o *Error) HasHostId() bool`

HasHostId returns a boolean if a field has been set.


