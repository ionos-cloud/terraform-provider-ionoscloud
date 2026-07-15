package ionoscloud

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/compute/v2"
)

var labelResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"key": {
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile("^[a-z0-9]*$"), "Invalid label key, please provide only lower case alphabets (a-z), or numeric (0-9) characters")),
		},
		"value": {
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile("^[a-z0-9]*$"), "Invalid label value, please provide only lower case alphabets (a-z), or numeric (0-9) characters")),
			Required:         true,
		},
	},
}

var labelDataSource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"key": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"value": {
			Type:     schema.TypeString,
			Computed: true,
		},
	},
}

type Label map[string]string

type LabelsService struct {
	ctx    context.Context
	client *ionoscloud.APIClient
}

func (ls *LabelsService) datacentersServersLabelsGet(datacenterID, serverID string, isDataSource bool) ([]Label, error) {
	labelsResponse, apiResponse, err := ls.client.LabelsApi.DatacentersServersLabelsGet(ls.ctx, datacenterID, serverID).Depth(1).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return nil, fmt.Errorf("error occurred while fetching labels for server with ID: %s, datacenter ID: %s, error: %w", serverID, datacenterID, err)
	}
	labels, err := processLabelsData(labelsResponse, isDataSource)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

func (ls *LabelsService) datacentersServersLabelsCreate(datacenterID, serverID string, labelsData any) error {
	if labelsData, ok := labelsData.(*schema.Set); ok {
		for _, labelData := range labelsData.List() {
			if label, ok := labelData.(map[string]any); ok {
				labelKey := label["key"].(string)
				labelValue := label["value"].(string)
				labelResource := ionoscloud.LabelResource{
					Properties: ionoscloud.LabelResourceProperties{Key: &labelKey, Value: &labelValue},
				}
				_, apiResponse, err := ls.client.LabelsApi.DatacentersServersLabelsPost(ls.ctx, datacenterID, serverID).Label(labelResource).Execute()
				apiResponse.LogInfo()
				if err != nil {
					return fmt.Errorf("error occurred while creating label for server with ID: %s, datacenter ID: %s, error: (%w)", serverID, datacenterID, err)
				}
			}
		}
	}
	return nil
}

func (ls *LabelsService) datacentersServersLabelsDelete(datacenterID, serverID string, labelsData any) error {
	if labelsData, ok := labelsData.(*schema.Set); ok {
		for _, labelData := range labelsData.List() {
			if label, ok := labelData.(map[string]any); ok {
				labelKey := label["key"].(string)
				apiResponse, err := ls.client.LabelsApi.DatacentersServersLabelsDelete(ls.ctx, datacenterID, serverID, labelKey).Execute()
				apiResponse.LogInfo()
				if err != nil {
					if httpNotFound(apiResponse) {
						tflog.Warn(ls.ctx, "label has been already removed from server", map[string]any{"key": labelKey, "server_id": serverID})
					} else {
						return fmt.Errorf("[label update] an error occurred while deleting label with key: %s, server ID: %s, error: %w", labelKey, serverID, err)
					}
				}
			}
		}
	}
	return nil
}

// Process the labels data fetched using the API and convert it a list of labels that can be
// used to set the resource data.
func processLabelsData(labelsData ionoscloud.LabelResources, isDataSource bool) ([]Label, error) {
	labels := make([]Label, 0, len(labelsData.Items))
	for _, labelData := range labelsData.Items {
		entry := make(Label)
		if labelData.Properties.Key == nil || labelData.Properties.Value == nil {
			return nil, errors.New("expected valid label properties from the API but received nil instead")
		}
		entry["key"] = *labelData.Properties.Key
		entry["value"] = *labelData.Properties.Value
		if isDataSource {
			if labelData.Id == nil {
				return nil, errors.New("expected valid label ID from the API but received nil instead")
			}
			entry["id"] = *labelData.Id
		}
		labels = append(labels, entry)
	}
	return labels, nil
}
