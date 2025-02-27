package cert

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	certmanager "github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

//nolint:golint
func (c *Client) GetCertificate(ctx context.Context, certId string) (certmanager.CertificateRead, *shared.APIResponse, error) {
	cert, apiResponse, err := c.sdkClient.CertificateApi.CertificatesFindById(ctx, certId).Execute()
	apiResponse.LogInfo()
	return cert, apiResponse, err
}

//nolint:golint
func (c *Client) ListCertificates(ctx context.Context) (certmanager.CertificateReadList, *shared.APIResponse, error) {
	certs, apiResponse, err := c.sdkClient.CertificateApi.CertificatesGet(ctx).Execute()
	apiResponse.LogInfo()
	return certs, apiResponse, err
}

//nolint:golint
func (c *Client) CreateCertificate(ctx context.Context, certPostDto certmanager.CertificateCreate) (certmanager.CertificateRead, *shared.APIResponse, error) {
	certResponse, apiResponse, err := c.sdkClient.CertificateApi.CertificatesPost(ctx).CertificateCreate(certPostDto).Execute()
	apiResponse.LogInfo()
	return certResponse, apiResponse, err
}

//nolint:golint
func (c *Client) UpdateCertificate(ctx context.Context, certId string, certPatch certmanager.CertificatePatch) (certmanager.CertificateRead, *shared.APIResponse, error) {
	certResponse, apiResponse, err := c.sdkClient.CertificateApi.CertificatesPatch(ctx, certId).CertificatePatch(certPatch).Execute()
	apiResponse.LogInfo()
	return certResponse, apiResponse, err
}

//nolint:golint
func (c *Client) DeleteCertificate(ctx context.Context, certId string) (*shared.APIResponse, error) {
	apiResponse, err := c.sdkClient.CertificateApi.CertificatesDelete(ctx, certId).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

//nolint:golint
func (c *Client) IsCertReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	cert, _, err := c.GetCertificate(ctx, d.Id())
	if err != nil {
		return true, fmt.Errorf("error checking certificate status: %w", err)
	}
	return strings.EqualFold(cert.Metadata.State, constant.Available), nil
}

//nolint:golint
func (c *Client) IsCertDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	_, apiResponse, err := c.GetCertificate(ctx, d.Id())
	if err != nil {
		if apiResponse.HttpNotFound() {
			return true, nil
		}
		return false, fmt.Errorf("error checking certificate deletion status: %w", err)
	}
	return false, nil
}

//nolint:golint
func SetCertificateData(d *schema.ResourceData, cert *certmanager.CertificateRead) error {
	d.SetId(cert.Id)

	err := d.Set("name", cert.Properties.Name)
	if err != nil {
		return fmt.Errorf("error while setting name property for certificate %s: %w", d.Id(), err)
	}
	err = d.Set("certificate", cert.Properties.Certificate)
	if err != nil {
		return fmt.Errorf("error while setting certificate property for certificate %s: %w", d.Id(), err)
	}
	err = d.Set("certificate_chain", cert.Properties.CertificateChain)
	if err != nil {
		return fmt.Errorf("error while setting certificate_chain property for certificate %s: %w", d.Id(), err)
	}
	return nil
}

//nolint:golint
func GetCertPostDto(d *schema.ResourceData) (*certmanager.CertificateCreate, error) {

	certificatePostDto := certmanager.CertificateCreate{
		Properties: certmanager.Certificate{},
	}

	if name, nameOk := d.GetOk("name"); nameOk {
		name := name.(string)
		certificatePostDto.Properties.Name = name
	} else {
		return nil, fmt.Errorf("name must be provided for the certificate")
	}

	if certField, certOk := d.GetOk("certificate"); certOk {
		certificate := certField.(string)
		certificatePostDto.Properties.Certificate = certificate
	} else {
		return nil, fmt.Errorf("certificate(body) must be provided for the certificate")
	}

	if certificateChain, ok := d.GetOk("certificate_chain"); ok {
		certChain := certificateChain.(string)
		certificatePostDto.Properties.CertificateChain = certChain
	} else {
		return nil, fmt.Errorf("certificate_chain must be provided for the certificate")
	}

	if privateKey, ok := d.GetOk("private_key"); ok {
		keyStr := privateKey.(string)
		certificatePostDto.Properties.PrivateKey = keyStr
	} else {
		return nil, fmt.Errorf("private_key must be provided for the certificate")
	}

	return &certificatePostDto, nil
}

//nolint:golint
func GetCertPatchDto(d *schema.ResourceData) *certmanager.CertificatePatch {
	certificatePatchDto := certmanager.CertificatePatch{
		Properties: certmanager.PatchName{},
	}

	if d.HasChange("name") {
		_, v := d.GetChange("name")
		vStr := v.(string)
		certificatePatchDto.Properties.Name = vStr
	}

	return &certificatePatchDto
}
