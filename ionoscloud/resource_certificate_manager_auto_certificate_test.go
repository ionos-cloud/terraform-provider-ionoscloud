package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	certSDK "github.com/ionos-cloud/sdk-go-cert-manager"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func TestAccCertificateManagerAutoCertificate(t *testing.T) {
	var autoCertificate certSDK.AutoCertificateRead

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCMAutoCertificateDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: CMAutoCertificateConfig,
				Check: resource.ComposeTestCheckFunc(
					CMAutoCertificateExistenceCheck(constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, &autoCertificate),
					resource.TestCheckResourceAttr(constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMProviderLocationAttr, CMProviderLocationVal),
					resource.TestCheckResourceAttr(constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMAutoCertCommonNameAttr, CMAutoCertCommonNameVal),
					resource.TestCheckResourceAttr(constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMAutoCertKeyAlgorithmAttr, CMAutoCertKeyAlgorithmVal),
					resource.TestCheckResourceAttr(constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMAutoCertNameAttr, CMAutoCertNameVal),
					resource.TestCheckResourceAttrSet(constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMAutoCertProviderIDAttr),
					resource.TestCheckResourceAttrSet(constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMAutoCertLastIssuedCertID),
					resource.TestCheckResourceAttr(constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMAutoCertSubjectAlternativeNamesAttr+".0", CMAutoCertSubjectAlternativeNamesVal),
				),
			},
			{
				Config: CMAutoCertificateConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					CMAutoCertificateExistenceCheck(constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, &autoCertificate),
					resource.TestCheckResourceAttr(constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMProviderLocationAttr, CMProviderLocationVal),
					resource.TestCheckResourceAttr(constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMAutoCertCommonNameAttr, CMAutoCertCommonNameVal),
					resource.TestCheckResourceAttr(constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMAutoCertKeyAlgorithmAttr, CMAutoCertKeyAlgorithmVal),
					resource.TestCheckResourceAttr(constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMAutoCertNameAttr, CMAutoCertNameUpdatedVal),
					resource.TestCheckResourceAttrSet(constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMAutoCertProviderIDAttr),
					resource.TestCheckResourceAttrSet(constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMAutoCertLastIssuedCertID),
					resource.TestCheckResourceAttr(constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMAutoCertSubjectAlternativeNamesAttr+".0", CMAutoCertSubjectAlternativeNamesVal),
				),
			},
			{
				Config: CMAutoCertificateDSByID,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMProviderLocationAttr, CMProviderLocationVal),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMAutoCertCommonNameAttr, CMAutoCertCommonNameVal),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMAutoCertKeyAlgorithmAttr, CMAutoCertKeyAlgorithmVal),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMAutoCertNameAttr, CMAutoCertNameVal),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMAutoCertProviderIDAttr),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMAutoCertLastIssuedCertID),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMAutoCertSubjectAlternativeNamesAttr+".0", CMAutoCertSubjectAlternativeNamesVal),
				),
			},
			{
				Config: CMAutoCertificateDSByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMProviderLocationAttr, CMProviderLocationVal),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMAutoCertCommonNameAttr, CMAutoCertCommonNameVal),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMAutoCertKeyAlgorithmAttr, CMAutoCertKeyAlgorithmVal),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMAutoCertNameAttr, CMAutoCertNameVal),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMAutoCertProviderIDAttr),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMAutoCertLastIssuedCertID),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.AutoCertificateResource+"."+constant.TestCMAutoCertificateName, CMAutoCertSubjectAlternativeNamesAttr+".0", CMAutoCertSubjectAlternativeNamesVal),
				),
			},
			{
				Config:      CMAutoCertificateDSInvalidConfBothIDAndName,
				ExpectError: regexp.MustCompile("ID and name cannot be provided at the same time"),
			},
			{
				Config:      CMAutoCertificateDSInvalidConfNoIDNoName,
				ExpectError: regexp.MustCompile("please provide either the auto-certificate ID or name"),
			},
			{
				Config:      CMAutoCertificateDSInvalidConfigWrongName,
				ExpectError: regexp.MustCompile("no auto-certificate found with the specified name:"),
			},
		},
	})
}

func testAccCMAutoCertificateDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).CertManagerClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.AutoCertificateResource {
			continue
		}
		autoCertificateID := rs.Primary.ID
		location := rs.Primary.Attributes["location"]
		_, apiResponse, err := client.GetAutoCertificate(ctx, autoCertificateID, location)
		if err != nil {
			if !apiResponse.HttpNotFound() {
				return fmt.Errorf("an error occurred while checking the destruction of auto-certificate with ID: %v, error: %w", autoCertificateID, err)
			}
		} else {
			return fmt.Errorf("auto-certificate with ID: %v still exists", autoCertificateID)
		}
	}
	return nil
}

