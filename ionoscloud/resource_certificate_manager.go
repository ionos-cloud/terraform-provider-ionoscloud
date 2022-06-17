package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	certmanager "github.com/ionos-cloud/sdk-cert-go"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cert"
	"log"
	"time"
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
				Type:         schema.TypeString,
				Description:  "The certificate name",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"certificate": {
				Type:         schema.TypeString,
				Description:  "The certificate body in PEM format. This attribute is immutable.",
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"certificate_chain": {
				Type:         schema.TypeString,
				Description:  "The certificate chain. This attribute is immutable.",
				Optional:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"private_key": {
				Type:         schema.TypeString,
				Description:  "The private key blob. This attribute is immutable.",
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
		},
		CustomizeDiff: checkCertImmutableFields,
		Timeouts:      &resourceDefaultTimeouts,
	}
}

func checkCertImmutableFields(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {

	//we do not want to check in case of resource creation
	if diff.Id() == "" {
		return nil
	}

	if diff.HasChange("certificate") {
		return fmt.Errorf("certificate %s", ImmutableError)
	}

	if diff.HasChange("certificate_chain") {
		return fmt.Errorf("certificate %s", ImmutableError)
	}

	if diff.HasChange("private_key") {
		return fmt.Errorf("certificate %s", ImmutableError)
	}

	return nil

}

func resourceCertificateManagerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CertManagerClient

	certificatePostDto := certmanager.CertificatePostDto{
		Properties: &certmanager.CertificatePostPropertiesDto{},
	}

	if name, nameOk := d.GetOk("name"); nameOk {
		name := name.(string)
		certificatePostDto.Properties.Name = &name
	} else {
		diags := diag.FromErr(fmt.Errorf("name must be provided for certificate manager"))
		return diags
	}

	if certField, certOk := d.GetOk("certificate"); certOk {
		certificate := certField.(string)
		certificatePostDto.Properties.Certificate = &certificate
	} else {
		diags := diag.FromErr(fmt.Errorf("certificate must be provided for certificate manager"))
		return diags
	}

	if certificateChain, ok := d.GetOk("certificate_chain"); ok {
		certChain := certificateChain.(string)
		certificatePostDto.Properties.CertificateChain = &certChain
	} else {
		diags := diag.FromErr(fmt.Errorf("certificateChain must be provided for certificate manager"))
		return diags
	}

	if privateKey, ok := d.GetOk("private_key"); ok {
		keyStr := privateKey.(string)
		certificatePostDto.Properties.PrivateKey = &keyStr
	} else {
		diags := diag.FromErr(fmt.Errorf("private key must be provided for certificate manager"))
		return diags
	}
	certificateDto, apiResponse, err := client.CreateCertificate(ctx, certificatePostDto)
	certManagerLogApiResponse(apiResponse)
	if err != nil {
		d.SetId("")
		diags := diag.FromErr(fmt.Errorf("error creating certificate: %w", err))
		return diags
	}

	d.SetId(*certificateDto.Id)

	diagErr := waitForCertToBeReady(ctx, d, client)
	if diagErr != nil {
		return diagErr
	}

	return resourceCertificateManagerRead(ctx, d, meta)
}

func waitForCertToBeReady(ctx context.Context, d *schema.ResourceData, client *cert.Client) diag.Diagnostics {
	for {
		log.Printf("[INFO] Waiting for certificate %s to be ready...", d.Id())

		certReady, rsErr := certReady(ctx, d, client)
		if rsErr != nil {
			diags := diag.FromErr(fmt.Errorf("error while checking readiness status of certificate %s: %w", d.Id(), rsErr))
			return diags
		}

		if certReady {
			log.Printf("[INFO] certificate ready: %s", d.Id())
			break
		}

		select {
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			diags := diag.FromErr(fmt.Errorf("certificate check timed out! WARNING: your certificate will still probably be created/updated " +
				"after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates"))
			return diags
		}
	}
	return nil
}

func certReady(ctx context.Context, d *schema.ResourceData, client *cert.Client) (bool, error) {
	backupUnit, apiResponse, err := client.GetCertificate(ctx, d.Id())
	certManagerLogApiResponse(apiResponse)
	if err != nil {
		return true, fmt.Errorf("error checking certificate status: %s", err)
	}
	return *backupUnit.Metadata.State == "AVAILABLE", nil
}

func resourceCertificateManagerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CertManagerClient

	certDto, apiResponse, err := client.GetCertificate(ctx, d.Id())
	certManagerLogApiResponse(apiResponse)
	if err != nil {
		if certManagerHttpNotFound(apiResponse) {
			log.Printf("[INFO] Resource %s not found: %+v", d.Id(), err)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully retreived certificate %s: %+v", d.Id(), certDto)

	if err := cert.SetCertificateData(d, &certDto); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceCertificateManagerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CertManagerClient

	certificatePatchDto := certmanager.CertificatePatchDto{
		Properties: &certmanager.CertificatePatchPropertiesDto{},
	}

	if d.HasChange("name") {
		_, v := d.GetChange("name")
		vStr := v.(string)
		certificatePatchDto.Properties.Name = &vStr
	}

	_, apiResponse, err := client.UpdateCertificate(ctx, d.Id(), certificatePatchDto)
	certManagerLogApiResponse(apiResponse)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while updating certificate ID %s %w", d.Id(), err))
		return diags
	}

	diagErr := waitForCertToBeReady(ctx, d, client)
	if diagErr != nil {
		return diagErr
	}

	return resourceCertificateManagerRead(ctx, d, meta)
}

func resourceCertificateManagerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CertManagerClient
	deleted := false
	for deleted != true {

		apiResponse, err := client.DeleteCertificate(ctx, d.Id())
		certManagerLogApiResponse(apiResponse)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occured while deleting an certificate %s %w", d.Id(), err))
			return diags
		}

		deleted, err = certDeleted(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
		select {
		case <-time.After(SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			diags := diag.FromErr(fmt.Errorf("certificate deletion timed out! WARNING: your certificate (%s) will still probably be deleted after some time "+
				"but the terraform state won't reflect that; check your Ionos Cloud account for updates", d.Id()))
			return diags
		}
	}

	log.Printf("[INFO] Successfully deleted certificate: %s", d.Id())

	d.SetId("")

	return nil
}

func certDeleted(ctx context.Context, d *schema.ResourceData, client *cert.Client) (bool, error) {

	_, apiResponse, err := client.GetCertificate(ctx, d.Id())

	if err != nil {
		if certManagerHttpNotFound(apiResponse) {
			return true, nil
		}
		return true, fmt.Errorf("error checking certificate deletion status: %w", err)
	}
	return false, nil
}

func resourceCertificateManagerImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(SdkBundle).CertManagerClient

	certId := d.Id()
	certDto, apiResponse, err := client.GetCertificate(ctx, d.Id())
	certManagerLogApiResponse(apiResponse)
	if err != nil {
		if certManagerHttpNotFound(apiResponse) {
			d.SetId("")
			return nil, fmt.Errorf("unable to find cert %q", certId)
		}
		return nil, fmt.Errorf("an error occured while retrieving the cert %q, %w", certId, err)
	}

	if err := cert.SetCertificateData(d, &certDto); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
