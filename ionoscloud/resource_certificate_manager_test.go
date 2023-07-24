//go:build compute || all || cert

package ionoscloud

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testCertificate = `-----BEGIN CERTIFICATE-----
MIIDazCCAlOgAwIBAgIUOH1cikhurIjCjm5Zxt7sfJmhIVAwDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yMjA2MDkxMjM0MzVaFw0zMjA2
MDYxMjM0MzVaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQDkVU596LWGR+/nC3r3MndJfimGMDvt9W4SwL0bOa+V
XxgVgUKYTTCPvwZaAQrtJRUjW2bGxwwj8/3uDEY6vwHJ1Yh+OrbrQHPFPKcBbRie
8mqwgjnAveqvlRKxi3VWwG0Bevki54ghwolmZ5GppvzeqLNYFF8nYuSAbseRoPFb
EJMLd5vuEkDytZl42eiZkv/aHEtUGXvcTY29K6G4yGOEr3Pr320ts8tVW4UNlBt4
0mDfBjtXAeSIcQfww/c69Pc3Xrfd3FVf4Qjo3bhMCvbg5shvRHmJrcbOPJO5kUn+
mwPU7DlJM9YeOMQBMgmw3NoKKI4dOU3HUBpBiN3M5tztAgMBAAGjUzBRMB0GA1Ud
DgQWBBRCecVMYml89VvfhBl+DTxzqcwWoTAfBgNVHSMEGDAWgBRCecVMYml89Vvf
hBl+DTxzqcwWoTAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQDJ
oF3c5MGbrg5Iu8TF9X8g3tgVANGUVWDHDxx0Fe3zQojW0NSfHPtQ+Qkf6BdYH6oc
OzQBUgWFGnrhPliUW9rD4/8c3BoVvT9ukYPOhwLDd2lPqTTbbfhdkzDSM/BKPP1g
7Ok2m/uk9jnsLQSCQE4zc8+X0M+zG9ZPyC0MJqM3d7gB+LVOE8PKIJz6fXCyoakz
18PV+e4RhL5daTFCdZ1XAL146kIorS4XX5iIyvCt1WBzSS8IUtAIgR/QLxk7ZqrL
BKEkcU1X0yvgyDUkpcJ1BS/++5q/EDEQCYP6gN0cvPFhFvQeNor5SId6EFFlEkMn
MYuea4TP5Gk2UkmDOxuJ
-----END CERTIFICATE-----`

