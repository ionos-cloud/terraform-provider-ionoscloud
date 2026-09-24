//go:build compute || all || server || cube

package ionoscloud

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
)

const importNoNicCubeAddr = "ionoscloud_cube_server.import_no_nic"

// TestAccCubeServerImportWithoutNic imports a cube server whose only NIC was deleted outside
// Terraform. The shared cube/GPU importer used to index the first NIC unchecked and panicked.
func TestAccCubeServerImportWithoutNic(t *testing.T) {
	var dcID, serverID, nicID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckCubeServerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCubeServerImportNoNicConfig,
				Check: func(s *terraform.State) error {
					rs, ok := s.RootModule().Resources[importNoNicCubeAddr]
					if !ok {
						return fmt.Errorf("resource %s not found in state", importNoNicCubeAddr)
					}
					dcID, serverID, nicID = rs.Primary.Attributes["datacenter_id"], rs.Primary.ID, rs.Primary.Attributes["primary_nic"]
					if nicID == "" {
						return fmt.Errorf("%s has no primary_nic", importNoNicCubeAddr)
					}
					return nil
				},
			},
			{
				PreConfig: func() {
					if err := testAccDeleteServerNic(dcID, serverID, nicID); err != nil {
						t.Fatal(err)
					}
				},
				ResourceName: importNoNicCubeAddr,
				ImportState:  true,
				ImportStateIdFunc: func(*terraform.State) (string, error) {
					return dcID + "/" + serverID, nil
				},
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected 1 imported state, got %d", len(states))
					}
					attrs := states[0].Attributes
					if attrs["primary_nic"] != "" || attrs["primary_ip"] != "" {
						return fmt.Errorf("imported a NIC-less cube with primary_nic=%q primary_ip=%q", attrs["primary_nic"], attrs["primary_ip"])
					}
					if count := attrs["nic.#"]; count != "" && count != "0" {
						return fmt.Errorf("imported a NIC-less cube with nic.# = %s", count)
					}
					return nil
				},
			},
		},
	})
}

// testAccDeleteServerNic deletes a NIC behind Terraform's back and waits until it is gone.
func testAccDeleteServerNic(dcID, serverID, nicID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	client, err := testAccProvider.Meta().(bundleclient.SdkBundle).NewCloudAPIClient(ctx, "")
	if err != nil {
		return err
	}

	apiResponse, err := client.NetworkInterfacesApi.DatacentersServersNicsDelete(ctx, dcID, serverID, nicID).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		return fmt.Errorf("deleting nic %s: %w", nicID, err)
	}

	for {
		_, apiResponse, err := client.NetworkInterfacesApi.DatacentersServersNicsFindById(ctx, dcID, serverID, nicID).Execute()
		logApiRequestTime(apiResponse)
		if httpNotFound(apiResponse) {
			return nil
		}
		if err != nil && ctx.Err() == nil && apiResponse.SafeStatusCode() < 500 {
			return fmt.Errorf("waiting for nic %s to be deleted: %w", nicID, err)
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for nic %s to be deleted: %w", nicID, ctx.Err())
		case <-time.After(10 * time.Second):
		}
	}
}

const testAccCubeServerImportNoNicConfig = `
data "ionoscloud_template" "import_no_nic" {
  name         = "Basic Cube XS"
  cores        = 1
  ram          = 2048
  storage_size = 60
}

resource "ionoscloud_datacenter" "import_no_nic" {
  name     = "tf-acctest-import-cube-no-nic"
  location = "de/txl"
}

resource "ionoscloud_lan" "import_no_nic" {
  datacenter_id = ionoscloud_datacenter.import_no_nic.id
  public        = false
  name          = "tf-acctest-import-cube-no-nic-lan"
}

resource "random_password" "import_no_nic" {
  length  = 16
  special = false
}

resource "ionoscloud_cube_server" "import_no_nic" {
  name              = "tf-acctest-import-cube-no-nic"
  datacenter_id     = ionoscloud_datacenter.import_no_nic.id
  availability_zone = "AUTO"
  image_name        = "ubuntu:latest"
  image_password    = random_password.import_no_nic.result
  template_uuid     = data.ionoscloud_template.import_no_nic.id

  volume {
    name         = "tf-acctest-import-cube-no-nic"
    licence_type = "LINUX"
    disk_type    = "DAS"
  }

  nic {
    lan  = ionoscloud_lan.import_no_nic.id
    name = "tf-acctest-import-cube-no-nic"
    dhcp = true
  }
}
`
