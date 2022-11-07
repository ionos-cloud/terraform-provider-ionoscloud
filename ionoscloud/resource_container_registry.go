package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	crService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/container-registry"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"log"
	"regexp"
	"time"
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
		CustomizeDiff: checkContainerRegistryImmutableFields,
		Schema: map[string]*schema.Schema{
			"garbage_collection_schedule": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"days": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.All(validation.IsDayOfTheWeek(true)),
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringMatch(regexp.MustCompile("^[a-z][-a-z0-9]{1,61}[a-z0-9]$"), "")),
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
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}
func checkContainerRegistryImmutableFields(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {

	//we do not want to check in case of resource creation
	if diff.Id() == "" {
		return nil
	}
	if diff.HasChange("name") {
		return fmt.Errorf("name %s", ImmutableError)
	}
	if diff.HasChange("location") {
		return fmt.Errorf("location %s", ImmutableError)
	}
	return nil

}
func resourceContainerRegistryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).ContainerClient

	containerRegistry := crService.GetRegistryDataCreate(d)

	containerRegistryResponse, _, err := client.CreateRegistry(ctx, *containerRegistry)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while creating the registry: %w", err))
		return diags
	}

	d.SetId(*containerRegistryResponse.Id)

	return resourceContainerRegistryRead(ctx, d, meta)
}

func resourceContainerRegistryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).ContainerClient

	registry, apiResponse, err := client.GetRegistry(ctx, d.Id())

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching registry %s: %w", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] Successfully retreived registry %s: %+v", d.Id(), registry)

	if err := crService.SetRegistryData(d, registry); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceContainerRegistryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).ContainerClient

	containerRegistry := crService.GetRegistryDataUpdate(d)

	registryId := d.Id()

	containerRegistryResponse, _, err := client.PatchRegistry(ctx, registryId, *containerRegistry)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating a registry: %s", err))
		return diags
	}

	d.SetId(*containerRegistryResponse.Id)

	time.Sleep(utils.SleepInterval)

	return resourceContainerRegistryRead(ctx, d, meta)
}

func resourceContainerRegistryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).ContainerClient

	registryId := d.Id()

	apiResponse, err := client.DeleteRegistry(ctx, registryId)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting registry %s: %s", registryId, err))
		return diags
	}

	for {
		log.Printf("[INFO] Waiting for cluster %s to be deleted...", d.Id())

		registryDeleted, dsErr := registryDeleted(ctx, client, d)

		if dsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking deletion status of registry %s: %s", d.Id(), dsErr))
			return diags
		}

		if registryDeleted {
			log.Printf("[INFO] Successfully deleted registry: %s", d.Id())
			break
		}

		select {
		case <-time.After(utils.SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			diags := diag.FromErr(fmt.Errorf("registry deletion timed out! WARNING: your k8s cluster (%s) will still probably be deleted after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates", d.Id()))
			return diags
		}
	}

	return nil
}

func resourceContainerRegistryImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(SdkBundle).ContainerClient

	registryId := d.Id()

	containerRegistry, apiResponse, err := client.GetRegistry(ctx, registryId)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("registry does not exist %q", registryId)
		}
		return nil, fmt.Errorf("an error occured while trying to fetch the import of registry %q", registryId)
	}

	log.Printf("[INFO] registry found: %+v", containerRegistry)

	if err := crService.SetRegistryData(d, containerRegistry); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func registryDeleted(ctx context.Context, client *crService.Client, d *schema.ResourceData) (bool, error) {

	_, apiResponse, err := client.GetRegistry(ctx, d.Id())

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			return true, nil
		}
		return true, fmt.Errorf("error checking registry deletion status: %s", err)
	}
	return false, nil
}
