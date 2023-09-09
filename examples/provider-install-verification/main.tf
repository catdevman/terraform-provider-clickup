terraform {
  required_providers {
    clickup = {
        version = "~> 0.0.1"
        source = "hashicorp.io/catdevman/clickup"
    }
  }
}

variable "CLICKUP_API_KEY" {
    type = string
}

provider "clickup" {
  api_token = var.CLICKUP_API_KEY
}

data "clickup_teams" "teams" {}

data "clickup_usergroups" "groups"{}

data "clickup_spaces" "spaces" {
    team_id = "9014024487"
}

data "clickup_space" "space" {
    space_id = "90140051562"
}

data "clickup_folders" "folders" {
    space_id = "90140051562"
}

data "clickup_folder" "folder" {
    folder_id = "90140096619"
}

data "clickup_lists" "lists" {
    folder_id = "90140103670"
}

output "outtie" {
    value = data.clickup_teams.teams
}

output "outtie2" {
    value = data.clickup_usergroups.groups
}

output "outtie3" {
    value = data.clickup_spaces.spaces
}

output "outtie4" {
    value = data.clickup_space.space
}

output "outtie5" {
    value = data.clickup_folders.folders
}

output "outtie6" {
    value = data.clickup_folder.folder
}
