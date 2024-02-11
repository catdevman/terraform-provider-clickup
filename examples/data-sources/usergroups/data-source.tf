# List all User Groups.
data "clickup_usergroups" "all" {
  team_id = 123
}

output "usergroups" {
  value = data.clickup_usergroups.all.groups
}
