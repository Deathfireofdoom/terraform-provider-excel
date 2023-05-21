package excel

import (
	"context"
	"time"

	"github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &workbookResource{}
	_ resource.ResourceWithConfigure = &workbookResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewWorkbookResource() resource.Resource {
	return &workbookResource{}
}

type workbookResource struct {
	client *exlcel_client.Client
}

type workbookResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Extension   types.String `tfsdk:"extension"`
	Folder      types.String `tfsdk:"folder"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

// Metadata returns the resource type name.
func (r *workbookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workbook"
}

// Schema defines the schema for the resource.
func (r *workbookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"extension": schema.StringAttribute{
				Computed: true,
			},
			"folder": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *workbookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// creates the model, and populates it with values from the plan
	var plan workbookResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasErrors() {
		return
	}

	// converts workbookResourceModel to excel.Workbook
	workbook := excel.Workbook{}

	// creates the workbook with help of the client
	workbook, err := r.client.CreateWorkbook(ctx, workbook)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create workbook",
			"failed to create workbook: %s",
			err.Error(),
		)
		return
	}

	// maps the values we got from the client to the terraform model
	plan.ID = types.String(workbook.ID)
	plan.Name = types.String(workbook.Name)
	plan.Extension = types.String(workbook.Extension)
	plan.Folder = types.String(workbook.Folder)

	// updates last_updated
	plan.LastUpdated = types.String(time.Now().Format(time.RFC850))

	// sets the state with the populated model
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasErrors() {
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *workbookResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*hashicups.Client)
}
