package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	certsdk "github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
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
		return utils.ToDiags(d, "ID and name cannot be provided at the same time", nil)
	}
	if !idOk && !nameOk {
		return utils.ToDiags(d, "please provide either the auto-certificate ID or name", nil)
	}

	var autoCertificate certsdk.AutoCertificateRead
	var apiResponse *shared.APIResponse
	var err error

	if idOk {
		id := id.(string)
		autoCertificate, apiResponse, err = client.GetAutoCertificate(ctx, id, location)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching the auto-certificate with ID: %v, error: %s", id, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
	} else {
		autoCertificates, apiResponse, err := client.ListAutoCertificates(ctx, location)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching auto-certificates: %s", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
		var results []certsdk.AutoCertificateRead
		if autoCertificates.Items != nil {
			for _, autoCertificateItem := range autoCertificates.Items {
				if strings.EqualFold(autoCertificateItem.Properties.Name, name.(string)) {
					results = append(results, autoCertificateItem)
				}
			}
		}

		if len(results) == 0 {
			return utils.ToDiags(d, fmt.Sprintf("no auto-certificate found with the specified name: %v", name), nil)
		}
		if len(results) > 1 {
			return utils.ToDiags(d, fmt.Sprintf("more than one auto-certificate found with the specified name: %v", name), nil)
		}
		autoCertificate = results[0]
	}

	if err := certService.SetAutoCertificateData(d, autoCertificate); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}
	return nil
}
