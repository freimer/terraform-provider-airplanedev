package main

import (
	"context"
	"terraform-provider-airplanedev/airplanedev"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name airplanedev

func main() {
	providerserver.Serve(context.Background(), airplanedev.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/freimer/airplanedev",
	})
}
