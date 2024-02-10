package usergroups

import (
	"context"
	"fmt"
	"io"
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
	tflog.Debug(ctx, "Starting Resource Create")

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
	tflog.Debug(ctx, "Starting Resource Read")

	// Get current state.
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

	// Get refreshed user group value from the API.
	groups, _, err := r.client.UserGroups.GetUserGroups(ctx, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading ClickUp User Groups",
			err.Error(),
		)
		return
	}
	if len(groups) != 1 {
		resp.Diagnostics.AddError(
			"Wrong number of ClickUp User Groups returned from API",
			"Expected 1, got: "+fmt.Sprint(len(groups))+" groups.",
		)
		return
	}

	// our API only returns a list, so we need to get out the first item.
	var group = groups[0]

	state.ID = types.StringValue(group.ID)
	state.Name = types.StringValue(group.Name)
	if group.Handle != "<null>" {
		state.Handle = types.StringValue(group.Handle)
	}

	var memberIDs []int
	for _, member := range group.Members {
		memberIDs = append(memberIDs, member.ID)
	}
	if len(memberIDs) > 0 {
		state.Members = memberIDs
	}

	// Set refreshed state.
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *ClickUpUserGroupsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Starting Resource Update")

	var plan ClickUpUserGroupCreateResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// get the old state, so we can figure out the changes required.
	var oldState ClickUpUserGroupCreateResourceModel
	req.State.Get(ctx, &oldState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var membersToAdd []int
	var membersToRemove []int
	tflog.Debug(ctx, fmt.Sprintf("Total state members: %d", len(oldState.Members)))
	tflog.Debug(ctx, fmt.Sprintf("Total plan members: %d", len(plan.Members)))

	// find the members in the plan that are not in the state.
	// eg: the ones we want to add.
	for _, member := range plan.Members {
		found := false

		for _, stateMember := range oldState.Members {
			if member == stateMember {
				found = true
				break
			}
		}
		if !found {
			tflog.Debug(ctx, fmt.Sprintf("Adding new member: %d", member))
			membersToAdd = append(membersToAdd, member)
		}
	}
	// find the members in the STATE that are not in the plan.
	// eg: the ones we want to remove.
	for _, member := range oldState.Members {
		found := false

		for _, stateMember := range plan.Members {
			if member == stateMember {
				found = true
				break
			}
		}

		if !found {
			tflog.Debug(ctx, fmt.Sprintf("Removing member: %d", member))
			membersToRemove = append(membersToRemove, member)
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Total members to add: %d", len(membersToAdd)))
	tflog.Debug(ctx, fmt.Sprintf("Total members to remove: %d", len(membersToRemove)))

	opts := &clickup.UpdateUserGroupRequest{
		Name: trimQuotes(plan.Name.String()),
	}

	if plan.Handle.String() != "<null>" {
		opts.Handle = plan.Handle.String()
	}
	if len(membersToAdd) > 0 {
		opts.Members.Add = membersToAdd
	}
	if len(membersToRemove) > 0 {
		opts.Members.Remove = membersToRemove
	}

	updatedGroup, updateResponse, err := r.client.UserGroups.UpdateUserGroup(ctx, trimQuotes(oldState.ID.String()), opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not update User Group, unexpected error: "+err.Error(),
			getResponseBody(ctx, updateResponse),
		)
		return
	}

	memberIDs := make([]int, len(updatedGroup.Members))
	for i, member := range updatedGroup.Members {
		memberIDs[i] = member.ID
	}
	tflog.Debug(ctx, fmt.Sprintf("Adding memebers to final state: %d", len(memberIDs)))
	// plan.Members = append(plan.Members, memberIDs...)
	plan.Members = memberIDs

	plan.Name = types.StringValue(updatedGroup.Name)
	// plan.Handle = types.StringValue(trimQuotes(updatedGroup.Handle))
	// if updatedGroup.Handle != "<null>" || updatedGroup.Handle != "" {
	// 	plan.Handle = types.StringValue(trimQuotes(updatedGroup.Handle))
	// }
	// if plan.Handle == types.StringValue("<null>") {
	// 	plan.Handle = types.StringValue("")
	// }

	plan.ID = types.StringValue(updatedGroup.ID)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ClickUpUserGroupsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Starting Resource Delete")

	var state ClickUpUserGroupCreateResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteResponse, err := r.client.UserGroups.DeleteUserGroup(ctx, trimQuotes(state.ID.String()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not delete User Group, unexpected error: "+err.Error(),
			getResponseBody(ctx, deleteResponse),
		)
		return
	}
}

func (r *ClickUpUserGroupsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Starting Resource Import")
	teamID, groupId, err := splitId(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error extracting values from import ID: "+err.Error(),
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
			"Error reading user group: "+err.Error(),
			getResponseBody(ctx, readResponse),
		)
	}

	if len(groups) == 0 {
		resp.Diagnostics.AddError(
			"No user group found with id: "+groupId,
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
			fmt.Errorf("Error setting state: %v", diags).Error(),
			"",
		)
	}

}

// 'id' looks like: 'team_id/id'.
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

func getResponseBody(_ context.Context, res *clickup.Response) string {
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "There was an error getting the response body"
	}

	return string(body)
}
