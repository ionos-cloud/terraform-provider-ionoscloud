//go:build all || dsaas
// +build all dsaas

package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	dsaas "github.com/ionos-cloud/sdk-go-autoscaling"
	"regexp"
	"testing"
)

func TestAccDSaaSClusterBasic(t *testing.T) {
	//var DSaaSCluster dsaas.ClusterResponseData

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDSaaSClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceDSaaSClusterWrongNameError,
				ExpectError: regexp.MustCompile("no DSaaS cluster found with the specified name"),
			},
		},
	})
}

func testAccCheckDSaaSClusterDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).DSaaSClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != DSaaSClusterResource {
			continue
		}

		_, apiResponse, err := client.GetCluster(ctx, rs.Primary.ID)

		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while checking the destruction of DSaaS cluster %s: %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("DSaaS cluster %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckDSaaSClusterExists(n string, cluster *dsaas.ClusterResponseData) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).DSaaSClient

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		foundCluster, _, err := client.GetCluster(ctx, rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("an error occured while fetching DSaaS Cluster %s: %s", rs.Primary.ID, err)
		}
		if *foundCluster.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		cluster = &foundCluster

		return nil
	}
}

const testAccDataSourceDSaaSClusterWrongNameError = `
data ` + DSaaSClusterResource + ` ` + DSaaSClusterTestDataSourceByName + ` {
	name = "wrong_name"
}
`
