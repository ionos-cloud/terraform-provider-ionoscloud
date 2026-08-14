package mariadbv2

import "github.com/hashicorp/terraform-plugin-framework/datasource"

// DataSources returns the list of data sources for the mariadbv2 package.
func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewClusterDataSource,
		NewClustersDataSource,
		NewBackupsDataSource,
		NewBackupLocationsDataSource,
		NewVersionsDataSource,
	}
}
