package kafka

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	kafkaSDK "github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"

	kafkaService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/kafka"
)

type userCredentialsModel struct {
	ClusterID            types.String `tfsdk:"cluster_id"`
	ID                   types.String `tfsdk:"id"`
	Username             types.String `tfsdk:"username"`
	Location             types.String `tfsdk:"location"`
	CertificateAuthority types.String `tfsdk:"certificate_authority"`
	PrivateKey           types.String `tfsdk:"private_key"`
	Certificate          types.String `tfsdk:"certificate"`
}

// hasMissingData verifies if the API response contains nil values.
func hasMissingData(userCredentials kafkaSDK.UserReadAccess) bool {
	return userCredentials.Metadata.CertificateAuthority == nil || userCredentials.Metadata.PrivateKey == nil || userCredentials.Metadata.Certificate == nil
}

// getUserCredentials gets the user credentials by making the proper API request depending on the provided parameters (GetByID or GetByName).
func getUserCredentials(ctx context.Context, client kafkaService.Client, data userCredentialsModel) (kafkaSDK.UserReadAccess, diag.Diagnostics) {
	var userCredentials kafkaSDK.UserReadAccess
	var err error
	var diags diag.Diagnostics

	clusterID := data.ClusterID.ValueString()
	userID := data.ID.ValueString()
	username := data.Username.ValueString()
	location := data.Location.ValueString()

	// The config validator ensures that 'id' and 'username' are mutually exclusive.
	if userID != "" {
		userCredentials, _, err = client.GetUserCredentialsByID(ctx, clusterID, userID, location)
		if err != nil {
			diags.AddError("API Error Reading Kafka User Credentials", fmt.Sprintf("Failed to retrieve user credentials for user with ID: %s, cluster ID: %s, error: %s", userID, clusterID, err))
			return userCredentials, diags
		}
	} else if username != "" {
		userCredentials, _, err = client.GetUserCredentialsByName(ctx, clusterID, username, location)
		if err != nil {
			diags.AddError("API Error Reading Kafka User Credentials", fmt.Sprintf("Failed to retrieve user credentials for user with name: %s, cluster ID: %s, error: %s", username, clusterID, err))
			return userCredentials, diags
		}
	}

	return userCredentials, diags
}

// populateUserCredentialsModel populates the user credentials model with information retrieved from the API.
func populateUserCredentialsModel(data *userCredentialsModel, userCredentials kafkaSDK.UserReadAccess) {
	data.CertificateAuthority = types.StringValue(*userCredentials.Metadata.CertificateAuthority)
	data.PrivateKey = types.StringValue(*userCredentials.Metadata.PrivateKey)
	data.Certificate = types.StringValue(*userCredentials.Metadata.Certificate)
	data.ID = types.StringValue(userCredentials.Id)
	data.Username = types.StringValue(userCredentials.Properties.Name)
}
