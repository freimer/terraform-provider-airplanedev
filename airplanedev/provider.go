package airplanedev

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/airplanedev/cli/pkg/api"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &airplanedevProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &airplanedevProvider{}
}

// airplanedevProvider is the provider implementation.
type airplanedevProvider struct{}

// airplanedevProviderModel maps provider schema data to a Go type.
type airplanedevProviderModel struct {
	Host   types.String `tfsdk:"host"`
	APIKey types.String `tfsdk:"api_key"`
	TeamID types.String `tfsdk:"team_id"`
}

// Metadata returns the provider type name.
func (p *airplanedevProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "airplanedev"
}

// GetSchema defines the provider-level schema for configuration data.
func (p *airplanedevProvider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"host": {
				Type:     types.StringType,
				Optional: true,
			},
			"api_key": {
				Type:      types.StringType,
				Optional:  true,
				Sensitive: true,
			},
			"team_id": {
				Type:      types.StringType,
				Optional:  true,
				Sensitive: true,
			},
		},
	}, nil
}

func (p *airplanedevProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config airplanedevProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Airplanedev API Host (host)",
			"The provider cannot create the Airplanedev API client as there is an unknown configuration value for the Airplanedev API Host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the AIRPLANEDEV_HOST environment variable.",
		)
	}

	if config.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Airplanedev APIKey (api_key)",
			"The provider cannot create the Airplanedev API client as there is an unknown configuration value for the Airplanedev APIKey. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the AIRPLANEDEV_APIKEY environment variable.",
		)
	}

	if config.TeamID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("team_id"),
			"Unknown Airplanedev TeamID (team_id)",
			"The provider cannot create the Airplanedev API client as there is an unknown configuration value for the Airplanedev TeamID. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the AIRPLANEDEV_TEAMID environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("AIRPLANEDEV_HOST")
	apiKey := os.Getenv("AIRPLANEDEV_APIKEY")
	teamID := os.Getenv("AIRPLANEDEV_TEAMID")

	if !config.Host.IsNull() {
		host = config.Host.Value
	}

	if !config.APIKey.IsNull() {
		apiKey = config.APIKey.Value
	}

	if !config.TeamID.IsNull() {
		teamID = config.TeamID.Value
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing Airplanedev APIKey (api_key)",
			"The provider cannot create the Airplanedev API client as there is a missing or empty value for the Airplanedev APIKey. "+
				"Set the api_key value in the configuration or use the AIRPLANEDEV_APIKEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if teamID == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("team_id"),
			"Missing Airplanedev TeamID (team_id)",
			"The provider cannot create the Airplanedev API client as there is a missing or empty value for the Airplanedev TeamID. "+
				"Set the team_id value in the configuration or use the AIRPLANEDEV_TEAMID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new Airplanedev client using the configuration values
	client := &api.Client{
		Host:   host,
		APIKey: apiKey,
		TeamID: teamID}

	// Make the Airplanedev client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *airplanedevProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewEnvironmentDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *airplanedevProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
