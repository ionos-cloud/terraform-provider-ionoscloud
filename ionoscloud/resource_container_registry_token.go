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
	crService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/containerregistry"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
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
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"expiry_date": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: DiffExpiryDate,
			},
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile("^[a-z][-a-z0-9]{1,61}[a-z0-9]$"), "")),
				ForceNew:         true,
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
							Description: "Example: [\"pull\", \"push\", \"delete\"]",
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
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"enabled", "disabled"}, true)),
				Description:      "Can be one of enabled, disabled",
			},
			"registry_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile("^[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}$"), "")),
			},
			"save_password_to_file": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Saves password to file. Only works on create. Takes as argument a file name, or a file path",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceContainerRegistryTokenCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).ContainerClient

	registryId := d.Get("registry_id").(string)
	fileStr := d.Get("save_password_to_file").(string)
	registryToken, err := crService.GetTokenDataCreate(d)

	if err != nil {
		return diag.FromErr(err)
	}
	registryTokenResponse, _, err := client.CreateTokens(ctx, registryId, *registryToken)
	if err != nil {
		return diag.FromErr(fmt.Errorf("an error occurred while creating the registry token: %w", err))
	}

	d.SetId(*registryTokenResponse.Id)

	if registryTokenResponse.Properties == nil {
		return diag.FromErr(fmt.Errorf("no token properties found with the specified id = %s", *registryTokenResponse.Id))
	}

	if fileStr != "" {
		if err := utils.WriteToFile(fileStr, *registryTokenResponse.Properties.Credentials.Password); err != nil {
			return diag.FromErr(err)
		}
	}

	if err = crService.SetTokenData(d, *registryTokenResponse.Properties); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceContainerRegistryTokenRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(bundleclient.SdkBundle).ContainerClient

	registryId := d.Get("registry_id").(string)
	registryTokenId := d.Id()
	registryToken, apiResponse, err := client.GetToken(ctx, registryId, registryTokenId)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching registry token %s: %w", d.Id(), err))
		return diags
	}

	log.Printf("[INFO] Successfully retrieved registry token %s: %+v", d.Id(), registryToken)

	if registryToken.Properties == nil {
		return diag.FromErr(fmt.Errorf("no token properties found with the specified id = %s", *registryToken.Id))
	}

	if err := crService.SetTokenData(d, *registryToken.Properties); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceContainerRegistryTokenUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).ContainerClient

	registryId := d.Get("registry_id").(string)
	registryTokenId := d.Id()
	registryToken, err := crService.GetTokenDataUpdate(d)
	if err != nil {
		return diag.FromErr(err)
	}

	_, _, err = client.PatchToken(ctx, registryId, registryTokenId, *registryToken)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while updating a registry token: %w", err))
		return diags
	}

	return resourceContainerRegistryTokenRead(ctx, d, meta)
}

func resourceContainerRegistryTokenDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).ContainerClient

	registryId := d.Get("registry_id").(string)
	registryTokenId := d.Id()

	apiResponse, err := client.DeleteToken(ctx, registryId, registryTokenId)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while deleting registry token %s: %w", registryTokenId, err))
		return diags
	}

	return nil
}

func resourceContainerRegistryTokenImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(bundleclient.SdkBundle).ContainerClient

	registryId := d.Get("registry_id").(string)
	registryTokenId := d.Id()

	registryToken, apiResponse, err := client.GetToken(ctx, registryId, registryTokenId)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("registry does not exist %q", registryTokenId)
		}
		return nil, fmt.Errorf("an error occurred while trying to fetch the import of registry token %q, error:%w", registryTokenId, err)
	}

	log.Printf("[INFO] registry token found: %+v", registryToken)

	if registryToken.Id != nil {
		d.SetId(*registryToken.Id)
	}
	if registryToken.Properties == nil {
		return nil, fmt.Errorf("no token properties found with the specified id = %s", *registryToken.Id)
	}

	if err := crService.SetTokenData(d, *registryToken.Properties); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
