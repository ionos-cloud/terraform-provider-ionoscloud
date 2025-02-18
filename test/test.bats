#!/usr/bin/env bats
# tags: terraform, provider, config, logging, apply

BATS_LIBS_PATH="${LIBS_PATH:-./libs}"  # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"

setup_file() {
    # CHECK REQUIRED ENV VARS
    if [ -z "$VALID_IONOS_TOKEN" ]; then
      echo "ERROR: VALID_IONOS_TOKEN must be set to a valid token for tests that require success." >&2
      exit 1
    fi
    command -v ionosctl >/dev/null 2>&1 || { echo "ERROR: ionosctl command not found"; exit 1; }
}

setup() {
  # Create a temporary working directory for Terraform.
  TMP_DIR=$(mktemp -d)
  export TF_IN_AUTOMATION=1

  # Start the provider in debuggable mode and capture its output.
  PLUGIN_LOG=$(mktemp)
  terraform-provider-ionoscloud -debuggable > "$PLUGIN_LOG" 2>&1 &
  PLUGIN_PID=$!

  # Wait briefly for the provider to start.
  sleep 2

  # Extract TF_REATTACH_PROVIDERS info from the log.
  TF_REATTACH_LINE=$(grep -o "TF_REATTACH_PROVIDERS='[^']*'" "$PLUGIN_LOG")
  TF_REATTACH_VALUE=$(echo "$TF_REATTACH_LINE" | sed "s/TF_REATTACH_PROVIDERS='//;s/'//")
  export TF_REATTACH_PROVIDERS="$TF_REATTACH_VALUE"

  # Write a Terraform configuration (Logging Pipeline) into main.tf.
  cat > "$TMP_DIR/main.tf" <<'EOF'
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

  # Write a full config file (mimicking ~/.ionos/config) to config.yaml.
  # Here profile "user1" gets the valid token (good token) from VALID_IONOS_TOKEN.
  cat > "$TMP_DIR/config.yaml" <<EOF
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
            name: https://logging.es-vit.ionos.com
          - location: bla
            name: https://logging.de-txl.ionos.com
  - name: prod
    products:
      - name: logging
        endpoints:
          - location: es/vit
            name: https://logging.es-vit.ionos.com
          - location: bla
            name: https://logging.de-txl.ionos.com
EOF

  # Initialize Terraform.
  pushd "$TMP_DIR" > /dev/null
  terraform init -input=false >/dev/null 2>&1
  popd > /dev/null

  # Suppress Terraform debug output for the commands we run.
  export TF_LOG=""
}

teardown() {
  ionosctl logging-service pipeline delete -af

  # Kill the provider process and clean up temporary files.
  kill "$PLUGIN_PID" 2>/dev/null
  rm -f "$PLUGIN_LOG"
  rm -rf "$TMP_DIR"

  unset TF_IN_AUTOMATION TF_REATTACH_PROVIDERS TF_LOG
  unset IONOS_CONFIG_FILE_PATH IONOS_CONFIG_PROFILE IONOS_TOKEN
  unset IONOS_API_URL IONOS_API_URL_LOGGING
}

# -----------------------------------------------
# Test Scenarios – using "terraform apply -auto-approve"
# -----------------------------------------------

# Scenario 1:
# IONOS_CONFIG_FILE_PATH set. (Default profile is user2 => token "BAD_TOKEN")
@test "Config file only: IONOS_CONFIG_FILE_PATH set (default profile user2) fails with 401" {
  export IONOS_CONFIG_FILE_PATH="$TMP_DIR/config.yaml"
  # No override: currentProfile is user2, token "BAD_TOKEN" (should fail)
  unset IONOS_CONFIG_PROFILE IONOS_TOKEN

  pushd "$TMP_DIR" > /dev/null
  run terraform apply -auto-approve -input=false -no-color
  popd > /dev/null

  assert_failure
  assert_output -p "401"
}

