package excel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &extensionsDataSource{}
	_ datasource.DataSourceWithConfigure = &extensionsDataSource{}
)

// NewCoffeesDataSource is a helper function to simplify the provider implementation.
func NewExtensionDataSource() datasource.DataSource {
	return &extensionsDataSource{}
}

type extensionsDataSourceModel struct {
	Extensions []exitensionModel `tfsdk:"extensions"`
}

type exitensionModel struct {
	Extension types.String `tfsdk:"extension"`
	Name      types.String `tfsdk:"name"`
}

type extensionsDataSource struct {
	client *excel_client.Client
}

// Metadata returns the data source type name.
func (d *extensionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_extensions"
}

func (d *extensionsDataSourceModel) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Schema{
			"extensions": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"extension": schema.StringAttribute{Computed: true},
						"name":      schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *extensionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// creates the model
	var state extensionsDataSourceModel

	// uses the client to get the extensions
	extensions, err := d.client.GetExtensions()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unabled to read extensions from excel client",
			err.Error())
		return
	}

	// maps the extensions to the model
	for _, extension := range extensions {
		extensionState := exitensionModel{
			Extension: types.StringValue(extension.Extension),
			Name:      types.StringValue(extension.Name),
		}
		state.Extensions = append(state.Extensions, extensionState)
	}

	// set the state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasErrors() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *extensionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*excel_client.Client)
}
