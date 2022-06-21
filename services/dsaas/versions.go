package dsaas

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dsaas "github.com/ionos-cloud/sdk-go-autoscaling"
)

type VersionService interface {
	GetVersions(ctx context.Context) ([]string, *dsaas.APIResponse, error)
}

func (c *Client) GetVersions(ctx context.Context) ([]string, *dsaas.APIResponse, error) {
	versions, apiResponse, err := c.DataPlatformMetaDataApi.VersionsGet(ctx).Execute()
	if apiResponse != nil {
		return versions, apiResponse, err

	}
	return versions, nil, err
}

func SetVersionsData(d *schema.ResourceData, versions []string) diag.Diagnostics {

	if versions != nil {
		err := d.Set("postgres_versions", versions)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting Data Stack API version: %s", err))
			return diags
		}
	}

	resourceId := uuid.New()
	d.SetId(resourceId.String())

	return nil
}
