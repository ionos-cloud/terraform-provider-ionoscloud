/*
 * IONOS Cloud VPN Gateway API
 *
 * The Managed VPN Gateway service provides secure and scalable connectivity, enabling encrypted communication between your IONOS cloud resources in a VDC and remote networks (on-premises, multi-cloud, private LANs in other VDCs etc).
 *
 * API version: 1.0.0
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package vpn

import (
	_context "context"
	"fmt"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"io"
	_nethttp "net/http"
	_neturl "net/url"
	"strings"
)

// Linger please
var (
	_ _context.Context
)

// IPSecGatewaysApiService IPSecGatewaysApi service
type IPSecGatewaysApiService service

type ApiIpsecgatewaysDeleteRequest struct {
	ctx        _context.Context
	ApiService *IPSecGatewaysApiService
	gatewayId  string
}

func (r ApiIpsecgatewaysDeleteRequest) Execute() (*shared.APIResponse, error) {
	return r.ApiService.IpsecgatewaysDeleteExecute(r)
}

/*
 * IpsecgatewaysDelete Delete IPSecGateway
 * Deletes the specified IPSecGateway.
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param gatewayId The ID (UUID) of the IPSecGateway.
 * @return ApiIpsecgatewaysDeleteRequest
 */
func (a *IPSecGatewaysApiService) IpsecgatewaysDelete(ctx _context.Context, gatewayId string) ApiIpsecgatewaysDeleteRequest {
	return ApiIpsecgatewaysDeleteRequest{
		ApiService: a,
		ctx:        ctx,
		gatewayId:  gatewayId,
	}
}

/*
 * Execute executes the request
 */
func (a *IPSecGatewaysApiService) IpsecgatewaysDeleteExecute(r ApiIpsecgatewaysDeleteRequest) (*shared.APIResponse, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodDelete
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "IPSecGatewaysApiService.IpsecgatewaysDelete")
	if err != nil {
		gerr := shared.GenericOpenAPIError{}
		gerr.SetError(err.Error())
		return nil, gerr
	}

	localVarPath := localBasePath + "/ipsecgateways/{gatewayId}"
	localVarPath = strings.Replace(localVarPath, "{"+"gatewayId"+"}", _neturl.PathEscape(parameterValueToString(r.gatewayId, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFormFileName, localVarFileName, localVarFileBytes)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, httpRequestTime, err := a.client.callAPI(req)

	localVarAPIResponse := &shared.APIResponse{
		Response:    localVarHTTPResponse,
		Method:      localVarHTTPMethod,
		RequestTime: httpRequestTime,
		RequestURL:  localVarPath,
		Operation:   "IpsecgatewaysDelete",
	}

	if err != nil || localVarHTTPResponse == nil {
		return localVarAPIResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	if err != nil {
		return localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := shared.GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)))
		if localVarHTTPResponse.StatusCode == 400 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 403 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 404 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 429 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 503 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		var v Error
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.SetError(err.Error())
			return localVarAPIResponse, newErr
		}
		newErr.SetModel(v)
		return localVarAPIResponse, newErr
	}

	return localVarAPIResponse, nil
}

type ApiIpsecgatewaysFindByIdRequest struct {
	ctx        _context.Context
	ApiService *IPSecGatewaysApiService
	gatewayId  string
}

func (r ApiIpsecgatewaysFindByIdRequest) Execute() (IPSecGatewayRead, *shared.APIResponse, error) {
	return r.ApiService.IpsecgatewaysFindByIdExecute(r)
}

/*
 * IpsecgatewaysFindById Retrieve IPSecGateway
 * Returns the IPSecGateway by ID.
 * @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param gatewayId The ID (UUID) of the IPSecGateway.
 * @return ApiIpsecgatewaysFindByIdRequest
 */
func (a *IPSecGatewaysApiService) IpsecgatewaysFindById(ctx _context.Context, gatewayId string) ApiIpsecgatewaysFindByIdRequest {
	return ApiIpsecgatewaysFindByIdRequest{
		ApiService: a,
		ctx:        ctx,
		gatewayId:  gatewayId,
	}
}

/*
 * Execute executes the request
 * @return IPSecGatewayRead
 */
func (a *IPSecGatewaysApiService) IpsecgatewaysFindByIdExecute(r ApiIpsecgatewaysFindByIdRequest) (IPSecGatewayRead, *shared.APIResponse, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodGet
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  IPSecGatewayRead
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "IPSecGatewaysApiService.IpsecgatewaysFindById")
	if err != nil {
		gerr := shared.GenericOpenAPIError{}
		gerr.SetError(err.Error())
		return localVarReturnValue, nil, gerr
	}

	localVarPath := localBasePath + "/ipsecgateways/{gatewayId}"
	localVarPath = strings.Replace(localVarPath, "{"+"gatewayId"+"}", _neturl.PathEscape(parameterValueToString(r.gatewayId, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFormFileName, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, httpRequestTime, err := a.client.callAPI(req)

	localVarAPIResponse := &shared.APIResponse{
		Response:    localVarHTTPResponse,
		Method:      localVarHTTPMethod,
		RequestTime: httpRequestTime,
		RequestURL:  localVarPath,
		Operation:   "IpsecgatewaysFindById",
	}

	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	if err != nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := shared.GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)))
		if localVarHTTPResponse.StatusCode == 400 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 403 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 404 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 429 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 503 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		var v Error
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.SetError(err.Error())
			return localVarReturnValue, localVarAPIResponse, newErr
		}
		newErr.SetModel(v)
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := shared.GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(err.Error())
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	return localVarReturnValue, localVarAPIResponse, nil
}

