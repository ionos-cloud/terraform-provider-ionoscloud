package dns

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dnssdk "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

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
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAlgorithmNumber(t *testing.T) {
	got, err := AlgorithmNumber(dnssdk.ALGORITHM_RSASHA256)
	require.NoError(t, err)
	assert.Equal(t, int64(8), got)

	_, err = AlgorithmNumber("ED25519")
	require.Error(t, err)
}

func TestBuildDSRecord(t *testing.T) {
	key := dnssdk.DnssecKey{
		KeyTag:                  new(int32(49057)),
		DigestAlgorithmMnemonic: new("SHA-256"),
		Digest:                  new("CF58B511B2D8EF99263704A112703586E542E4FA"),
	}

	ds, err := BuildDSRecord(key, dnssdk.ALGORITHM_RSASHA256)
	require.NoError(t, err)
	assert.Equal(t, "49057 8 2 CF58B511B2D8EF99263704A112703586E542E4FA", ds.String())

	t.Run("missing fields", func(t *testing.T) {
		_, err := BuildDSRecord(dnssdk.DnssecKey{KeyTag: new(int32(1))}, dnssdk.ALGORITHM_RSASHA256)
		require.Error(t, err)
	})
	t.Run("unsupported digest", func(t *testing.T) {
		bad := key
		bad.DigestAlgorithmMnemonic = new("MD5")
		_, err := BuildDSRecord(bad, dnssdk.ALGORITHM_RSASHA256)
		require.Error(t, err)
	})
}

func TestSigningKey(t *testing.T) {
	_, found := SigningKey(dnssdk.DnssecKeyReadList{})
	assert.False(t, found, "no metadata")

	_, found = SigningKey(dnssdk.DnssecKeyReadList{Metadata: &dnssdk.DnssecKeyReadListMetadata{Items: []dnssdk.DnssecKey{{KeyTag: new(int32(1))}}}})
	assert.False(t, found, "no digest")

	keys := dnssdk.DnssecKeyReadList{Metadata: &dnssdk.DnssecKeyReadListMetadata{Items: []dnssdk.DnssecKey{
		{KeyTag: new(int32(1))},
		{KeyTag: new(int32(2)), Digest: new("ABCD")},
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

func TestIsZoneAvailable(t *testing.T) {
	tests := map[string]struct {
		status    int
		body      string
		wantErr   bool
		permanent bool
	}{
		"available":        {status: http.StatusOK, body: `{"metadata":{"state":"AVAILABLE"}}`},
		"still deploying":  {status: http.StatusOK, body: `{"metadata":{"state":"PROVISIONING"}}`, wantErr: true},
		"zone not found":   {status: http.StatusNotFound, body: `{"messages":[{"message":"not found"}]}`, wantErr: true, permanent: true},
		"invalid token":    {status: http.StatusUnauthorized, body: `{"messages":[{"message":"unauthorized"}]}`, wantErr: true, permanent: true},
		"other client err": {status: http.StatusBadRequest, body: `{"messages":[{"message":"bad id"}]}`, wantErr: true, permanent: true},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer server.Close()

			config := shared.NewConfiguration("", "", "token", server.URL)
			config.MaxRetries = 0
			client := &Client{sdkClient: *dnssdk.NewAPIClient(config)}

			err := client.IsZoneAvailable(context.Background(), "zone")
			if !tc.wantErr {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			var permanent *backoff.PermanentError
			require.Equal(t, tc.permanent, errors.As(err, &permanent))

			// The retry loop must stop at once for permanent errors instead of waiting for the timeout.
			calls := 0
			start := time.Now()
			retryErr := backoff.Retry(func() error {
				calls++
				return client.IsZoneAvailable(context.Background(), "zone")
			}, backoff.WithContext(backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(2*time.Second)), context.Background()))
			require.Error(t, retryErr)
			if tc.permanent {
				assert.Equal(t, 1, calls)
				assert.Less(t, time.Since(start), time.Second)
			} else {
				assert.Greater(t, calls, 1)
			}
		})
	}
}

func TestIsZoneBusy(t *testing.T) {
	resp := func(status int, body string) *shared.APIResponse {
		return &shared.APIResponse{Response: &http.Response{StatusCode: status}, Payload: []byte(body)}
	}
	busy := `{"httpStatus":409,"messages":[{"errorCode":"paas-dns-rest-0513","message":"the zone has too many operations in progress, please retry later"}]}`

	assert.True(t, IsZoneBusy(resp(409, busy)))
	assert.True(t, IsZoneBusy(resp(409, `{"messages":[{"message":"Too many operations in progress"}]}`)))
	assert.False(t, IsZoneBusy(resp(409, `{"messages":[{"message":"key already exists"}]}`)), "other conflict")
	assert.False(t, IsZoneBusy(resp(400, busy)), "other status")
	assert.False(t, IsZoneBusy(nil), "nil response")
}
