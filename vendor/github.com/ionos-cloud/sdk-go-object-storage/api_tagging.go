/*
 * IONOS Object Storage API for contract-owned buckets
 *
 * ## Overview The IONOS Object Storage API for contract-owned buckets is a REST-based API that allows developers and applications to interact directly with IONOS' scalable storage solution, leveraging the S3 protocol for object storage operations. Its design ensures seamless compatibility with existing tools and libraries tailored for S3 systems.  ### API References - [S3 API Reference for contract-owned buckets](https://api.ionos.com/docs/s3-contract-owned-buckets/v2/) ### User documentation [IONOS Object Storage User Guide](https://docs.ionos.com/cloud/managed-services/s3-object-storage) * [Documentation on user-owned and contract-owned buckets](https://docs.ionos.com/cloud/managed-services/s3-object-storage/concepts/buckets) * [Documentation on S3 API Compatibility](https://docs.ionos.com/cloud/managed-services/s3-object-storage/concepts/s3-api-compatibility) * [S3 Tools](https://docs.ionos.com/cloud/managed-services/s3-object-storage/s3-tools)  ## Endpoints for contract-owned buckets | Location | Region Name | Bucket Type | Endpoint | | --- | --- | --- | --- | | **Berlin, Germany** | **eu-central-3** | Contract-owned | `https://s3.eu-central-3.ionoscloud.com` |  ## Changelog - 30.05.2024 Initial version
 *
 * API version: 2.0.2
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// TaggingApiService TaggingApi service
type TaggingApiService service

type ApiDeleteBucketTaggingRequest struct {
	ctx        context.Context
	ApiService *TaggingApiService
	bucket     string
}

func (r ApiDeleteBucketTaggingRequest) Execute() (*APIResponse, error) {
	return r.ApiService.DeleteBucketTaggingExecute(r)
}

/*
DeleteBucketTagging DeleteBucketTagging

<p>Deletes the tags from the bucket.</p> <p>To use this operation, you must have permission to perform the `PutBucketTagging` operation. By default, the bucket owner has this permission and can grant this permission to others.</p>

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param bucket
	@return ApiDeleteBucketTaggingRequest
*/
func (a *TaggingApiService) DeleteBucketTagging(ctx context.Context, bucket string) ApiDeleteBucketTaggingRequest {
	return ApiDeleteBucketTaggingRequest{
		ApiService: a,
		ctx:        ctx,
		bucket:     bucket,
	}
}

// Execute executes the request
func (a *TaggingApiService) DeleteBucketTaggingExecute(r ApiDeleteBucketTaggingRequest) (*APIResponse, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   interface{}
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TaggingApiService.DeleteBucketTagging")
	if err != nil {
		gerr := GenericOpenAPIError{}
		gerr.SetError(err.Error())
		return nil, gerr
	}

	localVarPath := localBasePath + "/{Bucket}?tagging"
	localVarPath = strings.Replace(localVarPath, "{"+"Bucket"+"}", parameterValueToString(r.bucket, "bucket"), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if Strlen(r.bucket) < 3 {
		return nil, reportError("bucket must have at least 3 elements")
	}
	if Strlen(r.bucket) > 63 {
		return nil, reportError("bucket must have less than 63 elements")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(ContextAPIKeys).(map[string]APIKey); ok {
			if apiKey, ok := auth["hmac"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["Authorization"] = key
			}
		}
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, httpRequestTime, err := a.client.callAPI(req)
	localVarAPIResponse := &APIResponse{
		Response:    localVarHTTPResponse,
		Method:      localVarHTTPMethod,
		RequestTime: httpRequestTime,
		RequestURL:  localVarPath,
		Operation:   "DeleteBucketTagging",
	}
	if err != nil || localVarHTTPResponse == nil {
		return localVarAPIResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)))
		return localVarAPIResponse, newErr
	}

	return localVarAPIResponse, nil
}

type ApiDeleteObjectTaggingRequest struct {
	ctx        context.Context
	ApiService *TaggingApiService
	bucket     string
	key        string
	versionId  *string
}

// The versionId of the object that the tag-set will be removed from.
func (r ApiDeleteObjectTaggingRequest) VersionId(versionId string) ApiDeleteObjectTaggingRequest {
	r.versionId = &versionId
	return r
}

func (r ApiDeleteObjectTaggingRequest) Execute() (map[string]interface{}, *APIResponse, error) {
	return r.ApiService.DeleteObjectTaggingExecute(r)
}