type ApiIpsecgatewaysGetRequest struct {
	ctx        _context.Context
	ApiService *IPSecGatewaysApiService
	offset     *int32
	limit      *int32
}

func (r ApiIpsecgatewaysGetRequest) Offset(offset int32) ApiIpsecgatewaysGetRequest {
	r.offset = &offset
	return r
}
func (r ApiIpsecgatewaysGetRequest) Limit(limit int32) ApiIpsecgatewaysGetRequest {
	r.limit = &limit
	return r
}

func (r ApiIpsecgatewaysGetRequest) Execute() (IPSecGatewayReadList, *shared.APIResponse, error) {
	return r.ApiService.IpsecgatewaysGetExecute(r)
}

/*
  - IpsecgatewaysGet Retrieve all IPSecGateways
  - This endpoint enables retrieving all IPSecGateways using

pagination and optional filters.

  - @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @return ApiIpsecgatewaysGetRequest
*/
func (a *IPSecGatewaysApiService) IpsecgatewaysGet(ctx _context.Context) ApiIpsecgatewaysGetRequest {
	return ApiIpsecgatewaysGetRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

/*
 * Execute executes the request
 * @return IPSecGatewayReadList
 */
func (a *IPSecGatewaysApiService) IpsecgatewaysGetExecute(r ApiIpsecgatewaysGetRequest) (IPSecGatewayReadList, *shared.APIResponse, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodGet
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  IPSecGatewayReadList
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "IPSecGatewaysApiService.IpsecgatewaysGet")
	if err != nil {
		gerr := shared.GenericOpenAPIError{}
		gerr.SetError(err.Error())
		return localVarReturnValue, nil, gerr
	}

	localVarPath := localBasePath + "/ipsecgateways"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}

	if r.offset != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "offset", r.offset, "")
	}
	if r.limit != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "limit", r.limit, "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFormFileName, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, httpRequestTime, err := a.client.callAPI(req)

	localVarAPIResponse := &shared.APIResponse{
		Response:    localVarHTTPResponse,
		Method:      localVarHTTPMethod,
		RequestTime: httpRequestTime,
		RequestURL:  localVarPath,
		Operation:   "IpsecgatewaysGet",
	}

	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	if err != nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := shared.GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)))
		if localVarHTTPResponse.StatusCode == 400 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 403 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 429 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 503 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		var v Error
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.SetError(err.Error())
			return localVarReturnValue, localVarAPIResponse, newErr
		}
		newErr.SetModel(v)
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := shared.GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(err.Error())
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	return localVarReturnValue, localVarAPIResponse, nil
}

type ApiIpsecgatewaysPostRequest struct {
	ctx                _context.Context
	ApiService         *IPSecGatewaysApiService
	iPSecGatewayCreate *IPSecGatewayCreate
}

func (r ApiIpsecgatewaysPostRequest) IPSecGatewayCreate(iPSecGatewayCreate IPSecGatewayCreate) ApiIpsecgatewaysPostRequest {
	r.iPSecGatewayCreate = &iPSecGatewayCreate
	return r
}

func (r ApiIpsecgatewaysPostRequest) Execute() (IPSecGatewayRead, *shared.APIResponse, error) {
	return r.ApiService.IpsecgatewaysPostExecute(r)
}

/*
  - IpsecgatewaysPost Create IPSecGateway
  - Creates a new IPSecGateway.

The full IPSecGateway needs to be provided to create the object.
Optional data will be filled with defaults or left empty.

  - @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @return ApiIpsecgatewaysPostRequest
*/
func (a *IPSecGatewaysApiService) IpsecgatewaysPost(ctx _context.Context) ApiIpsecgatewaysPostRequest {
	return ApiIpsecgatewaysPostRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

/*
 * Execute executes the request
 * @return IPSecGatewayRead
 */
func (a *IPSecGatewaysApiService) IpsecgatewaysPostExecute(r ApiIpsecgatewaysPostRequest) (IPSecGatewayRead, *shared.APIResponse, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodPost
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  IPSecGatewayRead
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "IPSecGatewaysApiService.IpsecgatewaysPost")
	if err != nil {
		gerr := shared.GenericOpenAPIError{}
		gerr.SetError(err.Error())
		return localVarReturnValue, nil, gerr
	}

	localVarPath := localBasePath + "/ipsecgateways"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}
	if r.iPSecGatewayCreate == nil {
		return localVarReturnValue, nil, reportError("iPSecGatewayCreate is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.iPSecGatewayCreate
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFormFileName, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, httpRequestTime, err := a.client.callAPI(req)

	localVarAPIResponse := &shared.APIResponse{
		Response:    localVarHTTPResponse,
		Method:      localVarHTTPMethod,
		RequestTime: httpRequestTime,
		RequestURL:  localVarPath,
		Operation:   "IpsecgatewaysPost",
	}

	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	if err != nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := shared.GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)))
		if localVarHTTPResponse.StatusCode == 400 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 403 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 415 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 422 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 429 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 503 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		var v Error
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.SetError(err.Error())
			return localVarReturnValue, localVarAPIResponse, newErr
		}
		newErr.SetModel(v)
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := shared.GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(err.Error())
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	return localVarReturnValue, localVarAPIResponse, nil
}

