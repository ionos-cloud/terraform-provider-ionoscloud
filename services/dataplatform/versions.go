package dataplatform

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dataplatform "github.com/ionos-cloud/sdk-go-dataplatform"
)

type VersionService interface {
	GetVersions(ctx context.Context) ([]string, *dataplatform.APIResponse, error)
}

func (c *Client) GetVersions(ctx context.Context) ([]string, *dataplatform.APIResponse, error) {
	versions, apiResponse, err := c.DataPlatformMetaDataApi.VersionsGet(ctx).Execute()
	if apiResponse != nil {
		return versions, apiResponse, err

	}
	return versions, nil, err
}

func SetVersionsData(d *schema.ResourceData, versions []string) diag.Diagnostics {

	if versions != nil {
		err := d.Set("versions", versions)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting Dataplatform API version: %s", err))
			return diags
		}
	}

	resourceId := uuid.New()
	d.SetId(resourceId.String())

	return nil
}
