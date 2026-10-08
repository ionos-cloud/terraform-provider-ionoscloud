package ionoscloud

import (
	"context"
	"encoding/json"
	"maps"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/clientoptions"
)

const (
	fwRuleTestTCP           = "TCP"
	fwRuleTestHTTP          = "HTTP"
	fwRuleTestClientTimeout = "client_timeout"
	fwRuleTestPort          = "port"
)

// patchBodyStub is a Cloud API stub that accepts every write, serves the given GET bodies and
// records the body of the last PATCH request.
type patchBodyStub struct {
	t    *testing.T
	url  string
	gets map[string]string

	mu    sync.Mutex
	patch map[string]any
}

func newPatchBodyStub(t *testing.T, gets map[string]string) *patchBodyStub {
	t.Helper()
	stub := &patchBodyStub{t: t, gets: gets}
	server := httptest.NewServer(http.HandlerFunc(stub.serve))
	t.Cleanup(server.Close)
	stub.url = server.URL
	return stub
}

func (s *patchBodyStub) serve(w http.ResponseWriter, r *http.Request) {
	// The SDK prefixes paths with /cloudapi/v6; only the resource path matters here.
	path := r.URL.Path[strings.Index(r.URL.Path, "/cloudapi/v6")+len("/cloudapi/v6"):]
	w.Header().Set("Content-Type", "application/json")
	switch {
	case path == "/requests/req-1/status":
		s.write(w, `{"metadata":{"status":"DONE"}}`)
	case r.Method == http.MethodPatch:
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			s.t.Errorf("failed to decode the PATCH body: %v", err)
		}
		s.mu.Lock()
		s.patch = body
		s.mu.Unlock()
		w.Header().Set("Location", s.url+"/cloudapi/v6/requests/req-1/status")
		w.WriteHeader(http.StatusAccepted)
		s.write(w, `{}`)
	case r.Method == http.MethodGet && s.gets[path] != "":
		s.write(w, s.gets[path])
	default:
		w.WriteHeader(http.StatusNotFound)
		s.write(w, `{"httpStatus":404,"messages":[{"message":"stubbed"}]}`)
	}
}

func (s *patchBodyStub) write(w http.ResponseWriter, body string) {
	if _, err := w.Write([]byte(body)); err != nil {
		s.t.Errorf("failed to write the stubbed response: %v", err)
	}
}

// applyChange plans and applies an update of a resource in the prior state that changes only the
// attributes in change, and returns the PATCH body that was sent.
func applyChange(t *testing.T, stub *patchBodyStub, r *schema.Resource, id string, prior, change map[string]any) map[string]any {
	t.Helper()
	ctx := context.Background()
	t.Setenv(shared.IonosApiUrlEnvVar, "")

	meta := *bundleclient.New(ctx, clientoptions.TerraformClientOptions{
		ClientOptions: shared.ClientOptions{
			Endpoint:    stub.url,
			Credentials: shared.Credentials{Token: "token-for-the-stub"},
		},
	}, nil)

	d := r.Data(&terraform.InstanceState{ID: id})
	for k, v := range prior {
		if err := d.Set(k, v); err != nil {
			t.Fatalf("setting %s in the prior state: %v", k, err)
		}
	}
	state := d.State()

	config := maps.Clone(prior)
	maps.Copy(config, change)
	diff, err := r.Diff(ctx, state, terraform.NewResourceConfigRaw(config), meta)
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	if _, diags := r.Apply(ctx, state, diff, meta); diags.HasError() {
		t.Fatalf("Apply: %v", diags)
	}

	stub.mu.Lock()
	defer stub.mu.Unlock()
	if stub.patch == nil {
		t.Fatal("no PATCH request was sent")
	}
	return stub.patch
}

func assertBodyFields(t *testing.T, body, want map[string]any) {
	t.Helper()
	for k, v := range want {
		if body[k] != v {
			t.Errorf("PATCH body %s = %v, want %v (body: %v)", k, body[k], v, body)
		}
	}
}

// The forwarding-rule PATCH body types require these fields and the SDK always sends them, so an
// update of another attribute must send their current values instead of "", 0 and null.
func TestNetworkLoadBalancerForwardingRuleUpdateSendsRequiredFields(t *testing.T) {
	stub := newPatchBodyStub(t, map[string]string{
		"/datacenters/dc-1/networkloadbalancers/nlb-1/forwardingrules/fr-1": `{"id":"fr-1","properties":{"name":"fr","algorithm":"ROUND_ROBIN","protocol":"TCP",` +
			`"listenerIp":"10.0.0.1","listenerPort":80,"healthCheck":{"clientTimeout":2000},"targets":[{"ip":"10.0.0.2","port":80,"weight":1}]}}`,
	})
	prior := map[string]any{
		"datacenter_id": "dc-1", "networkloadbalancer_id": "nlb-1", "name": "fr", "algorithm": "ROUND_ROBIN", "protocol": fwRuleTestTCP,
		"listener_ip": "10.0.0.1", "listener_port": 80, "health_check": []any{map[string]any{fwRuleTestClientTimeout: 1000}},
		"targets": []any{map[string]any{"ip": "10.0.0.2", fwRuleTestPort: 80, "weight": 1}},
	}

	body := applyChange(t, stub, resourceNetworkLoadBalancerForwardingRule(), "fr-1", prior,
		map[string]any{"health_check": []any{map[string]any{fwRuleTestClientTimeout: 2000}}})

	assertBodyFields(t, body, map[string]any{
		"name": "fr", "algorithm": "ROUND_ROBIN", "protocol": fwRuleTestTCP, "listenerIp": "10.0.0.1", "listenerPort": float64(80),
	})
	targets, ok := body["targets"].([]any)
	if !ok || len(targets) != 1 {
		t.Fatalf("PATCH body targets = %v, want the configured target", body["targets"])
	}
	if ip := targets[0].(map[string]any)["ip"]; ip != "10.0.0.2" {
		t.Errorf("PATCH body target ip = %v, want 10.0.0.2", ip)
	}
}

func TestApplicationLoadBalancerForwardingRuleUpdateSendsRequiredFields(t *testing.T) {
	stub := newPatchBodyStub(t, map[string]string{
		"/datacenters/dc-1/applicationloadbalancers/alb-1/forwardingrules/fr-1": `{"id":"fr-1","properties":{"name":"fr","protocol":"HTTP",` +
			`"listenerIp":"10.0.0.1","listenerPort":80,"clientTimeout":2000}}`,
	})
	prior := map[string]any{
		"datacenter_id": "dc-1", "application_loadbalancer_id": "alb-1", "name": "fr", "protocol": fwRuleTestHTTP,
		"listener_ip": "10.0.0.1", "listener_port": 80, fwRuleTestClientTimeout: 1000,
	}

	body := applyChange(t, stub, resourceApplicationLoadBalancerForwardingRule(), "fr-1", prior,
		map[string]any{fwRuleTestClientTimeout: 2000})

	assertBodyFields(t, body, map[string]any{
		"name": "fr", "protocol": fwRuleTestHTTP, "listenerIp": "10.0.0.1", "listenerPort": float64(80), "clientTimeout": float64(2000),
	})
	if _, ok := body["httpRules"]; ok {
		t.Errorf("PATCH body carries httpRules although they did not change: %v", body)
	}
}
