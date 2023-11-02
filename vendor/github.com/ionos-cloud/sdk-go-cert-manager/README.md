![Alt text](.github/IONOS.CLOUD.BLU.svg?raw=true "Title")

# sdk-go-cert-manager
IONOS Cloud GO SDK for Certificate Manager Service

Using the Certificate Manager Service, you can conveniently provision and manage SSL certificates with IONOS services and your internal connected resources. For the [Application Load Balancer](https://api.ionos.com/docs/cloud/v6/#Application-Load-Balancers-get-datacenters-datacenterId-applicationloadbalancers), you usually need a certificate to encrypt your HTTPS traffic. The service provides the basic functions of uploading and deleting your certificates for this purpose.

### Example

```golang
package main

import (
	"context"
	"fmt"
	"os"

	ionoscloud "github.com/ionos-cloud/sdk-go-cert-manager"
)

func main() {
	//either provide username and password, or token.
	configuration := ionoscloud.NewConfiguration("USERNAME", "PASSWORD", "TOKEN", "HOST_URL")
	apiClient := ionoscloud.NewAPIClient(configuration)
	resources, resp, err := apiClient.CertificatesApi.CertificatesGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CertificatesApi.CertificatesGet`: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
	}
	fmt.Fprintf(os.Stdout, "Response from `CertificatesApi.CertificatesGet`: %v\n", resources)
}
```