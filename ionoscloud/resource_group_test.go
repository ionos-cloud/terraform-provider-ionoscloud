//go:build compute || all || group

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccGroupBasic(t *testing.T) {
	var group ionoscloud.Group

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGroupConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupExists(constant.GroupResource+"."+constant.GroupTestResource, &group),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "name", constant.GroupTestResource),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_datacenter", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_snapshot", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "reserve_ip", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_activity_log", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_pcc", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "s3_privilege", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_backup_unit", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_internet_access", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_k8s_cluster", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_flow_log", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_monitoring", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_certificates", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_dns", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "manage_registry", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "manage_dataplatform", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_logging", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_cdn", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_vpn", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_api_gateway", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_kaas", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_network_file_storage", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_ai_model_hub", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_iam_resources", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_network_security_groups", "true"),
					utils.TestNotEmptySlice(constant.GroupResource, "users")),
			},
			{
				Config: testAccCheckGroupConfigBasicNoUsersData,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupExists(constant.GroupResource+"."+constant.GroupTestResource, &group),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "name", constant.GroupTestResource),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_datacenter", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_snapshot", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "reserve_ip", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_activity_log", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_pcc", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "s3_privilege", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_backup_unit", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_internet_access", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_k8s_cluster", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_flow_log", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_monitoring", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_certificates", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_dns", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "manage_registry", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "manage_dataplatform", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_logging", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_cdn", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_vpn", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_api_gateway", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_kaas", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_network_file_storage", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_ai_model_hub", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_iam_resources", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_network_security_groups", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "users.#", "0"),
				),
			},
			{
				Config: testAccCheckGroupConfigBasicGetUsersData,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupExists(constant.GroupResource+"."+constant.GroupTestResource, &group),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "name", constant.GroupTestResource),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_datacenter", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_snapshot", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "reserve_ip", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_activity_log", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_pcc", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "s3_privilege", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_backup_unit", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_internet_access", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_k8s_cluster", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_flow_log", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_monitoring", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_certificates", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_dns", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "manage_registry", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "manage_dataplatform", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_logging", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_cdn", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_vpn", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_api_gateway", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_kaas", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_network_file_storage", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_ai_model_hub", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_iam_resources", "true"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_network_security_groups", "true"),
					utils.TestNotEmptySlice(constant.GroupResource, "users")),
			},
			{
				Config: testAccDataSourceGroupMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "name",
						constant.GroupResource+"."+constant.GroupTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "create_datacenter", constant.GroupResource+"."+constant.GroupTestResource, "create_datacenter"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "create_snapshot", constant.GroupResource+"."+constant.GroupTestResource, "create_snapshot"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "reserve_ip", constant.GroupResource+"."+constant.GroupTestResource, "reserve_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_activity_log", constant.GroupResource+"."+constant.GroupTestResource, "access_activity_log"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "create_pcc", constant.GroupResource+"."+constant.GroupTestResource, "create_pcc"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "s3_privilege", constant.GroupResource+"."+constant.GroupTestResource, "s3_privilege"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "create_backup_unit", constant.GroupResource+"."+constant.GroupTestResource, "create_backup_unit"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "create_internet_access", constant.GroupResource+"."+constant.GroupTestResource, "create_internet_access"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "create_k8s_cluster", constant.GroupResource+"."+constant.GroupTestResource, "create_k8s_cluster"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_dns", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_dns"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "manage_registry", constant.GroupResource+"."+constant.GroupTestResource, "manage_registry"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "manage_dataplatform", constant.GroupResource+"."+constant.GroupTestResource, "manage_dataplatform"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_logging", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_logging"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_cdn", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_cdn"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_vpn", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_vpn"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_api_gateway", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_api_gateway"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_kaas", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_kaas"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_network_file_storage", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_network_file_storage"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_ai_model_hub", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_ai_model_hub"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_iam_resources", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_iam_resources"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "create_network_security_groups", constant.GroupResource+"."+constant.GroupTestResource, "create_network_security_groups"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "manage_dbaas", constant.GroupResource+"."+constant.GroupTestResource, "manage_dbaas"),
					utils.TestNotEmptySlice(constant.DataSource+"."+constant.GroupResource, "users"),
				),
			},
			{
				Config: testAccDataSourceGroupMatchIdGetUsersData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "name",
						constant.GroupResource+"."+constant.GroupTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "create_datacenter", constant.GroupResource+"."+constant.GroupTestResource, "create_datacenter"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "create_snapshot", constant.GroupResource+"."+constant.GroupTestResource, "create_snapshot"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "reserve_ip", constant.GroupResource+"."+constant.GroupTestResource, "reserve_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_activity_log", constant.GroupResource+"."+constant.GroupTestResource, "access_activity_log"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "create_pcc", constant.GroupResource+"."+constant.GroupTestResource, "create_pcc"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "s3_privilege", constant.GroupResource+"."+constant.GroupTestResource, "s3_privilege"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "create_backup_unit", constant.GroupResource+"."+constant.GroupTestResource, "create_backup_unit"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "create_internet_access", constant.GroupResource+"."+constant.GroupTestResource, "create_internet_access"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "create_k8s_cluster", constant.GroupResource+"."+constant.GroupTestResource, "create_k8s_cluster"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_dns", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_dns"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "manage_registry", constant.GroupResource+"."+constant.GroupTestResource, "manage_registry"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "manage_dataplatform", constant.GroupResource+"."+constant.GroupTestResource, "manage_dataplatform"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_logging", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_logging"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_cdn", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_cdn"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_vpn", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_vpn"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_api_gateway", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_api_gateway"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_kaas", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_kaas"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_network_file_storage", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_network_file_storage"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_ai_model_hub", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_ai_model_hub"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_iam_resources", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_iam_resources"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "create_network_security_groups", constant.GroupResource+"."+constant.GroupTestResource, "create_network_security_groups"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "manage_dbaas", constant.GroupResource+"."+constant.GroupTestResource, "manage_dbaas"),
					utils.TestNotEmptySlice(constant.DataSource+"."+constant.GroupResource, "users"),
				),
			},
			{
				Config: testAccDataSourceGroupMatchIdNoUsersData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "name",
						constant.GroupResource+"."+constant.GroupTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "create_datacenter", constant.GroupResource+"."+constant.GroupTestResource, "create_datacenter"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "create_snapshot", constant.GroupResource+"."+constant.GroupTestResource, "create_snapshot"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "reserve_ip", constant.GroupResource+"."+constant.GroupTestResource, "reserve_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_activity_log", constant.GroupResource+"."+constant.GroupTestResource, "access_activity_log"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "create_pcc", constant.GroupResource+"."+constant.GroupTestResource, "create_pcc"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "s3_privilege", constant.GroupResource+"."+constant.GroupTestResource, "s3_privilege"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "create_backup_unit", constant.GroupResource+"."+constant.GroupTestResource, "create_backup_unit"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "create_internet_access", constant.GroupResource+"."+constant.GroupTestResource, "create_internet_access"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "create_k8s_cluster", constant.GroupResource+"."+constant.GroupTestResource, "create_k8s_cluster"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_dns", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_dns"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "manage_registry", constant.GroupResource+"."+constant.GroupTestResource, "manage_registry"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "manage_dataplatform", constant.GroupResource+"."+constant.GroupTestResource, "manage_dataplatform"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_logging", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_logging"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_cdn", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_cdn"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_vpn", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_vpn"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_api_gateway", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_api_gateway"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_kaas", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_kaas"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_network_file_storage", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_network_file_storage"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_ai_model_hub", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_ai_model_hub"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "access_and_manage_iam_resources", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_iam_resources"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "create_network_security_groups", constant.GroupResource+"."+constant.GroupTestResource, "create_network_security_groups"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "manage_dbaas", constant.GroupResource+"."+constant.GroupTestResource, "manage_dbaas"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceById, "users.#", "0"),
				),
			},
			{
				Config: testAccDataSourceGroupMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "name", constant.GroupResource+"."+constant.GroupTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "create_datacenter", constant.GroupResource+"."+constant.GroupTestResource, "create_datacenter"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "create_snapshot", constant.GroupResource+"."+constant.GroupTestResource, "create_snapshot"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "reserve_ip", constant.GroupResource+"."+constant.GroupTestResource, "reserve_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_activity_log", constant.GroupResource+"."+constant.GroupTestResource, "access_activity_log"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "create_pcc", constant.GroupResource+"."+constant.GroupTestResource, "create_pcc"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "s3_privilege", constant.GroupResource+"."+constant.GroupTestResource, "s3_privilege"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "create_backup_unit", constant.GroupResource+"."+constant.GroupTestResource, "create_backup_unit"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "create_internet_access", constant.GroupResource+"."+constant.GroupTestResource, "create_internet_access"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "create_k8s_cluster", constant.GroupResource+"."+constant.GroupTestResource, "create_k8s_cluster"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "create_flow_log", constant.GroupResource+"."+constant.GroupTestResource, "create_flow_log"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_monitoring", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_monitoring"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_certificates", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_certificates"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_dns", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_dns"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "manage_registry", constant.GroupResource+"."+constant.GroupTestResource, "manage_registry"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "manage_dataplatform", constant.GroupResource+"."+constant.GroupTestResource, "manage_dataplatform"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_logging", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_logging"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_cdn", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_cdn"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_vpn", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_vpn"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_api_gateway", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_api_gateway"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_kaas", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_kaas"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_network_file_storage", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_network_file_storage"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_ai_model_hub", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_ai_model_hub"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_iam_resources", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_iam_resources"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "create_network_security_groups", constant.GroupResource+"."+constant.GroupTestResource, "create_network_security_groups"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "manage_dbaas", constant.GroupResource+"."+constant.GroupTestResource, "manage_dbaas"),
					utils.TestNotEmptySlice(constant.DataSource+"."+constant.GroupResource, "users"),
				),
			},
			{
				Config: testAccDataSourceGroupMatchNameGetUsersData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "name", constant.GroupResource+"."+constant.GroupTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "create_datacenter", constant.GroupResource+"."+constant.GroupTestResource, "create_datacenter"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "create_snapshot", constant.GroupResource+"."+constant.GroupTestResource, "create_snapshot"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "reserve_ip", constant.GroupResource+"."+constant.GroupTestResource, "reserve_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_activity_log", constant.GroupResource+"."+constant.GroupTestResource, "access_activity_log"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "create_pcc", constant.GroupResource+"."+constant.GroupTestResource, "create_pcc"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "s3_privilege", constant.GroupResource+"."+constant.GroupTestResource, "s3_privilege"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "create_backup_unit", constant.GroupResource+"."+constant.GroupTestResource, "create_backup_unit"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "create_internet_access", constant.GroupResource+"."+constant.GroupTestResource, "create_internet_access"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "create_k8s_cluster", constant.GroupResource+"."+constant.GroupTestResource, "create_k8s_cluster"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "create_flow_log", constant.GroupResource+"."+constant.GroupTestResource, "create_flow_log"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_monitoring", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_monitoring"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_certificates", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_certificates"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_dns", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_dns"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "manage_registry", constant.GroupResource+"."+constant.GroupTestResource, "manage_registry"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "manage_dataplatform", constant.GroupResource+"."+constant.GroupTestResource, "manage_dataplatform"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_logging", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_logging"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_cdn", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_cdn"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_vpn", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_vpn"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_api_gateway", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_api_gateway"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_kaas", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_kaas"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_network_file_storage", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_network_file_storage"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_ai_model_hub", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_ai_model_hub"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_iam_resources", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_iam_resources"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "create_network_security_groups", constant.GroupResource+"."+constant.GroupTestResource, "create_network_security_groups"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "manage_dbaas", constant.GroupResource+"."+constant.GroupTestResource, "manage_dbaas"),
					utils.TestNotEmptySlice(constant.DataSource+"."+constant.GroupResource, "users"),
				),
			},
			{
				Config: testAccDataSourceGroupMatchNameNoUsersData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "name", constant.GroupResource+"."+constant.GroupTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "create_datacenter", constant.GroupResource+"."+constant.GroupTestResource, "create_datacenter"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "create_snapshot", constant.GroupResource+"."+constant.GroupTestResource, "create_snapshot"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "reserve_ip", constant.GroupResource+"."+constant.GroupTestResource, "reserve_ip"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_activity_log", constant.GroupResource+"."+constant.GroupTestResource, "access_activity_log"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "create_pcc", constant.GroupResource+"."+constant.GroupTestResource, "create_pcc"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "s3_privilege", constant.GroupResource+"."+constant.GroupTestResource, "s3_privilege"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "create_backup_unit", constant.GroupResource+"."+constant.GroupTestResource, "create_backup_unit"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "create_internet_access", constant.GroupResource+"."+constant.GroupTestResource, "create_internet_access"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "create_k8s_cluster", constant.GroupResource+"."+constant.GroupTestResource, "create_k8s_cluster"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "create_flow_log", constant.GroupResource+"."+constant.GroupTestResource, "create_flow_log"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_monitoring", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_monitoring"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_certificates", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_certificates"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_dns", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_dns"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "manage_registry", constant.GroupResource+"."+constant.GroupTestResource, "manage_registry"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "manage_dataplatform", constant.GroupResource+"."+constant.GroupTestResource, "manage_dataplatform"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_logging", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_logging"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_cdn", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_cdn"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_vpn", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_vpn"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_api_gateway", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_api_gateway"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_kaas", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_kaas"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_network_file_storage", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_network_file_storage"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_ai_model_hub", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_ai_model_hub"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "access_and_manage_iam_resources", constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_iam_resources"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "create_network_security_groups", constant.GroupResource+"."+constant.GroupTestResource, "create_network_security_groups"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "manage_dbaas", constant.GroupResource+"."+constant.GroupTestResource, "manage_dbaas"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.GroupResource+"."+constant.GroupDataSourceByName, "users.#", "0"),
				),
			},
			{
				Config:      testAccDataSourceGroupMultipleResultsError,
				ExpectError: regexp.MustCompile("more than one group found with the specified criteria"),
			},
			{
				Config:      testAccDataSourceGroupWrongNameError,
				ExpectError: regexp.MustCompile("no group found with the specified name"),
			},
			{
				Config: testAccCheckGroupConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupExists(constant.GroupResource+"."+constant.GroupTestResource, &group),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_datacenter", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_snapshot", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "reserve_ip", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_activity_log", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_pcc", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "s3_privilege", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_backup_unit", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_internet_access", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_k8s_cluster", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_flow_log", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_monitoring", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_certificates", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_dns", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "manage_registry", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "manage_dataplatform", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_logging", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_cdn", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_vpn", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_api_gateway", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_kaas", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_network_file_storage", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_ai_model_hub", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_iam_resources", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_network_security_groups", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "manage_dbaas", "false"),
					resource.TestCheckResourceAttrPair(constant.GroupResource+".test_user_id", "users.0.id", constant.UserResource+"."+constant.UserTestResource+"3", "id"),
					utils.TestNotEmptySlice(constant.GroupResource, "users")),
			},
			{
				Config: testAccCheckGroupUpdateMigrateToUserIds,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupExists(constant.GroupResource+"."+constant.GroupTestResource, &group),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_datacenter", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_snapshot", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "reserve_ip", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_activity_log", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_pcc", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "s3_privilege", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_backup_unit", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_internet_access", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_k8s_cluster", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_flow_log", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_monitoring", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_certificates", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_dns", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "manage_registry", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "manage_dataplatform", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_logging", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_cdn", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_vpn", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_api_gateway", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_kaas", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_network_file_storage", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_ai_model_hub", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "access_and_manage_iam_resources", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "create_network_security_groups", "false"),
					resource.TestCheckResourceAttr(constant.GroupResource+"."+constant.GroupTestResource, "manage_dbaas", "false"),
					resource.TestCheckResourceAttrPair(constant.GroupResource+".test_user_id", "users.0.id", constant.UserResource+"."+constant.UserTestResource+"3", "id"),
					utils.TestNotEmptySlice(constant.GroupResource, "users")),
			},
			{
				Config:      testAccCheckGroupBothUserArgumentsError,
				ExpectError: regexp.MustCompile("Conflicting configuration arguments"),
			},
		},
	})
}

