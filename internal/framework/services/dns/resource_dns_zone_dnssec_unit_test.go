package dns

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dnssdk "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	dnsservice "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/dns"
)

const notSignedBody = `{"httpStatus":400,"messages":[{"errorCode":"paas-dns-rest-0438","message":"zone is not signed"}]}`

func apiResp(status int, body string) *shared.APIResponse {
	return &shared.APIResponse{Response: &http.Response{StatusCode: status}, Payload: []byte(body)}
}

func signedKeys() dnssdk.DnssecKeyReadList {
	digest := "ABCD"
	return dnssdk.DnssecKeyReadList{Metadata: &dnssdk.DnssecKeyReadListMetadata{Items: []dnssdk.DnssecKey{{Digest: &digest}}}}
}

// fetchSequence returns a keysFetcher that answers with the given results in order, repeating the last one.
type fetchResult struct {
	keys dnssdk.DnssecKeyReadList
	resp *shared.APIResponse
	err  error
}

func fetchSequence(results ...fetchResult) (keysFetcher, *int) {
	calls := 0
	return func(context.Context, string) (dnssdk.DnssecKeyReadList, *shared.APIResponse, error) {
		r := results[min(calls, len(results)-1)]
		calls++
		return r.keys, r.resp, r.err
	}, &calls
}

func available(context.Context, string) error { return nil }

func TestIsKeyAbsent(t *testing.T) {
	assert.True(t, isKeyAbsent(apiResp(http.StatusNotFound, "")))
	assert.True(t, isKeyAbsent(apiResp(http.StatusBadRequest, notSignedBody)))
	assert.False(t, isKeyAbsent(apiResp(http.StatusBadRequest, `{"messages":[{"message":"bad zone id"}]}`)))
	assert.False(t, isKeyAbsent(apiResp(http.StatusInternalServerError, "")))
	assert.False(t, isKeyAbsent(nil))
}

func TestKeyPending(t *testing.T) {
	tests := map[string]struct {
		model zoneDNSSECResourceModel
		want  bool
	}{
		"written by create before the key was read": {
			model: zoneDNSSECResourceModel{KskBits: types.Int64Value(2048)},
			want:  true,
		},
		"key read": {
			model: zoneDNSSECResourceModel{KskBits: types.Int64Value(2048), dnssecKeyModel: dnssecKeyModel{Digest: types.StringValue("ABCD")}},
		},
		"imported, creation arguments unknown": {
			model: zoneDNSSECResourceModel{},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.model.keyPending())
		})
	}
}

func TestWaitForKey(t *testing.T) {
	t.Run("retries until the key appears", func(t *testing.T) {
		fetch, calls := fetchSequence(
			fetchResult{resp: apiResp(http.StatusBadRequest, notSignedBody), err: errors.New("zone is not signed")},
			fetchResult{keys: signedKeys(), resp: apiResp(http.StatusOK, "")},
		)
		keys, err := waitForKey(context.Background(), time.Minute, "zone", available, fetch)
		require.NoError(t, err)
		require.Equal(t, 2, *calls)
		_, found := dnsservice.SigningKey(keys)
		require.True(t, found)
	})

	t.Run("other API errors are not retried", func(t *testing.T) {
		fetch, calls := fetchSequence(fetchResult{resp: apiResp(http.StatusInternalServerError, ""), err: errors.New("boom")})
		_, err := waitForKey(context.Background(), time.Minute, "zone", available, fetch)
		require.Error(t, err)
		require.Equal(t, 1, *calls)
	})

	t.Run("times out while the key never appears", func(t *testing.T) {
		fetch, _ := fetchSequence(fetchResult{resp: apiResp(http.StatusBadRequest, notSignedBody), err: errors.New("zone is not signed")})
		_, err := waitForKey(context.Background(), time.Millisecond, "zone", available, fetch)
		require.Error(t, err)
	})

	t.Run("waits for the zone to be available", func(t *testing.T) {
		zoneCalls := 0
		zoneAvailable := func(context.Context, string) error {
			zoneCalls++
			if zoneCalls < 2 {
				return errors.New("zone is not AVAILABLE")
			}
			return nil
		}
		fetch, _ := fetchSequence(fetchResult{keys: signedKeys()})
		_, err := waitForKey(context.Background(), time.Minute, "zone", zoneAvailable, fetch)
		require.NoError(t, err)
		require.Equal(t, 2, zoneCalls)
	})
}
