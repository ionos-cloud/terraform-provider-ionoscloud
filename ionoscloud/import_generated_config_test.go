//go:build compute || all

package ionoscloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
)

const (
	importGenServerAddr   = "ionoscloud_server.import_gen"
	importGenOtherSrvAddr = "ionoscloud_server.import_gen_other"
	importGenVolumeAddr   = "ionoscloud_volume.import_gen_unattached"
	importGenGeneratedSrv = "ionoscloud_server.generated"
	importGenGeneratedVol = "ionoscloud_volume.generated"
	importGenLocationAttr = "location"
	importGenSSHKeysAttr  = "ssh_key_path.#"
	importGenLicenceAttr  = "licence_type"
)

// TestAccImportGeneratedConfigNoOp imports resources the way a whole-datacenter import does: an
// import block plus generated config, which must plan as a pure no-op import. It covers a server
// without a boot volume whose NIC is a separate ionoscloud_nic, and a volume attached to no server,
// then attached, moved to another server and detached again.
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
					testAccCheckImportGenAttachedVolumes(importGenServerAddr, 1),
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
			{
				// Moving the volume to another server detaches it from the first one before
				// attaching it, in place, and the next plan is empty.
				Config: testAccImportGenConfig("server_id = ionoscloud_server.import_gen_other.id"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(importGenVolumeAddr, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(importGenVolumeAddr, "server_id", importGenOtherSrvAddr, "id"),
					testAccCheckImportGenAttachedVolumes(importGenServerAddr, 0),
					testAccCheckImportGenAttachedVolumes(importGenOtherSrvAddr, 1),
				),
			},
			{
				// Removing server_id only detaches the volume: it is updated in place, the API no
				// longer lists it on any server, and the next plan is empty.
				Config: testAccImportGenConfig(""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(importGenVolumeAddr, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					// SDKv2 stores an optional string removed from config as "", not null, so the
					// API check below is what proves the volume is detached.
					resource.TestCheckResourceAttr(importGenVolumeAddr, "server_id", ""),
					testAccCheckImportGenAttachedVolumes(importGenServerAddr, 0),
					testAccCheckImportGenAttachedVolumes(importGenOtherSrvAddr, 0),
				),
			},
		},
	})
}

// testAccCheckImportGenAttachedVolumes checks on the API, not in state, how many volumes are
// attached to the server at serverAddr.
func testAccCheckImportGenAttachedVolumes(serverAddr string, want int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[serverAddr]
		if !ok {
			return fmt.Errorf("resource %s not found in state", serverAddr)
		}
		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		defer cancel()

		client, err := testAccProvider.Meta().(bundleclient.SdkBundle).NewCloudAPIClient(ctx, rs.Primary.Attributes["location"])
		if err != nil {
			return err
		}
		volumes, apiResponse, err := client.ServersApi.DatacentersServersVolumesGet(ctx, rs.Primary.Attributes["datacenter_id"], rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return fmt.Errorf("listing volumes attached to server %s: %w", rs.Primary.ID, err)
		}
		got := 0
		if volumes.Items != nil {
			got = len(*volumes.Items)
		}
		if got != want {
			return fmt.Errorf("server %s has %d attached volume(s), want %d", rs.Primary.ID, got, want)
		}
		return nil
	}
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
// volumeAttachment is inserted into the data volume, so a step can attach it to a server.
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

# A second server, so a step can move the volume between servers.
resource "ionoscloud_server" "import_gen_other" {
  name          = "tf-acctest-import-gen-srv2"
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

const (
	importGenFailoverLanAddr = "ionoscloud_lan.import_gen_failover"
	importGenFailoverIPBAddr = "ionoscloud_ipblock.import_gen_failover"
	importGenGeneratedLan    = "ionoscloud_lan.generated"
	importGenGeneratedIPB    = "ionoscloud_ipblock.generated"
)

// TestAccImportGeneratedConfigComputedBlocksNoOp imports a LAN with an IP failover and the IP
// block it uses. Generated config writes one empty ip_failover {} / ip_consumers {} block per
// entry, because every nested attribute is computed. Dropping Optional would stop that but break
// configurations that already hold such blocks, so the blocks stay and must plan as a no-op.
func TestAccImportGeneratedConfigComputedBlocksNoOp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckIPBlockDestroyCheck,
			testAccCheckDatacenterDestroyCheck,
		),
		Steps: []resource.TestStep{
			{
				// The LAN and IP block are read before the failover exists; the import steps
				// below read them afresh and assert one entry each.
				Config: testAccImportGenFailoverConfig,
			},
			{
				ResourceName:    importGenFailoverLanAddr,
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithID,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[importGenFailoverLanAddr]
					if !ok {
						return "", fmt.Errorf("resource %s not found in state", importGenFailoverLanAddr)
					}
					return fmt.Sprintf("%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.ID), nil
				},
				GenerateConfig: true,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(importGenGeneratedLan, plancheck.ResourceActionNoop),
						plancheck.ExpectKnownValue(importGenGeneratedLan, tfjsonpath.New("ip_failover"), knownvalue.ListSizeExact(1)),
					},
				},
			},
			{
				ResourceName:    importGenFailoverIPBAddr,
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithID,
				GenerateConfig:  true,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(importGenGeneratedIPB, plancheck.ResourceActionNoop),
						plancheck.ExpectKnownValue(importGenGeneratedIPB, tfjsonpath.New("ip_consumers"), knownvalue.ListSizeExact(1)),
					},
				},
			},
		},
	})
}

const testAccImportGenFailoverConfig = `
resource "ionoscloud_datacenter" "import_gen_failover" {
  name     = "tf-acctest-import-gen-failover"
  location = "de/txl"
}

resource "ionoscloud_ipblock" "import_gen_failover" {
  location = ionoscloud_datacenter.import_gen_failover.location
  size     = 1
  name     = "tf-acctest-import-gen-failover"
}

resource "ionoscloud_lan" "import_gen_failover" {
  datacenter_id = ionoscloud_datacenter.import_gen_failover.id
  public        = true
  name          = "tf-acctest-import-gen-failover"
}

resource "ionoscloud_server" "import_gen_failover" {
  name          = "tf-acctest-import-gen-failover"
  datacenter_id = ionoscloud_datacenter.import_gen_failover.id
  cores         = 1
  ram           = 1024
}

resource "ionoscloud_nic" "import_gen_failover" {
  datacenter_id = ionoscloud_datacenter.import_gen_failover.id
  server_id     = ionoscloud_server.import_gen_failover.id
  lan           = ionoscloud_lan.import_gen_failover.id
  name          = "tf-acctest-import-gen-failover"
  dhcp          = true
  ips           = [ionoscloud_ipblock.import_gen_failover.ips[0]]
}

resource "ionoscloud_ipfailover" "import_gen_failover" {
  datacenter_id = ionoscloud_datacenter.import_gen_failover.id
  lan_id        = ionoscloud_lan.import_gen_failover.id
  ip            = ionoscloud_ipblock.import_gen_failover.ips[0]
  nicuuid       = ionoscloud_nic.import_gen_failover.id
}
`
