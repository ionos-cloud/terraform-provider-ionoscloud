package cert

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ionos-cloud/sdk-go-bundle/common"
	certmanager "github.com/ionos-cloud/sdk-go-bundle/products/cert"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"strings"
)

func (c *Client) GetCertificate(ctx context.Context, certId string) (certmanager.CertificateDto, *common.APIResponse, error) {
	cert, apiResponse, err := c.sdkClient.CertificatesApi.CertificatesGetById(ctx, certId).Execute()
	apiResponse.LogInfo()
	return cert, apiResponse, err
}

func (c *Client) ListCertificates(ctx context.Context) (certmanager.CertificateCollectionDto, *common.APIResponse, error) {
	certs, apiResponse, err := c.sdkClient.CertificatesApi.CertificatesGet(ctx).Execute()
	apiResponse.LogInfo()
	return certs, apiResponse, err
}

func (c *Client) CreateCertificate(ctx context.Context, certPostDto certmanager.CertificatePostDto) (certmanager.CertificateDto, *common.APIResponse, error) {
	certResponse, apiResponse, err := c.sdkClient.CertificatesApi.CertificatesPost(ctx).CertificatePostDto(certPostDto).Execute()
	apiResponse.LogInfo()
	return certResponse, apiResponse, err
}

func (c *Client) UpdateCertificate(ctx context.Context, certId string, certPatch certmanager.CertificatePatchDto) (certmanager.CertificateDto, *common.APIResponse, error) {
	certResponse, apiResponse, err := c.sdkClient.CertificatesApi.CertificatesPatch(ctx, certId).CertificatePatchDto(certPatch).Execute()
	apiResponse.LogInfo()
	return certResponse, apiResponse, err
}

func (c *Client) DeleteCertificate(ctx context.Context, certId string) (*common.APIResponse, error) {
	apiResponse, err := c.sdkClient.CertificatesApi.CertificatesDelete(ctx, certId).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

func (c *Client) IsCertReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	cert, _, err := c.GetCertificate(ctx, d.Id())
	if err != nil {
		return true, fmt.Errorf("error checking certificate status: %w", err)
	}
	if cert.Metadata == nil || cert.Metadata.State == nil {
		return false, fmt.Errorf("cert metadata or state is empty for id %s", d.Id())
	}
	return strings.EqualFold(*cert.Metadata.State, utils.Available), nil
}

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

func SetCertificateData(d *schema.ResourceData, cert *certmanager.CertificateDto) error {
	if cert.Id != nil {
		d.SetId(*cert.Id)
	}

	if cert.Properties != nil {
		if cert.Properties.Name != nil {
			err := d.Set("name", *cert.Properties.Name)
			if err != nil {
				return fmt.Errorf("error while setting name property for certificate %s: %w", d.Id(), err)
			}
		}
		if cert.Properties.Certificate != nil {
			err := d.Set("certificate", *cert.Properties.Certificate)
			if err != nil {
				return fmt.Errorf("error while setting certificate property for certificate %s: %w", d.Id(), err)
			}
		}
		if cert.Properties.CertificateChain != nil {
			err := d.Set("certificate_chain", *cert.Properties.CertificateChain)
			if err != nil {
				return fmt.Errorf("error while setting certificate_chain property for certificate %s: %w", d.Id(), err)
			}
		}
	}
	return nil
}

func GetCertPostDto(d *schema.ResourceData) (*certmanager.CertificatePostDto, error) {

	certificatePostDto := certmanager.CertificatePostDto{
		Properties: &certmanager.CertificatePostPropertiesDto{},
	}

	if name, nameOk := d.GetOk("name"); nameOk {
		name := name.(string)
		certificatePostDto.Properties.Name = &name
	} else {
		return nil, fmt.Errorf("name must be provided for the certificate")
	}

	if certField, certOk := d.GetOk("certificate"); certOk {
		certificate := certField.(string)
		certificatePostDto.Properties.Certificate = &certificate
	} else {
		return nil, fmt.Errorf("certificate(body) must be provided for the certificate")
	}

	if certificateChain, ok := d.GetOk("certificate_chain"); ok {
		certChain := certificateChain.(string)
		certificatePostDto.Properties.CertificateChain = &certChain
	} else {
		return nil, fmt.Errorf("certificate_chain must be provided for the certificate")
	}

	if privateKey, ok := d.GetOk("private_key"); ok {
		keyStr := privateKey.(string)
		certificatePostDto.Properties.PrivateKey = &keyStr
	} else {
		return nil, fmt.Errorf("private_key must be provided for the certificate")
	}

	return &certificatePostDto, nil
}

func GetCertPatchDto(d *schema.ResourceData) *certmanager.CertificatePatchDto {
	certificatePatchDto := certmanager.CertificatePatchDto{
		Properties: &certmanager.CertificatePatchPropertiesDto{},
	}

	if d.HasChange("name") {
		_, v := d.GetChange("name")
		vStr := v.(string)
		certificatePatchDto.Properties.Name = &vStr
	}

	return &certificatePatchDto
}
