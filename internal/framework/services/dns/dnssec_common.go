package dns

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	dnssdk "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"

	dnsservice "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dns"
)

// dnssecKeyModel holds the attributes describing the signing key of a zone, as returned by the API.
// It is embedded by both the resource and the data source model.
type dnssecKeyModel struct {
	KeyTag                  types.Int64  `tfsdk:"key_tag"`
	Flags                   types.Int64  `tfsdk:"flags"`
	PublicKey               types.String `tfsdk:"public_key"`
	ComposedKeyData         types.String `tfsdk:"composed_key_data"`
	Digest                  types.String `tfsdk:"digest"`
	DigestAlgorithmMnemonic types.String `tfsdk:"digest_algorithm_mnemonic"`
	DigestType              types.Int64  `tfsdk:"digest_type"`
	AlgorithmNumber         types.Int64  `tfsdk:"algorithm_number"`
	DSRecord                types.String `tfsdk:"ds_record"`
}

// dnssecKeyAttributeDescriptions are shared by the resource and the data source schemas.
const (
	descKeyTag          = "The key tag of the zone signing key, as used in the DS record."
	descFlags           = "The flags of the DNSKEY record."
	descPublicKey       = "The Base64 encoded public key of the DNSKEY record."
	descComposedKeyData = "The DNSKEY RDATA: flags, protocol, algorithm and public key."
	descDigest          = "The digest of the DNSKEY, as used in the DS record."
	descDigestMnemonic  = "The mnemonic of the digest algorithm, e.g. SHA256."
	descDigestType      = "The IANA DS digest type number, derived from `digest_algorithm_mnemonic`."
	descAlgorithmNumber = "The IANA DNS Security Algorithm Number, derived from the signing algorithm."
	descDSRecord        = "The DS record to hand over to the registrar, formatted as `<key_tag> <algorithm_number> <digest_type> <digest>`."
)

// setFromKeys fills the key attributes from the API response. The algorithm is the one the key was created with.
// The DS record is left null if it cannot be derived, so that it does not break the read of a valid key.
func (m *dnssecKeyModel) setFromKeys(keys dnssdk.DnssecKeyReadList, fallbackAlgorithm dnssdk.Algorithm) (found bool, algorithm dnssdk.Algorithm, nsecMode *dnssdk.NsecMode) {
	key, found := dnsservice.SigningKey(keys)
	if !found {
		return false, "", nil
	}

	algorithm = fallbackAlgorithm
	if keys.Properties != nil {
		if keys.Properties.KeyParameters.Algorithm != nil {
			algorithm = *keys.Properties.KeyParameters.Algorithm
		}
		nsecMode = keys.Properties.NsecParameters.NsecMode
	}

	m.KeyTag = int64Value(key.KeyTag)
	m.ComposedKeyData = types.StringPointerValue(key.ComposedKeyData)
	m.Digest = types.StringPointerValue(key.Digest)
	m.DigestAlgorithmMnemonic = types.StringPointerValue(key.DigestAlgorithmMnemonic)
	m.Flags = types.Int64Null()
	m.PublicKey = types.StringNull()
	if key.KeyData != nil {
		m.Flags = int64Value(key.KeyData.Flags)
		m.PublicKey = types.StringPointerValue(key.KeyData.PubKey)
	}

	m.DigestType = types.Int64Null()
	if key.DigestAlgorithmMnemonic != nil {
		if digestType, err := dnsservice.DigestTypeFromMnemonic(*key.DigestAlgorithmMnemonic); err == nil {
			m.DigestType = types.Int64Value(digestType)
		}
	}
	m.AlgorithmNumber = types.Int64Null()
	if algorithmNumber, err := dnsservice.AlgorithmNumber(algorithm); err == nil {
		m.AlgorithmNumber = types.Int64Value(algorithmNumber)
	}

	m.DSRecord = types.StringNull()
	if ds, err := dnsservice.BuildDSRecord(key, algorithm); err == nil {
		m.DSRecord = types.StringValue(ds.String())
	}

	return true, algorithm, nsecMode
}

func int64Value(v *int32) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*v))
}

func zoneNotFoundMessage(zoneID string) string {
	return fmt.Sprintf("DNSSEC is not enabled for zone %s", zoneID)
}