func CMAutoCertificateExistenceCheck(path string, autoCertificate *certSDK.AutoCertificateRead) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).CertManagerClient
		rs, ok := s.RootModule().Resources[path]

		if !ok {
			return fmt.Errorf("not found: %s", path)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set for the auto-certificate")
		}
		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		defer cancel()
		autoCertificateID := rs.Primary.ID
		location := rs.Primary.Attributes["location"]
		autoCertificateResponse, _, err := client.GetAutoCertificate(ctx, autoCertificateID, location)
		if err != nil {
			return fmt.Errorf("an error occurred while fetching auto-certificate with ID: %v, error: %w", autoCertificateID, err)
		}
		autoCertificate = &autoCertificateResponse
		return nil
	}
}

// Attributes

const CMAutoCertProviderIDAttr = "provider_id"
const CMAutoCertCommonNameAttr = "common_name"
const CMAutoCertKeyAlgorithmAttr = "key_algorithm"
const CMAutoCertNameAttr = "name"
const CMAutoCertSubjectAlternativeNamesAttr = "subject_alternative_names"
const CMAutoCertLastIssuedCertID = "last_issued_certificate_id"

// Values

const CMAutoCertCommonNameVal = "devsdkionos.net"
const CMAutoCertKeyAlgorithmVal = "rsa4096"
const CMAutoCertNameVal = "testAutoCertificate"
const CMAutoCertNameUpdatedVal = "updatedTestAutoCertificate"
const CMAutoCertSubjectAlternativeNamesVal = "devsdkionos.net"

// Configurations

const CMAutoCertificateConfig = CMProviderConfig + `
resource ` + constant.AutoCertificateResource + ` ` + constant.TestCMAutoCertificateName + ` {
	` + CMAutoCertProviderIDAttr + ` = ` + constant.AutoCertificateProviderResource + `.` + constant.TestCMProviderName + `.id
	` + CMAutoCertCommonNameAttr + ` = "` + CMAutoCertCommonNameVal + `"
	` + CMProviderLocationAttr + ` = "` + CMProviderLocationVal + `"
	` + CMAutoCertKeyAlgorithmAttr + ` = "` + CMAutoCertKeyAlgorithmVal + `"
	` + CMAutoCertNameAttr + ` = "` + CMAutoCertNameVal + `"
	` + CMAutoCertSubjectAlternativeNamesAttr + ` = ` + `["` + CMAutoCertSubjectAlternativeNamesVal + `"]
}
`
const CMAutoCertificateConfigUpdate = CMProviderConfig + `
resource ` + constant.AutoCertificateResource + ` ` + constant.TestCMAutoCertificateName + ` {
	` + CMAutoCertProviderIDAttr + ` = ` + constant.AutoCertificateProviderResource + `.` + constant.TestCMProviderName + `.id
	` + CMAutoCertCommonNameAttr + ` = "` + CMAutoCertCommonNameVal + `"
	` + CMProviderLocationAttr + ` = "` + CMProviderLocationVal + `"
	` + CMAutoCertKeyAlgorithmAttr + ` = "` + CMAutoCertKeyAlgorithmVal + `"
	` + CMAutoCertNameAttr + ` = "` + CMAutoCertNameUpdatedVal + `"
	` + CMAutoCertSubjectAlternativeNamesAttr + ` = ` + `["` + CMAutoCertSubjectAlternativeNamesVal + `"]
}
`

const CMAutoCertificateDSByID = CMAutoCertificateConfig + `
` + constant.DataSource + ` ` + constant.AutoCertificateResource + ` ` + constant.TestCMAutoCertificateName + `{
	id = ` + constant.AutoCertificateResource + `.` + constant.TestCMAutoCertificateName + `.id
	` + CMProviderLocationAttr + ` = "` + CMProviderLocationVal + `"
}
`

const CMAutoCertificateDSByName = CMAutoCertificateConfig + `
` + constant.DataSource + ` ` + constant.AutoCertificateResource + ` ` + constant.TestCMAutoCertificateName + `{
	name = ` + constant.AutoCertificateResource + `.` + constant.TestCMAutoCertificateName + `.name
	` + CMProviderLocationAttr + ` = "` + CMProviderLocationVal + `"
}
`

const CMAutoCertificateDSInvalidConfBothIDAndName = CMAutoCertificateConfig + `
` + constant.DataSource + ` ` + constant.AutoCertificateResource + ` ` + constant.TestCMAutoCertificateName + `{
	id = ` + constant.AutoCertificateResource + `.` + constant.TestCMAutoCertificateName + `.id
	name = ` + constant.AutoCertificateResource + `.` + constant.TestCMAutoCertificateName + `.name
	` + CMProviderLocationAttr + ` = "` + CMProviderLocationVal + `"
}
`

const CMAutoCertificateDSInvalidConfNoIDNoName = CMAutoCertificateConfig + `
` + constant.DataSource + ` ` + constant.AutoCertificateResource + ` ` + constant.TestCMAutoCertificateName + `{
	` + CMProviderLocationAttr + ` = "` + CMProviderLocationVal + `"
}
`

const CMAutoCertificateDSInvalidConfigWrongName = CMAutoCertificateConfig + `
` + constant.DataSource + ` ` + constant.AutoCertificateResource + ` ` + constant.TestCMAutoCertificateName + `{
	` + CMProviderLocationAttr + ` = "` + CMProviderLocationVal + `"
	name = "wrongName"
}
`
