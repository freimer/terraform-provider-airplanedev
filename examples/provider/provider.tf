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
