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
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
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
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile("^[a-z][-a-z0-9]{1,61}[a-z0-9]$"), "must start with a lowercase letter, contain only lowercase alphanumeric characters or hyphens, and be 3-63 characters long")),
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
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile("^[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}$"), "must be a valid UUID")),
			},
			"save_password_to_file": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Saves password to file. Only works on create. Takes as argument a file name, or a file path",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The location of the resource. This field should be used only if you are also using a file configuration and should not be configured otherwise.",
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceContainerRegistryTokenCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewContainerRegistryClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	registryId := d.Get("registry_id").(string)
	fileStr := d.Get("save_password_to_file").(string)
	registryToken, err := crService.GetTokenDataCreate(d)

	if err != nil {
		return diagutil.ToDiags(d, err, nil)
	}
	registryTokenResponse, apiResponse, err := client.CreateToken(ctx, registryId, *registryToken)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while creating the registry token: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	d.SetId(*registryTokenResponse.Id)

	if fileStr != "" {
		if err := utils.WriteToFile(fileStr, registryTokenResponse.Properties.Credentials.Password); err != nil {
			return diagutil.ToDiags(d, err, nil)
		}
	}

	if err = crService.SetTokenData(d, registryTokenResponse.Properties); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	var credentials []any
	credentialsEntry := crService.SetCredentials(registryTokenResponse.Properties.Credentials)
	credentials = append(credentials, credentialsEntry)
	if err := d.Set("credentials", credentials); err != nil {
		return diagutil.ToDiags(d, utils.GenerateSetError("token", "credentials", err), nil)
	}
	return nil
}

func resourceContainerRegistryTokenRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewContainerRegistryClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	registryId := d.Get("registry_id").(string)
	registryTokenId := d.Id()
	registryToken, apiResponse, err := client.GetToken(ctx, registryId, registryTokenId)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, fmt.Errorf("error while fetching registry token: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	log.Printf("[INFO] Successfully retrieved registry token %s: %+v", d.Id(), registryToken)

	if err := crService.SetTokenData(d, registryToken.Properties); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	return nil
}

func resourceContainerRegistryTokenUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewContainerRegistryClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	registryId := d.Get("registry_id").(string)
	registryTokenId := d.Id()
	registryToken, err := crService.GetTokenDataUpdate(d)
	if err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	_, apiResponse, err := client.PatchToken(ctx, registryId, registryTokenId, *registryToken)
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while updating a registry token: %w", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	return resourceContainerRegistryTokenRead(ctx, d, meta)
}

func resourceContainerRegistryTokenDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewContainerRegistryClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	registryId := d.Get("registry_id").(string)
	registryTokenId := d.Id()

	apiResponse, err := client.DeleteToken(ctx, registryId, registryTokenId)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, fmt.Errorf("error while deleting registry token %s: %w", registryTokenId, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	return nil
}

func resourceContainerRegistryTokenImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	importID := d.Id()
	location, parts := splitImportID(importID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import identifier: expected one of <location>:<registry-id>/<token-id> or <registry-id>/<token-id>, got: %s", importID)
	}

	if err := validateImportIDParts(parts); err != nil {
		return nil, fmt.Errorf("failed validating import identifier %q: %w", importID, err)
	}

	client, err := meta.(bundleclient.SdkBundle).NewContainerRegistryClient(location)
	if err != nil {
		return nil, err
	}

	registryId := parts[0]
	registryTokenId := parts[1]

	registryToken, apiResponse, err := client.GetToken(ctx, registryId, registryTokenId)

	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, diagutil.ToError(d, fmt.Errorf("registry does not exist %q", registryTokenId), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}
		return nil, diagutil.ToError(d, fmt.Errorf("an error occurred while trying to fetch the import of registry token %q, error:%w", registryTokenId, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
	}

	log.Printf("[INFO] registry token found: %+v", registryToken)

	if registryToken.Id != nil {
		d.SetId(*registryToken.Id)
	}

	err = d.Set("registry_id", registryId)
	if err != nil {
		return nil, err
	}

	err = d.Set("location", location)
	if err != nil {
		return nil, err
	}

	if err := crService.SetTokenData(d, registryToken.Properties); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}

	return []*schema.ResourceData{d}, nil
}
