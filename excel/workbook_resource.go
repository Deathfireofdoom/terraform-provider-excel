package excel

import (
	"context"
	"time"

	"github.com/Deathfireofdoom/excel-client-go/pkg/client"
	"github.com/Deathfireofdoom/excel-client-go/pkg/models"

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
	client *client.ExcelClient
}

type workbookResourceModel struct {
	ID          types.String `tfsdk:"id"`
	LastUpdated types.String `tfsdk:"last_updated"`
	FileName    types.String `tfsdk:"file_name"`
	Extension   types.String `tfsdk:"extension"`
	FolderPath  types.String `tfsdk:"folder_path"`
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
			"file_name": schema.StringAttribute{
				Required: true,
			},
			"folder_path": schema.StringAttribute{
				Required: true,
			},
			"extension": schema.StringAttribute{
				Required: true,
			},
			"last_updated": schema.StringAttribute{
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
	if resp.Diagnostics.HasError() {
		return
	}

	// creates the workbook with help of the client
	workbook, err := r.client.CreateWorkbook(plan.FolderPath.ValueString(), plan.FileName.ValueString(), plan.Extension.ValueString(), "")
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create workbook: %s",
			err.Error(),
		)
		return
	}

	// maps the values we got from the client to the terraform model
	plan.ID = types.StringValue(workbook.ID)
	plan.FileName = types.StringValue(workbook.FileName)
	plan.Extension = types.StringValue(string(workbook.Extension))
	plan.FolderPath = types.StringValue(workbook.FolderPath)

	// updates last_updated
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// sets the state with the populated model
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *workbookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state workbookResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from HashiCups
	workbook, err := r.client.ReadWorkbook(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Workbook",
			"Could not read workbook with ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.ID = types.StringValue(workbook.ID)
	state.FileName = types.StringValue(workbook.FileName)
	state.Extension = types.StringValue(string(workbook.Extension))
	state.FolderPath = types.StringValue(workbook.FolderPath)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *workbookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state workbookResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DeleteWorkbook(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Workbook",
			"Could not delete workbook, unexpected error: "+err.Error(),
		)
		return
	}
}

// update the workbook
func (r *workbookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// old state
	var state workbookResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Retrieve values from plan
	var plan workbookResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Converts tf-workbook-model to excel.Workbook
	workbook := &models.Workbook{
		ID:         state.ID.ValueString(),
		FileName:   plan.FileName.ValueString(),
		Extension:  models.Extension(plan.Extension.ValueString()),
		FolderPath: plan.FolderPath.ValueString(),
	}

	// Update existing order
	_, err := r.client.UpdateWorkbook(workbook)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Workbook",
			"Could not update workbook, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated items from GetOrder as UpdateOrder items are not
	// populated.
	workbook, err = r.client.ReadWorkbook(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Workbook",
			"Could not read Workbook ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Update resource state with updated items and timestamp
	plan.ID = types.StringValue(workbook.ID)
	plan.FileName = types.StringValue(workbook.FileName)
	plan.Extension = types.StringValue(string(workbook.Extension))
	plan.FolderPath = types.StringValue(workbook.FolderPath)

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *workbookResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.ExcelClient)
}
