/*
 * IONOS S3 Object Storage API for contract-owned buckets
 *
 * ## Overview The IONOS S3 Object Storage API for contract-owned buckets is a REST-based API that allows developers and applications to interact directly with IONOS' scalable storage solution, leveraging the S3 protocol for object storage operations. Its design ensures seamless compatibility with existing tools and libraries tailored for S3 systems.  ### API References - [S3 API Reference for contract-owned buckets](https://api.ionos.com/docs/s3-contract-owned-buckets/v2/) ### User documentation [IONOS S3 Object Storage User Guide](https://docs.ionos.com/cloud/managed-services/s3-object-storage) * [Documentation on user-owned and contract-owned buckets](https://docs.ionos.com/cloud/managed-services/s3-object-storage/concepts/buckets) * [Documentation on S3 API Compatibility](https://docs.ionos.com/cloud/managed-services/s3-object-storage/concepts/s3-api-compatibility) * [S3 Tools](https://docs.ionos.com/cloud/managed-services/s3-object-storage/s3-tools)  ## Endpoints for contract-owned buckets | Location | Region Name | Bucket Type | Endpoint | | --- | --- | --- | --- | | **Berlin, Germany** | **eu-central-3** | Contract-owned | `https://s3.eu-central-3.ionoscloud.com` |  ## Changelog - 30.05.2024 Initial version
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

// WebsiteApiService WebsiteApi service
type WebsiteApiService service

type ApiDeleteBucketWebsiteRequest struct {
	ctx        context.Context
	ApiService *WebsiteApiService
	bucket     string
}

func (r ApiDeleteBucketWebsiteRequest) Execute() (*APIResponse, error) {
	return r.ApiService.DeleteBucketWebsiteExecute(r)
}

/*
DeleteBucketWebsite DeleteBucketWebsite

<p>This operation removes the website configuration for a bucket. IONOS S3 Object Storage returns a `200 OK` response upon successfully deleting a website configuration on the specified bucket. You will get a `200 OK` response if the website configuration you are trying to delete does not exist on the bucket. IONOS S3 Object Storage returns a `404` response if the bucket specified in the request does not exist.</p> <p>This DELETE operation requires the `DeleteBucketWebsite` permission. By default, only the bucket owner can delete the website configuration attached to a bucket. However, bucket owners can grant other users permission to delete the website configuration by writing a bucket policy granting them the `DeleteBucketWebsite` permission.</p>

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param bucket
	@return ApiDeleteBucketWebsiteRequest
*/
func (a *WebsiteApiService) DeleteBucketWebsite(ctx context.Context, bucket string) ApiDeleteBucketWebsiteRequest {
	return ApiDeleteBucketWebsiteRequest{
		ApiService: a,
		ctx:        ctx,
		bucket:     bucket,
	}
}

// Execute executes the request
func (a *WebsiteApiService) DeleteBucketWebsiteExecute(r ApiDeleteBucketWebsiteRequest) (*APIResponse, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   interface{}
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "WebsiteApiService.DeleteBucketWebsite")
	if err != nil {
		gerr := GenericOpenAPIError{}
		gerr.SetError(err.Error())
		return nil, gerr
	}

	localVarPath := localBasePath + "/{Bucket}?website"
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
		Operation:   "DeleteBucketWebsite",
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

type ApiGetBucketWebsiteRequest struct {
	ctx        context.Context
	ApiService *WebsiteApiService
	bucket     string
}

func (r ApiGetBucketWebsiteRequest) Execute() (*GetBucketWebsiteOutput, *APIResponse, error) {
	return r.ApiService.GetBucketWebsiteExecute(r)
}

/*
GetBucketWebsite GetBucketWebsite

<p>Returns the website configuration for a bucket. </p> <p>This GET operation requires the `GetBucketWebsite` permission. By default, only the bucket owner can read the bucket website configuration. However, bucket owners can allow other users to read the website configuration by writing a bucket policy granting them the `GetBucketWebsite` permission.</p>

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param bucket
	@return ApiGetBucketWebsiteRequest
*/
func (a *WebsiteApiService) GetBucketWebsite(ctx context.Context, bucket string) ApiGetBucketWebsiteRequest {
	return ApiGetBucketWebsiteRequest{
		ApiService: a,
		ctx:        ctx,
		bucket:     bucket,
	}
}