func testAccCheckGroupDestroyCheck(s *terraform.State) error {
	client, err := testAccProvider.Meta().(bundleclient.SdkBundle).NewCloudAPIClient("")
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}
	for _, rs := range s.RootModule().Resources {

		if rs.Type != constant.GroupResource {
			continue
		}
		_, apiResponse, err := client.UserManagementApi.UmGroupsFindById(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if apiResponse == nil || apiResponse.Response != nil && apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while checking the destruction of group %s: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("group %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckGroupExists(n string, group *ionoscloud.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := testAccProvider.Meta().(bundleclient.SdkBundle).NewCloudAPIClient("")
		if err != nil {
			return err
		}

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckGroupExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		foundgroup, apiResponse, err := client.UserManagementApi.UmGroupsFindById(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("error occurred while fetching Group: %s", rs.Primary.ID)
		}
		if *foundgroup.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		group = &foundgroup

		return nil
	}
}

var testAccCheckGroupCreateUsers = `
resource ` + constant.UserResource + ` ` + constant.UserTestResource + ` {
  first_name = "user"
  last_name = "test"
  email = "` + utils.GenerateEmail() + `"
  password = ` + constant.RandomPassword + `.user1_password.result
  administrator = false
  force_sec_auth= false
  active = false
}

resource ` + constant.UserResource + ` ` + constant.UserTestResource + `2 {
  first_name = "user"
  last_name = "test"
  email = "` + utils.GenerateEmail() + `"
  password = ` + constant.RandomPassword + `.user2_password.result
  administrator = false
  force_sec_auth= false
  active = false
}

resource ` + constant.RandomPassword + ` "user1_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource ` + constant.RandomPassword + ` "user2_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

`

var testAccCheckGroupConfigBasic = testAccCheckGroupCreateUsers + `
resource ` + constant.GroupResource + ` ` + constant.GroupTestResource + ` {
  name = "` + constant.GroupTestResource + `"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = true
  create_pcc = true
  s3_privilege = true
  create_backup_unit = true
  create_internet_access = true
  create_k8s_cluster = true
  create_flow_log = true
  access_and_manage_monitoring = true
  access_and_manage_certificates = true
  access_and_manage_dns = true
  manage_registry = true
  manage_dataplatform = true
  access_and_manage_logging = true
  access_and_manage_cdn = true
  access_and_manage_vpn = true
  access_and_manage_api_gateway = true
  access_and_manage_kaas = true
  access_and_manage_network_file_storage = true
  access_and_manage_ai_model_hub = true
  create_network_security_groups = true
  access_and_manage_iam_resources = true
  manage_dbaas = true


  user_ids = [` + constant.UserResource + `.` + constant.UserTestResource + `.id, ` + constant.UserResource + `.` + constant.UserTestResource + `2.id]
}
`

var testAccCheckGroupConfigBasicNoUsersData = testAccCheckGroupCreateUsers + `
resource ` + constant.GroupResource + ` ` + constant.GroupTestResource + ` {
  name = "` + constant.GroupTestResource + `"
  get_users_data = false
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = true
  create_pcc = true
  s3_privilege = true
  create_backup_unit = true
  create_internet_access = true
  create_k8s_cluster = true
  create_flow_log = true
  access_and_manage_monitoring = true
  access_and_manage_certificates = true
  access_and_manage_dns = true
  manage_registry = true
  manage_dataplatform = true
  access_and_manage_logging = true
  access_and_manage_cdn = true
  access_and_manage_vpn = true
  access_and_manage_api_gateway = true
  access_and_manage_kaas = true
  access_and_manage_network_file_storage = true
  access_and_manage_ai_model_hub = true
  create_network_security_groups = true
  access_and_manage_iam_resources = true
  manage_dbaas = true


  user_ids = [` + constant.UserResource + `.` + constant.UserTestResource + `.id, ` + constant.UserResource + `.` + constant.UserTestResource + `2.id]
}
`

var testAccCheckGroupConfigBasicGetUsersData = testAccCheckGroupCreateUsers + `
resource ` + constant.GroupResource + ` ` + constant.GroupTestResource + ` {
  name = "` + constant.GroupTestResource + `"
  get_users_data = true
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = true
  create_pcc = true
  s3_privilege = true
  create_backup_unit = true
  create_internet_access = true
  create_k8s_cluster = true
  create_flow_log = true
  access_and_manage_monitoring = true
  access_and_manage_certificates = true
  access_and_manage_dns = true
  manage_registry = true
  manage_dataplatform = true
  access_and_manage_logging = true
  access_and_manage_cdn = true
  access_and_manage_vpn = true
  access_and_manage_api_gateway = true
  access_and_manage_kaas = true
  access_and_manage_network_file_storage = true
  access_and_manage_ai_model_hub = true
  create_network_security_groups = true
  access_and_manage_iam_resources = true
  manage_dbaas = true


  user_ids = [` + constant.UserResource + `.` + constant.UserTestResource + `.id, ` + constant.UserResource + `.` + constant.UserTestResource + `2.id]
}
`

var testAccDataSourceGroupMatchId = testAccCheckGroupConfigBasic + `
data ` + constant.GroupResource + ` ` + constant.GroupDataSourceById + ` {
  id			= ` + constant.GroupResource + `.` + constant.GroupTestResource + `.id
}
`

var testAccDataSourceGroupMatchIdGetUsersData = testAccCheckGroupConfigBasic + `
data ` + constant.GroupResource + ` ` + constant.GroupDataSourceById + ` {
  id			= ` + constant.GroupResource + `.` + constant.GroupTestResource + `.id
  get_users_data = true
}
`

var testAccDataSourceGroupMatchIdNoUsersData = testAccCheckGroupConfigBasic + `
data ` + constant.GroupResource + ` ` + constant.GroupDataSourceById + ` {
  id			= ` + constant.GroupResource + `.` + constant.GroupTestResource + `.id
  get_users_data = false
}
`

var testAccDataSourceGroupMatchName = testAccCheckGroupConfigBasic + `
resource ` + constant.GroupResource + ` ` + constant.GroupTestResource + `similar {
  name = "similar` + constant.GroupTestResource + `"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = true
  create_pcc = true
  s3_privilege = true
  create_backup_unit = true
  create_internet_access = true
  create_k8s_cluster = true
}
data ` + constant.GroupResource + ` ` + constant.GroupDataSourceByName + ` {
  name			= "` + constant.GroupTestResource + `"
}
`

var testAccDataSourceGroupMatchNameGetUsersData = testAccCheckGroupConfigBasic + `
resource ` + constant.GroupResource + ` ` + constant.GroupTestResource + `similar {
  name = "similar` + constant.GroupTestResource + `"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = true
  create_pcc = true
  s3_privilege = true
  create_backup_unit = true
  create_internet_access = true
  create_k8s_cluster = true
}
data ` + constant.GroupResource + ` ` + constant.GroupDataSourceByName + ` {
  name			= "` + constant.GroupTestResource + `"
  get_users_data = true
}
`

var testAccDataSourceGroupMatchNameNoUsersData = testAccCheckGroupConfigBasic + `
resource ` + constant.GroupResource + ` ` + constant.GroupTestResource + `similar {
  name = "similar` + constant.GroupTestResource + `"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = true
  create_pcc = true
  s3_privilege = true
  create_backup_unit = true
  create_internet_access = true
  create_k8s_cluster = true
}
data ` + constant.GroupResource + ` ` + constant.GroupDataSourceByName + ` {
  name			= "` + constant.GroupTestResource + `"
  get_users_data = false
}
`

var testAccDataSourceGroupMultipleResultsError = testAccCheckGroupConfigBasic + `
resource ` + constant.GroupResource + ` ` + constant.GroupTestResource + `_multiple_results {
  name = "` + constant.GroupTestResource + `"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = true
  create_pcc = true
  s3_privilege = true
  create_backup_unit = true
  create_internet_access = true
  create_k8s_cluster = true
}

data ` + constant.GroupResource + ` ` + constant.GroupDataSourceByName + ` {
  name			= "` + constant.GroupTestResource + `"
}
`

var testAccDataSourceGroupWrongNameError = testAccCheckGroupConfigBasic + `
data ` + constant.GroupResource + ` ` + constant.GroupDataSourceByName + ` {
  name			= "wrong_name"
}
`

var testAccCheckGroupConfigUpdate = testAccCheckGroupCreateUsers + `
resource ` + constant.UserResource + ` ` + constant.UserTestResource + `3 {
  first_name = "user"
  last_name = "test"
  email = "` + utils.GenerateEmail() + `"
  password = ` + constant.RandomPassword + `.user3_password.result
  administrator = false
  force_sec_auth= false
  active = false
}

resource ` + constant.GroupResource + ` "test_user_id" {
  name = "` + constant.GroupTestResource + `"
  create_datacenter = false
  create_snapshot = false
  reserve_ip = false
  access_activity_log = false
  create_pcc = false
  s3_privilege = false
  create_backup_unit = false
  create_internet_access = false
  create_k8s_cluster = false
  create_flow_log = false
  access_and_manage_monitoring = false
  access_and_manage_certificates = false
  access_and_manage_dns = false
  manage_registry = false
  manage_dataplatform = false
  access_and_manage_logging = false
  access_and_manage_cdn = false
  access_and_manage_vpn = false
  access_and_manage_api_gateway = false
  access_and_manage_kaas = false
  access_and_manage_network_file_storage = false
  access_and_manage_ai_model_hub = false
  create_network_security_groups = false
  access_and_manage_iam_resources = false
  manage_dbaas = false
  user_id = ` + constant.UserResource + `.` + constant.UserTestResource + `3.id
}

resource ` + constant.GroupResource + ` ` + constant.GroupTestResource + ` {
  name = "` + constant.UpdatedResources + `"
  create_datacenter = false
  create_snapshot = false
  reserve_ip = false
  access_activity_log = false
  create_pcc = false
  s3_privilege = false
  create_backup_unit = false
  create_internet_access = false
  create_k8s_cluster = false
  user_ids = [` + constant.UserResource + `.` + constant.UserTestResource + `.id, ` + constant.UserResource + `.` + constant.UserTestResource + `3.id]
}

resource ` + constant.RandomPassword + ` "user3_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
`

var testAccCheckGroupUpdateMigrateToUserIds = testAccCheckGroupCreateUsers + `
resource ` + constant.UserResource + ` ` + constant.UserTestResource + `3 {
  first_name = "user"
  last_name = "test"
  email = "` + utils.GenerateEmail() + `"
  password = ` + constant.RandomPassword + `.user3_password.result
  administrator = false
  force_sec_auth= false
  active = false
}

resource ` + constant.GroupResource + ` "test_user_id" {
  name = "` + constant.GroupTestResource + `"
  create_datacenter = false
  create_snapshot = false
  reserve_ip = false
  access_activity_log = false
  create_pcc = false
  s3_privilege = false
  create_backup_unit = false
  create_internet_access = false
  create_k8s_cluster = false
  user_ids = [` + constant.UserResource + `.` + constant.UserTestResource + `3.id]
}

resource ` + constant.GroupResource + ` ` + constant.GroupTestResource + ` {
  name = "` + constant.UpdatedResources + `"
  create_datacenter = false
  create_snapshot = false
  reserve_ip = false
  access_activity_log = false
  create_pcc = false
  s3_privilege = false
  create_backup_unit = false
  create_internet_access = false
  create_k8s_cluster = false
  user_ids = [` + constant.UserResource + `.` + constant.UserTestResource + `.id, ` + constant.UserResource + `.` + constant.UserTestResource + `3.id]
}
resource ` + constant.RandomPassword + ` "user3_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
`

var testAccCheckGroupBothUserArgumentsError = testAccCheckGroupCreateUsers + `
resource ` + constant.UserResource + ` ` + constant.UserTestResource + `3 {
  first_name = "user"
  last_name = "test"
  email = "` + utils.GenerateEmail() + `"
  password = ` + constant.RandomPassword + `.user3_password.result
  administrator = false
  force_sec_auth= false
  active = false
}

resource ` + constant.GroupResource + ` "test_user_id" {
  name = "` + constant.GroupTestResource + `"
  create_datacenter = false
  create_snapshot = false
  reserve_ip = false
  access_activity_log = false
  create_pcc = false
  s3_privilege = false
  create_backup_unit = false
  create_internet_access = false
  create_k8s_cluster = false
  user_ids = [` + constant.UserResource + `.` + constant.UserTestResource + `3.id]
}

resource ` + constant.GroupResource + ` ` + constant.GroupTestResource + ` {
  name = "` + constant.UpdatedResources + `"
  create_datacenter = false
  create_snapshot = false
  reserve_ip = false
  access_activity_log = false
  create_pcc = false
  s3_privilege = false
  create_backup_unit = false
  create_internet_access = false
  create_k8s_cluster = false
  user_ids = [` + constant.UserResource + `.` + constant.UserTestResource + `.id, ` + constant.UserResource + `.` + constant.UserTestResource + `3.id]
  user_id = ` + constant.UserResource + `.` + constant.UserTestResource + `.id
}
resource ` + constant.RandomPassword + ` "user3_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
`