const privateKey = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDkVU596LWGR+/n
C3r3MndJfimGMDvt9W4SwL0bOa+VXxgVgUKYTTCPvwZaAQrtJRUjW2bGxwwj8/3u
DEY6vwHJ1Yh+OrbrQHPFPKcBbRie8mqwgjnAveqvlRKxi3VWwG0Bevki54ghwolm
Z5GppvzeqLNYFF8nYuSAbseRoPFbEJMLd5vuEkDytZl42eiZkv/aHEtUGXvcTY29
K6G4yGOEr3Pr320ts8tVW4UNlBt40mDfBjtXAeSIcQfww/c69Pc3Xrfd3FVf4Qjo
3bhMCvbg5shvRHmJrcbOPJO5kUn+mwPU7DlJM9YeOMQBMgmw3NoKKI4dOU3HUBpB
iN3M5tztAgMBAAECggEBANJZQFlAA9Kz/O+VpO+L/1amMmzbjKo4evItu0kUiIwM
MezFyurx2XXjnl9WLJGxotqSvokLIEUS5vDhP+Wox2YAIKFhR9hL5RtkN9pZfeAY
JW98WOYWT9j3dWQ4vJ1x4joF5vRf5gpr5BaB/TAUlUoukiHnio2HTkh/Rb0ETrT7
Pvl9hYFO50xmaxwd5Vy+726ZLwOkkraDpXB1jZC9Kp7EfnMi7ekZ8LfBYmEdl87Y
VvBghjSsRL4VdY/WTOpWM1DnOIBrUmM/0UfYW1uaV4upSPScOjFeBllY+lSpyO38
B+L1eQSJghIULOntN5XUGnrTpMSXW8C67qaEFPfa9qECgYEA/UyAQCEXuFiEDpia
CkZ0Ykh6xxY2sA6jMB52RvjpWxqbrVUE6yGMM2UJxNplwdZk9lmpzU9KPfDgulKX
Uq34O94wDSXKJQYI0GsrXs3IgheXKVT/4s3t9oCc2hH0F3/o37jYOkYP571e2sdd
yQd0aTZqG1qp7bZlRKWahKrB/FUCgYEA5sSmc7dIwxgX/kp/4VUdjgFUy/GU5xr9
6xnioGsdnb6rBpicklri/h2E4eLLzgVbuzVQMLIAG/MpwrIxWspUR64yAaPEAVm8
3GIg771JZHl06lYYjAqaSy1qC8v3/3T+masWwa/MNCxXB6YN9ptohAf4M3hlEL/J
jlR8Qp5M9jkCgYAqcPgIRsM3szUlUPJ2iEmV8jkIRLOTGlDDjkcZKznGdxXgnB8/
2pYoQmS5pDJqoSa3lFx8Ny3kZQjyj0Ylp1qxhVAd09gkDffKHDrfHrHbAmLknQZn
FUQrCm+9pkZ07Yyyd8FbOkQN+0/6bm9LcMFTo7dxr+ZLG0Wqk+jpE8d/JQKBgCFf
s0rs6OL+KwolsBTggGO3IZJVH9nEd5B2r+XPV/smRgmwLISmDEn/7uXULPFgqQGM
FkrUk1t3cUStDKI6vLGZKbY+/uvLFJsyvdyuHV0gi54QUYB/UA0rRjjqiLUzMFb8
/U+JoxiwiO2cQEy38QeXN3gKI2OmuPmSkl34Et1RAoGBAIqPAdXfyoMb6stylvC1
N20fcwpG3aiTESteYpnXCNFW8XrMnoBWL6bK6st4eBSUbvOfjTJrSVC/KBLR6awV
i+U582LTWq8y6WA8tdqfeZO+TUl+8DBk6k6aDbA8a3+X/D+sTsRfSavEVyEeV7EO
wkv+4ThHJ677Dpi/P8F8iOJp
-----END PRIVATE KEY-----`

func TestAccCertificateResAndDataSource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCertificateDestroyCheck,
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
	client := testAccProvider.Meta().(services.SdkBundle).CertManagerClient
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
` + testCertificate + `
EOT
	certificate_chain = <<EOT
` + testCertificate + `
EOT
	private_key 	  = <<EOT
` + privateKey + `
EOT
}
`
	testAccCheckCertUpdateName = `
resource ` + constant.CertificateResource + ` ` + constant.TestCertName + ` {
	name        	  = "` + constant.TestCertName + `1"
	certificate 	  = <<EOT
` + testCertificate + `
EOT
	certificate_chain = <<EOT
` + testCertificate + `
EOT
	private_key 	  = <<EOT
` + privateKey + `
EOT
}
`
	testAccCheckDataSourceByName = `
resource ` + constant.CertificateResource + ` ` + constant.TestCertName + ` {
	name        	  = "` + constant.TestCertName + `1"
	certificate 	  = <<EOT
` + testCertificate + `
EOT
	certificate_chain = <<EOT
` + testCertificate + `
EOT
	private_key 	  = <<EOT
` + privateKey + `
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
` + testCertificate + `
EOT
	certificate_chain = <<EOT
` + testCertificate + `
EOT
	private_key 	  = <<EOT
` + privateKey + `
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
` + testCertificate + `
EOT
	certificate_chain = <<EOT
` + testCertificate + `
EOT
	private_key 	  = <<EOT
` + privateKey + `
EOT
}
` + constant.DataSource + ` ` + constant.CertificateResource + ` ` + constant.TestCertName + ` {
name ="should_not_work"
}
`
)
