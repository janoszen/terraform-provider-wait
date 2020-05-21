package wait

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"wait_tcp":  dataSourceTcp(),
			"wait_http": dataSourceHttp(),
		},
	}
}
