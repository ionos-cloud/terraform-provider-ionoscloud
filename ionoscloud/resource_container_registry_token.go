package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	crService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/container-registry"
	"log"
	"regexp"
	"time"
)

func resourceContainerRegistryToken() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContainerRegistryTokenCreate,
		ReadContext:   resourceContainerRegistryTokenRead,
		UpdateContext: resourceContainerRegistryTokenUpdate,
		DeleteContext: resourceContainerRegistryTokenDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceContainerRegistryTokenImport,
		},
		CustomizeDiff: checkContainerRegistryTokenImmutableFields,
		Schema: map[string]*schema.Schema{
			"credentials": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Required: true,
						},
						"password": {
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
			"expiry_date": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringMatch(regexp.MustCompile("^[a-z][-a-z0-9]{1,61}[a-z0-9]$"), "")),
			},
			"scopes": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"actions": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.All(validation.StringInSlice([]string{"enabled", "disabled"}, true)),
			},
			"registry_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}
func checkContainerRegistryTokenImmutableFields(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {

	//we do not want to check in case of resource creation
	if diff.Id() == "" {
		return nil
	}
	if diff.HasChange("name") {
		return fmt.Errorf("name %s", ImmutableError)
	}
	return nil

}
func resourceContainerRegistryTokenCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).ContainerClient

	registryId := d.Get("registry_id").(string)

	registryToken, err := crService.GetTokenDataCreate(d)

	if err != nil {
		return diag.FromErr(err)
	}
	registryTokenResponse, _, err := client.CreateTokens(ctx, registryId, *registryToken)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while creating the registry token: %w", err))
		return diags
	}

	d.SetId(*registryTokenResponse.Id)

	for {
		log.Printf("[INFO] Waiting for registry token %s to be ready...", d.Id())

		registryTokenReady, rsErr := registryTokenReady(ctx, client, d)

		if rsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking readiness status of registry token %s: %w", d.Id(), rsErr))
			return diags
		}

		if registryTokenReady {
			log.Printf("[INFO] registry token ready: %s", d.Id())
			break
		}

		select {
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			log.Printf("[INFO] create timed out")
			diags := diag.FromErr(fmt.Errorf("registry creation timed out! WARNING: your registry token (%s) will still probably be created after some time but the terraform state wont reflect that; check your Ionos Cloud account for updates", d.Id()))
			return diags
		}

	}

	return resourceContainerRegistryTokenRead(ctx, d, meta)
}

func resourceContainerRegistryTokenRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).ContainerClient

	registryId := d.Get("registry_id").(string)
	registryTokenId := d.Id()
	registryToken, apiResponse, err := client.GetToken(ctx, registryId, registryTokenId)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching registry token %s: %w", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] Successfully retreived registry token %s: %+v", d.Id(), registryToken)

	if err := crService.SetTokenData(d, registryToken); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceContainerRegistryTokenUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).ContainerClient

	registryId := d.Get("registry_id").(string)
	registryTokenId := d.Id()
	registryToken, err := crService.GetTokenDataUpdate(d)

	if err != nil {
		return diag.FromErr(err)
	}

	registryTokenResponse, _, err := client.PatchToken(ctx, registryId, registryTokenId, *registryToken)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating a registry token: %s", err))
		return diags
	}

	d.SetId(*registryTokenResponse.Id)

	time.Sleep(SleepInterval)

	for {
		log.Printf("[INFO] Waiting for cluster %s to be ready...", d.Id())

		clusterReady, rsErr := registryTokenReady(ctx, client, d)

		if rsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking readiness status of registry token %s: %w", d.Id(), rsErr))
			return diags
		}

		if clusterReady {
			log.Printf("[INFO]registry token ready: %s", d.Id())
			break
		}

		select {
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			log.Printf("[INFO] create timed out")
			diags := diag.FromErr(fmt.Errorf("registry update timed out! WARNING: your registry token (%s) will still probably be updated after some time but the terraform state wont reflect that; check your Ionos Cloud account for updates", d.Id()))
			return diags
		}

	}

	return resourceContainerRegistryTokenRead(ctx, d, meta)
}

func resourceContainerRegistryTokenDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).ContainerClient

	registryId := d.Get("registry_id").(string)
	registryTokenId := d.Id()

	apiResponse, err := client.DeleteToken(ctx, registryId, registryTokenId)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting registry token %s: %s", registryTokenId, err))
		return diags
	}

	for {
		log.Printf("[INFO] Waiting for cluster %s to be deleted...", d.Id())

		registryTokenDeleted, dsErr := registryTokenDeleted(ctx, client, d)

		if dsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking deletion status of registry token %s: %s", d.Id(), dsErr))
			return diags
		}

		if registryTokenDeleted {
			log.Printf("[INFO] Successfully deleted registry token: %s", d.Id())
			break
		}

		select {
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			diags := diag.FromErr(fmt.Errorf("registry deletion timed out! WARNING: your k8s cluster (%s) will still probably be deleted after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates", d.Id()))
			return diags
		}
	}

	// wait 15 seconds after the deletion of the cluster, for the lan to be freed
	time.Sleep(SleepInterval * 3)

	return nil
}

func resourceContainerRegistryTokenImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(SdkBundle).ContainerClient

	registryId := d.Get("registry_id").(string)
	registryTokenId := d.Id()

	registryToken, apiResponse, err := client.GetToken(ctx, registryId, registryTokenId)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("registry does not exist %q", registryTokenId)
		}
		return nil, fmt.Errorf("an error occured while trying to fetch the import of registry token %q", registryTokenId)
	}

	log.Printf("[INFO] registry token found: %+v", registryToken)

	if err := crService.SetTokenData(d, registryToken); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func registryTokenReady(ctx context.Context, client *crService.Client, d *schema.ResourceData) (bool, error) {
	registryId := d.Get("registry_id").(string)
	registryTokenId := d.Id()

	subjectToken, _, err := client.GetToken(ctx, registryId, registryTokenId)

	if err != nil {
		return true, fmt.Errorf("error checking registry token status: %s", err)
	}

	if *subjectToken.Metadata.State == "FAILED" {

		time.Sleep(time.Second * 3)

		subjectToken, _, err = client.GetToken(ctx, registryId, registryTokenId)

		if err != nil {
			return true, fmt.Errorf("error checking registry token status: %s", err)
		}

		if *subjectToken.Metadata.State == "FAILED" {
			return false, fmt.Errorf("registry has failed. WARNING: your registry token may still recover after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates")
		}
	}
	return *subjectToken.Metadata.State == "AVAILABLE", nil
}

func registryTokenDeleted(ctx context.Context, client *crService.Client, d *schema.ResourceData) (bool, error) {

	registryId := d.Get("registry_id").(string)
	registryTokenId := d.Id()

	_, apiResponse, err := client.GetToken(ctx, registryId, registryTokenId)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			return true, nil
		}
		return true, fmt.Errorf("error checking registry token deletion status: %s", err)
	}
	return false, nil
}
