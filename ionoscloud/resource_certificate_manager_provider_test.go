package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	certSDK "github.com/ionos-cloud/sdk-go-cert-manager"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	"regexp"
	"testing"
)

func TestAccCertificateManagerProvider(t *testing.T) {
	var provider certSDK.ProviderRead

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCMProviderDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: CMProviderConfig,
				Check: resource.ComposeTestCheckFunc(
					CMProviderExistenceCheck(constant.AutoCertificateProviderResource+"."+constant.TestCMProviderName, &provider),
					resource.TestCheckResourceAttr(constant.AutoCertificateProviderResource+"."+constant.TestCMProviderName, CMProviderNameAttr, CMProviderNameVal),
					resource.TestCheckResourceAttr(constant.AutoCertificateProviderResource+"."+constant.TestCMProviderName, CMProviderEmailAttr, CMProviderEmailVal),
					resource.TestCheckResourceAttr(constant.AutoCertificateProviderResource+"."+constant.TestCMProviderName, CMProviderLocationAttr, CMProviderLocationVal),
					resource.TestCheckResourceAttr(constant.AutoCertificateProviderResource+"."+constant.TestCMProviderName, CMProviderServerAttr, CMProviderServerVal),
					resource.TestCheckTypeSetElemNestedAttrs(constant.AutoCertificateProviderResource+"."+constant.TestCMProviderName, CMProviderExternalAccountBindingAttr+".*", map[string]string{
						CMProviderEABKeyIDAttr: CMProviderEABKeyIDVal,
					}),
					resource.TestCheckResourceAttrSet(constant.AutoCertificateProviderResource+"."+constant.TestCMProviderName, CMProviderExternalAccountBindingAttr+".0."+CMProviderEABKeySecretAttr),
				),
			},
			{
				Config: CMProviderConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					CMProviderExistenceCheck(constant.AutoCertificateProviderResource+"."+constant.TestCMProviderName, &provider),
					resource.TestCheckResourceAttr(constant.AutoCertificateProviderResource+"."+constant.TestCMProviderName, CMProviderNameAttr, CMProviderNameUpdatedVal),
					resource.TestCheckResourceAttr(constant.AutoCertificateProviderResource+"."+constant.TestCMProviderName, CMProviderEmailAttr, CMProviderEmailVal),
					resource.TestCheckResourceAttr(constant.AutoCertificateProviderResource+"."+constant.TestCMProviderName, CMProviderLocationAttr, CMProviderLocationVal),
					resource.TestCheckResourceAttr(constant.AutoCertificateProviderResource+"."+constant.TestCMProviderName, CMProviderServerAttr, CMProviderServerVal),
					resource.TestCheckTypeSetElemNestedAttrs(constant.AutoCertificateProviderResource+"."+constant.TestCMProviderName, CMProviderExternalAccountBindingAttr+".*", map[string]string{
						CMProviderEABKeyIDAttr: CMProviderEABKeyIDVal,
					}),
				),
			},
			{
				Config: CMProviderDSByID,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.AutoCertificateProviderResource+"."+constant.TestCMProviderName, CMProviderNameAttr, CMProviderNameVal),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.AutoCertificateProviderResource+"."+constant.TestCMProviderName, CMProviderEmailAttr, CMProviderEmailVal),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.AutoCertificateProviderResource+"."+constant.TestCMProviderName, CMProviderLocationAttr, CMProviderLocationVal),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.AutoCertificateProviderResource+"."+constant.TestCMProviderName, CMProviderServerAttr, CMProviderServerVal),
					resource.TestCheckTypeSetElemNestedAttrs(constant.DataSource+"."+constant.AutoCertificateProviderResource+"."+constant.TestCMProviderName, CMProviderExternalAccountBindingAttr+".*", map[string]string{
						CMProviderEABKeyIDAttr: CMProviderEABKeyIDVal,
					}),
				),
			},
			{
				Config: CMProviderDSByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.AutoCertificateProviderResource+"."+constant.TestCMProviderName, CMProviderNameAttr, CMProviderNameVal),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.AutoCertificateProviderResource+"."+constant.TestCMProviderName, CMProviderEmailAttr, CMProviderEmailVal),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.AutoCertificateProviderResource+"."+constant.TestCMProviderName, CMProviderLocationAttr, CMProviderLocationVal),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.AutoCertificateProviderResource+"."+constant.TestCMProviderName, CMProviderServerAttr, CMProviderServerVal),
					resource.TestCheckTypeSetElemNestedAttrs(constant.DataSource+"."+constant.AutoCertificateProviderResource+"."+constant.TestCMProviderName, CMProviderExternalAccountBindingAttr+".*", map[string]string{
						CMProviderEABKeyIDAttr: CMProviderEABKeyIDVal,
					}),
				),
			},
			{
				Config:      CMProviderDSInvalidConfBothIDAndName,
				ExpectError: regexp.MustCompile("ID and name cannot be provided at the same time"),
			},
			{
				Config:      CMProviderDSInvalidConfNoIDNoName,
				ExpectError: regexp.MustCompile("please provide either the auto-certificate provider ID or name"),
			},
			{
				Config:      CMProviderDSInvalidConfigWrongName,
				ExpectError: regexp.MustCompile("no auto-certificate provider found with the specified name:"),
			},
		},
	})
}

func testAccCMProviderDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).CertManagerClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.AutoCertificateProviderResource {
			continue
		}
		providerID := rs.Primary.ID
		location := rs.Primary.Attributes["location"]
		_, apiResponse, err := client.GetProvider(ctx, providerID, location)
		if err != nil {
			if !apiResponse.HttpNotFound() {
				return fmt.Errorf("an error occurred while checking the destruction of auto-certificate provider with ID: %v, error: %w", providerID, err)
			}
		} else {
			return fmt.Errorf("auto-certificate provider with ID: %v still exists", providerID)
		}
	}
	return nil
}

func CMProviderExistenceCheck(path string, provider *certSDK.ProviderRead) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).CertManagerClient
		rs, ok := s.RootModule().Resources[path]

		if !ok {
			return fmt.Errorf("not found: %s", path)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set for the auto-certificate provider")
		}
		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		defer cancel()
		providerID := rs.Primary.ID
		location := rs.Primary.Attributes["location"]
		providerResponse, _, err := client.GetProvider(ctx, providerID, location)
		if err != nil {
			return fmt.Errorf("an error occurred while fetching auto-certificate provider with ID: %v, error: %w", providerID, err)
		}
		provider = &providerResponse
		return nil
	}
}

// Attributes

const CMProviderNameAttr = "name"
const CMProviderEmailAttr = "email"
const CMProviderLocationAttr = "location"
const CMProviderServerAttr = "server"
const CMProviderExternalAccountBindingAttr = "external_account_binding"
const CMProviderEABKeyIDAttr = "key_id"
const CMProviderEABKeySecretAttr = "key_secret"

// Values

const CMProviderNameVal = "testCMProvider"
const CMProviderNameUpdatedVal = "updatedTestCMProvider"
const CMProviderEmailVal = "sdk-go-v6@cloud.ionos.com"
const CMProviderLocationVal = "de/fra"
const CMProviderServerVal = "https://acme-staging-v02.api.letsencrypt.org/directory"
const CMProviderEABKeyIDVal = "some-key-id"
const CMProviderEABKeySecretVal = "secret"
const CMProviderExternalAccountBinding = CMProviderExternalAccountBindingAttr + `{
	` + CMProviderEABKeyIDAttr + ` = "` + CMProviderEABKeyIDVal + `"
	` + CMProviderEABKeySecretAttr + ` = "` + CMProviderEABKeySecretVal + `"
}`

// Configurations

const CMProviderConfig = `
resource ` + constant.AutoCertificateProviderResource + ` ` + constant.TestCMProviderName + ` {
	` + CMProviderNameAttr + ` = "` + CMProviderNameVal + `"
	` + CMProviderEmailAttr + ` = "` + CMProviderEmailVal + `"
	` + CMProviderLocationAttr + ` = "` + CMProviderLocationVal + `"
	` + CMProviderServerAttr + ` = "` + CMProviderServerVal + `"
	` + CMProviderExternalAccountBinding + `
}
`
const CMProviderConfigUpdate = `
resource ` + constant.AutoCertificateProviderResource + ` ` + constant.TestCMProviderName + ` {
	` + CMProviderNameAttr + ` = "` + CMProviderNameUpdatedVal + `"
	` + CMProviderEmailAttr + ` = "` + CMProviderEmailVal + `"
	` + CMProviderLocationAttr + ` = "` + CMProviderLocationVal + `"
	` + CMProviderServerAttr + ` = "` + CMProviderServerVal + `"
	` + CMProviderExternalAccountBinding + `
}
`

const CMProviderDSByID = CMProviderConfig + `
` + constant.DataSource + ` ` + constant.AutoCertificateProviderResource + ` ` + constant.TestCMProviderName + `{
	id = ` + constant.AutoCertificateProviderResource + `.` + constant.TestCMProviderName + `.id
	` + CMProviderLocationAttr + ` = "` + CMProviderLocationVal + `"
}
`

const CMProviderDSByName = CMProviderConfig + `
` + constant.DataSource + ` ` + constant.AutoCertificateProviderResource + ` ` + constant.TestCMProviderName + `{
	name = ` + constant.AutoCertificateProviderResource + `.` + constant.TestCMProviderName + `.name
	` + CMProviderLocationAttr + ` = "` + CMProviderLocationVal + `"
}
`

const CMProviderDSInvalidConfBothIDAndName = CMProviderConfig + `
` + constant.DataSource + ` ` + constant.AutoCertificateProviderResource + ` ` + constant.TestCMProviderName + `{
	id = ` + constant.AutoCertificateProviderResource + `.` + constant.TestCMProviderName + `.id
	name = ` + constant.AutoCertificateProviderResource + `.` + constant.TestCMProviderName + `.name
	` + CMProviderLocationAttr + ` = "` + CMProviderLocationVal + `"
}
`

const CMProviderDSInvalidConfNoIDNoName = CMProviderConfig + `
` + constant.DataSource + ` ` + constant.AutoCertificateProviderResource + ` ` + constant.TestCMProviderName + `{
	` + CMProviderLocationAttr + ` = "` + CMProviderLocationVal + `"
}
`

const CMProviderDSInvalidConfigWrongName = CMProviderConfig + `
` + constant.DataSource + ` ` + constant.AutoCertificateProviderResource + ` ` + constant.TestCMProviderName + `{
	` + CMProviderLocationAttr + ` = "` + CMProviderLocationVal + `"
	name = "wrongName"

}
`
