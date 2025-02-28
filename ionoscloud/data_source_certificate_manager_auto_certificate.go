package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	certSDK "github.com/ionos-cloud/sdk-go-cert-manager"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	certService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cert"
)

func dataSourceCertificateManagerAutoCertificate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAutoCertificateRead,
		Schema: map[string]*schema.Schema{
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The location of the auto-certificate",
			},
			"id": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "The ID of the auto-certificate",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the auto-certificate",
			},
			"common_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The common name (DNS) of the certificate to issue. The common name is a part of a zone in IONOS Cloud DNS",
			},
			"key_algorithm": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The key algorithm used to generate the certificate",
			},
			"subject_alternative_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Additional names added to the issued certificate. The additional names are part of a zone in IONOS Cloud DNS",
			},
			"last_issued_certificate_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the last certificate that was issued",
			},
			"provider_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The certificate provider used to issue the certificates",
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceAutoCertificateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CertManagerClient
	location := d.Get("location").(string)
	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if idOk && nameOk {
		return diag.FromErr(fmt.Errorf("ID and name cannot be provided at the same time"))
	}
	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf("please provide either the auto-certificate ID or name"))
	}

	var autoCertificate certSDK.AutoCertificateRead
	var err error

	if idOk {
		id := id.(string)
		autoCertificate, _, err = client.GetAutoCertificate(ctx, id, location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching the auto-certificate with ID: %v, error: %w", id, err))
		}
	} else {
		autoCertificates, _, err := client.ListAutoCertificates(ctx, location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching auto-certificates: %w", err))
		}
		var results []certSDK.AutoCertificateRead
		if autoCertificates.Items != nil {
			for _, autoCertificateItem := range *autoCertificates.Items {
				if autoCertificateItem.Properties != nil && autoCertificateItem.Properties.Name != nil && strings.EqualFold(*autoCertificateItem.Properties.Name, name.(string)) {
					results = append(results, autoCertificateItem)
				}
			}
		}

		if len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no auto-certificate found with the specified name: %v", name))
		}
		if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one auto-certificate found with the specified name: %v", name))
		}
		autoCertificate = results[0]
	}

	if err := certService.SetAutoCertificateData(d, autoCertificate); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