/*
DeleteObjectTagging DeleteObjectTagging

<p>Removes the entire tag set from the specified object.</p>  <p>To use this operation, you must have permission to perform the `DeleteObjectTagging` operation.</p> <p>To delete tags of a specific object version, add the `versionId` query parameter in the request. You will need permission for the `DeleteObjectVersionTagging` operation.</p>

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param bucket
	@param key The key that identifies the object in the bucket from which to remove all tags.
	@return ApiDeleteObjectTaggingRequest
*/
func (a *TaggingApiService) DeleteObjectTagging(ctx context.Context, bucket string, key string) ApiDeleteObjectTaggingRequest {
	return ApiDeleteObjectTaggingRequest{
		ApiService: a,
		ctx:        ctx,
		bucket:     bucket,
		key:        key,
	}
}

// Execute executes the request
//
//	@return map[string]interface{}
func (a *TaggingApiService) DeleteObjectTaggingExecute(r ApiDeleteObjectTaggingRequest) (map[string]interface{}, *APIResponse, error) {
	var (
		localVarHTTPMethod  = http.MethodDelete
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue map[string]interface{}
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TaggingApiService.DeleteObjectTagging")
	if err != nil {
		gerr := GenericOpenAPIError{}
		gerr.SetError(err.Error())
		return localVarReturnValue, nil, gerr
	}

	localVarPath := localBasePath + "/{Bucket}/{Key}?tagging"
	localVarPath = strings.Replace(localVarPath, "{"+"Bucket"+"}", parameterValueToString(r.bucket, "bucket"), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"Key"+"}", parameterValueToString(r.key, "key"), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if Strlen(r.bucket) < 3 {
		return localVarReturnValue, nil, reportError("bucket must have at least 3 elements")
	}
	if Strlen(r.bucket) > 63 {
		return localVarReturnValue, nil, reportError("bucket must have less than 63 elements")
	}
	if Strlen(r.key) < 1 {
		return localVarReturnValue, nil, reportError("key must have at least 1 elements")
	}

	if r.versionId != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "versionId", r.versionId, "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/xml"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(ContextAPIKeys).(map[string]APIKey); ok {
			if apiKey, ok := auth["hmac"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["Authorization"] = key
			}
		}
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, httpRequestTime, err := a.client.callAPI(req)
	localVarAPIResponse := &APIResponse{
		Response:    localVarHTTPResponse,
		Method:      localVarHTTPMethod,
		RequestTime: httpRequestTime,
		RequestURL:  localVarPath,
		Operation:   "DeleteObjectTagging",
	}
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)))
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(err.Error())
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	return localVarReturnValue, localVarAPIResponse, nil
}

type ApiGetBucketTaggingRequest struct {
	ctx        context.Context
	ApiService *TaggingApiService
	bucket     string
}

func (r ApiGetBucketTaggingRequest) Execute() (*GetBucketTaggingOutput, *APIResponse, error) {
	return r.ApiService.GetBucketTaggingExecute(r)
}

/*
GetBucketTagging GetBucketTagging

<p>Returns the tag set associated with the bucket.</p> <p>To use this operation, you must have permission to perform the `GetBucketTagging` operation. By default, the bucket owner has this permission and can grant this permission to others.</p> <p> `GetBucketTagging` has the following special error:</p> <ul> <li> <p>Error code: `NoSuchTagSetError` </p> <ul> <li> <p>Description: There is no tag set associated with the bucket.</p> </li> </ul> </li> </ul>

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param bucket
	@return ApiGetBucketTaggingRequest
*/
func (a *TaggingApiService) GetBucketTagging(ctx context.Context, bucket string) ApiGetBucketTaggingRequest {
	return ApiGetBucketTaggingRequest{
		ApiService: a,
		ctx:        ctx,
		bucket:     bucket,
	}
}

// Execute executes the request
//
//	@return GetBucketTaggingOutput
func (a *TaggingApiService) GetBucketTaggingExecute(r ApiGetBucketTaggingRequest) (*GetBucketTaggingOutput, *APIResponse, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *GetBucketTaggingOutput
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TaggingApiService.GetBucketTagging")
	if err != nil {
		gerr := GenericOpenAPIError{}
		gerr.SetError(err.Error())
		return localVarReturnValue, nil, gerr
	}

	localVarPath := localBasePath + "/{Bucket}?tagging"
	localVarPath = strings.Replace(localVarPath, "{"+"Bucket"+"}", parameterValueToString(r.bucket, "bucket"), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if Strlen(r.bucket) < 3 {
		return localVarReturnValue, nil, reportError("bucket must have at least 3 elements")
	}
	if Strlen(r.bucket) > 63 {
		return localVarReturnValue, nil, reportError("bucket must have less than 63 elements")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/xml"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(ContextAPIKeys).(map[string]APIKey); ok {
			if apiKey, ok := auth["hmac"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["Authorization"] = key
			}
		}
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, httpRequestTime, err := a.client.callAPI(req)
	localVarAPIResponse := &APIResponse{
		Response:    localVarHTTPResponse,
		Method:      localVarHTTPMethod,
		RequestTime: httpRequestTime,
		RequestURL:  localVarPath,
		Operation:   "GetBucketTagging",
	}
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)))
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(err.Error())
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	return localVarReturnValue, localVarAPIResponse, nil
}

