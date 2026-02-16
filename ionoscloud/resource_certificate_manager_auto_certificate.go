package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cert"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func resourceCertificateManagerAutoCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: autoCertificateCreate,
		ReadContext:   autoCertificateRead,
		UpdateContext: autoCertificateUpdate,
		DeleteContext: autoCertificateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: autoCertificateImport,
		},
		Schema: map[string]*schema.Schema{
			"provider_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The certificate provider used to issue the certificates",
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The location of the auto-certificate",
			},
			"common_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The common name (DNS) of the certificate to issue. The common name needs to be part of a zone in IONOS Cloud DNS",
			},
			"key_algorithm": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The key algorithm used to generate the certificate",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A certificate name used for management purposes",
			},
			"subject_alternative_names": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Optional additional names to be added to the issued certificate. The additional names needs to be part of a zone in IONOS Cloud DNS",
			},
			"last_issued_certificate_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the last certificate that was issued",
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func autoCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CertManagerClient
	location := d.Get("location").(string)

	autoCertificateCreateData := cert.GetAutoCertificateDataCreate(d)
	response, apiResponse, err := client.CreateAutoCertificate(ctx, location, *autoCertificateCreateData)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while creating an auto-certificate: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	autoCertificateID := response.Id
	d.SetId(autoCertificateID)

	err = utils.WaitForResourceToBeReady(ctx, d, client.IsAutoCertificateReady)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while checking the creation status for the auto-certificate with ID: %v, error: %s", autoCertificateID, err), nil)
	}
	// Return with another read call because 'last_issued_certificate_id' is not provided in the
	// creation response.
	return autoCertificateRead(ctx, d, meta)
}

func autoCertificateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CertManagerClient
	autoCertificateID := d.Id()
	location := d.Get("location").(string)
	autoCertificate, apiResponse, err := client.GetAutoCertificate(ctx, autoCertificateID, location)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return utils.ToDiags(d, fmt.Sprintf("error while fetching auto-certificate: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	log.Printf("[INFO] Successfully retrieved auto-certificate with ID: %v, auto-certificate info: %+v", autoCertificateID, autoCertificate)
	if err := cert.SetAutoCertificateData(d, autoCertificate); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}
	return nil
}

func autoCertificateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CertManagerClient
	autoCertificateID := d.Id()
	location := d.Get("location").(string)

	autoCertificateUpdateData := cert.GetAutoCertificateDataUpdate(d)
	autoCertificate, apiResponse, err := client.UpdateAutoCertificate(ctx, autoCertificateID, location, *autoCertificateUpdateData)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while updating auto-certificate: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	if err = utils.WaitForResourceToBeReady(ctx, d, client.IsAutoCertificateReady); err != nil {
		return utils.ToDiags(d, fmt.Sprintf("an error occurred while checking the update status for the auto-certificate: %s", err), nil)
	}
	if err = cert.SetAutoCertificateData(d, autoCertificate); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}
	return nil
}

func autoCertificateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CertManagerClient
	autoCertificateID := d.Id()
	location := d.Get("location").(string)
	apiResponse, err := client.DeleteAutoCertificate(ctx, autoCertificateID, location)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil
		}
		return utils.ToDiags(d, fmt.Sprintf("error while deleting auto-certificate: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsAutoCertificateDeleted)
	if err != nil {
		return utils.ToDiags(d, fmt.Sprintf("deletion check failed for auto-certificate: %s", err), &utils.DiagsOpts{Timeout: schema.TimeoutDelete})
	}
	return nil
}

func autoCertificateImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(bundleclient.SdkBundle).CertManagerClient
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return nil, utils.ToError(d, "invalid import, expected ID in the format: '<location>:<auto_certificate_ID>", nil)
	}
	location := parts[0]
	autoCertificateID := parts[1]
	autoCertificate, apiResponse, err := client.GetAutoCertificate(ctx, autoCertificateID, location)
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, utils.ToError(d, fmt.Sprintf("auto-certificate with ID: %v does not exist", autoCertificateID), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
		return nil, utils.ToError(d, fmt.Sprintf("an error occurred while trying to import auto-certificate with ID: %v, error: %s", autoCertificateID, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
	}
	log.Printf("[INFO] auto-certificate found: %+v", autoCertificate)
	if err := d.Set("location", location); err != nil {
		return nil, utils.GenerateSetError("Auto-certificate", "location", err)
	}
	if err := cert.SetAutoCertificateData(d, autoCertificate); err != nil {
		return nil, utils.ToError(d, err.Error(), nil)
	}
	return []*schema.ResourceData{d}, nil
}
