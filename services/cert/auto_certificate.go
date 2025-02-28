package cert

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	certmanager "github.com/ionos-cloud/sdk-go-cert-manager"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

// Setting nolint:golint for all these functions since there is no need to add function docs (the functions
// are suggestive enough).
//
//nolint:golint
func (c *Client) GetAutoCertificate(ctx context.Context, autoCertificateID, location string) (certmanager.AutoCertificateRead, *certmanager.APIResponse, error) {
	c.ChangeConfigURL(location)
	autoCertificate, apiResponse, err := c.sdkClient.AutoCertificateApi.AutoCertificatesFindById(ctx, autoCertificateID).Execute()
	apiResponse.LogInfo()
	return autoCertificate, apiResponse, err
}

//nolint:golint
func (c *Client) ListAutoCertificates(ctx context.Context, location string) (certmanager.AutoCertificateReadList, *certmanager.APIResponse, error) {
	c.ChangeConfigURL(location)
	autoCertificates, apiResponse, err := c.sdkClient.AutoCertificateApi.AutoCertificatesGet(ctx).Execute()
	apiResponse.LogInfo()
	return autoCertificates, apiResponse, err
}

//nolint:golint
func (c *Client) CreateAutoCertificate(ctx context.Context, location string, autoCertificatePostData certmanager.AutoCertificateCreate) (certmanager.AutoCertificateRead, *certmanager.APIResponse, error) {
	c.ChangeConfigURL(location)
	autoCertificate, apiResponse, err := c.sdkClient.AutoCertificateApi.AutoCertificatesPost(ctx).AutoCertificateCreate(autoCertificatePostData).Execute()
	apiResponse.LogInfo()
	return autoCertificate, apiResponse, err
}

//nolint:golint
func (c *Client) UpdateAutoCertificate(ctx context.Context, autoCertificateID, location string, autoCertificatePatchData certmanager.AutoCertificatePatch) (certmanager.AutoCertificateRead, *certmanager.APIResponse, error) {
	c.ChangeConfigURL(location)
	autoCertificate, apiResponse, err := c.sdkClient.AutoCertificateApi.AutoCertificatesPatch(ctx, autoCertificateID).AutoCertificatePatch(autoCertificatePatchData).Execute()
	apiResponse.LogInfo()
	return autoCertificate, apiResponse, err
}

//nolint:golint
func (c *Client) DeleteAutoCertificate(ctx context.Context, autoCertificateID, location string) (*certmanager.APIResponse, error) {
	c.ChangeConfigURL(location)
	apiResponse, err := c.sdkClient.AutoCertificateApi.AutoCertificatesDelete(ctx, autoCertificateID).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

//nolint:golint
func (c *Client) IsAutoCertificateReady(ctx context.Context, d *schema.ResourceData) (bool, error) {
	autoCertificateID := d.Id()
	location := d.Get("location").(string)
	autoCertificate, _, err := c.GetAutoCertificate(ctx, autoCertificateID, location)
	if err != nil {
		return false, fmt.Errorf("error checking auto-certificate status: %w", err)
	}
	if autoCertificate.Metadata == nil || autoCertificate.Metadata.State == nil {
		return false, fmt.Errorf("metadata or state is empty for auto-certificate with ID: %v", autoCertificateID)
	}
	if utils.IsStateFailed(*autoCertificate.Metadata.State) {
		return false, fmt.Errorf("error while checking if auto-certificate is ready, auto-certificate ID: %v, state: %v", autoCertificateID, *autoCertificate.Metadata.State)
	}
	return strings.EqualFold(*autoCertificate.Metadata.State, constant.Available), nil
}

//nolint:golint
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
	if autoCertificate.Metadata != nil && autoCertificate.Metadata.State != nil && utils.IsStateFailed(*autoCertificate.Metadata.State) {
		return false, fmt.Errorf("error while checking if auto-certificate is deleted properly, auto-certificate ID: %v, state: %v", autoCertificateID, *autoCertificate.Metadata.State)
	}
	return false, nil
}

