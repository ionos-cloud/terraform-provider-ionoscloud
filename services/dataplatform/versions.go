package dataplatform

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dataplatform "github.com/ionos-cloud/sdk-go-dataplatform"
)

func (c *Client) GetVersions(ctx context.Context) ([]string, *dataplatform.APIResponse, error) {
	versions, apiResponse, err := c.sdkClient.DataPlatformMetaDataApi.VersionsGet(ctx).Execute()
	apiResponse.LogInfo()
	if versions.Items == nil {
		return nil, nil, fmt.Errorf("expected a list of Dataplatform versions but received 'nil' instead")
	}
	return *versions.Items, apiResponse, err
}

func SetVersionsData(d *schema.ResourceData, versions []string) diag.Diagnostics {

	if versions != nil {
		err := d.Set("versions", versions)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting Dataplatform API version: %w", err))
			return diags
		}
	}

	resourceId := uuid.New()
	d.SetId(resourceId.String())

	return nil
}
