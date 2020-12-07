/* This file contains helpers to check whether given error
is specific http status code or not.

*/
package profitbricks

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	SnapshotInUseErrorCode = "200"
)

type ClientErrorType int

const (
	RequestFailed ClientErrorType = iota
	UnexpectedResponse
	HttpClientError
	InvalidInput
)

type ClientError struct {
	errType ClientErrorType
	msg     string
}

func (c ClientError) Error() string {
	return c.msg
}

func NewClientError(errType ClientErrorType, msg string) ClientError {
	return ClientError{
		errType: errType,
		msg:     msg,
	}
}

func IsClientErrorType(err error, errType ClientErrorType) bool {
	var clientErr ClientError
	if errors.As(err, &clientErr) {
		return clientErr.errType == errType
	}
	return false
}

func IsHttpStatus(err error, status int) bool {
	var apiErr ApiError
	if errors.As(err, &apiErr) {
		return apiErr.HttpStatusCode() == status
	}
	return false
}

// IsStatusOK - (200)
func IsStatusOK(err error) bool {
	return IsHttpStatus(err, http.StatusOK)
}

// IsStatusAccepted - (202) Used for asynchronous requests using PUT, DELETE, POST and PATCH methods.
// The response will also include a Location header pointing to a resource. This can be used for polling.
func IsStatusAccepted(err error) bool {
	return IsHttpStatus(err, http.StatusAccepted)
}

// IsStatusNotModified - (304) Response for GETs on resources that have not been changed. (based on ETag values).
func IsStatusNotModified(err error) bool {
	return IsHttpStatus(err, http.StatusNotModified)
}

// IsStatusBadRequest - (400) Response to malformed requests or general client errors.
func IsStatusBadRequest(err error) bool {
	return IsHttpStatus(err, http.StatusBadRequest)
}

// IsStatusUnauthorized - (401) Response to an unauthenticated connection.
// You will need to use your API username and password to be authenticated.
func IsStatusUnauthorized(err error) bool {
	return IsHttpStatus(err, http.StatusUnauthorized)
}

// IsStatusForbidden - (403) Forbidden
func IsStatusForbidden(err error) bool {
	return IsHttpStatus(err, http.StatusForbidden)
}

// IsStatusNotFound - (404) if resource does not exist
func IsStatusNotFound(err error) bool {
	return IsHttpStatus(err, http.StatusNotFound)
}

// IsStatusMethodNotAllowed - (405) Use for any POST, PUT, PATCH, or DELETE performed
// on read-only resources. This is also the response to PATCH requests
// on resources that do not support partial updates.
func IsStatusMethodNotAllowed(err error) bool {
	return IsHttpStatus(err, http.StatusMethodNotAllowed)
}

// IsStatusUnsupportedMediaType - (415) The content-type is incorrect for the payload.
func IsStatusUnsupportedMediaType(err error) bool {
	return IsHttpStatus(err, http.StatusUnsupportedMediaType)
}

// IsStatusUnprocessableEntity - (422) Validation errors.
func IsStatusUnprocessableEntity(err error) bool {
	return IsHttpStatus(err, http.StatusUnprocessableEntity)
}

// IsStatusTooManyRequests - (429) The number of requests exceeds the rate limit.
func IsStatusTooManyRequests(err error) bool {
	return IsHttpStatus(err, http.StatusTooManyRequests)
}

// IsRequestFailed - returns true if the error reason was that the request status was failed
func IsRequestFailed(err error) bool {
	return IsClientErrorType(err, RequestFailed)
}

type ApiError struct {
	HTTPStatus int `json:"httpStatus"`
	Messages   []struct {
		ErrorCode string `json:"errorCode"`
		Message   string `json:"message"`
	} `json:"messages"`
	RawBody []byte
}

func (e ApiError) Error() string {
	return e.String()
}

func (e ApiError) String() string {
	toReturn := fmt.Sprintf("HTTP Status: %d\nError Messages:", e.HTTPStatus)
	for _, m := range e.Messages {
		toReturn = toReturn + fmt.Sprintf("Error Code: %s Message: %s\n", m.ErrorCode, m.Message)
	}
	return toReturn
}

func (e ApiError) HttpStatusCode() int {
	return e.HTTPStatus
}

func (e ApiError) Body() []byte {
	return e.RawBody
}

func (e ApiError) HasErrorCode(code string) bool {
	for _, m := range e.Messages {
		if m.ErrorCode == code {
			return true
		}
	}
	return false
}