//nolint:golint
func GetAutoCertificateDataCreate(d *schema.ResourceData) *certmanager.AutoCertificateCreate {
	autoCertificate := certmanager.AutoCertificateCreate{
		Properties: &certmanager.AutoCertificate{},
	}

	providerID := d.Get("provider_id").(string)
	autoCertificate.Properties.Provider = &providerID
	commonName := d.Get("common_name").(string)
	autoCertificate.Properties.CommonName = &commonName
	name := d.Get("name").(string)
	autoCertificate.Properties.Name = &name
	keyAlgorithm := d.Get("key_algorithm").(string)
	autoCertificate.Properties.KeyAlgorithm = &keyAlgorithm
	if subjectAlternativeNames, subjectAlternativeNamesOk := d.GetOk("subject_alternative_names"); subjectAlternativeNamesOk {
		subjectAlternativeNames := subjectAlternativeNames.([]interface{})
		var subjectAlternativeNamesList []string
		for _, subjectAlternativeName := range subjectAlternativeNames {
			subjectAlternativeName := subjectAlternativeName.(string)
			subjectAlternativeNamesList = append(subjectAlternativeNamesList, subjectAlternativeName)
		}
		autoCertificate.Properties.SubjectAlternativeNames = &subjectAlternativeNamesList
	}
	return &autoCertificate
}

//nolint:golint
func GetAutoCertificateDataUpdate(d *schema.ResourceData) *certmanager.AutoCertificatePatch {
	autoCertificate := certmanager.AutoCertificatePatch{
		Properties: &certmanager.PatchName{},
	}
	if d.HasChange("name") {
		_, newValue := d.GetChange("name")
		newValueStr := newValue.(string)
		autoCertificate.Properties.Name = &newValueStr
	}
	return &autoCertificate
}

//nolint:golint
func SetAutoCertificateData(d *schema.ResourceData, autoCertificate certmanager.AutoCertificateRead) error {
	resourceName := "Auto-certificate"
	if autoCertificate.Id != nil {
		d.SetId(*autoCertificate.Id)
	}
	if autoCertificate.Metadata == nil || autoCertificate.Properties == nil {
		return fmt.Errorf("response properties/metadata should not be empty for auto-certificate with ID: %v", *autoCertificate.Id)
	}
	if autoCertificate.Metadata.LastIssuedCertificate != nil {
		if err := d.Set("last_issued_certificate_id", *autoCertificate.Metadata.LastIssuedCertificate); err != nil {
			return utils.GenerateSetError(resourceName, "last_issued_certificate_id", err)
		}
	}
	if autoCertificate.Properties.Provider != nil {
		if err := d.Set("provider_id", *autoCertificate.Properties.Provider); err != nil {
			return utils.GenerateSetError(resourceName, "provider_id", err)
		}
	}
	if autoCertificate.Properties.CommonName != nil {
		if err := d.Set("common_name", *autoCertificate.Properties.CommonName); err != nil {
			return utils.GenerateSetError(resourceName, "common_name", err)
		}
	}
	if autoCertificate.Properties.KeyAlgorithm != nil {
		if err := d.Set("key_algorithm", *autoCertificate.Properties.KeyAlgorithm); err != nil {
			return utils.GenerateSetError(resourceName, "key_algorithm", err)
		}
	}
	if autoCertificate.Properties.Name != nil {
		if err := d.Set("name", *autoCertificate.Properties.Name); err != nil {
			return utils.GenerateSetError(resourceName, "name", err)
		}
	}
	if autoCertificate.Properties.SubjectAlternativeNames != nil {
		subjectAlternativeNames := []string{}
		subjectAlternativeNames = append(subjectAlternativeNames, *autoCertificate.Properties.SubjectAlternativeNames...)
		if err := d.Set("subject_alternative_names", subjectAlternativeNames); err != nil {
			return utils.GenerateSetError(resourceName, "subject_alternative_names", err)
		}
	}
	return nil
}
