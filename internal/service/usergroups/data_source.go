package usergroups

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/raksul/go-clickup/clickup"
)

var _ datasource.DataSource = &ClickUpUserGroupsDataSource{}

func NewDataSource() datasource.DataSource {
	return &ClickUpUserGroupsDataSource{}
}

type ClickUpUserGroupsDataSource struct {
	client *clickup.Client
}

func (c *ClickUpUserGroupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_usergroups"
}

func (c *ClickUpUserGroupsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*clickup.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"unexpected resource configure type",
			fmt.Sprintf("expect clickup.Client, got: %T. Please report this issue to the provider developer.", req.ProviderData),
		)

		return
	}
	c.client = client
}

func (c *ClickUpUserGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ClickUpUserGroupsDataSourceModel
	var opts clickup.GetUserGroupsOptions
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !data.TeamId.IsNull() {
		opts.TeamID = data.TeamId.String()
	}

	groups, _, err := c.client.UserGroups.GetUserGroups(ctx, &opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to make call to ClickUp API",
			fmt.Sprintf("err: %s", err),
		)
	}

	var group ClickUpUserGroupDataSourceModel

	for _, g := range groups {
		group = ClickUpUserGroupDataSourceModel{
			Id:          types.StringValue(g.ID),
			UserId:      types.StringValue(fmt.Sprint(g.UserID)),
			Name:        types.StringValue(g.Name),
			Handle:      types.StringValue(g.Handle),
			DateCreated: types.StringValue(g.DateCreated),
			Initials:    types.StringValue(g.Initials),
			Members:     getMembers(ctx, g.Members),
			Avatar:      getAvatar(ctx, g.Avatar),
		}
		data.Groups = append(data.Groups, group)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func getMembers(ctx context.Context, members []clickup.GroupMember) []ClickUpUserGroupMemberSourceModel {
	mems := []ClickUpUserGroupMemberSourceModel{}

	for _, m := range members {
		mem := ClickUpUserGroupMemberSourceModel{
			ID:             types.StringValue(fmt.Sprint(m.ID)),
			Username:       types.StringValue(m.Username),
			Email:          types.StringValue(m.Email),
			Color:          types.StringValue(m.Color),
			Intials:        types.StringValue(m.Initials),
			ProfilePicture: types.StringValue(m.ProfilePicture),
		}
		mems = append(mems, mem)
	}

	return mems
}

// TODO: Figure out why avatar comes back as any
func getAvatar(ctx context.Context, avatar any) ClickUpUserGroupAvatarSourceModel {
	return ClickUpUserGroupAvatarSourceModel{}
}
