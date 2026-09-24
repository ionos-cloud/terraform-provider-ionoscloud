//go:build compute || all

package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

const (
	importGenServerAddr   = "ionoscloud_server.import_gen"
	importGenVolumeAddr   = "ionoscloud_volume.import_gen_unattached"
	importGenGeneratedSrv = "ionoscloud_server.generated"
	importGenGeneratedVol = "ionoscloud_volume.generated"
	importGenLocationAttr = "location"
	importGenSSHKeysAttr  = "ssh_key_path.#"
	importGenLicenceAttr  = "licence_type"
)

// TestAccImportGeneratedConfigNoOp imports resources the way a whole-datacenter import does: an
// import block plus generated config, which must plan as a pure no-op import. It covers a server
// without a boot volume whose NIC is a separate ionoscloud_nic, and a volume attached to no server.
func TestAccImportGeneratedConfigNoOp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckVolumeDestroyCheck,
			testAccCheckServerDestroyCheck,
			testAccCheckDatacenterDestroyCheck,
		),
		Steps: []resource.TestStep{
			{
				// The post-apply plan must be empty: a server without a boot volume used to plan
				// inline_volume_ids as "known after apply" forever.
				Config: testAccImportGenConfig(""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(importGenServerAddr, "inline_volume_ids.#", "0"),
					resource.TestCheckResourceAttr(importGenServerAddr, "boot_volume", ""),
					resource.TestCheckNoResourceAttr(importGenVolumeAddr, "server_id"),
				),
			},
			{
				ResourceName:      importGenServerAddr,
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateIdFunc: testAccImportGenServerID,
				GenerateConfig:    true,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(importGenGeneratedSrv, plancheck.ResourceActionNoop),
						// dc/srv leaves the NIC to its own ionoscloud_nic resource.
						plancheck.ExpectKnownValue(importGenGeneratedSrv, tfjsonpath.New("nic"), knownvalue.ListSizeExact(0)),
						plancheck.ExpectKnownValue(importGenGeneratedSrv, tfjsonpath.New("primary_nic"), knownvalue.Null()),
						plancheck.ExpectKnownValue(importGenGeneratedSrv, tfjsonpath.New("inline_volume_ids"), knownvalue.ListSizeExact(0)),
					},
				},
			},
			{
				ResourceName:      importGenVolumeAddr,
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateIdFunc: testAccImportGenUnattachedVolumeID,
				GenerateConfig:    true,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(importGenGeneratedVol, plancheck.ResourceActionNoop),
						plancheck.ExpectKnownValue(importGenGeneratedVol, tfjsonpath.New("server_id"), knownvalue.Null()),
					},
				},
			},
			{
				ResourceName:      importGenVolumeAddr,
				ImportState:       true,
				ImportStateIdFunc: testAccImportGenUnattachedVolumeID,
				ImportStateVerify: true,
				// Create writes ssh_key_path = [] and the configured licence_type; neither is read
				// back from the API, and import always sets location.
				ImportStateVerifyIgnore: []string{importGenSSHKeysAttr, importGenLicenceAttr, importGenLocationAttr},
			},
			{
				// Attaching a previously unattached volume must not need a volume PATCH or replace.
				Config: testAccImportGenConfig("server_id = ionoscloud_server.import_gen.id"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(importGenVolumeAddr, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(importGenVolumeAddr, "server_id", importGenServerAddr, "id"),
				),
			},
			{
				ResourceName:      importGenVolumeAddr,
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateIdFunc: testAccImportGenAttachedVolumeID,
				GenerateConfig:    true,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(importGenGeneratedVol, plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

func testAccImportGenServerID(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources[importGenServerAddr]
	if !ok {
		return "", fmt.Errorf("resource %s not found in state", importGenServerAddr)
	}
	return fmt.Sprintf("%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.ID), nil
}

func testAccImportGenUnattachedVolumeID(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources[importGenVolumeAddr]
	if !ok {
		return "", fmt.Errorf("resource %s not found in state", importGenVolumeAddr)
	}
	return fmt.Sprintf("%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.ID), nil
}

func testAccImportGenAttachedVolumeID(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources[importGenVolumeAddr]
	if !ok {
		return "", fmt.Errorf("resource %s not found in state", importGenVolumeAddr)
	}
	return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["server_id"], rs.Primary.ID), nil
}

// testAccImportGenConfig is the smallest datacenter holding the shapes the import has to handle.
// volumeAttachment is inserted into the data volume, so a step can attach it to the server.
func testAccImportGenConfig(volumeAttachment string) string {
	return fmt.Sprintf(`
resource "ionoscloud_datacenter" "import_gen" {
  name     = "tf-acctest-import-gen"
  location = "de/txl"
}

resource "ionoscloud_lan" "import_gen" {
  datacenter_id = ionoscloud_datacenter.import_gen.id
  public        = false
  name          = "tf-acctest-import-gen-lan"
}

# No volume block: the server has no boot volume.
resource "ionoscloud_server" "import_gen" {
  name          = "tf-acctest-import-gen-srv"
  datacenter_id = ionoscloud_datacenter.import_gen.id
  cores         = 1
  ram           = 1024
}

resource "ionoscloud_nic" "import_gen" {
  datacenter_id = ionoscloud_datacenter.import_gen.id
  server_id     = ionoscloud_server.import_gen.id
  lan           = ionoscloud_lan.import_gen.id
  name          = "tf-acctest-import-gen-nic"
  dhcp          = true
}

resource "ionoscloud_volume" "import_gen_unattached" {
  datacenter_id = ionoscloud_datacenter.import_gen.id
  name          = "tf-acctest-import-gen-vol"
  size          = 1
  disk_type     = "HDD"
  licence_type  = "UNKNOWN"
  %s
}
`, volumeAttachment)
}
