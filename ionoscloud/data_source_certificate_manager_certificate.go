package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	certmanager "github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
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
			if !strings.EqualFold(certificate.Properties.Name, name) {
				return diag.FromErr(fmt.Errorf("name of cert (UUID=%s, name=%s) does not match expected name: %s",
					certificate.Id, certificate.Properties.Name, name))
			}
		}
		log.Printf("[INFO] Got certificate [Name=%s]", certificate.Properties.Name)

	} else {
		log.Printf("[INFO] Using data source for certificate with name: %s", name)

		certificates, _, err := client.ListCertificates(ctx)
		if err != nil {
			return diag.FromErr(fmt.Errorf("an error occurred while fetching certificates: %w ", err))
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
