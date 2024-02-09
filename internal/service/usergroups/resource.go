package usergroups

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
			"handle": schema.StringAttribute{
				MarkdownDescription: "User Group handle",
				Optional:            true,
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
	var userGroup ClickUpUserGroupCreateResourceModel
	diags := req.Plan.Get(ctx, &userGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var userGroupRequest clickup.CreateUserGroupRequest

	tflog.Debug(ctx, "Preparing to make members list")
	for _, member := range userGroup.Members {
		memberId := member
		userGroupRequest.Members = append(userGroupRequest.Members, memberId)
	}
	userGroupRequest.Name = userGroup.Name.String()

	tflog.Debug(ctx, "Preparing to send API request")
	newUserGroup, createResponse, err := c.client.UserGroups.CreateUserGroup(ctx, trimQuotes(userGroup.TeamID.String()), &userGroupRequest)

	if err != nil {
		resp.Diagnostics.AddError(
			"Could not create user group: "+err.Error(),
			getResponseBody(ctx, createResponse),
		)
		return
	}

	tflog.Debug(ctx, "Processing response")
	userGroup.ID = types.StringValue(newUserGroup.ID)
	userGroup.Name = types.StringValue(trimQuotes(newUserGroup.Name))
	for userGroupIndex, userGroupMember := range newUserGroup.Members {
		userGroup.Members[userGroupIndex] = userGroupMember.ID
	}

	tflog.Debug(ctx, "Setting final state")
	diags = resp.State.Set(ctx, userGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *ClickUpUserGroupsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state ClickUpUserGroupCreateResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := &clickup.GetUserGroupsOptions{
		TeamID:   trimQuotes(state.TeamID.String()),
		GroupIDs: []string{trimQuotes(state.ID.String())},
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
		if group.ID == state.ID.String() {
			state.Name = types.StringValue(group.Name)
			break
		}
	}

	// throw error if group is not found
	if state.ID.String() == "" {
		resp.Diagnostics.AddError(
			"Error Reading ClickUp User Groups",
			"Group not found: "+state.ID.String(),
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
	var plan clickup.UserGroup
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var oldState clickup.UserGroup
	req.Config.Get(ctx, &oldState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var membersToAdd []int64

	// find the members in the plan that are not in the state
	for _, member := range plan.Members {
		found := false
		for _, stateMember := range oldState.Members {
			if member == stateMember {
				found = true
				break
			}
		}
		if !found {
			membersToAdd = append(membersToAdd, int64(member.ID))
		}
	}

	opts := &clickup.UpdateUserGroupRequest{
		Name:   plan.Name,
		Handle: plan.Handle,
		Members: clickup.UpdateUserGroupMember{
			Add:    make([]int, len(membersToAdd)),
			Remove: []int{},
		},
	}

	updatedGroup, _, err := r.client.UserGroups.UpdateUserGroup(ctx, plan.TeamID, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating User Group",
			"Could not update user group, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Members = append(plan.Members, updatedGroup.Members...)
	plan.Name = updatedGroup.Name
	plan.Handle = updatedGroup.Handle

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ClickUpUserGroupsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state clickup.UserGroup
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UserGroups.DeleteUserGroup(ctx, state.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting User Group",
			"Could not delete user group, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *ClickUpUserGroupsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	teamID, groupId, err := splitId(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Errorf("error extracting values from import ID: %w", err).Error(),
			"",
		)
	}

	// call read with id and teamID
	opts := &clickup.GetUserGroupsOptions{
		TeamID:   teamID,
		GroupIDs: []string{groupId},
	}

	groups, readResponse, err := r.client.UserGroups.GetUserGroups(ctx, opts)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Errorf("error reading user group: %w", err).Error(),
			getResponseBody(ctx, readResponse),
		)
	}

	if len(groups) == 0 {
		resp.Diagnostics.AddError(
			fmt.Errorf("no user group found with id: %s", groupId).Error(),
			getResponseBody(ctx, readResponse),
		)
	}

	// Set the state
	var userGroup ClickUpUserGroupCreateResourceModel
	var group = groups[0]
	userGroup.ID = types.StringValue(group.ID)
	userGroup.Name = types.StringValue(group.Name)
	userGroup.Members = make([]int, len(group.Members))
	for i, member := range group.Members {
		userGroup.Members[i] = member.ID
	}
	userGroup.TeamID = types.StringValue(group.TeamID)
	diags := resp.State.Set(ctx, &userGroup)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			fmt.Errorf("error setting state: %v", diags).Error(),
			"",
		)
	}

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

func trimQuotes(s string) string {
	return strings.Trim(s, "\"")
}

func getResponseBody(ctx context.Context, res *clickup.Response) string {
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}
