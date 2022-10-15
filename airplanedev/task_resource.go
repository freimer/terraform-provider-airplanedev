package airplanedev

import (
	"context"
	"time"

	"github.com/airplanedev/cli/pkg/api"
	libapi "github.com/airplanedev/lib/pkg/api"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &taskResource{}
	_ resource.ResourceWithConfigure = &taskResource{}
)

// NewTaskResource is a helper function to simplify the provider implementation.
func NewTaskResource() resource.Resource {
	return &taskResource{}
}

// taskResource is the resource implementation.
type taskResource struct {
	client *api.Client
}

// taskResourceModel maps the resource schema data.
type taskResourceModel struct {
	ID          types.String         `tfsdk:"id"`
	Revision    types.String         `tfsdk:"revision"`
	LastUpdated types.String         `tfsdk:"last_updated"`
	Slug        types.String         `tfsdk:"slug"`
	Name        types.String         `tfsdk:"name"`
	Description types.String         `tfsdk:"description"`
	Command     []string             `tfsdk:"command"`
	Arguments   []string             `tfsdk:"arguments"`
	Parameters  []taskParameterModel `tfsdk:"parameters"`
}

// taskParameterModel maps task parameter data.
type taskParameterModel struct {
	Slug      types.String `tfsdk:"slug"`
	Name      types.String `tfsdk:"name"`
	Desc      types.String `tfsdk:"desc"`
	Type      types.String `tfsdk:"type"`
	Component types.String `tfsdk:"component"`
}

// Metadata returns the resource type name.
func (r *taskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_task"
}

// Configure adds the provider configured client to the resource.
func (r *taskResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*api.Client)
}

// GetSchema defines the schema for the resource.
func (r *taskResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"revision": {
				Type:     types.StringType,
				Computed: true,
			},
			"last_updated": {
				Type:     types.StringType,
				Computed: true,
			},
			"slug": {
				Type:     types.StringType,
				Required: true,
			},
			"name": {
				Type:     types.StringType,
				Required: true,
			},
			"description": {
				Type:     types.StringType,
				Optional: true,
			},
			"command": {
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Required: true,
			},
			"arguments": {
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Required: true,
			},
			"parameters": {
				Required: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"name": {
						Type:     types.StringType,
						Required: true,
					},
					"slug": {
						Type:     types.StringType,
						Required: true,
					},
					"type": {
						Type:     types.StringType,
						Required: true,
						/* valid types:
						// All Parameter types.
						const (
							TypeString    Type = "string"
							TypeBoolean   Type = "boolean"
							TypeUpload    Type = "upload"
							TypeInteger   Type = "integer"
							TypeFloat     Type = "float"
							TypeDate      Type = "date"
							TypeDatetime  Type = "datetime"
							TypeConfigVar Type = "configvar"
						)
						*/
					},
					"desc": {
						Type:     types.StringType,
						Optional: true,
					},
					"component": {
						Type:     types.StringType,
						Optional: true,
					},
					// Default:     nil,
					// Constraints: libapi.Constraints{},
				}),
			},
		},
	}, nil
}

// Create a new resource
func (r *taskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan taskResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert plan structures to API request structues
	parameters := make([]libapi.Parameter, len(plan.Parameters))
	for i := 0; i < len(plan.Parameters); i++ {
		parameters[i] = libapi.Parameter{
			Name:        plan.Parameters[i].Name.Value,
			Slug:        plan.Parameters[i].Slug.Value,
			Type:        libapi.Type(plan.Parameters[i].Type.Value),
			Desc:        plan.Parameters[i].Desc.Value,
			Component:   libapi.Component(plan.Parameters[i].Component.Value),
			Constraints: libapi.Constraints{},
		}
	}

	// Create new order
	createTaskResponse, err := r.client.CreateTask(ctx, api.CreateTaskRequest{
		Slug:             plan.Slug.Value,
		Name:             plan.Name.Value,
		Description:      plan.Description.Value,
		Image:            new(string), // fix
		Command:          plan.Command,
		Arguments:        plan.Arguments,
		Parameters:       parameters,
		Configs:          []libapi.ConfigAttachment{},
		Constraints:      libapi.RunConstraints{},
		EnvVars:          map[string]libapi.EnvVarValue{},
		ResourceRequests: map[string]string{},
		Resources:        map[string]string{},
		Kind:             "",
		KindOptions:      map[string]interface{}{},
		Runtime:          "",
		Repo:             "",
		Timeout:          0,
		EnvSlug:          "",
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating task",
			"Could not create task, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.String{Value: createTaskResponse.TaskID}
	plan.Revision = types.String{Value: createTaskResponse.TaskRevisionID}
	plan.LastUpdated = types.String{Value: string(time.Now().Format(time.RFC850))}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *taskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *taskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *taskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
