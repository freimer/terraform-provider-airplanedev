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
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
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
			"username": {
				Type:     types.StringType,
				Optional: true,
			},
			"password": {
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
			"Unknown Airplanedev API Host",
			"The provider cannot create the Airplanedev API client as there is an unknown configuration value for the Airplanedev API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the AIRPLANEDEV_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Airplanedev API Username",
			"The provider cannot create the Airplanedev API client as there is an unknown configuration value for the Airplanedev API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the AIRPLANEDEV_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown Airplanedev API Password",
			"The provider cannot create the Airplanedev API client as there is an unknown configuration value for the Airplanedev API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the AIRPLANEDEV_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("AIRPLANEDEV_HOST")
	username := os.Getenv("AIRPLANEDEV_USERNAME")
	password := os.Getenv("AIRPLANEDEV_PASSWORD")

	if !config.Host.IsNull() {
		host = config.Host.Value
	}

	if !config.Username.IsNull() {
		username = config.Username.Value
	}

	if !config.Password.IsNull() {
		password = config.Password.Value
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Airplanedev API Host",
			"The provider cannot create the Airplanedev API client as there is a missing or empty value for the Airplanedev API host. "+
				"Set the host value in the configuration or use the AIRPLANEDEV_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Airplanedev API Username",
			"The provider cannot create the Airplanedev API client as there is a missing or empty value for the Airplanedev API username. "+
				"Set the username value in the configuration or use the AIRPLANEDEV_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Airplanedev API Password",
			"The provider cannot create the Airplanedev API client as there is a missing or empty value for the Airplanedev API password. "+
				"Set the password value in the configuration or use the AIRPLANEDEV_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new Airplanedev client using the configuration values
	client, err := airplanedev.NewClient(&host, &username, &password)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Airplanedev API Client",
			"An unexpected error occurred when creating the Airplanedev API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Airplanedev Client Error: "+err.Error(),
		)
		return
	}

	// Make the Airplanedev client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *airplanedevProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *airplanedevProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
