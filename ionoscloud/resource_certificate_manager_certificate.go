package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cert"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func resourceCertificateManager() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCertificateManagerCreate,
		ReadContext:   resourceCertificateManagerRead,
		UpdateContext: resourceCertificateManagerUpdate,
		DeleteContext: resourceCertificateManagerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCertificateManagerImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Description:      "The certificate name",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"certificate": {
				Type:                  schema.TypeString,
				Description:           "The certificate body in PEM format. This attribute is immutable.",
				Required:              true,
				ValidateDiagFunc:      validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				DiffSuppressFunc:      utils.DiffWithoutNewLines,
				DiffSuppressOnRefresh: true,
			},
			"certificate_chain": {
				Type:                  schema.TypeString,
				Description:           "The certificate chain. This attribute is immutable.",
				Optional:              true,
				ValidateDiagFunc:      validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
				DiffSuppressFunc:      utils.DiffWithoutNewLines,
				DiffSuppressOnRefresh: true,
			},
			"private_key": {
				Type:             schema.TypeString,
				Description:      "The private key blob. This attribute is immutable.",
				Required:         true,
				Sensitive:        true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
		},
		CustomizeDiff: checkCertImmutableFields,
		Timeouts:      &resourceDefaultTimeouts,
	}
}

func checkCertImmutableFields(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {

	// we do not want to check in case of resource creation
	if diff.Id() == "" {
		return nil
	}

	if diff.HasChange("certificate") {
		oldV, newV := diff.GetChange("certificate")
		old := utils.RemoveNewLines(oldV.(string))
		newStr := utils.RemoveNewLines(newV.(string))
		// we get extraneous newlines in the certificate, so we must remove them before checking equality
		if !strings.EqualFold(old, newStr) {
			return fmt.Errorf("certificate %s", ImmutableError)
		}
	}

	if diff.HasChange("certificate_chain") {
		oldV, newV := diff.GetChange("certificate_chain")
		old := utils.RemoveNewLines(oldV.(string))
		newStr := utils.RemoveNewLines(newV.(string))
		if !strings.EqualFold(old, newStr) {
			return fmt.Errorf("certificate_chain %s", ImmutableError)
		}
	}

	if diff.HasChange("private_key") {
		return fmt.Errorf("private_key %s", ImmutableError)
	}

	return nil

}

func resourceCertificateManagerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CertManagerClient

	certPostDto, err := cert.GetCertPostDto(d)
	if err != nil {
		return diag.FromErr(err)
	}
	certificateDto, _, err := client.CreateCertificate(ctx, *certPostDto)
	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating certificate: %w", err))
		return diags
	}

	d.SetId(certificateDto.Id)

	if err = utils.WaitForResourceToBeReady(ctx, d, client.IsCertReady); err != nil {
		return diag.FromErr(err)
	}

	return resourceCertificateManagerRead(ctx, d, meta)
}

func resourceCertificateManagerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CertManagerClient

	certDto, apiResponse, err := client.GetCertificate(ctx, d.Id())
	if err != nil {
		if apiResponse.HttpNotFound() {
			log.Printf("[INFO] Resource %s not found: %+v", d.Id(), err)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully retrieved certificate %s: %+v", d.Id(), certDto)

	if err := cert.SetCertificateData(d, &certDto); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceCertificateManagerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CertManagerClient

	certPatchDto := cert.GetCertPatchDto(d)

	_, _, err := client.UpdateCertificate(ctx, d.Id(), *certPatchDto)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while updating certificate with ID %s, %w", d.Id(), err))
		return diags
	}

	if err = utils.WaitForResourceToBeReady(ctx, d, client.IsCertReady); err != nil {
		return diag.FromErr(err)
	}

	return resourceCertificateManagerRead(ctx, d, meta)
}

func resourceCertificateManagerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CertManagerClient

	_, err := client.DeleteCertificate(ctx, d.Id())
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occurred while deleting the certificate %s %w", d.Id(), err))
		return diags
	}

	err = utils.WaitForResourceToBeDeleted(ctx, d, client.IsCertDeleted)
	if err != nil {
		return diag.FromErr(fmt.Errorf("deleting %w", err))
	}

	log.Printf("[INFO] Successfully deleted certificate: %s", d.Id())

	d.SetId("")

	return nil
}

func resourceCertificateManagerImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(bundleclient.SdkBundle).CertManagerClient

	certId := d.Id()
	certDto, apiResponse, err := client.GetCertificate(ctx, d.Id())
	if err != nil {
		if apiResponse.HttpNotFound() {
			d.SetId("")
			return nil, fmt.Errorf("unable to find cert %q", certId)
		}
		return nil, fmt.Errorf("an error occurred while retrieving the cert %q, %w", certId, err)
	}

	if err := cert.SetCertificateData(d, &certDto); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
