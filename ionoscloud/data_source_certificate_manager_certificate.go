package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	certmanager "github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cert"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCertificate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCertificateRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"certificate": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"certificate_chain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func dataSourceCertificateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(bundleclient.SdkBundle).CertManagerClient

	var name, idStr string
	id, idOk := d.GetOk("id")
	if idOk {
		idStr = id.(string)
	}
	t, nameOk := d.GetOk("name")
	if nameOk {
		name = t.(string)
	}

	var certificate certmanager.CertificateRead
	var apiResponse *shared.APIResponse
	var err error

	if !idOk && !nameOk {
		return utils.ToDiags(d, "either id, or name must be set", nil)
	}

	if idOk {
		certificate, apiResponse, err = client.GetCertificate(ctx, idStr)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("error getting certificate with id %s %s", idStr, err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}
		if nameOk {
			if !strings.EqualFold(certificate.Properties.Name, name) {
				return utils.ToDiags(d, fmt.Sprintf("name of cert (UUID=%s, name=%s) does not match expected name: %s",
					certificate.Id, certificate.Properties.Name, name), nil)
			}
		}
		log.Printf("[INFO] Got certificate [Name=%s]", certificate.Properties.Name)

	} else {
		log.Printf("[INFO] Using data source for certificate with name: %s", name)

		certificates, apiResponse, err := client.ListCertificates(ctx)
		if err != nil {
			return utils.ToDiags(d, fmt.Sprintf("an error occurred while fetching certificates: %s ", err), &utils.DiagsOpts{StatusCode: apiResponse.StatusCode})
		}

		var results []certmanager.CertificateRead

		if certificates.Items != nil {
			var certsFound []certmanager.CertificateRead
			for _, certItem := range certificates.Items {
				if certItem.Properties.Name == name {
					certsFound = append(certsFound, certItem)
				}
			}

			if certsFound == nil {
				return utils.ToDiags(d, fmt.Sprintf("no certificate found with the specified criteria: name = %s", name), nil)
			} else {
				results = certsFound
			}
		}

		if results == nil || len(results) == 0 {
			return utils.ToDiags(d, fmt.Sprintf("no certificate found with the specified criteria: name = %s", name), nil)
		} else if len(results) > 1 {
			return utils.ToDiags(d, fmt.Sprintf("more than one certificate found with the specified criteria: name = %s", name), nil)
		} else {
			certificate = results[0]
		}

	}

	if err := cert.SetCertificateData(d, &certificate); err != nil {
		return utils.ToDiags(d, err.Error(), nil)
	}

	return nil
}
