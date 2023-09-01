terraform {
  required_providers {
    clickup = {
        version = "0.0.1"
        source = "hashicorp.io/catdevman/clickup"
    }
  }
}

provider "clickup" {
  api_token = "pk_48124324_LUX47IOWYEC42RRWLJNDQ2BJUIU5MGG9"
}

data "clickup_teams" "teams" {}

output "outtie" {
    value = data.clickup_teams.teams.teams[*].id
}

