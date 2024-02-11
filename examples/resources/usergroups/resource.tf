resource "clickup_usergroup" "test_group" {
  name    = "some-test-name"
  team_id = "123"
  members = [456]
}
