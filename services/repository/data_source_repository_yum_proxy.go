package repository

import (
	"github.com/SimCubeLtd/terraform-provider-nexus/schema/common"
	repositorySchema "github.com/SimCubeLtd/terraform-provider-nexus/schema/repository"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceRepositoryYumProxy() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing yum proxy repository.",

		Read: dataSourceRepositoryYumProxyRead,

		Schema: map[string]*schema.Schema{
			// Common schemas
			"id":     common.DataSourceID,
			"name":   repositorySchema.DataSourceName,
			"online": repositorySchema.DataSourceOnline,
			// Proxy schemas
			"cleanup":        repositorySchema.DataSourceCleanup,
			"http_client":    repositorySchema.DataSourceHTTPClient,
			"negative_cache": repositorySchema.DataSourceNegativeCache,
			"proxy":          repositorySchema.DataSourceProxy,
			"routing_rule":   repositorySchema.DataSourceRoutingRule,
			"storage":        repositorySchema.DataSourceStorage,
			// Yum proxy schemas
			"yum_signing": repositorySchema.DataSourceYumSigning,
		},
	}
}

func dataSourceRepositoryYumProxyRead(resourceData *schema.ResourceData, m interface{}) error {
	resourceData.SetId(resourceData.Get("name").(string))

	return resourceYumProxyRepositoryRead(resourceData, m)
}
