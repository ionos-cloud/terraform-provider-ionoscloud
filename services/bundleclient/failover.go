package bundleclient

import (
	"fmt"

	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/failover"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

// configureFailover sets config.Servers from the given cloud endpoints and, per the
// failover strategy from the file config, sets config.HTTPClient.Transport:
//   - roundRobin: install a failover round-tripper over all endpoints, in order.
//   - none / "": pin the first endpoint and apply its dedicated TLS transport if needed.
//
// It returns an error when endpoints is empty or the strategy is unknown.
func (c SdkBundle) configureFailover(
	config *ionoscloud.Configuration,
	endpoints []fileconfiguration.Endpoint,
) error {
	product := fileconfiguration.Cloud
	opts := c.fileConfig.GetFailoverOptions()
	if opts == nil {
		opts = &failover.Options{Strategy: failover.None}
	}
	if len(endpoints) == 0 {
		return fmt.Errorf("no endpoints configured for %q to build a failover client", product)
	}

	failoverEndpoints := make([]failover.Endpoint, 0, len(endpoints))
	servers := make(ionoscloud.ServerConfigurations, 0, len(endpoints))
	for _, ep := range endpoints {
		failoverEndpoints = append(failoverEndpoints, failover.Endpoint{
			URL:                 ep.Name,
			SkipTLSVerify:       ep.SkipTLSVerify,
			CertificateAuthData: ep.CertificateAuthData,
		})
		description := fmt.Sprintf("global %s", shared.EndpointOverridden)
		if ep.Location != "" {
			description = fmt.Sprintf("regional %s, region: %s", shared.EndpointOverridden, ep.Location)
		}
		servers = append(servers, ionoscloud.ServerConfiguration{
			URL:         ep.Name,
			Description: description,
		})
	}

	//nolint:exhaustive // cases are NormalizeStrategy() results, not bare Strategy constants
	switch failover.NormalizeStrategy(opts.Strategy) {
	case failover.NormalizeStrategy(failover.RoundRobin):
		config.HTTPClient.Transport = failover.NewRoundTripper(failoverEndpoints, *opts, config.HTTPClient.Transport)
	case failover.NormalizeStrategy(failover.None), "":
		servers = servers[0:1]
		ep := failoverEndpoints[0]
		if ep.SkipTLSVerify || ep.CertificateAuthData != "" {
			config.HTTPClient.Transport = shared.CreateTransport(ep.SkipTLSVerify, ep.CertificateAuthData)
		}
	default:
		return fmt.Errorf(
			"invalid failover strategy %q defined in file config, only %q, %q or an empty value are supported",
			opts.Strategy, failover.RoundRobin, failover.None,
		)
	}

	config.Servers = servers
	return nil
}
