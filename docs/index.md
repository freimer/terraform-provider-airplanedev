---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "airplanedev Provider"
subcategory: ""
description: |-
  Provision Airplane.dev resources.
---

# airplanedev Provider

Provision Airplane.dev resources.

## Example Usage

```terraform
# Configuration-based authentication
# It is highly recommended you use the environment variables
# for the API Key (AIRPLANEDEV_APIKEY) and if possible the Team
# ID (AIRPLANEDEV_TEAMID) instead of checking credentials into
# a version control system (git).  The host typically does not
# need to be specified, and can use the AIRPLANEDEV_HOST env
# var also.  So most provider configs are simply:
# provider "airplanedev" {}
provider "airplanedev" {
  team_id = "tea20220101abcdefg12345"
  api_key = "tkn_averylongtokenwithmanycharacters"
  host     = "api.airplane.dev"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `api_key` (String, Sensitive) API Key for Airplane.dev.  May also be provided via AIRPLANEDEV_APIKEY environment variable.
- `host` (String) API Host for Airplane.dev.  May also be provided via AIRPLANEDEV_HOST environment variable.  If not specified the provider will use the standard API host for Airplane.dev.
- `team_id` (String, Sensitive) Team ID for Airplane.dev.  May also be provided via AIRPLANEDEV_TEAMID environment variable.