package ionoscloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	certmanager "github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cert"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
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

func dataSourceCertificateRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		return diagutil.ToDiags(d, fmt.Errorf("either id, or name must be set"), nil)
	}

	if idOk {
		certificate, apiResponse, err = client.GetCertificate(ctx, idStr)
		if err != nil {
			return diagutil.ToDiags(d, fmt.Errorf("error getting certificate with id %s %w", idStr, err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
		}
		if nameOk {
			if !strings.EqualFold(certificate.Properties.Name, name) {
				return diagutil.ToDiags(d, fmt.Errorf("name of cert (UUID=%s, name=%s) does not match expected name: %s",
					certificate.Id, certificate.Properties.Name, name), nil)
			}
		}
		tflog.Info(ctx, "got certificate", map[string]interface{}{"name": certificate.Properties.Name})

	} else {
		tflog.Info(ctx, "searching certificate by name", map[string]interface{}{"name": name})

		certificates, apiResponse, err := client.ListCertificates(ctx)
		if err != nil {
			return diagutil.ToDiags(d, fmt.Errorf("an error occurred while fetching certificates: %w ", err), &diagutil.ErrorContext{StatusCode: apiResponse.SafeStatusCode()})
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
				return diagutil.ToDiags(d, fmt.Errorf("no certificate found with the specified criteria: name = %s", name), nil)
			} else {
				results = certsFound
			}
		}

		if results == nil || len(results) == 0 {
			return diagutil.ToDiags(d, fmt.Errorf("no certificate found with the specified criteria: name = %s", name), nil)
		} else if len(results) > 1 {
			return diagutil.ToDiags(d, fmt.Errorf("more than one certificate found with the specified criteria: name = %s", name), nil)
		} else {
			certificate = results[0]
		}

	}

	if err := cert.SetCertificateData(d, &certificate); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	return nil
}
