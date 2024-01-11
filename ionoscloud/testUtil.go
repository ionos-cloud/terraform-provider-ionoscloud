package ionoscloud

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

// mocked server, that responds with jsonResponse to whatever we send to it
func createMockServer(jsonResponse string) *httptest.Server {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := fmt.Fprintln(w, jsonResponse)
		if err != nil {
			log.Println("error while writing to server", err)
		}
	}))
	return ts
}

// we need to connect the mocked server to our httpClient, so we inject a mocked client into cfg
func getMockedClient(jsonResponse string) interface{} {
	ts := createMockServer(jsonResponse)

	cfg := ionoscloud.NewConfiguration("", "", "", ts.URL)
	cfg.HTTPClient = ts.Client()

	return services.SdkBundle{
		CloudApiClient: ionoscloud.NewAPIClient(cfg),
	}
}

func getEmptyTestResourceData(t *testing.T, resourceSchema map[string]*schema.Schema) *schema.ResourceData {
	testMap := map[string]interface{}{}
	var testSchema = resourceSchema
	return schema.TestResourceDataRaw(t, testSchema, testMap)
}

func getConfigurationFromTemplate(templateConfig string, data interface{}) string {
	T, _ := template.New("cfg-template").Parse(templateConfig)
	b := &bytes.Buffer{}
	if err := T.Execute(b, data); err != nil {
		panic(err)
	}
	return b.String()
}