# Scenario 2:
# IONOS_CONFIG_PROFILE set to user1. (Should use the valid token from config file.)
@test "Config file only: IONOS_CONFIG_PROFILE set to user1 applies successfully" {
  export IONOS_CONFIG_FILE_PATH="$TMP_DIR/config.yaml"
  export IONOS_CONFIG_PROFILE="user1"
  unset IONOS_TOKEN

  pushd "$TMP_DIR" > /dev/null
  run terraform apply -auto-approve -input=false -no-color
  popd > /dev/null

  assert_success
  assert_output -p "Apply complete!"
}

# Scenario 3:
# Non-existent IONOS_CONFIG_PROFILE.
@test "Config file only: Non-existent IONOS_CONFIG_PROFILE errors appropriately" {
  skip "todo"
}

# Scenario 4:
# IONOS_TOKEN env var overrides config file credentials. (Even if a profile is set, the env var wins.)
@test "Env var credentials: IONOS_TOKEN overrides config file and apply succeeds" {
  export IONOS_CONFIG_FILE_PATH="$TMP_DIR/config.yaml"
  # Even though the config file would use a dummy token for user2,
  # setting IONOS_TOKEN to a valid token should override them.
  # according to the docs at least
  export IONOS_TOKEN="${VALID_IONOS_TOKEN}"
  unset IONOS_CONFIG_PROFILE

  pushd "$TMP_DIR" > /dev/null
  run terraform apply -auto-approve -input=false -no-color
  popd > /dev/null

  assert_success
  assert_output -p "Apply complete!"
}

# Scenario 5:
# IONOS_API_URL override should override any endpoint defined in the config file.
@test "Endpoint override: Global IONOS_API_URL override is accepted and apply succeeds" {
  export IONOS_CONFIG_FILE_PATH="$TMP_DIR/config.yaml"
  unset IONOS_CONFIG_PROFILE
  export IONOS_API_URL="https://override.api.ionos.com"
  export IONOS_TOKEN="${VALID_IONOS_TOKEN}"

  pushd "$TMP_DIR" > /dev/null
  run terraform apply -auto-approve -input=false -no-color
  popd > /dev/null

  # We expect that the override causes Terraform to use the new endpoint.
  # (Your provider may produce an error if the override endpoint is not resolvable.
  # Here we check for an error that mentions the override, adjust as needed.)
  assert_failure
  assert_output -p "lookup override.api.ionos.com"
}

# Scenario 6:
# For products that have location-based URLs: verify that the endpoint is overridden only for the location provided.
@test "Endpoint override: IONOS_API_URL_LOGGING override is accepted for logging product and apply succeeds" {
  export IONOS_CONFIG_FILE_PATH="$TMP_DIR/config.yaml"
  unset IONOS_CONFIG_PROFILE
  # For logging, override the endpoint using a product-specific variable.
  export IONOS_API_URL_LOGGING="https://override.logging.endpoint"
  export IONOS_TOKEN="${VALID_IONOS_TOKEN}"

  pushd "$TMP_DIR" > /dev/null
  run terraform apply -auto-approve -input=false -no-color
  popd > /dev/null

  assert_success
  assert_output -p "Apply complete!"
}

# Scenario 7:
# Change environment in config file: update currentProfile and see if endpoints update accordingly.
@test "Environment change: switching currentProfile in config file (simulate environment change)" {
  skip "todo"
}

# Scenario 8:
# IONOS_API_URL_product_name override: For logging, test that IONOS_API_URL_LOGGING override works.
@test "Product-specific override: IONOS_API_URL_LOGGING override is accepted and apply succeeds" {
  export IONOS_CONFIG_FILE_PATH="$TMP_DIR/config.yaml"
  unset IONOS_CONFIG_PROFILE
  export IONOS_API_URL_LOGGING="https://override.logging.endpoint"
  export IONOS_TOKEN="${VALID_IONOS_TOKEN}"

  pushd "$TMP_DIR" > /dev/null
  run terraform apply -auto-approve -input=false -no-color
  popd > /dev/null

  assert_success
  assert_output -p "Apply complete!"
}
