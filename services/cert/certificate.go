package cert

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	certmanager "github.com/ionos-cloud/sdk-cert-go"
)

type CertificateService interface {
	GetCertificate(ctx context.Context, certId string) (certmanager.CertificateDto, *certmanager.APIResponse, error)
	ListCertificates(ctx context.Context, filterName string) (certmanager.CertificateDto, *certmanager.APIResponse, error)
	CreateCertificate(ctx context.Context, certPostDto certmanager.CertificatePostDto) (certmanager.CertificateDto, *certmanager.APIResponse, error)
	UpdateCertificate(ctx context.Context, certPostDto certmanager.CertificatePatchDto) (certmanager.CertificateDto, *certmanager.APIResponse, error)
	DeleteCertificate(ctx context.Context, certIf string) (*certmanager.APIResponse, error)
}

func (c *Client) GetCertificate(ctx context.Context, certId string) (certmanager.CertificateDto, *certmanager.APIResponse, error) {
	cert, apiResponse, err := c.CertificateApi.GetCertificateByUuid(ctx, certId).Execute()
	if apiResponse != nil {
		return cert, apiResponse, err

	}
	return cert, nil, err
}

func (c *Client) ListCertificates(ctx context.Context) (certmanager.CertificateCollectionDto, *certmanager.APIResponse, error) {
	certs, apiResponse, err := c.CertificateApi.GetCertificates(ctx).Execute()
	if apiResponse != nil {
		return certs, apiResponse, err
	}
	return certs, nil, err
}

func (c *Client) CreateCertificate(ctx context.Context, certPostDto certmanager.CertificatePostDto) (certmanager.CertificateDto, *certmanager.APIResponse, error) {
	certResponse, apiResponse, err := c.CertificateApi.AddCertificate(ctx).CertificatePostDto(certPostDto).Execute()
	if apiResponse != nil {
		return certResponse, apiResponse, err
	}
	return certResponse, nil, err
}

func (c *Client) UpdateCertificate(ctx context.Context, certId string, certPatch certmanager.CertificatePatchDto) (certmanager.CertificateDto, *certmanager.APIResponse, error) {
	certResponse, apiResponse, err := c.CertificateApi.PatchCertificateByUuid(ctx, certId).CertificatePatchDto(certPatch).Execute()
	if apiResponse != nil {
		return certResponse, apiResponse, err
	}
	return certResponse, nil, err
}

func (c *Client) DeleteCertificate(ctx context.Context, certId string) (*certmanager.APIResponse, error) {
	apiResponse, err := c.CertificateApi.DeleteCertificateByUuid(ctx, certId).Execute()
	if apiResponse != nil {
		return apiResponse, err
	}
	return nil, err
}

func SetCertificateData(d *schema.ResourceData, applicationLoadBalancer *certmanager.CertificateDto) error {

	if applicationLoadBalancer.Id != nil {
		d.SetId(*applicationLoadBalancer.Id)
	}

	if applicationLoadBalancer.Properties != nil {
		if applicationLoadBalancer.Properties.Name != nil {
			err := d.Set("name", *applicationLoadBalancer.Properties.Name)
			if err != nil {
				return fmt.Errorf("error while setting name property for certificate %s: %w", d.Id(), err)
			}
		}
		if applicationLoadBalancer.Properties.Certificate != nil {
			err := d.Set("certificate", *applicationLoadBalancer.Properties.Certificate)
			if err != nil {
				return fmt.Errorf("error while setting certificate property for certificate %s: %w", d.Id(), err)
			}
		}
		if applicationLoadBalancer.Properties.CertificateChain != nil {
			err := d.Set("certificate_chain", *applicationLoadBalancer.Properties.CertificateChain)
			if err != nil {
				return fmt.Errorf("error while setting certificate_chain property for certificate %s: %w", d.Id(), err)
			}
		}
	}
	return nil
}
