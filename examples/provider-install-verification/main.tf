terraform {
  required_providers {
    airplanedev = {
      source = "freimer/airplanedev"
    }
  }
}

provider "airplanedev" {}

data "airplanedev_environment" "dev" {}