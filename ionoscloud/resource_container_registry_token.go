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
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
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
		return utils.ToDiags(d, err.Error(), nil)
	}
	registryTokenResponse, apiResponse, err := client.CreateToken(ctx, registryId, *registryToken)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while creating the registry token: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	d.SetId(*registryTokenResponse.Id)

	if fileStr != "" {
		if err := utils.WriteToFile(fileStr, registryTokenResponse.Properties.Credentials.Password); err != nil {
			return utils.ToDiags(d, err.Error(), nil)
		}
	}

	if err = crService.SetTokenData(d, registryTokenResponse.Properties); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	var credentials []any
	credentialsEntry := crService.SetCredentials(registryTokenResponse.Properties.Credentials)
	credentials = append(credentials, credentialsEntry)
	if err := d.Set("credentials", credentials); err != nil {
		return utils.ToDiags(d, utils.GenerateSetError("token", "credentials", err).Error(), nil)
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
		return utils.ToDiags(d, fmt.Sprintf("error while fetching registry token: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	log.Printf("[INFO] Successfully retrieved registry token %s: %+v", d.Id(), registryToken)

	if err := crService.SetTokenData(d, registryToken.Properties); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	return nil
}

func resourceContainerRegistryTokenUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).ContainerClient

	registryId := d.Get("registry_id").(string)
	registryTokenId := d.Id()
	registryToken, err := crService.GetTokenDataUpdate(d)
	if err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	_, apiResponse, err := client.PatchToken(ctx, registryId, registryTokenId, *registryToken)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while updating a registry token: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
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
		return utils.ToDiags(d, fmt.Sprintf("error while deleting registry token %s: %s", registryTokenId, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
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
			return nil, utils.ToError(d, fmt.Sprintf("registry does not exist %q", registryTokenId), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
		return nil, utils.ToError(d, fmt.Sprintf("an error occurred while trying to fetch the import of registry token %q, error:%s", registryTokenId, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}

	log.Printf("[INFO] registry token found: %+v", registryToken)

	if registryToken.Id != nil {
		d.SetId(*registryToken.Id)
	}

	if err := crService.SetTokenData(d, registryToken.Properties); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}

	return []*schema.ResourceData{d}, nil
}
