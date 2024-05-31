//go:build all || kafka
// +build all kafka

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestKafkaClusterResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testCheckKafkaClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testKafkaClusterBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.KafkaClusterResource+"."+kafkaTestResource, "name", kafkaClusterTestResource),
					resource.TestCheckResourceAttr(constant.KafkaClusterResource+"."+kafkaTestResource, "version", kafkaClusterVersion),
					resource.TestCheckResourceAttr(constant.KafkaClusterResource+"."+kafkaTestResource, "size", kafkaClusterSize),
				),
			},
			{
				Config: testKafkaClusterMatchID,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.KafkaClusterResource+"."+kafkaDSTestResource, "name", constant.KafkaClusterResource+"."+kafkaTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.KafkaClusterResource+"."+kafkaDSTestResource, "version", constant.KafkaClusterResource+"."+kafkaTestResource, "version"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.KafkaClusterResource+"."+kafkaDSTestResource, "size", constant.KafkaClusterResource+"."+kafkaTestResource, "size"),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.KafkaClusterResource+"."+kafkaDSTestResource, "server_address"),
				),
			},
			{
				Config: testKafkaClusterMatchNameExact,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.KafkaClusterResource+"."+kafkaDSTestResource, "name", constant.KafkaClusterResource+"."+kafkaTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.KafkaClusterResource+"."+kafkaDSTestResource, "version", constant.KafkaClusterResource+"."+kafkaTestResource, "version"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.KafkaClusterResource+"."+kafkaDSTestResource, "size", constant.KafkaClusterResource+"."+kafkaTestResource, "size"),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.KafkaClusterResource+"."+kafkaDSTestResource, "server_address"),
				),
			},
			{
				Config: testKafkaClusterMatchNamePartial,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.KafkaClusterResource+"."+kafkaDSTestResource, "name", constant.KafkaClusterResource+"."+kafkaTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.KafkaClusterResource+"."+kafkaDSTestResource, "version", constant.KafkaClusterResource+"."+kafkaTestResource, "version"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.KafkaClusterResource+"."+kafkaDSTestResource, "size", constant.KafkaClusterResource+"."+kafkaTestResource, "size"),
					resource.TestCheckResourceAttrSet(constant.DataSource+"."+constant.KafkaClusterResource+"."+kafkaDSTestResource, "server_address"),
				),
			},
			{
				Config:      testKafkaClusterMatchNamePartialNotFound,
				ExpectError: regexp.MustCompile(`no kafka clusters found with the specified name: willnotwork`),
			},
		},
	})
}

func testCheckKafkaClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).KafkaClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.KafkaClusterResource {
			continue
		}

		_, apiResponse, err := client.GetClusterById(ctx, rs.Primary.ID)
		apiResponse.LogInfo()
		if err != nil {
			if !apiResponse.HttpNotFound() {
				return fmt.Errorf("an error occurred while checking the destruction of resource %s: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("resource %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

const testKafkaClusterBasic = `resource ` + constant.KafkaClusterResource + ` ` + kafkaTestResource + ` {
	  name = "` + kafkaClusterTestResource + `"
	  version = "` + kafkaClusterVersion + `"
	  size = "` + kafkaClusterSize + `"
	}`

const testKafkaClusterMatchID = testKafkaClusterBasic +
	`
data ` + constant.KafkaClusterResource + ` ` + kafkaDSTestResource + ` {
	  id = ` + constant.KafkaClusterResource + `.` + kafkaTestResource + `.id
	}`

const testKafkaClusterMatchNameExact = testKafkaClusterBasic +
	`
data ` + constant.KafkaClusterResource + ` ` + kafkaDSTestResource + ` {
	  name = ` + constant.KafkaClusterResource + `.` + kafkaTestResource + `.name
	}`

const testKafkaClusterMatchNamePartial = testKafkaClusterBasic + `
data ` + constant.KafkaClusterResource + ` ` + kafkaDSTestResource + ` {
	  name = "kafka_cluster"
	  partial_match = true
	}`

const testKafkaClusterMatchNamePartialNotFound = testKafkaClusterBasic + `
data ` + constant.KafkaClusterResource + ` ` + kafkaDSTestResource + ` {
	  name = "willnotwork"
	  partial_match = true
	}`