type ApiGetObjectTaggingRequest struct {
	ctx        context.Context
	ApiService *TaggingApiService
	bucket     string
	key        string
	versionId  *string
}

// The versionId of the object for which to get the tagging information.
func (r ApiGetObjectTaggingRequest) VersionId(versionId string) ApiGetObjectTaggingRequest {
	r.versionId = &versionId
	return r
}

func (r ApiGetObjectTaggingRequest) Execute() (*GetObjectTaggingOutput, *APIResponse, error) {
	return r.ApiService.GetObjectTaggingExecute(r)
}

/*
GetObjectTagging GetObjectTagging

<p>Returns the tag-set of an object. You send the GET request against the tagging subresource associated with the object.</p> <p>To use this operation, you must have permission to perform the `GetObjectTagging` operation. By default, the GET operation returns information about current version of an object. For a versioned bucket, you can have multiple versions of an object in your bucket. To retrieve tags of any other version, use the versionId query parameter. You also need permission for the `GetObjectVersionTagging` operation.</p> <p> By default, the bucket owner has this permission and can grant this permission to others.</p>

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param bucket
	@param key Object key for which to get the tagging information.
	@return ApiGetObjectTaggingRequest
*/
func (a *TaggingApiService) GetObjectTagging(ctx context.Context, bucket string, key string) ApiGetObjectTaggingRequest {
	return ApiGetObjectTaggingRequest{
		ApiService: a,
		ctx:        ctx,
		bucket:     bucket,
		key:        key,
	}
}

// Execute executes the request
//
//	@return GetObjectTaggingOutput
func (a *TaggingApiService) GetObjectTaggingExecute(r ApiGetObjectTaggingRequest) (*GetObjectTaggingOutput, *APIResponse, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *GetObjectTaggingOutput
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TaggingApiService.GetObjectTagging")
	if err != nil {
		gerr := GenericOpenAPIError{}
		gerr.SetError(err.Error())
		return localVarReturnValue, nil, gerr
	}

	localVarPath := localBasePath + "/{Bucket}/{Key}?tagging"
	localVarPath = strings.Replace(localVarPath, "{"+"Bucket"+"}", parameterValueToString(r.bucket, "bucket"), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"Key"+"}", parameterValueToString(r.key, "key"), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if Strlen(r.bucket) < 3 {
		return localVarReturnValue, nil, reportError("bucket must have at least 3 elements")
	}
	if Strlen(r.bucket) > 63 {
		return localVarReturnValue, nil, reportError("bucket must have less than 63 elements")
	}
	if Strlen(r.key) < 1 {
		return localVarReturnValue, nil, reportError("key must have at least 1 elements")
	}

	if r.versionId != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "versionId", r.versionId, "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/xml"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(ContextAPIKeys).(map[string]APIKey); ok {
			if apiKey, ok := auth["hmac"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["Authorization"] = key
			}
		}
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, httpRequestTime, err := a.client.callAPI(req)
	localVarAPIResponse := &APIResponse{
		Response:    localVarHTTPResponse,
		Method:      localVarHTTPMethod,
		RequestTime: httpRequestTime,
		RequestURL:  localVarPath,
		Operation:   "GetObjectTagging",
	}
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)))
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(err.Error())
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	return localVarReturnValue, localVarAPIResponse, nil
}

type ApiPutBucketTaggingRequest struct {
	ctx                     context.Context
	ApiService              *TaggingApiService
	bucket                  string
	putBucketTaggingRequest *PutBucketTaggingRequest
	contentMD5              *string
}

func (r ApiPutBucketTaggingRequest) PutBucketTaggingRequest(putBucketTaggingRequest PutBucketTaggingRequest) ApiPutBucketTaggingRequest {
	r.putBucketTaggingRequest = &putBucketTaggingRequest
	return r
}