type ApiIpsecgatewaysPutRequest struct {
	ctx                _context.Context
	ApiService         *IPSecGatewaysApiService
	gatewayId          string
	iPSecGatewayEnsure *IPSecGatewayEnsure
}

func (r ApiIpsecgatewaysPutRequest) IPSecGatewayEnsure(iPSecGatewayEnsure IPSecGatewayEnsure) ApiIpsecgatewaysPutRequest {
	r.iPSecGatewayEnsure = &iPSecGatewayEnsure
	return r
}

func (r ApiIpsecgatewaysPutRequest) Execute() (IPSecGatewayRead, *shared.APIResponse, error) {
	return r.ApiService.IpsecgatewaysPutExecute(r)
}

/*
  - IpsecgatewaysPut Ensure IPSecGateway
  - Ensures that the IPSecGateway with the provided ID is created or modified.

The full IPSecGateway needs to be provided to ensure
(either update or create) the IPSecGateway. Non present data will
only be filled with defaults or left empty, but not take
previous values into consideration.

  - @param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
  - @param gatewayId The ID (UUID) of the IPSecGateway.
  - @return ApiIpsecgatewaysPutRequest
*/
func (a *IPSecGatewaysApiService) IpsecgatewaysPut(ctx _context.Context, gatewayId string) ApiIpsecgatewaysPutRequest {
	return ApiIpsecgatewaysPutRequest{
		ApiService: a,
		ctx:        ctx,
		gatewayId:  gatewayId,
	}
}

/*
 * Execute executes the request
 * @return IPSecGatewayRead
 */
func (a *IPSecGatewaysApiService) IpsecgatewaysPutExecute(r ApiIpsecgatewaysPutRequest) (IPSecGatewayRead, *shared.APIResponse, error) {
	var (
		localVarHTTPMethod   = _nethttp.MethodPut
		localVarPostBody     interface{}
		localVarFormFileName string
		localVarFileName     string
		localVarFileBytes    []byte
		localVarReturnValue  IPSecGatewayRead
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "IPSecGatewaysApiService.IpsecgatewaysPut")
	if err != nil {
		gerr := shared.GenericOpenAPIError{}
		gerr.SetError(err.Error())
		return localVarReturnValue, nil, gerr
	}

	localVarPath := localBasePath + "/ipsecgateways/{gatewayId}"
	localVarPath = strings.Replace(localVarPath, "{"+"gatewayId"+"}", _neturl.PathEscape(parameterValueToString(r.gatewayId, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}
	if r.iPSecGatewayEnsure == nil {
		return localVarReturnValue, nil, reportError("iPSecGatewayEnsure is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.iPSecGatewayEnsure
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFormFileName, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, httpRequestTime, err := a.client.callAPI(req)

	localVarAPIResponse := &shared.APIResponse{
		Response:    localVarHTTPResponse,
		Method:      localVarHTTPMethod,
		RequestTime: httpRequestTime,
		RequestURL:  localVarPath,
		Operation:   "IpsecgatewaysPut",
	}

	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarAPIResponse.Payload = localVarBody
	if err != nil {
		return localVarReturnValue, localVarAPIResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := shared.GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(fmt.Sprintf("%s: %s", localVarHTTPResponse.Status, string(localVarBody)))
		if localVarHTTPResponse.StatusCode == 400 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 403 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 404 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 409 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 415 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 422 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 429 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		if localVarHTTPResponse.StatusCode == 503 {
			var v Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.SetError(err.Error())
				return localVarReturnValue, localVarAPIResponse, newErr
			}
			newErr.SetModel(v)
		}
		var v Error
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.SetError(err.Error())
			return localVarReturnValue, localVarAPIResponse, newErr
		}
		newErr.SetModel(v)
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := shared.GenericOpenAPIError{}
		newErr.SetStatusCode(localVarHTTPResponse.StatusCode)
		newErr.SetBody(localVarBody)
		newErr.SetError(err.Error())
		return localVarReturnValue, localVarAPIResponse, newErr
	}

	return localVarReturnValue, localVarAPIResponse, nil
}
