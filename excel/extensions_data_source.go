package excel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type extensionsDataSourceModel struct {
	Extensions []exitensionModel `json:"extensions"`
}

type exitensionModel struct {
	Extension types.String `json:"extension"`
	Name      types.String `json:"name"`
}

type extensionsDataSource struct {
	//client *excel_client.Client
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
