//go:build all || kafka || user

package kafka_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
)

func TestAccUsersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccUsersDataSourceConfig_nonexistentCluster,
				ExpectError: regexp.MustCompile("404 Not Found"),
			},
			{
				Config:      testAccUsersDataSourceConfig_invalidClusterID,
				ExpectError: regexp.MustCompile("String must be a valid UUID"),
			},
			{
				Config:      testAccUsersDataSourceConfig_invalidLocation,
				ExpectError: regexp.MustCompile("Attribute location value must be one of"),
			},
			{
				Config:      testAccUsersDataSourceConfig_smallTimeoutValue,
				ExpectError: regexp.MustCompile("context deadline exceeded"),
			},
			{
				Config: testAccUsersDataSourceConfig_valid,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.ionoscloud_kafka_users.kafka_users_ds",
						tfjsonpath.New("users"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"id":       knownvalue.NotNull(),
								"username": knownvalue.StringExact("admin"),
							}),
						}),
					),
				},
			},
		},
	})
}

const (
	// TODO - Write the configuration for the Kafka Cluster and reuse it in other configurations.
	testAccUsersDataSourceConfig_valid = `
		data "ionoscloud_kafka_users" "kafka_users_ds" {
		  cluster_id = "92ca35f4-5fb8-438a-9c94-1806e76b63dd"
		  location = "de/fra"
          timeouts = {
			read = "1m"
		  }
		}
	`
	testAccUsersDataSourceConfig_nonexistentCluster = `
		data "ionoscloud_kafka_users" "kafka_users_ds" {
		  cluster_id = "92ca35f4-5fb8-438a-9c94-1806e76b63de"
		  location = "de/fra"
		}
	`

	testAccUsersDataSourceConfig_invalidClusterID = `
		data "ionoscloud_kafka_users" "kafka_users_ds" {
		  cluster_id = "invalid-UUID"
		  location = "de/fra"
		}
	`

	testAccUsersDataSourceConfig_invalidLocation = `
		data "ionoscloud_kafka_users" "kafka_users_ds" {
		  cluster_id = "92ca35f4-5fb8-438a-9c94-1806e76b63dd"
		  location = "invalid"
		}
	`
	testAccUsersDataSourceConfig_smallTimeoutValue = `
		data "ionoscloud_kafka_users" "kafka_users_ds" {
		  cluster_id = "92ca35f4-5fb8-438a-9c94-1806e76b63dd"
		  location = "de/fra"
          timeouts = {
			read = "1ms"
		  }
		}
	`
)
