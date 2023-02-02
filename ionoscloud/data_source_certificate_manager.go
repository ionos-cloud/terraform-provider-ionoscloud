package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	certmanager "github.com/ionos-cloud/sdk-go-bundle/products/cert"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cert"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCertificate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCertificateRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
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
	client := meta.(SdkBundle).CertManagerClient

	var name, idStr string
	id, idOk := d.GetOk("id")
	if idOk {
		idStr = id.(string)
	}
	t, nameOk := d.GetOk("name")
	if nameOk {
		name = t.(string)
	}

	var certificate certmanager.CertificateDto
	var err error

	if !idOk && !nameOk {
		return diag.FromErr(fmt.Errorf("either id, or name must be set"))
	}

	if idOk {
		certificate, _, err = client.GetCertificate(ctx, idStr)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error getting certificate with id %s %w", idStr, err))
		}
		if nameOk {
			if certificate.Properties != nil && certificate.Properties.Name != nil &&
				strings.EqualFold(*certificate.Properties.Name, name) {
				return diag.FromErr(fmt.Errorf("name of cert (UUID=%s, name=%s) does not match expected name: %s",
					*certificate.Id, *certificate.Properties.Name, name))
			}
		}
		if certificate.Properties != nil {
			log.Printf("[INFO] Got certificate [Name=%s]", *certificate.Properties.Name)
		}

	} else {
		log.Printf("[INFO] Using data source for certificate with name: %s", name)

		certificates, _, err := client.ListCertificates(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occured while fetching certificates: %w ", err))
		}

		var results []certmanager.CertificateDto

		if certificates.Items != nil {
			var certsFound []certmanager.CertificateDto
			for _, certItem := range *certificates.Items {
				if certItem.Properties != nil && certItem.Properties.Name != nil && *certItem.Properties.Name == name {
					certsFound = append(certsFound, certItem)
				}
			}

			if certsFound == nil {
				return diag.FromErr(fmt.Errorf("no certificate found with the specified criteria: name = %s", name))
			} else {
				results = certsFound
			}
		}

		if results == nil || len(results) == 0 {
			return diag.FromErr(fmt.Errorf("no certificate found with the specified criteria: name = %s", name))
		} else if len(results) > 1 {
			return diag.FromErr(fmt.Errorf("more than one certificate found with the specified criteria: name = %s", name))
		} else {
			certificate = results[0]
		}

	}

	if err := cert.SetCertificateData(d, &certificate); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
