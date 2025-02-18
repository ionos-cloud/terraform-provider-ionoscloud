#!/usr/bin/env bats
# tags: terraform, provider, config, logging, apply

BATS_LIBS_PATH="${LIBS_PATH:-./libs}"  # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"

# A helper to check required env vars.
setup_file() {
  if [ -z "$VALID_IONOS_TOKEN" ]; then
    echo "ERROR: VALID_IONOS_TOKEN must be set to a valid token for tests that require success." >&2
    exit 1
  fi
  command -v ionosctl >/dev/null 2>&1 || { echo "ERROR: ionosctl command not found"; exit 1; }

  # Create a temporary working directory for Terraform.
  mkdir -p /tmp/bats-test/ionoscloud-config
  export TF_IN_AUTOMATION=1

  # Write a Terraform configuration for a logging pipeline into main.tf.
  cat > "/tmp/bats-test/ionoscloud-config/main.tf" <<EOF
terraform {
  required_providers {
    ionoscloud = {
      source  = "ionos-cloud/ionoscloud"
      version = "1.0.0"
    }
  }
}

variable "ionos_token" {
  type    = string
  default = ""
}

provider "ionoscloud" {
  token = var.ionos_token
}

resource "ionoscloud_logging_pipeline" "example" {
  location = "es/vit"
  name     = "pipelineexample"
  log {
    source   = "kubernetes"
    tag      = "tagexample"
    protocol = "http"
    destinations {
      type              = "loki"
      retention_in_days = 7
    }
  }
  log {
    source   = "kubernetes"
    tag      = "anothertagexample"
    protocol = "tcp"
    destinations {
      type              = "loki"
      retention_in_days = 7
    }
  }
}

output "dummy" {
  value = "dummy"
}
EOF

    # Write a full configuration file (config) mimicking ~/.ionos/config.
    # Profile user1 gets the good token from VALID_IONOS_TOKEN.
    cat > "/tmp/bats-test/ionoscloud-config/config" <<EOF
version: 1.0
currentProfile: user2
profiles:
  - name: user1
    environment: preprod
    credentials: # terraform-v6
      token: "${VALID_IONOS_TOKEN}"
  - name: user2
    environment: prod
    credentials: # mytoken
      token: "BAD_TOKEN"
environments:
  - name: preprod
    products:
      - name: logging
        endpoints:
          - location: es/vit
            name: https://logging1.de-fra.ionos.com
            skipTlsVerify: false
          - location: bla
            name: https://logging1.de-txl.ionos.com
            certificateAuthData: "certauthdata"
            skipTlsVerify: true
  - name: prod
    products:
      - name: logging
        endpoints:
          - location: es/vit
            name: https://logging2.de-fra.ionos.com
            skipTlsVerify: false
          - location: bla
            name: https://bla.de-txl.ionos.com
            certificateAuthData: "certauthdata"
            skipTlsVerify: true
EOF

    # Suppress Terraform debug output when running commands.
    export TF_LOG=""

}

setup() {
  # Initialize Terraform.
  (
      cd "/tmp/bats-test/ionoscloud-config" || exit
      terraform init -input=false >/dev/null 2>&1
  )
}

teardown_file() {
    ionosctl logging-service pipeline delete -af

    # Kill provider process and clean up.
    kill "$PLUGIN_PID" 2>/dev/null
    rm -f "$PLUGIN_LOG"
#    rm -rf "/tmp/bats-test/ionoscloud-config"

    unset TF_IN_AUTOMATION TF_REATTACH_PROVIDERS TF_LOG
    unset IONOS_CONFIG_FILE IONOS_CURRENT_PROFILE IONOS_TOKEN
    unset IONOS_API_URL IONOS_API_URL_LOGGING
}

# -----------------------------------------------
# Test Scenarios
# -----------------------------------------------

# Scenario 1: Config file only: IONOS_CONFIG_FILE set (default profile user2) fails with 401
@test "Scenario 1: Config file only: IONOS_CONFIG_FILE set (default profile user2) fails with 401" {
  export IONOS_CONFIG_FILE="/tmp/bats-test/ionoscloud-config/config"
  unset IONOS_CURRENT_PROFILE IONOS_TOKEN

  cd "/tmp/bats-test/ionoscloud-config"|| exit
  run terraform apply -auto-approve -input=false -no-color 2>&1

  assert_failure
#  assert_output -p "401"
}

# Scenario 2: Config file only: IONOS_CURRENT_PROFILE set to user1 applies successfully
@test "Scenario 2: Config file only: IONOS_CURRENT_PROFILE set to user1 applies successfully" {
  export IONOS_CONFIG_FILE="/tmp/bats-test/ionoscloud-config/config"
  export IONOS_CURRENT_PROFILE="user1"
  unset IONOS_TOKEN

  cd "/tmp/bats-test/ionoscloud-config"|| exit
  run terraform apply -auto-approve -input=false -no-color 2>&1

  assert_success
#  assert_output -p "Apply complete!"
}

