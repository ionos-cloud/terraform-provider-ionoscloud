package acctest

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	objstorage "github.com/ionos-cloud/sdk-go-object-storage"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/envar"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/provider"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/ionoscloud"
	monitoringService "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/monitoring"
	objstorageservice "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/objectstorage"
)

const (
	// ProviderName is the name of the provider
	ProviderName = "ionoscloud"
)

// testAccProviderConfigure ensures Provider is only configured once
//
// The PreCheck(t) function is invoked for every test and this prevents
// extraneous reconfiguration to the same values each time. However, this does
// not prevent reconfiguration that may happen should the address of
// Provider be errantly reused in ProviderFactories.
var testAccProviderConfigure sync.Once

var (
	// TestAccProtoV6ProviderFactories is a map of provider names to provider factories
	TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		ProviderName: func() (tfprotov6.ProviderServer, error) {
			ctx := context.Background()
			upgradedSdkServer, err := tf5to6server.UpgradeServer(
				ctx,
				ionoscloud.Provider().GRPCProvider, // Example terraform-plugin-sdk provider
			)
			if err != nil {
				log.Fatal(err)
			}

			providers := []func() tfprotov6.ProviderServer{
				providerserver.NewProtocol6(provider.New()), // Example terraform-plugin-framework provider
				func() tfprotov6.ProviderServer {
					return upgradedSdkServer
				}, // Example terraform-plugin-sdk provider
			}

			muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
			if err != nil {
				return nil, err
			}

			return muxServer.ProviderServer(), nil
		},
	}
)

// PreCheck ensures the necessary environment variables are set for acceptance testing
func PreCheck(t *testing.T) {
	t.Helper()

	// Since we are outside the scope of the Terraform configuration we must
	// call Configure() to properly initialize the provider configuration.
	testAccProviderConfigure.Do(func() {
		envar.FailIfAllEmpty(t, []string{envar.IonosToken, envar.IonosUsername, envar.IonosPassword, envar.IonosS3AccessKey, envar.IonosS3SecretKey}, "credentials for running acceptance testing")

		username := os.Getenv(envar.IonosUsername)
		password := os.Getenv(envar.IonosPassword)
		token := os.Getenv(envar.IonosToken)

		if token == "" {
			if username == "" || password == "" {
				t.Fatalf("%s/%s or %s must be set for acceptance tests", envar.IonosUsername, envar.IonosPassword, envar.IonosToken)
			}
		}
	})
}

// ObjectStorageClient returns a new S3 client for acceptance testing
func ObjectStorageClient() (*objstorage.APIClient, error) {
	accessKey := os.Getenv(envar.IonosS3AccessKey)
	secretKey := os.Getenv(envar.IonosS3SecretKey)
	insecureStr := os.Getenv("IONOS_ALLOW_INSECURE")
	insecureBool := false
	if insecureStr != "" {
		boolValue, err := strconv.ParseBool(insecureStr)
		if err != nil {
			log.Fatal(err)
		}
		insecureBool = boolValue

	}
	if accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("%s and %s must be set for acceptance tests", envar.IonosS3AccessKey, envar.IonosS3SecretKey)
	}

	return objstorageservice.NewClient(accessKey, secretKey, "", "", insecureBool).GetBaseClient(), nil
}

func MonitoringClient() *monitoringService.MonitoringClient {
	token := os.Getenv(envar.IonosToken)
	username := os.Getenv(envar.IonosUsername)
	password := os.Getenv(envar.IonosPassword)
	insecureStr := os.Getenv(envar.IonosInsecure)

	insecureBool := false
	if insecureStr != "" {
		boolValue, err := strconv.ParseBool(insecureStr)
		if err != nil {
			log.Fatal(err)
		}
		insecureBool = boolValue
	}

	return monitoringService.NewClient(username, password, token, "", "", insecureBool)
}
