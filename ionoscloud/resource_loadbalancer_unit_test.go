package ionoscloud

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// The balanced NICs are only expanded with depth >= 1. Read and import used to fetch the load
// balancer without depth, so nic_ids stayed null and generated configuration was invalid.
func TestLoadbalancerReadAndImportSetNicIDs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/datacenters/dc-1/loadbalancers/lb-1") {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		body := `{"id":"lb-1","properties":{"name":"lb"},"entities":{"balancednics":{"id":"lb-1/balancednics"}}}`
		if depth := r.URL.Query().Get("depth"); depth != "" && depth != "0" {
			body = `{"id":"lb-1","properties":{"name":"lb"},"entities":{"balancednics":{"items":[{"id":"nic-1"},{"id":"nic-2"}]}}}`
		}
		if _, err := w.Write([]byte(body)); err != nil {
			t.Errorf("writing response: %v", err)
		}
	}))
	t.Cleanup(server.Close)
	meta := stubMeta(t, server.URL)
	want := []any{"nic-1", "nic-2"}

	t.Run("read", func(t *testing.T) {
		d := resourceLoadbalancer().Data(&terraform.InstanceState{ID: "lb-1", Attributes: map[string]string{"datacenter_id": "dc-1"}})
		if diags := resourceLoadbalancerRead(context.Background(), d, meta); diags.HasError() {
			t.Fatalf("Read: %v", diags)
		}
		if got := d.Get("nic_ids"); !reflect.DeepEqual(got, want) {
			t.Errorf("nic_ids = %v, want %v", got, want)
		}
	})

	t.Run("read keeps the state's order of the same NICs", func(t *testing.T) {
		d := resourceLoadbalancer().Data(&terraform.InstanceState{ID: "lb-1", Attributes: map[string]string{
			"datacenter_id": "dc-1", "nic_ids.#": "2", "nic_ids.0": "nic-2", "nic_ids.1": "nic-1",
		}})
		if diags := resourceLoadbalancerRead(context.Background(), d, meta); diags.HasError() {
			t.Fatalf("Read: %v", diags)
		}
		if got, want := d.Get("nic_ids"), []any{"nic-2", "nic-1"}; !reflect.DeepEqual(got, want) {
			t.Errorf("nic_ids = %v, want %v", got, want)
		}
	})

	t.Run("read keeps the state's order and appends NICs it lacks", func(t *testing.T) {
		d := resourceLoadbalancer().Data(&terraform.InstanceState{ID: "lb-1", Attributes: map[string]string{
			"datacenter_id": "dc-1", "nic_ids.#": "2", "nic_ids.0": "nic-2", "nic_ids.1": "nic-3",
		}})
		if diags := resourceLoadbalancerRead(context.Background(), d, meta); diags.HasError() {
			t.Fatalf("Read: %v", diags)
		}
		if got, want := d.Get("nic_ids"), []any{"nic-2", "nic-1"}; !reflect.DeepEqual(got, want) {
			t.Errorf("nic_ids = %v, want %v", got, want)
		}
	})

	t.Run("read replaces NICs that differ", func(t *testing.T) {
		d := resourceLoadbalancer().Data(&terraform.InstanceState{ID: "lb-1", Attributes: map[string]string{
			"datacenter_id": "dc-1", "nic_ids.#": "2", "nic_ids.0": "nic-3", "nic_ids.1": "nic-1",
		}})
		if diags := resourceLoadbalancerRead(context.Background(), d, meta); diags.HasError() {
			t.Fatalf("Read: %v", diags)
		}
		if got := d.Get("nic_ids"); !reflect.DeepEqual(got, want) {
			t.Errorf("nic_ids = %v, want %v", got, want)
		}
	})

	t.Run("import", func(t *testing.T) {
		d := resourceLoadbalancer().Data(&terraform.InstanceState{ID: "dc-1/lb-1"})
		if _, err := resourceLoadbalancerImporter(context.Background(), d, meta); err != nil {
			t.Fatalf("Import: %v", err)
		}
		if got := d.Get("nic_ids"); !reflect.DeepEqual(got, want) {
			t.Errorf("nic_ids = %v, want %v", got, want)
		}
	})
}

// Update used to detach every NIC in the old list before attaching the new one, so repairing one
// backend, or replacing one, left the load balancer without any for a while.
func TestLoadbalancerUpdateTouchesOnlyChangedNICs(t *testing.T) {
	var mu sync.Mutex
	var calls []string
	var url string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[strings.Index(r.URL.Path, "/cloudapi/v6")+len("/cloudapi/v6"):]
		w.Header().Set("Content-Type", "application/json")
		switch {
		case path == "/requests/req-1/status":
			_, _ = w.Write([]byte(`{"metadata":{"status":"DONE"}}`))
		case r.Method == http.MethodGet:
			_, _ = w.Write([]byte(`{"id":"lb-1","properties":{"name":"lb"},"entities":{"balancednics":{"items":[{"id":"nic-2"},{"id":"nic-3"}]}}}`))
		default:
			mu.Lock()
			calls = append(calls, r.Method+" "+path)
			mu.Unlock()
			w.Header().Set("Location", url+"/cloudapi/v6/requests/req-1/status")
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write([]byte(`{"id":"nic-3"}`))
		}
	}))
	t.Cleanup(server.Close)
	url = server.URL
	ctx := context.Background()
	meta := stubMeta(t, url)
	r := resourceLoadbalancer()

	d := r.Data(&terraform.InstanceState{ID: "lb-1"})
	for k, v := range map[string]any{"datacenter_id": "dc-1", "name": "lb", "nic_ids": []any{"nic-1", "nic-2"}} {
		if err := d.Set(k, v); err != nil {
			t.Fatalf("setting %s: %v", k, err)
		}
	}
	state := d.State()
	config := terraform.NewResourceConfigRaw(map[string]any{"datacenter_id": "dc-1", "name": "lb", "nic_ids": []any{"nic-2", "nic-3"}})
	diff, err := r.Diff(ctx, state, config, meta)
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	if _, diags := r.Apply(ctx, state, diff, meta); diags.HasError() {
		t.Fatalf("Apply: %v", diags)
	}

	want := []string{
		"DELETE /datacenters/dc-1/loadbalancers/lb-1/balancednics/nic-1",
		"POST /datacenters/dc-1/loadbalancers/lb-1/balancednics",
	}
	if !reflect.DeepEqual(calls, want) {
		t.Errorf("calls = %v, want %v", calls, want)
	}
}
