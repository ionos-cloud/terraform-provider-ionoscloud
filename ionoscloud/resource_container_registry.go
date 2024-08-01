package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
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
	client := meta.(services.SdkBundle).ContainerClient

	containerRegistry, err := crService.GetRegistryDataCreate(d)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error occurred while getting container registry from schema: %w", err))
		return diags
	}

	containerRegistryFeatures, warnings := crService.GetRegistryFeatures(d)
	containerRegistry.Properties.Features = containerRegistryFeatures

	containerRegistryResponse, _, err := client.CreateRegistry(ctx, *containerRegistry)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while creating the registry: %w", err))
		return diags
	}
	d.SetId(*containerRegistryResponse.Id)

	if err := utils.WaitForResourceToBeReady(ctx, d, client.IsRegistryReady); err != nil {
		return diag.FromErr(fmt.Errorf("error waiting for registry to be ready: %w", err))
	}
	return append(warnings, resourceContainerRegistryRead(ctx, d, meta)...)
}

func resourceContainerRegistryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(services.SdkBundle).ContainerClient

	registry, apiResponse, err := client.GetRegistry(ctx, d.Id())
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching registry %s: %w", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] Successfully retrieved registry %s: %+v", d.Id(), registry)

	if err := crService.SetRegistryData(d, registry); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceContainerRegistryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).ContainerClient

	containerRegistry, err := crService.GetRegistryDataUpdate(d)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error occurred while getting container registry from schema: %w", err))
		return diags
	}
	containerRegistryFeatures, warnings := crService.GetRegistryFeatures(d)
	containerRegistry.Features = containerRegistryFeatures
	// suppress warnings if there are no changes to the features set
	if !d.HasChange("features") {
		warnings = diag.Diagnostics{}
	}

	registryId := d.Id()

	_, _, err = client.PatchRegistry(ctx, registryId, *containerRegistry)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while updating a registry: %w", err))
		return diags
	}

	if err := utils.WaitForResourceToBeReady(ctx, d, client.IsRegistryReady); err != nil {
		return diag.FromErr(fmt.Errorf("error waiting for registry to be ready: %w", err))
	}

	return append(warnings, resourceContainerRegistryRead(ctx, d, meta)...)
}

func resourceContainerRegistryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(services.SdkBundle).ContainerClient

	registryId := d.Id()

	apiResponse, err := client.DeleteRegistry(ctx, registryId)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting registry %s: %w", registryId, err))
		return diags
	}

	return diag.FromErr(utils.WaitForResourceToBeDeleted(ctx, d, client.IsRegistryDeleted))
}

func resourceContainerRegistryImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(services.SdkBundle).ContainerClient

	registryId := d.Id()

	containerRegistry, apiResponse, err := client.GetRegistry(ctx, registryId)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("registry does not exist %q", registryId)
		}
		return nil, fmt.Errorf("an error occurred while trying to fetch the import of registry %q, error:%w", registryId, err)
	}

	log.Printf("[INFO] registry found: %+v", containerRegistry)

	if err := crService.SetRegistryData(d, containerRegistry); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
