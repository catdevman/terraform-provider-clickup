terraform {
  required_providers {
    clickup = {
        version = "0.0.1"
        source = "hashicorp.io/catdevman/clickup"
    }
  }
}

provider "clickup" {
  api_token = "pk_48124324_H5YETRF03658ON0G6IF9F859Y00P1MEW"
}

data "clickup_teams" "teams" {}

output "outtie" {
    value = data.clickup_teams.teams
}
