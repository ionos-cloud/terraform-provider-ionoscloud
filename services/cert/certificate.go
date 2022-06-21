package cert

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	certmanager "github.com/ionos-cloud/sdk-cert-go"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"log"
	"net/http"
	"time"
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
	LogApiResponse(apiResponse)
	if apiResponse != nil {
		return cert, apiResponse, err

	}
	return cert, nil, err
}

func (c *Client) ListCertificates(ctx context.Context) (certmanager.CertificateCollectionDto, *certmanager.APIResponse, error) {
	certs, apiResponse, err := c.CertificateApi.GetCertificates(ctx).Execute()
	LogApiResponse(apiResponse)
	if apiResponse != nil {
		return certs, apiResponse, err
	}
	return certs, nil, err
}

func (c *Client) CreateCertificate(ctx context.Context, certPostDto certmanager.CertificatePostDto) (certmanager.CertificateDto, *certmanager.APIResponse, error) {
	certResponse, apiResponse, err := c.CertificateApi.AddCertificate(ctx).CertificatePostDto(certPostDto).Execute()
	LogApiResponse(apiResponse)
	if apiResponse != nil {
		return certResponse, apiResponse, err
	}
	return certResponse, nil, err
}

func (c *Client) UpdateCertificate(ctx context.Context, certId string, certPatch certmanager.CertificatePatchDto) (certmanager.CertificateDto, *certmanager.APIResponse, error) {
	certResponse, apiResponse, err := c.CertificateApi.PatchCertificateByUuid(ctx, certId).CertificatePatchDto(certPatch).Execute()
	LogApiResponse(apiResponse)
	if apiResponse != nil {
		return certResponse, apiResponse, err
	}
	return certResponse, nil, err
}

func (c *Client) DeleteCertificate(ctx context.Context, certId string) (*certmanager.APIResponse, error) {
	apiResponse, err := c.CertificateApi.DeleteCertificateByUuid(ctx, certId).Execute()
	LogApiResponse(apiResponse)
	if apiResponse != nil {
		return apiResponse, err
	}
	return nil, err
}

func (c *Client) IsCertReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	cert, _, err := c.GetCertificate(ctx, d.Id())
	if err != nil {
		return true, fmt.Errorf("error checking certificate status: %w", err)
	}
	return *cert.Metadata.State == utils.Available, nil
}

func (c *Client) WaitForCertToBeReady(ctx context.Context, d *schema.ResourceData) error {
	for {
		log.Printf("[INFO] Waiting for certificate %s to be ready...", d.Id())

		certReady, rsErr := c.IsCertReady(ctx, d)
		if rsErr != nil {
			return fmt.Errorf("error while checking readiness status of certificate %s: %w", d.Id(), rsErr)
		}

		if certReady {
			log.Printf("[INFO] certificate ready: %s", d.Id())
			break
		}

		select {
		case <-time.After(utils.SleepInterval):
			log.Printf("[INFO] trying again ...")
		case <-ctx.Done():
			return fmt.Errorf("certificate check timed out! WARNING: your certificate will still probably be created/updated " +
				"after some time but the terraform state won't reflect that; check your Ionos Cloud account for updates")
		}
	}
	return nil
}

func (c *Client) IsCertDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	_, apiResponse, err := c.GetCertificate(ctx, d.Id())
	if err != nil {
		if HttpNotFound(apiResponse) {
			return true, nil
		}
		return true, fmt.Errorf("error checking certificate deletion status: %w", err)
	}
	return false, nil
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
		return nil, fmt.Errorf("certificateChain must be provided for the certificate")
	}

	if privateKey, ok := d.GetOk("private_key"); ok {
		keyStr := privateKey.(string)
		certificatePostDto.Properties.PrivateKey = &keyStr
	} else {
		return nil, fmt.Errorf("private key must be provided for the certificate")
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

func LogApiResponse(resp *certmanager.APIResponse) {
	if resp != nil {
		log.Printf("[DEBUG] Operation : %s",
			resp.Operation)
		if resp.Response != nil {
			log.Printf("[DEBUG] response status code : %d\n", resp.StatusCode)
		}
	}
}
func HttpNotFound(resp *certmanager.APIResponse) bool {
	if resp != nil && resp.Response != nil && resp.StatusCode == http.StatusNotFound {
		return true
	}
	return false
}
