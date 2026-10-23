package mongov2

import "github.com/hashicorp/terraform-plugin-framework/datasource"

// DataSources returns the list of data sources for the mongov2 package.
func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewClusterDataSource,
		NewClustersDataSource,
		NewTemplatesDataSource,
		NewSnapshotsDataSource,
		NewSnapshotLocationsDataSource,
		NewVersionsDataSource,
	}
}
