# Terraform Provider ClickUp

# Data Sources
- [x] Teams (Workspaces)
  - [x] Authorized Teams
  - [x] Seats
  - [x] Plan
- [ ] User Groups
  - [ ] Groups
- [ ] Users
- [ ] User
- [x] Spaces
- [x] Space
- [x] Folders
- [x] Folder
- [x] Lists
- [ ] List
- [x] Folderless Lists
- [ ] Folderless List
- [ ] Tasks
  - [ ] List
  - [ ] Team (Workspace)
  - [ ] By Task Ids
- [ ] Task

# Resources
- [x] User Group
  - [x] Create
  - [x] Read
  - [x] Update
  - [x] Delete
- [ ] User
  - [ ] Create
  - [ ] Read
  - [ ] Update
  - [ ] Delete
- [ ] Space
  - [ ] Create
  - [ ] Read
  - [ ] Update
  - [ ] Delete
- [ ] Folder
  - [ ] Create
  - [ ] Read
  - [ ] Update
  - [ ] Delete
- [ ] List
  - [ ] Create
    - [ ] Foldered
    - [ ] Folderless
  - [ ] Read
  - [ ] Update
  - [ ] Delete
- [ ] Task
  - [ ] Create
  - [ ] Read
  - [ ] Update
  - [ ] Delete

# Development Notes

## Local Go CLI

When working on this provider, you can add this to the top of `go.mod` to use your local dev version of the CLI to test any changes:
```
replace github.com/raksul/go-clickup => ../go-clickup
```

## Log levels

```sh
TF_LOG=DEBUG terraform apply
```
