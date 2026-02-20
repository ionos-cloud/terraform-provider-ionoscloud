package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"

	crService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/containerregistry"
)

func resourceContainerRegistry() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContainerRegistryCreate,
		ReadContext:   resourceContainerRegistryRead,
		UpdateContext: resourceContainerRegistryUpdate,
		DeleteContext: resourceContainerRegistryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceContainerRegistryImport,
		},
		Schema: map[string]*schema.Schema{
			"garbage_collection_schedule": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "UTC time of day e.g. 01:00:00 - as defined by partial-time - RFC3339",
						},
						"days": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validation.ToDiagFunc(validation.IsDayOfTheWeek(true)),
							},
						},
					},
				},
			},
			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				ForceNew:         true,
			},
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile("^[a-z][-a-z0-9]{1,61}[a-z0-9]$"), "")),
				ForceNew:         true,
			},
			"api_subnet_allow_list": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "The subnet CIDRs that are allowed to connect to the registry. Specify 'a.b.c.d/32' for an individual IP address. __Note__: If this list is empty or not set, there are no restrictions.",
			},
			"storage_usage": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bytes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"features": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vulnerability_scanning": {
							Type:        schema.TypeBool,
							Description: "Enables vulnerability scanning for images in the container registry. Note: this feature can incur additional charges",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceContainerRegistryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).ContainerClient

	containerRegistry, err := crService.GetRegistryDataCreate(d)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error occurred while getting container registry from schema: %s", err), nil)
	}

	containerRegistryFeatures, warnings := crService.GetRegistryFeatures(d)
	containerRegistry.Properties.Features = containerRegistryFeatures

	containerRegistryResponse, apiResponse, err := client.CreateRegistry(ctx, *containerRegistry)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while creating the registry: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	d.SetId(*containerRegistryResponse.Id)

	if err := utils.WaitForResourceToBeReady(ctx, d, client.IsRegistryReady); err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error waiting for registry to be ready: %s", err), nil)
	}
	return append(warnings, resourceContainerRegistryRead(ctx, d, meta)...)
}

func resourceContainerRegistryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(bundleclient.SdkBundle).ContainerClient

	registry, apiResponse, err := client.GetRegistry(ctx, d.Id())
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return utils.ToDiags(d, fmt.Sprintf("error while fetching registry: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	log.Printf("[INFO] Successfully retrieved registry %s: %+v", d.Id(), registry)

	if err := crService.SetRegistryData(d, registry); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	return nil
}

func resourceContainerRegistryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).ContainerClient

	containerRegistry, err := crService.GetRegistryDataUpdate(d)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error occurred while getting container registry from schema: %s", err), nil)
	}
	containerRegistryFeatures, warnings := crService.GetRegistryFeatures(d)
	containerRegistry.Features = containerRegistryFeatures
	// suppress warnings if there are no changes to the features set
	if !d.HasChange("features") {
		warnings = diag.Diagnostics{}
	}

	registryId := d.Id()

	_, apiResponse, err := client.PatchRegistry(ctx, registryId, *containerRegistry)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while updating a registry: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	if err := utils.WaitForResourceToBeReady(ctx, d, client.IsRegistryReady); err != nil {
		return utils.ToDiags(d, fmt.Sprintf("error waiting for registry to be ready: %s", err), nil)
	}

	return append(warnings, resourceContainerRegistryRead(ctx, d, meta)...)
}

func resourceContainerRegistryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).ContainerClient

	registryId := d.Id()

	apiResponse, err := client.DeleteRegistry(ctx, registryId)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return utils.ToDiags(d, fmt.Sprintf("error while deleting registry %s: %s", registryId, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	return utils.ToDiags(d, utils.WaitForResourceToBeDeleted(ctx, d, client.IsRegistryDeleted).Error(), &utils.DiagsOpts{Timeout: schema.TimeoutDelete})
}

func resourceContainerRegistryImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(bundleclient.SdkBundle).ContainerClient

	registryId := d.Id()

	containerRegistry, apiResponse, err := client.GetRegistry(ctx, registryId)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, utils.ToError(d, fmt.Sprintf("registry does not exist %q", registryId), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
		return nil, utils.ToError(d, fmt.Sprintf("an error occurred while trying to fetch the import of registry %q, error:%s", registryId, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	log.Printf("[INFO] registry found: %+v", containerRegistry)

	if err := crService.SetRegistryData(d, containerRegistry); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}

	return []*schema.ResourceData{d}, nil
}
