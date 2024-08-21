package usergroups

import "github.com/hashicorp/terraform-plugin-framework/types"

type ClickUpUserGroupsDataSourceModel struct {
	TeamId types.String                      `tfsdk:"team_id"`
	Groups []ClickUpUserGroupDataSourceModel `tfsdk:"groups"`
}

type ClickUpUserGroupDataSourceModel struct {
	Id          types.String                        `tfsdk:"id"`
	UserId      types.String                        `tfsdk:"userid"`
	Name        types.String                        `tfsdk:"name"`
	Handle      types.String                        `tfsdk:"handle"`
	DateCreated types.String                        `tfsdk:"date_created"`
	Initials    types.String                        `tfsdk:"initials"`
	Members     []ClickUpUserGroupMemberSourceModel `tfsdk:"members"`
	Avatar      ClickUpUserGroupAvatarSourceModel   `tfsdk:"avatar"`
}

type ClickUpUserGroupMemberSourceModel struct {
	ID             types.String `tfsdk:"id"`
	Username       types.String `tfsdk:"username"`
	Email          types.String `tfsdk:"email"`
	Color          types.String `tfsdk:"color"`
	Intials        types.String `tfsdk:"initials"`
	ProfilePicture types.String `tfsdk:"profilePicture"`
}

type ClickUpUserGroupAvatarSourceModel struct {
	AttachmentId types.String `tfsdk:"attachment_id"`
	Color        types.String `tfsdk:"color"`
	Source       types.String `tfsdk:"source"`
	Icon         types.String `tfsdk:"icon"`
}
