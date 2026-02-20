package ionoscloud

import (
	"bytes"
	"html/template"
)

// getConfigurationFromTemplate use text templates to insert values that need to be generated at test runtime (such as expiry dates). Used only for testing.
func getConfigurationFromTemplate(templateConfig string, data interface{}) string {
	T, _ := template.New("").Parse(templateConfig)
	b := &bytes.Buffer{}
	if err := T.Execute(b, data); err != nil {
		panic(err)
	}
	return b.String()
}
