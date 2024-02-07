package usergroups

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raksul/go-clickup/clickup"
)

var (
	_ resource.Resource                = &ClickUpUserGroupsResource{}
	_ resource.ResourceWithConfigure   = &ClickUpUserGroupsResource{}
	_ resource.ResourceWithImportState = &ClickUpUserGroupsResource{}
)

func NewResource() resource.Resource {
	return &ClickUpUserGroupsResource{}
}

type ClickUpUserGroupsResource struct {
	client *clickup.Client
}

func (r *ClickUpUserGroupsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_usergroup"
}

// Configure adds the provider configured client to the resource.
func (r *ClickUpUserGroupsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*clickup.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *clickup.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Schema defines the schema for the resource.
func (c *ClickUpUserGroupsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "User Group resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "User Group name",
				Required:            true,
			},
			"team_id": schema.StringAttribute{
				MarkdownDescription: "Team ID (Workspace)",
				Required:            true,
			},
			// "user_id": schema.StringAttribute{
			// 	MarkdownDescription: "User ID - Owner of this group",
			// 	Required:            true,
			// },
			"members": schema.ListAttribute{
				MarkdownDescription: "User Group Member Ids",
				ElementType:         types.Int64Type,
				Required:            true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (c *ClickUpUserGroupsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var userGroup clickup.UserGroup
	diags := req.Plan.Get(ctx, &userGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var userGroupRequest clickup.CreateUserGroupRequest

	for _, member := range userGroup.Members {
		memberId := member.ID
		userGroupRequest.Members = append(userGroupRequest.Members, memberId)
	}
	userGroupRequest.Name = userGroup.Name

	newUserGroup, _, err := c.client.UserGroups.CreateUserGroup(ctx, userGroup.TeamID, &userGroupRequest)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user group",
			"Could not create user group, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	userGroup.ID = newUserGroup.ID
	userGroup.Name = newUserGroup.Name
	for userGroupIndex, userGroupMember := range newUserGroup.Members {
		userGroup.Members[userGroupIndex] = clickup.GroupMember{
			ID:             userGroupMember.ID,
			Email:          userGroupMember.Email,
			Color:          userGroupMember.Color,
			Username:       userGroupMember.Username,
			Initials:       userGroupMember.Initials,
			ProfilePicture: userGroupMember.ProfilePicture,
		}
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, userGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *ClickUpUserGroupsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state clickup.UserGroup
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := &clickup.GetUserGroupsOptions{
		TeamID:   state.TeamID,
		GroupIDs: []string{state.ID},
	}

	// Get refreshed use group value from the API
	groups, _, err := r.client.UserGroups.GetUserGroups(ctx, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading ClickUp User Groups",
			err.Error(),
		)
		return
	}

	// loop through groups and get the group with the same id as the state
	for _, group := range groups {
		if group.ID == state.ID {
			state = group
			break
		}
	}

	// throw error if group is not found
	if state.ID == "" {
		resp.Diagnostics.AddError(
			"Error Reading ClickUp User Groups",
			"Group not found: "+state.ID,
		)
		return
	}

	// Overwrite values with refreshed state
	// state.Name = group.name
	// state.Items = []ClickUpUserGroupResourceModel{}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *ClickUpUserGroupsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ClickUpUserGroupsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *ClickUpUserGroupsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, teamID, err := splitId(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Errorf("error extracting values from import ID: %w", err).Error(),
			"",
		)
	}

	// call read with id and teamID
	opts := &clickup.GetUserGroupsOptions{
		TeamID:   teamID,
		GroupIDs: []string{id},
	}

	groups, _, err := r.client.UserGroups.GetUserGroups(ctx, opts)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Errorf("error reading user group: %w", err).Error(),
			"",
		)
	}

	if len(groups) == 0 {
		resp.Diagnostics.AddError(
			fmt.Errorf("no user group found with id: %s", id).Error(),
			"",
		)
	}

	// Set the state
	// var group = groups[0]
	// group.TeamID = teamID
	// group.ID = id
	// diags := resp.State.Set(ctx, &group)
	// if diags.HasError() {
	// 	resp.Diagnostics.AddError(
	// 		fmt.Errorf("error setting state: %v", diags).Error(),
	// 		"",
	// 	)
	// }

}

// id looks like: team_id/id
func splitId(id string) (string, string, error) {
	splitLine := strings.Split(id, "/")
	if len(splitLine) != 2 {
		return "", "", fmt.Errorf("invalid ID. Use the format: team_id/group_id")
	}
	var team_id = splitLine[0]
	var group_id = splitLine[1]
	return team_id, group_id, nil
}
