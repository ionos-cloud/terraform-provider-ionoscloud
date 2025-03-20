package bundleclient

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GetStateChangeConf gets the default configuration for tracking a request progress
func GetStateChangeConf(meta interface{}, d *schema.ResourceData, location string, timeoutType string) *retry.StateChangeConf {
	stateConf := &retry.StateChangeConf{
		Pending:        resourcePendingStates,
		Target:         resourceTargetStates,
		Refresh:        resourceStateRefreshFunc(meta, location),
		Timeout:        d.Timeout(timeoutType),
		MinTimeout:     5 * time.Second,
		Delay:          0,    // Don't delay the start
		NotFoundChecks: 9999, // Setting high number, to support long timeouts
	}

	return stateConf
}

// resourceStateRefreshFunc tracks progress of a request
func resourceStateRefreshFunc(meta interface{}, path string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		client := meta.(SdkBundle).CloudApiClient

		log.Printf("[INFO] Checking PATH %s\n", path)
		if path == "" {
			return nil, "", fmt.Errorf("can not check a state when path is empty")
		}

		request, apiResponse, err := client.GetRequestStatus(context.Background(), path)
		apiResponse.LogInfo()
		if err != nil {
			return nil, "", fmt.Errorf("request failed with following error: %w", err)
		}
		if request != nil && request.Metadata != nil && request.Metadata.Status != nil {
			if *request.Metadata.Status == "FAILED" {
				var msg string
				if request.Metadata.Message != nil {
					msg = fmt.Sprintf("Request failed with following error: %s", *request.Metadata.Message)
				} else {
					msg = "Request failed with an unknown error"
				}
				return nil, "", RequestFailedError{msg}
			}

			if *request.Metadata.Status == "DONE" {
				return request, "DONE", nil
			}
		} else {
			if request == nil {
				log.Printf("[DEBUG] request is nil")
			} else if request.Metadata == nil {
				log.Printf("[DEBUG] request metadata is nil")
			}
			if request != nil && request.Metadata != nil && request.Metadata.Message != nil {
				log.Printf("[DEBUG] request failed with following error: %s", *request.Metadata.Message)
			}
			if apiResponse != nil {
				log.Printf("[DEBUG] response message %s", apiResponse.Message)
			}
			return nil, "", RequestFailedError{fmt.Sprintf("request metadata status is nil for path %s", path)}
		}

		return nil, *request.Metadata.Status, nil
	}
}

// WaitForStateChange tracks state change progress of a resource
func WaitForStateChange(ctx context.Context, meta interface{}, d *schema.ResourceData, apiResponse *ionoscloud.APIResponse, opTimeout string) error {
	var err error
	var loc *url.URL

	if apiResponse == nil {
		return fmt.Errorf("cannot track resource state change, apiResponse was nil")
	}

	if loc, err = apiResponse.Location(); err != nil {
		return fmt.Errorf("error retrieving 'location' header: %w", err)
	}
	_, errState := GetStateChangeConf(meta, d, loc.String(), opTimeout).WaitForStateContext(ctx)
	return errState
}

// resourcePendingStates defines states of working in progress
var resourcePendingStates = []string{
	"RUNNING",
	"QUEUED",
}

// resourceTargetStates defines states of completion
var resourceTargetStates = []string{
	"DONE",
}

type RequestFailedError struct {
	Msg string
}

func (e RequestFailedError) Error() string {
	return e.Msg
}

func IsRequestFailed(err error) bool {
	var requestFailedError RequestFailedError
	ok := errors.As(err, &requestFailedError)
	return ok
}
