
terraform {
  required_providers {
    clickup = {
      source = "terraform.cntxt.ai/infra/clickup"
    }
  }
}

provider "clickup" {
  api_token = "API_KEY"
}

data "clickup_usergroups" "all" {
  team_id = 123
}

output "usergroups" {
  value = data.clickup_usergroups.all.groups
}