// Execute executes the request
//
//	@return GetBucketWebsiteOutput
func (a *WebsiteApiService) GetBucketWebsiteExecute(r ApiGetBucketWebsiteRequest) (*GetBucketWebsiteOutput, *APIResponse, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *GetBucketWebsiteOutput
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "WebsiteApiService.GetBucketWebsite")
	if err != nil {
		gerr := GenericOpenAPIError{}
		gerr.SetError(err.Error())
		return localVarReturnValue, nil, gerr
	}

	localVarPath := localBasePath + "/{Bucket}?website"
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
		Operation:   "GetBucketWebsite",
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

type ApiPutBucketWebsiteRequest struct {
	ctx                     context.Context
	ApiService              *WebsiteApiService
	bucket                  string
	putBucketWebsiteRequest *PutBucketWebsiteRequest
	contentMD5              *string
}

func (r ApiPutBucketWebsiteRequest) PutBucketWebsiteRequest(putBucketWebsiteRequest PutBucketWebsiteRequest) ApiPutBucketWebsiteRequest {
	r.putBucketWebsiteRequest = &putBucketWebsiteRequest
	return r
}

func (r ApiPutBucketWebsiteRequest) ContentMD5(contentMD5 string) ApiPutBucketWebsiteRequest {
	r.contentMD5 = &contentMD5
	return r
}

func (r ApiPutBucketWebsiteRequest) Execute() (*APIResponse, error) {
	return r.ApiService.PutBucketWebsiteExecute(r)
}

/*
PutBucketWebsite PutBucketWebsite

<p>Sets the configuration of the website that is specified in the `website` subresource. To configure a bucket as a website, you can add this subresource on the bucket with website configuration information such as the file name of the index document and any redirect rules. </p>                  <p>This PUT operation requires the `PutBucketWebsite` permission. By default, only the bucket owner can configure the website attached to a bucket; however, bucket owners can allow other users to set the website configuration by writing a bucket policy that grants them the `PutBucketWebsite` permission.</p> <p>To redirect all website requests sent to the bucket's website endpoint, you add a website configuration with the following elements. Because all requests are sent to another website, you don't need to provide index document name for the bucket.</p> <ul> <li> <p> `WebsiteConfiguration` </p> </li> <li> <p> `RedirectAllRequestsTo` </p> </li> <li> <p> `HostName` </p> </li> <li> <p> `Protocol` </p> </li> </ul> <p>If you want granular control over redirects, you can use the following elements to add routing rules that describe conditions for redirecting requests and information about the redirect destination. In this case, the website configuration must provide an index document for the bucket, because some requests might not be redirected. </p> <ul> <li> <p> `WebsiteConfiguration` </p> </li> <li> <p> `IndexDocument` </p> </li> <li> <p> `Suffix` </p> </li> <li> <p> `ErrorDocument` </p> </li> <li> <p> `Key` </p> </li> <li> <p> `RoutingRules` </p> </li> <li> <p> `RoutingRule` </p> </li> <li> <p> `Condition` </p> </li> <li> <p> `HttpErrorCodeReturnedEquals` </p> </li> <li> <p> `KeyPrefixEquals` </p> </li> <li> <p> `Redirect` </p> </li> <li> <p> `Protocol` </p> </li> <li> <p> `HostName` </p> </li> <li> <p> `ReplaceKeyPrefixWith` </p> </li> <li> <p> `ReplaceKeyWith` </p> </li> <li> <p> `HttpRedirectCode` </p> </li> </ul>

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param bucket
	@return ApiPutBucketWebsiteRequest
*/
func (a *WebsiteApiService) PutBucketWebsite(ctx context.Context, bucket string) ApiPutBucketWebsiteRequest {
	return ApiPutBucketWebsiteRequest{
		ApiService: a,
		ctx:        ctx,
		bucket:     bucket,
	}
}

// Execute executes the request
func (a *WebsiteApiService) PutBucketWebsiteExecute(r ApiPutBucketWebsiteRequest) (*APIResponse, error) {
	var (
		localVarHTTPMethod = http.MethodPut
		localVarPostBody   interface{}
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "WebsiteApiService.PutBucketWebsite")
	if err != nil {
		gerr := GenericOpenAPIError{}
		gerr.SetError(err.Error())
		return nil, gerr
	}

	localVarPath := localBasePath + "/{Bucket}?website"
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
	if r.putBucketWebsiteRequest == nil {
		return nil, reportError("putBucketWebsiteRequest is required and must be specified")
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
	localVarPostBody = r.putBucketWebsiteRequest
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
		Operation:   "PutBucketWebsite",
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