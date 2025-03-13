//go:build compute || all || cert

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCertificateResAndDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckCertificateDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCertConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(constant.CertificateResource+"."+constant.TestCertName, "certificate"),
					resource.TestCheckResourceAttrSet(constant.CertificateResource+"."+constant.TestCertName, "certificate_chain"),
					resource.TestCheckResourceAttrSet(constant.CertificateResource+"."+constant.TestCertName, "private_key"),
					resource.TestCheckResourceAttr(constant.CertificateResource+"."+constant.TestCertName, "name", constant.TestCertName),
				),
			},
			{
				Config: testAccCheckCertUpdateName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(constant.CertificateResource+"."+constant.TestCertName, "certificate"),
					resource.TestCheckResourceAttrSet(constant.CertificateResource+"."+constant.TestCertName, "certificate_chain"),
					resource.TestCheckResourceAttrSet(constant.CertificateResource+"."+constant.TestCertName, "private_key"),
					resource.TestCheckResourceAttr(constant.CertificateResource+"."+constant.TestCertName, "name", constant.TestCertName+"1"),
				),
			},
			{
				Config: testAccCheckDataSourceByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.CertificateResource+"."+constant.TestCertName, "certificate", constant.DataSource+"."+constant.CertificateResource+"."+constant.TestCertName, "certificate")),
			},
			{
				Config: testAccCheckDataSourceById,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.CertificateResource+"."+constant.TestCertName, "certificate", constant.DataSource+"."+constant.CertificateResource+"."+constant.TestCertName, "certificate")),
			},
			{
				Config:      testAccCheckDataSourceWrongName,
				ExpectError: regexp.MustCompile(`no certificate found with the specified criteria: name = should_not_work`),
			},
		},
	})
}

func testAccCheckCertificateDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(bundleclient.SdkBundle).CertManagerClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {

		if rs.Type != constant.CertificateResource {
			continue
		}

		_, apiResponse, err := client.GetCertificate(ctx, rs.Primary.ID)

		if err != nil {
			if !apiResponse.HttpNotFound() {
				return fmt.Errorf("an error occurred while checking for the destruction of certificate %s: %w",
					rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("certificate %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const (
	testAccCheckCertConfigBasic = `
resource ` + constant.CertificateResource + ` ` + constant.TestCertName + ` {
	name        	  = "` + constant.TestCertName + `"
	certificate 	  = <<EOT
` + constant.TestCertificate + `
EOT
	certificate_chain = <<EOT
` + constant.TestCertificate + `
EOT
	private_key 	  = <<EOT
` + constant.PrivateKey + `
EOT
}
`
	testAccCheckCertUpdateName = `
resource ` + constant.CertificateResource + ` ` + constant.TestCertName + ` {
	name        	  = "` + constant.TestCertName + `1"
	certificate 	  = <<EOT
` + constant.TestCertificate + `
EOT
	certificate_chain = <<EOT
` + constant.TestCertificate + `
EOT
	private_key 	  = <<EOT
` + constant.PrivateKey + `
EOT
}
`
	testAccCheckDataSourceByName = `
resource ` + constant.CertificateResource + ` ` + constant.TestCertName + ` {
	name        	  = "` + constant.TestCertName + `1"
	certificate 	  = <<EOT
` + constant.TestCertificate + `
EOT
	certificate_chain = <<EOT
` + constant.TestCertificate + `
EOT
	private_key 	  = <<EOT
` + constant.PrivateKey + `
EOT
}
` + constant.DataSource + ` ` + constant.CertificateResource + ` ` + constant.TestCertName + ` {
name ="` + constant.TestCertName + `1"
}
`
	testAccCheckDataSourceById = `
resource ` + constant.CertificateResource + ` ` + constant.TestCertName + ` {
	name        	  = "` + constant.TestCertName + `1"
	certificate 	  = <<EOT
` + constant.TestCertificate + `
EOT
	certificate_chain = <<EOT
` + constant.TestCertificate + `
EOT
	private_key 	  = <<EOT
` + constant.PrivateKey + `
EOT
}
` + constant.DataSource + ` ` + constant.CertificateResource + ` ` + constant.TestCertName + ` {
id =` + constant.CertificateResource + `.` + constant.TestCertName + `.id
}
`
	testAccCheckDataSourceWrongName = `
resource ` + constant.CertificateResource + ` ` + constant.TestCertName + ` {
	name        	  = "` + constant.TestCertName + `1"
	certificate 	  = <<EOT
` + constant.TestCertificate + `
EOT
	certificate_chain = <<EOT
` + constant.TestCertificate + `
EOT
	private_key 	  = <<EOT
` + constant.PrivateKey + `
EOT
}
` + constant.DataSource + ` ` + constant.CertificateResource + ` ` + constant.TestCertName + ` {
name ="should_not_work"
}
`
)