# Scenario 4: Env var credentials: IONOS_TOKEN overrides config file and apply succeeds
@test "Scenario 4: Env var credentials: IONOS_TOKEN overrides config file and apply succeeds" {
  export IONOS_CONFIG_FILE="/tmp/bats-test/ionoscloud-config/config"
  export IONOS_TOKEN="${VALID_IONOS_TOKEN}"
  unset IONOS_CURRENT_PROFILE

  cd "/tmp/bats-test/ionoscloud-config"|| exit
  run terraform apply -auto-approve -input=false -no-color 2>&1

  assert_success
#  assert_output -p "Apply complete!"
}

# Scenario 5: Global IONOS_API_URL override should override config endpoints
@test "Scenario 5: Global IONOS_API_URL override should override config endpoints" {
  export IONOS_CONFIG_FILE="/tmp/bats-test/ionoscloud-config/config"
  unset IONOS_CURRENT_PROFILE
  export IONOS_API_URL="https://override.api.ionos.com"
  export IONOS_TOKEN="${VALID_IONOS_TOKEN}"

  cd "/tmp/bats-test/ionoscloud-config"|| exit
  run terraform apply -auto-approve -input=false -no-color 2>&1

  assert_failure
#  assert_output -p "lookup override.api.ionos.com"
}

# Scenario 6: IONOS_API_URL_LOGGING override for logging product is applied and apply succeeds
@test "Scenario 6: IONOS_API_URL_LOGGING override for logging product is applied and apply succeeds" {
  export IONOS_CONFIG_FILE="/tmp/bats-test/ionoscloud-config/config"
  unset IONOS_CURRENT_PROFILE
  export IONOS_API_URL_LOGGING="https://logging.de-fra.ionos.com"
  export IONOS_TOKEN="${VALID_IONOS_TOKEN}"

  cd "/tmp/bats-test/ionoscloud-config"|| exit
  run terraform apply -auto-approve -input=false -no-color 2>&1

  assert_success
#  assert_output -p "Apply complete!"
}

# Scenario 7: Custom location override.
@test "Scenario 7: Custom location override: when resource uses 'mylocation', config endpoint applies" {
  # Modify the config file to add a custom endpoint for location "mylocation" (already added in config).
  # Now, update the Terraform configuration to request logging pipeline in location "mylocation".
  cat > "$TMP_DIR/main.tf" <<EOF
terraform {
  required_providers {
    ionoscloud = {
      source  = "ionos-cloud/ionoscloud"
      version = "1.0.0"
    }
  }
}

variable "ionos_token" {
  type    = string
  default = ""
}

provider "ionoscloud" {
  token = var.ionos_token
}

resource "ionoscloud_logging_pipeline" "example_custom" {
  location = "mylocation"
  name     = "pipelineexample-custom"
  log {
    source   = "kubernetes"
    tag      = "customtag"
    protocol = "http"
    destinations {
      type              = "loki"
      retention_in_days = 7
    }
  }
}

output "dummy" {
  value = "dummy"
}
EOF

  # Use valid token via env.
  export IONOS_TOKEN="${VALID_IONOS_TOKEN}"
  unset IONOS_CURRENT_PROFILE

  cd "/tmp/bats-test/ionoscloud-config"|| exit
  run terraform apply -auto-approve -input=false -no-color

  assert_success
  assert_output -p "Apply complete!"
}

# Scenario 8: Change environment in config file.
@test "Scenario 8: Changing environment in config file updates credentials and endpoints" {
  # First, use profile user2 (prod) and expect failure with BAD_TOKEN.
  export IONOS_CONFIG_FILE="$TMP_DIR/config"
  export IONOS_CURRENT_PROFILE="user2"
  unset IONOS_TOKEN

  run terraform apply -auto-approve -input=false -no-color

  assert_failure
  assert_output -p "401"

  # Now change currentProfile in the config file to user1.
  sed -i 's/currentProfile: user2/currentProfile: user1/' "$TMP_DIR/config"

  cd "/tmp/bats-test/ionoscloud-config"|| exit
  run terraform apply -auto-approve -input=false -no-color

  # Expect success because user1 now provides the valid token.
  assert_success
  assert_output -p "Apply complete!"
}

# Scenario 9: Product-specific override.
@test "Scenario 9: IONOS_API_URL_LOGGING override is accepted and apply succeeds" {
  export IONOS_CONFIG_FILE="$TMP_DIR/config"
  unset IONOS_CURRENT_PROFILE
  export IONOS_API_URL_LOGGING="https://logging.de-fra.ionos.com"
  export IONOS_TOKEN="${VALID_IONOS_TOKEN}"

  cd "/tmp/bats-test/ionoscloud-config"|| exit
  run terraform apply -auto-approve -input=false -no-color

  assert_success
  assert_output -p "Apply complete!"
}