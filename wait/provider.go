package wait

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"wait_tcp":  dataSourceTcp(),
			"wait_http": dataSourceHttp(),
		},
	}
}
