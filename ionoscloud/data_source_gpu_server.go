package ionoscloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/serverutil"
)

func dataSourceGpuServer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCubeServerRead,
		Schema:      serverutil.SchemaTemplatedServerDatasource,
		Timeouts:    &resourceDefaultTimeouts,
	}
}
