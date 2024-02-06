package usergroups

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/raksul/go-clickup/clickup"
)

var (
	_ resource.Resource              = &ClickUpUserGroupResourceModel{}
	_ resource.ResourceWithConfigure = &ClickUpUserGroupResourceModel{}
)

func NewResource() resource.Resource {
	return &ClickUpUserGroupResourceModel{}
}

type ClickUpUserGroupResourceModel struct {
	client *clickup.Client
}

func (r *ClickUpUserGroupResourceModel) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_usergroup"
}

// Configure adds the provider configured client to the resource.
func (r *ClickUpUserGroupResourceModel) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (c *ClickUpUserGroupResourceModel) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

// Create creates the resource and sets the initial Terraform state.
func (c *ClickUpUserGroupResourceModel) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	// var userGroup ClickUpUserGroupResourceModel
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
			// Email:          types.StringValue(userGroupMember.Email),
			// Username:       types.StringValue(userGroupMember.Username),
			// Color:          types.StringValue(userGroupMember.Color),
			// Initials:       types.StringValue(userGroupMember.Initials),
			// ProfilePicture: types.StringValue(userGroupMember.ProfilePicture),
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
func (r *ClickUpUserGroupResourceModel) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *ClickUpUserGroupResourceModel) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ClickUpUserGroupResourceModel) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
