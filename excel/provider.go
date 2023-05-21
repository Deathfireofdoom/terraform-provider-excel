package excel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &excelProvider{}
)

func New() provider.Provider {
	return &excelProvider{}
}

type excelProvider struct {
}

type excelProviderModel struct {
}

func (p *excelProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "excel"
}

// Schema defines the provider-level schema for configuration data.
func (p *excelProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	// creates a client
	client := excel_client.New()

	// make client available for DataSource and Resource type in configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *excelProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExtensionsDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *excelProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}
