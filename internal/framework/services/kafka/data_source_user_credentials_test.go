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

func TestAccUserCredentialsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccUserCredentialsDataSourceConfig_nonexistentCluster,
				ExpectError: regexp.MustCompile("not found"),
			},
			{
				Config:      testAccUserCredentialsDataSourceConfig_nonexistentUserID,
				ExpectError: regexp.MustCompile("not found"),
			},
			{
				Config:      testAccUserCredentialsDataSourceConfig_nonexistentUsername,
				ExpectError: regexp.MustCompile("no Kafka user was found"),
			},
			{
				Config:      testAccUserCredentialsDataSourceConfig_invalidClusterID,
				ExpectError: regexp.MustCompile("String must be a valid UUID"),
			},
			{
				Config:      testAccUserCredentialsDataSourceConfig_invalidUserID,
				ExpectError: regexp.MustCompile("String must be a valid UUID"),
			},
			{
				Config:      testAccUserCredentialsDataSourceConfig_smallTimeoutValue,
				ExpectError: regexp.MustCompile("context deadline exceeded"),
			},
			{
				Config:      testAccUserCredentialsDataSourceConfig_invalidLocation,
				ExpectError: regexp.MustCompile("Attribute location value must be one of"),
			},
			{
				Config:      testAccUserCredentialsDataSourceConfig_bothUserIDUsername,
				ExpectError: regexp.MustCompile("Exactly one of these attributes must be configured"),
			},
			{
				Config: testAccUserCredentialsDataSourceConfig_validByID,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.ionoscloud_kafka_user_credentials.user_credentials_ds",
						tfjsonpath.New("certificate_authority"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.ionoscloud_kafka_user_credentials.user_credentials_ds",
						tfjsonpath.New("private_key"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.ionoscloud_kafka_user_credentials.user_credentials_ds",
						tfjsonpath.New("certificate"),
						knownvalue.NotNull(),
					),
				},
			},
			{
				Config: testAccUserCredentialsDataSourceConfig_validByName,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.ionoscloud_kafka_user_credentials.user_credentials_ds",
						tfjsonpath.New("certificate_authority"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.ionoscloud_kafka_user_credentials.user_credentials_ds",
						tfjsonpath.New("private_key"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.ionoscloud_kafka_user_credentials.user_credentials_ds",
						tfjsonpath.New("certificate"),
						knownvalue.NotNull(),
					),
				},
			},
			{
				Config: testAccUserCredentialsDataSourceConfig_defaultLocation,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.ionoscloud_kafka_user_credentials.user_credentials_ds",
						tfjsonpath.New("certificate_authority"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.ionoscloud_kafka_user_credentials.user_credentials_ds",
						tfjsonpath.New("private_key"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.ionoscloud_kafka_user_credentials.user_credentials_ds",
						tfjsonpath.New("certificate"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

const (
	// TODO -- Change this with an actual kafka resource that will be created.
	kafkaUsersDataSource = `
		data "ionoscloud_kafka_users" "kafka_users_ds" {
		  cluster_id = "92ca35f4-5fb8-438a-9c94-1806e76b63dd"
		  location = "de/fra"
		}
	`

	// INVALID CONFIGURATIONS
	testAccUserCredentialsDataSourceConfig_nonexistentCluster = kafkaUsersDataSource + `
		data "ionoscloud_kafka_user_credentials" "user_credentials_ds" {
		  cluster_id = "92ca35f4-5fb8-438a-9c94-1806e76b63de"
		  id = data.ionoscloud_kafka_users.kafka_users_ds.users[0].id
		  location = "de/fra"
		  timeouts = {
			read = "2m"
		  }
		}
	`

	testAccUserCredentialsDataSourceConfig_invalidClusterID = kafkaUsersDataSource + `
		data "ionoscloud_kafka_user_credentials" "user_credentials_ds" {
		  cluster_id = "invalid UUID"
		  id = data.ionoscloud_kafka_users.kafka_users_ds.users[0].id
		  location = "de/fra"
		  timeouts = {
			read = "2m"
		  }
		}
	`

	testAccUserCredentialsDataSourceConfig_invalidUserID = kafkaUsersDataSource + `
		data "ionoscloud_kafka_user_credentials" "user_credentials_ds" {
		  cluster_id = data.ionoscloud_kafka_users.kafka_users_ds.cluster_id
		  id = "invalid UUID"
		  location = "de/fra"
		  timeouts = {
			read = "2m"
		  }
		}
	`

	testAccUserCredentialsDataSourceConfig_nonexistentUserID = kafkaUsersDataSource + `
		data "ionoscloud_kafka_user_credentials" "user_credentials_ds" {
		  cluster_id = data.ionoscloud_kafka_users.kafka_users_ds.cluster_id
		  id = "92ca35f4-5fb8-438a-9c94-1806e76b63dd"
		  location = "de/fra"
		  timeouts = {
			read = "2m"
		  }
		}
	`

	testAccUserCredentialsDataSourceConfig_nonexistentUsername = kafkaUsersDataSource + `
		data "ionoscloud_kafka_user_credentials" "user_credentials_ds" {
		  cluster_id = data.ionoscloud_kafka_users.kafka_users_ds.cluster_id
		  username = "nonexistent"
		  location = "de/fra"
		  timeouts = {
			read = "2m"
		  }
		}
	`

	testAccUserCredentialsDataSourceConfig_bothUserIDUsername = kafkaUsersDataSource + `
		data "ionoscloud_kafka_user_credentials" "user_credentials_ds" {
		  cluster_id = data.ionoscloud_kafka_users.kafka_users_ds.cluster_id
		  id = data.ionoscloud_kafka_users.kafka_users_ds.users[0].id
          username = data.ionoscloud_kafka_users.kafka_users_ds.users[0].username
		  location = "de/fra"
		  timeouts = {
			read = "2m"
		  }
		}
	`

	testAccUserCredentialsDataSourceConfig_smallTimeoutValue = kafkaUsersDataSource + `
		data "ionoscloud_kafka_user_credentials" "user_credentials_ds" {
		  cluster_id = data.ionoscloud_kafka_users.kafka_users_ds.cluster_id
		  id = data.ionoscloud_kafka_users.kafka_users_ds.users[0].id
		  location = "de/fra"
		  timeouts = {
			read = "1ms"
		  }
		}
	`

	testAccUserCredentialsDataSourceConfig_invalidLocation = kafkaUsersDataSource + `
		data "ionoscloud_kafka_user_credentials" "user_credentials_ds" {
		  cluster_id = data.ionoscloud_kafka_users.kafka_users_ds.cluster_id
		  id = data.ionoscloud_kafka_users.kafka_users_ds.users[0].id
		  location = "nonexistent"
		  timeouts = {
			read = "2m"
		  }
		}
	`

	// VALID CONFIGURATIONS
	testAccUserCredentialsDataSourceConfig_validByID = kafkaUsersDataSource + `
		data "ionoscloud_kafka_user_credentials" "user_credentials_ds" {
		  cluster_id = data.ionoscloud_kafka_users.kafka_users_ds.cluster_id
		  id = data.ionoscloud_kafka_users.kafka_users_ds.users[0].id
		  location = "de/fra"
		  timeouts = {
			read = "2m"
		  }
		}
	`

	testAccUserCredentialsDataSourceConfig_validByName = kafkaUsersDataSource + `
		data "ionoscloud_kafka_user_credentials" "user_credentials_ds" {
		  cluster_id = data.ionoscloud_kafka_users.kafka_users_ds.cluster_id
		  username = data.ionoscloud_kafka_users.kafka_users_ds.users[0].username
		  location = "de/fra"
		  timeouts = {
			read = "2m"
		  }
		}
	`

	// Omit 'location' from configuration and make sure the default location is used.
	testAccUserCredentialsDataSourceConfig_defaultLocation = kafkaUsersDataSource + `
		data "ionoscloud_kafka_user_credentials" "user_credentials_ds" {
		  cluster_id = data.ionoscloud_kafka_users.kafka_users_ds.cluster_id
		  username = data.ionoscloud_kafka_users.kafka_users_ds.users[0].username
		  timeouts = {
			read = "2m"
		  }
		}
	`
)
