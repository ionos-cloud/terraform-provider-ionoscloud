package dns

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	dnssdk "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

func ptr[T any](v T) *T { return &v }

func TestDigestTypeFromMnemonic(t *testing.T) {
	tests := map[string]struct {
		in      string
		want    int64
		wantErr bool
	}{
		"sha1 dashed":      {in: "SHA-1", want: 1},
		"sha256 plain":     {in: "SHA256", want: 2},
		"sha256 dashed":    {in: "SHA-256", want: 2},
		"sha384 lowercase": {in: "sha-384", want: 4},
		"gost":             {in: "GOST_R_34.11-94", wantErr: true},
		"unknown":          {in: "MD5", wantErr: true},
		"empty":            {in: "", wantErr: true},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := DigestTypeFromMnemonic(tc.in)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAlgorithmNumber(t *testing.T) {
	got, err := AlgorithmNumber(dnssdk.ALGORITHM_RSASHA256)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, int64(8), got)

	_, err = AlgorithmNumber("ED25519")
	assert.Error(t, err)
}

func TestBuildDSRecord(t *testing.T) {
	key := dnssdk.DnssecKey{
		KeyTag:                  ptr(int32(49057)),
		DigestAlgorithmMnemonic: ptr("SHA-256"),
		Digest:                  ptr("CF58B511B2D8EF99263704A112703586E542E4FA"),
	}

	ds, err := BuildDSRecord(key, dnssdk.ALGORITHM_RSASHA256)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "49057 8 2 CF58B511B2D8EF99263704A112703586E542E4FA", ds.String())

	t.Run("missing fields", func(t *testing.T) {
		_, err := BuildDSRecord(dnssdk.DnssecKey{KeyTag: ptr(int32(1))}, dnssdk.ALGORITHM_RSASHA256)
		assert.Error(t, err)
	})
	t.Run("unsupported digest", func(t *testing.T) {
		bad := key
		bad.DigestAlgorithmMnemonic = ptr("MD5")
		_, err := BuildDSRecord(bad, dnssdk.ALGORITHM_RSASHA256)
		assert.Error(t, err)
	})
}

func TestSigningKey(t *testing.T) {
	_, found := SigningKey(dnssdk.DnssecKeyReadList{})
	assert.False(t, found, "no metadata")

	_, found = SigningKey(dnssdk.DnssecKeyReadList{Metadata: &dnssdk.DnssecKeyReadListMetadata{Items: []dnssdk.DnssecKey{{KeyTag: ptr(int32(1))}}}})
	assert.False(t, found, "no digest")

	keys := dnssdk.DnssecKeyReadList{Metadata: &dnssdk.DnssecKeyReadListMetadata{Items: []dnssdk.DnssecKey{
		{KeyTag: ptr(int32(1))},
		{KeyTag: ptr(int32(2)), Digest: ptr("ABCD")},
	}}}
	key, found := SigningKey(keys)
	if !assert.True(t, found) {
		return
	}
	assert.Equal(t, int32(2), *key.KeyTag)
}

func TestIsZoneNotSigned(t *testing.T) {
	resp := func(status int, body string) *shared.APIResponse {
		return &shared.APIResponse{Response: &http.Response{StatusCode: status}, Payload: []byte(body)}
	}
	notSigned := `{"httpStatus":400,"messages":[{"errorCode":"paas-dns-rest-0438","message":"zone is not signed"}]}`

	assert.True(t, IsZoneNotSigned(resp(400, notSigned)))
	assert.True(t, IsZoneNotSigned(resp(400, `{"messages":[{"message":"Zone is not signed"}]}`)))
	assert.False(t, IsZoneNotSigned(resp(400, `{"messages":[{"message":"bad zone id"}]}`)), "other 400")
	assert.False(t, IsZoneNotSigned(resp(500, notSigned)), "other status")
	assert.False(t, IsZoneNotSigned(nil), "nil response")
}
