package cloudapilocation

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

// newMockAPIClient returns an APIClient whose requests are served by handler. The SDK
// keeps the /cloudapi/v6 base path from its default server and only swaps scheme+host,
// so handler paths must include the /cloudapi/v6 prefix.
func newMockAPIClient(t *testing.T, handler http.HandlerFunc) *ionoscloud.APIClient {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return ionoscloud.NewAPIClient(ionoscloud.NewConfiguration("", "", "", srv.URL))
}

func TestResolveParentLocation(t *testing.T) {
	locationJSON := func(id, metroRegion string) string {
		if metroRegion == "" {
			return fmt.Sprintf(`{"id":%q,"type":"location","properties":{"name":"loc"}}`, id)
		}
		return fmt.Sprintf(`{"id":%q,"type":"location","properties":{"name":"loc","metroRegion":%q}}`, id, metroRegion)
	}

	tests := []struct {
		name            string
		locationID      string
		metroRegion     string
		status          int
		wantParentID    string
		wantLocationIDs []string
	}{
		{
			name: "child location returns its parent", locationID: "pc/txl/1", metroRegion: "de/txl", status: http.StatusOK,
			wantParentID: "de/txl", wantLocationIDs: []string{"pc/txl/1", "de/txl"},
		},
		{
			name: "three-segment child", locationID: "de/fra/2", metroRegion: "de/fra", status: http.StatusOK,
			wantParentID: "de/fra", wantLocationIDs: []string{"de/fra/2", "de/fra"},
		},
		{
			name: "self-referential metroRegion is not a parent", locationID: "de/fra", metroRegion: "de/fra", status: http.StatusOK,
			wantParentID: "", wantLocationIDs: []string{"de/fra"},
		},
		{
			name: "no metroRegion", locationID: "us/las", metroRegion: "", status: http.StatusOK,
			wantParentID: "", wantLocationIDs: []string{"us/las"},
		},
		{
			name: "fetch failure degrades to the requested location", locationID: "pc/txl/1", status: http.StatusNotFound,
			wantParentID: "", wantLocationIDs: []string{"pc/txl/1"},
		},
		{
			name: "invalid location id", locationID: "wrong",
			wantParentID: "", wantLocationIDs: []string{"wrong"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newMockAPIClient(t, func(w http.ResponseWriter, r *http.Request) {
				if tt.status != http.StatusOK {
					w.WriteHeader(tt.status)
					fmt.Fprint(w, `{"messages":[{"errorCode":"309","message":"Resource does not exist"}]}`)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, locationJSON(tt.locationID, tt.metroRegion))
			})

			info := ResolveParentLocation(context.Background(), client, tt.locationID)
			if info.ParentID != tt.wantParentID {
				t.Errorf("ParentID = %q, want %q", info.ParentID, tt.wantParentID)
			}
			if !slices.Equal(info.LocationIDs, tt.wantLocationIDs) {
				t.Errorf("LocationIDs = %v, want %v", info.LocationIDs, tt.wantLocationIDs)
			}
		})
	}
}

func TestLocationInSet(t *testing.T) {
	tests := []struct {
		name      string
		locations []string
		loc       string
		expect    bool
	}{
		{name: "present", locations: []string{"de/txl", "de/fra"}, loc: "de/fra", expect: true},
		{name: "case-insensitive", locations: []string{"DE/TXL"}, loc: "de/txl", expect: true},
		{name: "absent", locations: []string{"de/txl"}, loc: "us/las", expect: false},
		{name: "empty set", locations: nil, loc: "de/txl", expect: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LocationInSet(tt.locations, tt.loc); got != tt.expect {
				t.Fatalf("LocationInSet(%v, %q) = %v, want %v", tt.locations, tt.loc, got, tt.expect)
			}
		})
	}
}
