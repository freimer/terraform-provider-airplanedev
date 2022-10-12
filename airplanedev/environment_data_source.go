package airplanedev

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/airplanedev/cli/pkg/api"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &environmentDataSource{}
	_ datasource.DataSourceWithConfigure = &environmentDataSource{}
)

type environmentDataSource struct {
	client *api.Client
}

// environmentDataSourceModel maps the data source schema data.
type environmentDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	Slug       types.String `tfsdk:"slug"`
	Name       types.String `tfsdk:"name"`
	TeamID     types.String `tfsdk:"team_id"`
	Default    types.Bool   `tfsdk:"default"`
	CreatedAt  types.String `tfsdk:"created_at"`
	CreatedBy  types.String `tfsdk:"created_by"`
	UpdatedAt  types.String `tfsdk:"updated_at"`
	UpdatedBy  types.String `tfsdk:"updated_by"`
	IsArchived types.Bool   `tfsdk:"is_archived"`
	ArchivedAt types.String `tfsdk:"archived_at"`
}

func NewEnvironmentDataSource() datasource.DataSource {
	return &environmentDataSource{}
}

// Configure adds the provider configured client to the data source.
func (d *environmentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*api.Client)
}

func (d *environmentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

// GetSchema defines the schema for the data source.
func (d *environmentDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"slug": {
				Type:     types.StringType,
				Required: true,
			},
			"name": {
				Type:     types.StringType,
				Computed: true,
			},
			"team_id": {
				Type:     types.StringType,
				Computed: true,
			},
			"default": {
				Type:     types.BoolType,
				Computed: true,
			},
			"created_at": {
				Type:     types.StringType,
				Computed: true,
			},
			"created_by": {
				Type:     types.StringType,
				Computed: true,
			},
			"updated_at": {
				Type:     types.StringType,
				Computed: true,
			},
			"updated_by": {
				Type:     types.StringType,
				Computed: true,
			},
			"is_archived": {
				Type:     types.BoolType,
				Computed: true,
			},
			"archived_at": {
				Type:     types.StringType,
				Computed: true,
			},
		},
	}, nil
}

// Read refreshes the Terraform state with the latest data.
func (d *environmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state environmentDataSourceModel

	// Retrieve values from config
	var slug string
	diags := req.Config.GetAttribute(ctx, path.Root("slug"), &slug)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	env, err := d.client.GetEnv(ctx, slug)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Airplane.dev Environment",
			err.Error(),
		)
		return
	}

	// Map response body to model
	state = environmentDataSourceModel{
		ID:         types.String{Value: env.ID},
		Slug:       types.String{Value: env.Slug},
		Name:       types.String{Value: env.Name},
		TeamID:     types.String{Value: env.TeamID},
		Default:    types.Bool{Value: env.Default},
		CreatedAt:  types.String{Value: env.CreatedAt.String()},
		CreatedBy:  types.String{Value: env.CreatedBy},
		UpdatedAt:  types.String{Value: env.UpdatedAt.String()},
		UpdatedBy:  types.String{Value: env.UpdatedBy},
		IsArchived: types.Bool{Value: env.IsArchived},
		ArchivedAt: types.String{Null: true},
	}
	if env.IsArchived {
		state.ArchivedAt = types.String{Null: false, Value: env.ArchivedAt.String()}
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
