terraform {
  required_providers {
    airplanedev = {
      source = "freimer/airplanedev"
    }
  }
}

provider "airplanedev" {}

data "airplanedev_environment" "dev" {
  slug = "dev"
}

output "dev_environment" {
  value = data.airplanedev_environment.dev
}
