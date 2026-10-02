package dns

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cenkalti/backoff/v4"
	dns "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// dnssecAlgorithmNumbers maps the algorithms supported by the API to their IANA DNS Security Algorithm Number
// (https://www.iana.org/assignments/dns-sec-alg-numbers), as used in the DS record.
var dnssecAlgorithmNumbers = map[dns.Algorithm]int64{
	dns.ALGORITHM_RSASHA256: 8,
}

// dnssecDigestTypes maps the normalized digest algorithm mnemonic to its IANA DS Digest Type
// (https://www.iana.org/assignments/ds-rr-types).
var dnssecDigestTypes = map[string]int64{
	"SHA1":        1,
	"SHA256":      2,
	"GOSTR341194": 3,
	"SHA384":      4,
}

// zoneNotSignedErrorCode is the API error code returned when a zone has no DNSSEC key.
const zoneNotSignedErrorCode = "paas-dns-rest-0438"

// zoneBusyErrorCode is the API error code (409) returned when a zone has too many operations in progress.
const zoneBusyErrorCode = "paas-dns-rest-0513"

// DNSSECDSRecord holds the values needed to publish a DS record at the registrar.
type DNSSECDSRecord struct {
	KeyTag          int64
	AlgorithmNumber int64
	DigestType      int64
	Digest          string
}

// String renders the DS record RDATA: "<key tag> <algorithm> <digest type> <digest>".
func (r DNSSECDSRecord) String() string {
	return fmt.Sprintf("%d %d %d %s", r.KeyTag, r.AlgorithmNumber, r.DigestType, r.Digest)
}

// CreateDNSSECKey enables DNSSEC for a zone by creating a signing key.
func (c *Client) CreateDNSSECKey(ctx context.Context, zoneID string, properties dns.DnssecKeyParameters) (dns.DnssecKeyReadCreation, *shared.APIResponse, error) {
	key, apiResponse, err := c.sdkClient.DNSSECApi.ZonesKeysPost(ctx, zoneID).DnssecKeyCreate(dns.DnssecKeyCreate{Properties: properties}).Execute()
	apiResponse.LogInfo()
	return key, apiResponse, err
}

// GetDNSSECKeys retrieves the DNSSEC key information of a zone.
func (c *Client) GetDNSSECKeys(ctx context.Context, zoneID string) (dns.DnssecKeyReadList, *shared.APIResponse, error) {
	keys, apiResponse, err := c.sdkClient.DNSSECApi.ZonesKeysGet(ctx, zoneID).Execute()
	apiResponse.LogInfo()
	return keys, apiResponse, err
}

// DeleteDNSSECKey disables DNSSEC for a zone by deleting its signing key.
func (c *Client) DeleteDNSSECKey(ctx context.Context, zoneID string) (*shared.APIResponse, error) {
	_, apiResponse, err := c.sdkClient.DNSSECApi.ZonesKeysDelete(ctx, zoneID).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// IsZoneAvailable returns an error unless the zone is in the AVAILABLE state. It is meant to be used with a retry loop:
// a zone in another state is a retryable error, while a failed request (e.g. missing zone, invalid credentials) is
// wrapped in backoff.Permanent, which stops the retry.
func (c *Client) IsZoneAvailable(ctx context.Context, zoneID string) error {
	zone, _, err := c.GetZoneById(ctx, zoneID)
	if err != nil {
		return backoff.Permanent(err)
	}
	if !strings.EqualFold(string(zone.Metadata.State), string(dns.PROVISIONINGSTATE_AVAILABLE)) {
		return fmt.Errorf("zone %s is in state %s, expected %s", zoneID, zone.Metadata.State, dns.PROVISIONINGSTATE_AVAILABLE)
	}
	return nil
}

// IsZoneNotSigned reports whether the API answered that the zone has no DNSSEC key (yet). The key endpoint returns
// 400 "zone is not signed" instead of 404 for a zone without a key, and also while a new key is still being created.
func IsZoneNotSigned(apiResponse *shared.APIResponse) bool {
	if apiResponse.SafeStatusCode() != http.StatusBadRequest {
		return false
	}
	return strings.Contains(string(apiResponse.Payload), zoneNotSignedErrorCode) ||
		strings.Contains(strings.ToLower(string(apiResponse.Payload)), "zone is not signed")
}

// IsZoneBusy reports whether the API rejected a request because the zone still has operations in progress, e.g. the
// asynchronous removal of a previous key. The request can be repeated later.
func IsZoneBusy(apiResponse *shared.APIResponse) bool {
	if apiResponse.SafeStatusCode() != http.StatusConflict {
		return false
	}
	return strings.Contains(string(apiResponse.Payload), zoneBusyErrorCode) ||
		strings.Contains(strings.ToLower(string(apiResponse.Payload)), "too many operations in progress")
}

// SigningKey returns the key that is published to the registrar, i.e. the first entry that carries a digest.
// It returns false if the zone has no such key.
func SigningKey(keys dns.DnssecKeyReadList) (dns.DnssecKey, bool) {
	if keys.Metadata == nil {
		return dns.DnssecKey{}, false
	}
	for _, item := range keys.Metadata.Items {
		if item.Digest != nil && *item.Digest != "" {
			return item, true
		}
	}
	return dns.DnssecKey{}, false
}

// DigestTypeFromMnemonic converts a digest algorithm mnemonic (e.g. "SHA256" or "SHA-256") to its IANA DS digest type.
func DigestTypeFromMnemonic(mnemonic string) (int64, error) {
	normalized := strings.NewReplacer("-", "", "_", "", " ", "").Replace(strings.ToUpper(mnemonic))
	digestType, ok := dnssecDigestTypes[normalized]
	if !ok {
		return 0, fmt.Errorf("unsupported DNSSEC digest algorithm %q", mnemonic)
	}
	return digestType, nil
}

// AlgorithmNumber converts a signing algorithm to its IANA DNS Security Algorithm Number.
func AlgorithmNumber(algorithm dns.Algorithm) (int64, error) {
	number, ok := dnssecAlgorithmNumbers[dns.Algorithm(strings.ToUpper(string(algorithm)))]
	if !ok {
		return 0, fmt.Errorf("unsupported DNSSEC algorithm %q", algorithm)
	}
	return number, nil
}

// BuildDSRecord assembles the DS record of a signing key.
func BuildDSRecord(key dns.DnssecKey, algorithm dns.Algorithm) (DNSSECDSRecord, error) {
	if key.KeyTag == nil || key.Digest == nil || key.DigestAlgorithmMnemonic == nil {
		return DNSSECDSRecord{}, fmt.Errorf("the DNSSEC key does not contain the key tag, digest and digest algorithm needed for a DS record")
	}
	algorithmNumber, err := AlgorithmNumber(algorithm)
	if err != nil {
		return DNSSECDSRecord{}, err
	}
	digestType, err := DigestTypeFromMnemonic(*key.DigestAlgorithmMnemonic)
	if err != nil {
		return DNSSECDSRecord{}, err
	}
	return DNSSECDSRecord{
		KeyTag:          int64(*key.KeyTag),
		AlgorithmNumber: algorithmNumber,
		DigestType:      digestType,
		Digest:          *key.Digest,
	}, nil
}
