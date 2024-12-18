package monitoring

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// DataSources returns the list of data sources for the objectstorage package.
func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewPipelineDataSource,
	}
}