func (r ApiPutBucketTaggingRequest) ContentMD5(contentMD5 string) ApiPutBucketTaggingRequest {
	r.contentMD5 = &contentMD5
	return r
}

func (r ApiPutBucketTaggingRequest) Execute() (*APIResponse, error) {
	return r.ApiService.PutBucketTaggingExecute(r)
}

/*
PutBucketTagging PutBucketTagging

<p>Sets the tags for a bucket.</p>          <note> <p> When this operation sets the tags for a bucket, it will overwrite any current tags the bucket already has. You cannot use this operation to add tags to an existing list of tags.</p> </note> <p>To use this operation, you must have permissions to perform the `PutBucketTagging` operation. The bucket owner has this permission by default and can grant this permission to others. </p> <p> `PutBucketTagging` has the following special errors:</p> <ul> <li> <p>Error code: `InvalidTagError` </p> <ul> <li> <p>Description: The tag provided was not a valid tag. This error can occur if the tag did not pass input validation. </p> </li> </ul> </li> <li> <p>Error code: `MalformedXMLError` </p> <ul> <li> <p>Description: The XML provided does not match the schema.</p> </li> </ul> </li> <li> <p>Error code: `OperationAbortedError ` </p> <ul> <li> <p>Description: A conflicting conditional operation is currently in progress against this resource. Please try again.</p> </li> </ul> </li> <li> <p>Error code: `InternalError` </p> <ul> <li> <p>Description: The service was unable to apply the provided tag to the bucket.</p> </li> </ul> </li> </ul>

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param bucket
	@return ApiPutBucketTaggingRequest
*/
func (a *TaggingApiService) PutBucketTagging(ctx context.Context, bucket string) ApiPutBucketTaggingRequest {
	return ApiPutBucketTaggingRequest{
		ApiService: a,
		ctx:        ctx,
		bucket:     bucket,
	}
}

