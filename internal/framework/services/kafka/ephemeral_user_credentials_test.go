//go:build all || kafka || user

package kafka_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/echoprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
)

func TestAccUserCredentialsEphemeral(t *testing.T) {
	testProviderFactories := make(map[string]func() (tfprotov6.ProviderServer, error))
	for k, v := range acctest.TestAccProtoV6ProviderFactories {
		testProviderFactories[k] = v
	}
	testProviderFactories["echo"] = echoprovider.NewProviderServer()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccUserCredentialsEphemeralConfig_nonexistentCluster,
				ExpectError: regexp.MustCompile("not found"),
			},
			{
				Config:      testAccUserCredentialsEphemeralConfig_nonexistentUserID,
				ExpectError: regexp.MustCompile("not found"),
			},
			{
				Config:      testAccUserCredentialsEphemeralConfig_nonexistentUsername,
				ExpectError: regexp.MustCompile("no Kafka user was found"),
			},
			{
				Config:      testAccUserCredentialsEphemeralConfig_invalidClusterID,
				ExpectError: regexp.MustCompile("String must be a valid UUID"),
			},
			{
				Config:      testAccUserCredentialsEphemeralConfig_invalidUserID,
				ExpectError: regexp.MustCompile("String must be a valid UUID"),
			},
			{
				Config:      testAccUserCredentialsEphemeralConfig_smallTimeoutValue,
				ExpectError: regexp.MustCompile("context deadline exceeded"),
			},
			{
				Config:      testAccUserCredentialsEphemeralConfig_invalidLocation,
				ExpectError: regexp.MustCompile("Attribute location value must be one of"),
			},
			{
				Config:      testAccUserCredentialsEphemeralConfig_bothUserIDUsername,
				ExpectError: regexp.MustCompile("Exactly one of these attributes must be configured"),
			},
			{
				Config: testAccUserCredentialsEphemeralConfig_validByID,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.verification",
						tfjsonpath.New("data").AtMapKey("certificate_authority"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"echo.verification",
						tfjsonpath.New("data").AtMapKey("private_key"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"echo.verification",
						tfjsonpath.New("data").AtMapKey("certificate"),
						knownvalue.NotNull(),
					),
				},
			},
			{
				Config: testAccUserCredentialsEphemeralConfig_validByName,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.verification",
						tfjsonpath.New("data").AtMapKey("certificate_authority"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"echo.verification",
						tfjsonpath.New("data").AtMapKey("private_key"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"echo.verification",
						tfjsonpath.New("data").AtMapKey("certificate"),
						knownvalue.NotNull(),
					),
				},
			},
			{
				Config: testAccUserCredentialsEphemeralConfig_defaultLocation,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.verification",
						tfjsonpath.New("data").AtMapKey("certificate_authority"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"echo.verification",
						tfjsonpath.New("data").AtMapKey("private_key"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"echo.verification",
						tfjsonpath.New("data").AtMapKey("certificate"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

const (
	// INVALID CONFIGURATIONS
	testAccUserCredentialsEphemeralConfig_nonexistentCluster = kafkaUsersDataSource + `
		ephemeral "ionoscloud_kafka_user_credentials" "user_credentials_ephemeral" {
		  cluster_id = "92ca35f4-5fb8-438a-9c94-1806e76b63de"
		  id = data.ionoscloud_kafka_users.kafka_users_ds.users[0].id
		  location = "de/fra"
		  timeouts = {
			open = "2m"
		  }
		}
	`

	testAccUserCredentialsEphemeralConfig_nonexistentUserID = kafkaUsersDataSource + `
		ephemeral "ionoscloud_kafka_user_credentials" "user_credentials_ephemeral" {
		  cluster_id = data.ionoscloud_kafka_users.kafka_users_ds.cluster_id
		  id = "92ca35f4-5fb8-438a-9c94-1806e76b63dd"
		  location = "de/fra"
		  timeouts = {
			open = "2m"
		  }
		}
	`

	testAccUserCredentialsEphemeralConfig_nonexistentUsername = kafkaUsersDataSource + `
		ephemeral "ionoscloud_kafka_user_credentials" "user_credentials_ephemeral" {
		  cluster_id = data.ionoscloud_kafka_users.kafka_users_ds.cluster_id
		  username = "nonexistent"
		  location = "de/fra"
		  timeouts = {
			open = "2m"
		  }
		}
	`

	testAccUserCredentialsEphemeralConfig_invalidClusterID = kafkaUsersDataSource + `
		ephemeral "ionoscloud_kafka_user_credentials" "user_credentials_ephemeral" {
		  cluster_id = "invalid UUID"
		  id = data.ionoscloud_kafka_users.kafka_users_ds.users[0].id
		  location = "de/fra"
		  timeouts = {
			open = "2m"
		  }
		}
	`

	testAccUserCredentialsEphemeralConfig_invalidUserID = kafkaUsersDataSource + `
		ephemeral "ionoscloud_kafka_user_credentials" "user_credentials_ephemeral" {
		  cluster_id = data.ionoscloud_kafka_users.kafka_users_ds.cluster_id
		  id = "invalid UUID"
		  location = "de/fra"
		  timeouts = {
			open = "2m"
		  }
		}
	`

	testAccUserCredentialsEphemeralConfig_smallTimeoutValue = kafkaUsersDataSource + `
		ephemeral "ionoscloud_kafka_user_credentials" "user_credentials_ephemeral" {
		  cluster_id = data.ionoscloud_kafka_users.kafka_users_ds.cluster_id
		  id = data.ionoscloud_kafka_users.kafka_users_ds.users[0].id
		  location = "de/fra"
		  timeouts = {
			open = "1ms"
		  }
		}
	`

	testAccUserCredentialsEphemeralConfig_invalidLocation = kafkaUsersDataSource + `
		ephemeral "ionoscloud_kafka_user_credentials" "user_credentials_ephemeral" {
		  cluster_id = data.ionoscloud_kafka_users.kafka_users_ds.cluster_id
		  id = data.ionoscloud_kafka_users.kafka_users_ds.users[0].id
		  location = "nonexistent"
		  timeouts = {
			open = "2m"
		  }
		}
	`

	testAccUserCredentialsEphemeralConfig_bothUserIDUsername = kafkaUsersDataSource + `
		ephemeral "ionoscloud_kafka_user_credentials" "user_credentials_ephemeral" {
		  cluster_id = data.ionoscloud_kafka_users.kafka_users_ds.cluster_id
		  id = data.ionoscloud_kafka_users.kafka_users_ds.users[0].id
          username = data.ionoscloud_kafka_users.kafka_users_ds.users[0].username
		  location = "de/fra"
		  timeouts = {
			open = "2m"
		  }
		}
	`

	// VALID CONFIGURATIONS
	testAccUserCredentialsEphemeralConfig_validByID = kafkaUsersDataSource + `
		ephemeral "ionoscloud_kafka_user_credentials" "user_credentials_ephemeral" {
		  cluster_id = data.ionoscloud_kafka_users.kafka_users_ds.cluster_id
		  id = data.ionoscloud_kafka_users.kafka_users_ds.users[0].id
		  location = "de/fra"
		  timeouts = {
			open = "2m"
		  }
		}

		provider "echo" {
		  data = {
			certificate_authority = ephemeral.ionoscloud_kafka_user_credentials.user_credentials_ephemeral.certificate_authority
			private_key           = ephemeral.ionoscloud_kafka_user_credentials.user_credentials_ephemeral.private_key
			certificate           = ephemeral.ionoscloud_kafka_user_credentials.user_credentials_ephemeral.certificate
		  }
		}

		resource "echo" "verification" {}
	`

	testAccUserCredentialsEphemeralConfig_validByName = kafkaUsersDataSource + `
		ephemeral "ionoscloud_kafka_user_credentials" "user_credentials_ephemeral" {
		  cluster_id = data.ionoscloud_kafka_users.kafka_users_ds.cluster_id
		  username = data.ionoscloud_kafka_users.kafka_users_ds.users[0].username
		  location = "de/fra"
		  timeouts = {
			open = "2m"
		  }
		}

		provider "echo" {
		  data = {
			certificate_authority = ephemeral.ionoscloud_kafka_user_credentials.user_credentials_ephemeral.certificate_authority
			private_key           = ephemeral.ionoscloud_kafka_user_credentials.user_credentials_ephemeral.private_key
			certificate           = ephemeral.ionoscloud_kafka_user_credentials.user_credentials_ephemeral.certificate
		  }
		}

		resource "echo" "verification" {}
	`

	// Omit 'location' from configuration and make sure the default location is used.
	testAccUserCredentialsEphemeralConfig_defaultLocation = kafkaUsersDataSource + `
		ephemeral "ionoscloud_kafka_user_credentials" "user_credentials_ephemeral" {
		  cluster_id = data.ionoscloud_kafka_users.kafka_users_ds.cluster_id
		  username = data.ionoscloud_kafka_users.kafka_users_ds.users[0].username
		  timeouts = {
			open = "2m"
		  }
		}

		provider "echo" {
		  data = {
			certificate_authority = ephemeral.ionoscloud_kafka_user_credentials.user_credentials_ephemeral.certificate_authority
			private_key           = ephemeral.ionoscloud_kafka_user_credentials.user_credentials_ephemeral.private_key
			certificate           = ephemeral.ionoscloud_kafka_user_credentials.user_credentials_ephemeral.certificate
		  }
		}

		resource "echo" "verification" {}
	`
)
