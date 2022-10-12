package main

import (
	"context"
	"terraform-provider-airplanedev/airplanedev"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	providerserver.Serve(context.Background(), airplanedev.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/freimer/airplanedev",
	})
}
