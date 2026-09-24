package ionoscloud

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/clientoptions"
)

const (
	attachTestDC     = "dc-1"
	attachTestVolume = "vol-1"
	attachTestOldSrv = "srv-old"
	attachTestNewSrv = "srv-new"
	attachTestDisk   = "SSD Standard"
)

// volumeAttachStub is a Cloud API stub that records every call and answers the detach with
// detachStatus, so the order of detach and attach calls can be asserted.
type volumeAttachStub struct {
	t            *testing.T
	url          string
	detachStatus int

	mu    sync.Mutex
	calls []string
}

func newVolumeAttachStub(t *testing.T, detachStatus int) *volumeAttachStub {
	t.Helper()
	stub := &volumeAttachStub{t: t, detachStatus: detachStatus}
	server := httptest.NewServer(http.HandlerFunc(stub.serve))
	t.Cleanup(server.Close)
	stub.url = server.URL
	return stub
}

func (s *volumeAttachStub) serve(w http.ResponseWriter, r *http.Request) {
	// The SDK prefixes paths with /cloudapi/v6; only the resource path matters here.
	path := r.URL.Path[strings.Index(r.URL.Path, "/cloudapi/v6")+len("/cloudapi/v6"):]
	call := r.Method + " " + path
	if path != "/requests/req-1/status" {
		s.mu.Lock()
		s.calls = append(s.calls, call)
		s.mu.Unlock()
	}

	accepted := func() {
		w.Header().Set("Location", s.url+"/cloudapi/v6/requests/req-1/status")
		w.WriteHeader(http.StatusAccepted)
	}
	volume := ionoscloud.Volume{
		Id:         new(attachTestVolume),
		Properties: &ionoscloud.VolumeProperties{Size: new(float32(1)), Type: new(attachTestDisk)},
	}

	w.Header().Set("Content-Type", "application/json")
	switch {
	case call == "GET /requests/req-1/status":
		s.write(w, map[string]any{"metadata": map[string]any{"status": "DONE"}})
	case r.Method == http.MethodDelete && strings.HasSuffix(path, "/volumes/"+attachTestVolume):
		if s.detachStatus == http.StatusAccepted {
			accepted()
			return
		}
		w.WriteHeader(s.detachStatus)
		s.write(w, map[string]any{"httpStatus": s.detachStatus, "messages": []map[string]string{{"message": "stubbed"}}})
	case r.Method == http.MethodPost && strings.HasSuffix(path, "/volumes"):
		accepted()
		s.write(w, volume)
	case r.Method == http.MethodGet && strings.HasSuffix(path, "/volumes/"+attachTestVolume):
		s.write(w, volume)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (s *volumeAttachStub) write(w http.ResponseWriter, body any) {
	if err := json.NewEncoder(w).Encode(body); err != nil {
		s.t.Errorf("failed to write the stubbed response: %v", err)
	}
}

func (s *volumeAttachStub) recorded() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.calls...)
}

// applyVolumeServerIDChange runs a plan and apply that moves an attached volume from
// attachTestOldSrv to newServerID ("" removes server_id from the config).
func applyVolumeServerIDChange(t *testing.T, stub *volumeAttachStub, newServerID string) (*terraform.InstanceState, diag.Diagnostics) {
	t.Helper()
	ctx := context.Background()
	t.Setenv(shared.IonosApiUrlEnvVar, "")

	meta := *bundleclient.New(ctx, clientoptions.TerraformClientOptions{
		ClientOptions: shared.ClientOptions{
			Endpoint:    stub.url,
			Credentials: shared.Credentials{Token: "token-for-the-stub"},
		},
	}, nil)

	r := resourceVolume()
	state := &terraform.InstanceState{
		ID: attachTestVolume,
		Attributes: map[string]string{
			"id":            attachTestVolume,
			"datacenter_id": attachTestDC,
			"server_id":     attachTestOldSrv,
			"size":          "1",
			"disk_type":     attachTestDisk,
		},
	}
	config := map[string]any{
		"datacenter_id": attachTestDC,
		"size":          1,
		"disk_type":     attachTestDisk,
	}
	if newServerID != "" {
		config["server_id"] = newServerID
	}

	diff, err := r.Diff(ctx, state, terraform.NewResourceConfigRaw(config), meta)
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	if diff == nil || diff.Attributes["server_id"] == nil {
		t.Fatalf("expected a server_id change in the diff, got %#v", diff)
	}

	return r.Apply(ctx, state, diff, meta)
}

