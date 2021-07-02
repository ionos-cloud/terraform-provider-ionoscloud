package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPrivateCrossConnect_Basic(t *testing.T) {
	var privateCrossConnect ionoscloud.PrivateCrossConnect
	privateCrossConnectName := "example"
	privateCrossConnectDescription := "example-description"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckprivateCrossConnectDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckprivateCrossConnectConfigBasic, privateCrossConnectName, privateCrossConnectDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckprivateCrossConnectExists("ionoscloud_private_crossconnect.example", &privateCrossConnect),
					resource.TestCheckResourceAttr("ionoscloud_private_crossconnect.example", "name", privateCrossConnectName),
					resource.TestCheckResourceAttr("ionoscloud_private_crossconnect.example", "description", "example-description"),
				),
			},
			{
				Config: testAccCheckprivateCrossConnectConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckprivateCrossConnectExists("ionoscloud_private_crossconnect.example", &privateCrossConnect),
					resource.TestCheckResourceAttr("ionoscloud_private_crossconnect.example", "name", "example-renamed"),
					resource.TestCheckResourceAttr("ionoscloud_private_crossconnect.example", "description", "example-description-updated"),
				),
			},
		},
	})
}

func testAccCheckprivateCrossConnectDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_private_crossconnect" {
			continue
		}

		_, apiResponse, err := client.PrivateCrossConnectsApi.PccsFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while checking private cross-connect %s: %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("unable to fetch private cross-connect %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckprivateCrossConnectExists(n string, privateCrossConnect *ionoscloud.PrivateCrossConnect) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.APIClient)

		rs, ok := s.RootModule().Resources[n]
		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		if cancel != nil {
			defer cancel()
		}

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		foundPrivateCrossConnect, _, err := client.PrivateCrossConnectsApi.PccsFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			return fmt.Errorf("error occured while fetching private cross-connect: %s", rs.Primary.ID)
		}
		if *foundPrivateCrossConnect.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		privateCrossConnect = &foundPrivateCrossConnect

		return nil
	}
}

const testAccCheckprivateCrossConnectConfigBasic = `
resource "ionoscloud_private_crossconnect" "example" {
  name        = "%s"
  description = "%s"
}`

const testAccCheckprivateCrossConnectConfigUpdate = `
resource "ionoscloud_private_crossconnect" "example" {
  name        = "example-renamed"
  description = "example-description-updated"
}`