// Execute executes the request
func (a *TaggingApiService) PutBucketTaggingExecute(r ApiPutBucketTaggingRequest) (*APIResponse, error) {
	var (
		localVarHTTPMethod = http.MethodPut
		localVarPostBody   interface{}
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TaggingApiService.PutBucketTagging")
	if err != nil {
		gerr := GenericOpenAPIError{}
		gerr.SetError(err.Error())
		return nil, gerr
	}

	localVarPath := localBasePath + "/{Bucket}?tagging"
	localVarPath = strings.Replace(localVarPath, "{"+"Bucket"+"}", parameterValueToString(r.bucket, "bucket"), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if Strlen(r.bucket) < 3 {
		return nil, reportError("bucket must have at least 3 elements")
	}
	if Strlen(r.bucket) > 63 {
		return nil, reportError("bucket must have less than 63 elements")
	}
	if r.putBucketTaggingRequest == nil {
		return nil, reportError("putBucketTaggingRequest is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/xml"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	if r.contentMD5 != nil {
		parameterAddToHeaderOrQuery(localVarHeaderParams, "Content-MD5", r.contentMD5, "")
	}
	// body params
	localVarPostBody = r.putBucketTaggingRequest
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(ContextAPIKeys).(map[string]APIKey); ok {
			if apiKey, ok := auth["hmac"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["Authorization"] = key
			}
		}
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, httpRequestTime, err := a.client.callAPI(req)
	localVarAPIResponse := &APIResponse{
		Response:    localVarHTTPResponse,
		Method:      localVarHTTPMethod,
		RequestTime: httpRequestTime,
		RequestURL:  localVarPath,
		Operation:   "PutBucketTagging",
	}
	if err != nil || localVarHTTPResponse == nil {
		return localVarAPIResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)))
		return localVarAPIResponse, newErr
	}

	return localVarAPIResponse, nil
}

type ApiPutObjectTaggingRequest struct {
	ctx                     context.Context
	ApiService              *TaggingApiService
	bucket                  string
	key                     string
	putObjectTaggingRequest *PutObjectTaggingRequest
	versionId               *string
	contentMD5              *string
}

func (r ApiPutObjectTaggingRequest) PutObjectTaggingRequest(putObjectTaggingRequest PutObjectTaggingRequest) ApiPutObjectTaggingRequest {
	r.putObjectTaggingRequest = &putObjectTaggingRequest
	return r
}

// The versionId of the object that the tag-set will be added to.
func (r ApiPutObjectTaggingRequest) VersionId(versionId string) ApiPutObjectTaggingRequest {
	r.versionId = &versionId
	return r
}

func (r ApiPutObjectTaggingRequest) ContentMD5(contentMD5 string) ApiPutObjectTaggingRequest {
	r.contentMD5 = &contentMD5
	return r
}

func (r ApiPutObjectTaggingRequest) Execute() (map[string]interface{}, *APIResponse, error) {
	return r.ApiService.PutObjectTaggingExecute(r)
}

/*
PutObjectTagging PutObjectTagging

<p>Sets the supplied tag-set to an object that already exists in a bucket.</p> <p>A tag is a key-value pair. You can associate tags with an object by sending a PUT request against the tagging subresource that is associated with the object. You can retrieve tags by sending a GET request.</p> <p>Note that IONOS Object Storage limits the maximum number of tags to 10 tags per object.</p> <p>To use this operation, you must have permission to perform the `PutObjectTagging` operation. By default, the bucket owner has this permission and can grant this permission to others.</p> <p>To put tags of any other version, use the `versionId` query parameter. You also need permission for the `PutObjectVersionTagging` operation.</p> <p class="title"> <b>Special Errors</b> </p> <ul> <li> <ul> <li> <p> <i>Code: InvalidTagError </i> </p> </li> <li> <p> <i>Cause: The tag provided was not a valid tag. This error can occur if the tag did not pass input validation.</i> </p> </li> </ul> </li> <li> <ul> <li> <p> <i>Code: MalformedXMLError </i> </p> </li> <li> <p> <i>Cause: The XML provided does not match the schema.</i> </p> </li> </ul> </li> <li> <ul> <li> <p> <i>Code: OperationAbortedError </i> </p> </li> <li> <p> <i>Cause: A conflicting conditional operation is currently in progress against this resource. Please try again.</i> </p> </li> </ul> </li> <li> <ul> <li> <p> <i>Code: InternalError</i> </p> </li> <li> <p> <i>Cause: The service was unable to apply the provided tag to the object.</i> </p> </li> </ul> </li> </ul>

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param bucket
	@param key Name of the object key.
	@return ApiPutObjectTaggingRequest
*/
func (a *TaggingApiService) PutObjectTagging(ctx context.Context, bucket string, key string) ApiPutObjectTaggingRequest {
	return ApiPutObjectTaggingRequest{
		ApiService: a,
		ctx:        ctx,
		bucket:     bucket,
		key:        key,
	}
}

// Execute executes the request
//
//	@return map[string]interface{}
func (a *TaggingApiService) PutObjectTaggingExecute(r ApiPutObjectTaggingRequest) (map[string]interface{}, *APIResponse, error) {
	var (
		localVarHTTPMethod  = http.MethodPut
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue map[string]interface{}
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "TaggingApiService.PutObjectTagging")
	if err != nil {
		gerr := GenericOpenAPIError{}
		gerr.SetError(err.Error())
		return localVarReturnValue, nil, gerr
	}

	localVarPath := localBasePath + "/{Bucket}/{Key}?tagging"
	localVarPath = strings.Replace(localVarPath, "{"+"Bucket"+"}", parameterValueToString(r.bucket, "bucket"), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"Key"+"}", parameterValueToString(r.key, "key"), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if Strlen(r.bucket) < 3 {
		return localVarReturnValue, nil, reportError("bucket must have at least 3 elements")
	}
	if Strlen(r.bucket) > 63 {
		return localVarReturnValue, nil, reportError("bucket must have less than 63 elements")
	}
	if Strlen(r.key) < 1 {
		return localVarReturnValue, nil, reportError("key must have at least 1 elements")
	}
	if r.putObjectTaggingRequest == nil {
		return localVarReturnValue, nil, reportError("putObjectTaggingRequest is required and must be specified")
	}

	if r.versionId != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "versionId", r.versionId, "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/xml"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/xml"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	if r.contentMD5 != nil {
		parameterAddToHeaderOrQuery(localVarHeaderParams, "Content-MD5", r.contentMD5, "")
	}
	// body params
	localVarPostBody = r.putObjectTaggingRequest
	if r.ctx != nil {
		// API Key Authentication
		if auth, ok := r.ctx.Value(ContextAPIKeys).(map[string]APIKey); ok {
			if apiKey, ok := auth["hmac"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["Authorization"] = key
			}
		}
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, httpRequestTime, err := a.client.callAPI(req)
	localVarAPIResponse := &APIResponse{
		Response:    localVarHTTPResponse,
		Method:      localVarHTTPMethod,
		RequestTime: httpRequestTime,
		RequestURL:  localVarPath,
		Operation:   "PutObjectTagging",
	}
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)))
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(err.Error())
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	return localVarReturnValue, localVarAPIResponse, nil
}