func assertCalls(t *testing.T, got, want []string) {
	t.Helper()
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("API calls:\n  got  %q\n  want %q", got, want)
	}
}

// Moving a volume between servers must detach it from the old server before attaching it to
// the new one, since a volume can only be attached to one server, and must not PATCH the volume.
func TestResourceVolumeUpdate_MoveDetachesBeforeAttach(t *testing.T) {
	stub := newVolumeAttachStub(t, http.StatusAccepted)

	state, diags := applyVolumeServerIDChange(t, stub, attachTestNewSrv)
	if diags.HasError() {
		t.Fatalf("apply: %v", diags)
	}

	assertCalls(t, stub.recorded(), []string{
		"DELETE /datacenters/dc-1/servers/srv-old/volumes/vol-1",
		"POST /datacenters/dc-1/servers/srv-new/volumes",
		"GET /datacenters/dc-1/volumes/vol-1",
		"GET /datacenters/dc-1/servers/srv-new/volumes/vol-1",
	})
	if got := state.Attributes["server_id"]; got != attachTestNewSrv {
		t.Fatalf("server_id = %q, want %q", got, attachTestNewSrv)
	}
}

// Removing server_id detaches the volume and leaves it unattached.
func TestResourceVolumeUpdate_UnsetServerIDOnlyDetaches(t *testing.T) {
	stub := newVolumeAttachStub(t, http.StatusAccepted)

	state, diags := applyVolumeServerIDChange(t, stub, "")
	if diags.HasError() {
		t.Fatalf("apply: %v", diags)
	}

	assertCalls(t, stub.recorded(), []string{
		"DELETE /datacenters/dc-1/servers/srv-old/volumes/vol-1",
		"GET /datacenters/dc-1/volumes/vol-1",
	})
	// SDKv2 keeps an optional string removed from config as "" in state; that is the unset value.
	if got := state.Attributes["server_id"]; got != "" {
		t.Fatalf("server_id = %q, want it unset", got)
	}
}

// A volume the old server no longer holds (detached outside Terraform, or the server is gone)
// counts as detached, so the move still attaches it to the new server.
func TestResourceVolumeUpdate_DetachNotFoundCountsAsDetached(t *testing.T) {
	stub := newVolumeAttachStub(t, http.StatusNotFound)

	state, diags := applyVolumeServerIDChange(t, stub, attachTestNewSrv)
	if diags.HasError() {
		t.Fatalf("apply: %v", diags)
	}

	assertCalls(t, stub.recorded(), []string{
		"DELETE /datacenters/dc-1/servers/srv-old/volumes/vol-1",
		"POST /datacenters/dc-1/servers/srv-new/volumes",
		"GET /datacenters/dc-1/volumes/vol-1",
		"GET /datacenters/dc-1/servers/srv-new/volumes/vol-1",
	})
	if got := state.Attributes["server_id"]; got != attachTestNewSrv {
		t.Fatalf("server_id = %q, want %q", got, attachTestNewSrv)
	}
}

// Any other detach failure stops the update before the attach, so the volume is never left
// attached to the new server while the old attachment is still in place.
func TestResourceVolumeUpdate_DetachFailureStopsBeforeAttach(t *testing.T) {
	stub := newVolumeAttachStub(t, http.StatusUnprocessableEntity)

	if _, diags := applyVolumeServerIDChange(t, stub, attachTestNewSrv); !diags.HasError() {
		t.Fatal("expected the failed detach to fail the update")
	}

	assertCalls(t, stub.recorded(), []string{
		"DELETE /datacenters/dc-1/servers/srv-old/volumes/vol-1",
	})
}
