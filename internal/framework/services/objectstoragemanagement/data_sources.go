package objectstoragemanagement

import "github.com/hashicorp/terraform-plugin-framework/datasource"

// DataSources returns the list of data sources for the objectstoragemanagement package.
func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewRegionDataSource,
		NewAccesskeyDataSource,
	}
}
