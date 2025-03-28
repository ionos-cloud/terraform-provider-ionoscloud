package cert

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	certmanager "github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func (c *Client) GetAutoCertificate(ctx context.Context, autoCertificateID, location string) (certmanager.AutoCertificateRead, *shared.APIResponse, error) {
	c.ChangeConfigURL(location)
	autoCertificate, apiResponse, err := c.sdkClient.AutoCertificateApi.AutoCertificatesFindById(ctx, autoCertificateID).Execute()
	apiResponse.LogInfo()
	return autoCertificate, apiResponse, err
}

func (c *Client) ListAutoCertificates(ctx context.Context, location string) (certmanager.AutoCertificateReadList, *shared.APIResponse, error) {
	c.ChangeConfigURL(location)
	autoCertificates, apiResponse, err := c.sdkClient.AutoCertificateApi.AutoCertificatesGet(ctx).Execute()
	apiResponse.LogInfo()
	return autoCertificates, apiResponse, err
}

func (c *Client) CreateAutoCertificate(ctx context.Context, location string, autoCertificatePostData certmanager.AutoCertificateCreate) (certmanager.AutoCertificateRead, *shared.APIResponse, error) {
	c.ChangeConfigURL(location)
	autoCertificate, apiResponse, err := c.sdkClient.AutoCertificateApi.AutoCertificatesPost(ctx).AutoCertificateCreate(autoCertificatePostData).Execute()
	apiResponse.LogInfo()
	return autoCertificate, apiResponse, err
}

func (c *Client) UpdateAutoCertificate(ctx context.Context, autoCertificateID, location string, autoCertificatePatchData certmanager.AutoCertificatePatch) (certmanager.AutoCertificateRead, *shared.APIResponse, error) {
	c.ChangeConfigURL(location)
	autoCertificate, apiResponse, err := c.sdkClient.AutoCertificateApi.AutoCertificatesPatch(ctx, autoCertificateID).AutoCertificatePatch(autoCertificatePatchData).Execute()
	apiResponse.LogInfo()
	return autoCertificate, apiResponse, err
}

func (c *Client) DeleteAutoCertificate(ctx context.Context, autoCertificateID, location string) (*shared.APIResponse, error) {
	c.ChangeConfigURL(location)
	apiResponse, err := c.sdkClient.AutoCertificateApi.AutoCertificatesDelete(ctx, autoCertificateID).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

func (c *Client) IsAutoCertificateReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	autoCertificateID := d.Id()
	location := d.Get("location").(string)
	autoCertificate, _, err := c.GetAutoCertificate(ctx, autoCertificateID, location)
	if err != nil {
		return false, fmt.Errorf("error checking auto-certificate status: %w", err)
	}
	if utils.IsStateFailed(autoCertificate.Metadata.State) {
		return false, fmt.Errorf("error while checking if auto-certificate is ready, auto-certificate ID: %v, state: %v", autoCertificateID, autoCertificate.Metadata.State)
	}
	return strings.EqualFold(autoCertificate.Metadata.State, constant.Available), nil
}

func (c *Client) IsAutoCertificateDeleted(ctx context.Context, d *schema.ResourceData) (bool, error) {
	autoCertificateID := d.Id()
	location := d.Get("location").(string)
	autoCertificate, apiResponse, err := c.GetAutoCertificate(ctx, autoCertificateID, location)
	if err != nil {
		if apiResponse.HttpNotFound() {
			return true, nil
		}
		return false, fmt.Errorf("error while checking deletion status for auto-certificate with ID: %v, error: %w", d.Id(), err)
	}
	if utils.IsStateFailed(autoCertificate.Metadata.State) {
		return false, fmt.Errorf("error while checking if auto-certificate is deleted properly, auto-certificate ID: %v, state: %v", autoCertificateID, autoCertificate.Metadata.State)
	}
	return false, nil
}

func GetAutoCertificateDataCreate(d *schema.ResourceData) *certmanager.AutoCertificateCreate {
	autoCertificate := certmanager.AutoCertificateCreate{
		Properties: certmanager.AutoCertificate{},
	}

	providerID := d.Get("provider_id").(string)
	autoCertificate.Properties.Provider = providerID
	commonName := d.Get("common_name").(string)
	autoCertificate.Properties.CommonName = commonName
	name := d.Get("name").(string)
	autoCertificate.Properties.Name = name
	keyAlgorithm := d.Get("key_algorithm").(string)
	autoCertificate.Properties.KeyAlgorithm = keyAlgorithm
	if subjectAlternativeNames, subjectAlternativeNamesOk := d.GetOk("subject_alternative_names"); subjectAlternativeNamesOk {
		subjectAlternativeNames := subjectAlternativeNames.([]interface{})
		var subjectAlternativeNamesList []string
		for _, subjectAlternativeName := range subjectAlternativeNames {
			subjectAlternativeName := subjectAlternativeName.(string)
			subjectAlternativeNamesList = append(subjectAlternativeNamesList, subjectAlternativeName)
		}
		autoCertificate.Properties.SubjectAlternativeNames = subjectAlternativeNamesList
	}
	return &autoCertificate
}

func GetAutoCertificateDataUpdate(d *schema.ResourceData) *certmanager.AutoCertificatePatch {
	autoCertificate := certmanager.AutoCertificatePatch{
		Properties: certmanager.PatchName{},
	}
	if d.HasChange("name") {
		_, newValue := d.GetChange("name")
		newValueStr := newValue.(string)
		autoCertificate.Properties.Name = newValueStr
	}
	return &autoCertificate
}

func SetAutoCertificateData(d *schema.ResourceData, autoCertificate certmanager.AutoCertificateRead) error {
	resourceName := "Auto-certificate"
	d.SetId(autoCertificate.Id)

	if autoCertificate.Metadata.LastIssuedCertificate != nil {
		if err := d.Set("last_issued_certificate_id", *autoCertificate.Metadata.LastIssuedCertificate); err != nil {
			return utils.GenerateSetError(resourceName, "last_issued_certificate_id", err)
		}
	}
	if err := d.Set("provider_id", autoCertificate.Properties.Provider); err != nil {
		return utils.GenerateSetError(resourceName, "provider_id", err)
	}
	if err := d.Set("common_name", autoCertificate.Properties.CommonName); err != nil {
		return utils.GenerateSetError(resourceName, "common_name", err)
	}
	if err := d.Set("key_algorithm", autoCertificate.Properties.KeyAlgorithm); err != nil {
		return utils.GenerateSetError(resourceName, "key_algorithm", err)
	}
	if err := d.Set("name", autoCertificate.Properties.Name); err != nil {
		return utils.GenerateSetError(resourceName, "name", err)
	}
	subjectAlternativeNames := []string{}
	subjectAlternativeNames = append(subjectAlternativeNames, autoCertificate.Properties.SubjectAlternativeNames...)
	if err := d.Set("subject_alternative_names", subjectAlternativeNames); err != nil {
		return utils.GenerateSetError(resourceName, "subject_alternative_names", err)
	}
	return nil
}
