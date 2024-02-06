package usergroups

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
	ID             int    `tfsdk:"id"`
	Username       string `tfsdk:"username"`
	Email          string `tfsdk:"email"`
	Color          string `tfsdk:"color"`
	Initials       string `tfsdk:"initials"`
	ProfilePicture string `tfsdk:"profile_picture"`
}

type ClickUpUserGroupAvatarSourceModel struct {
	AttachmentId *string `tfsdk:"attachment_id"`
	Color        *string `tfsdk:"color"`
	Source       *string `tfsdk:"source"`
	Icon         *string `tfsdk:"icon"`
}

type ClickUpUserGroupResourceModel struct {
	Id          types.String                        `tfsdk:"id"`
	TeamId      types.String                        `tfsdk:"team_id"`
	UserId      types.String                        `tfsdk:"user_id"`
	Name        types.String                        `tfsdk:"name"`
	Handle      types.String                        `tfsdk:"handle"`
	DateCreated types.String                        `tfsdk:"date_created"`
	Initials    types.String                        `tfsdk:"initials"`
	Members     []ClickUpUserGroupMemberSourceModel `tfsdk:"members"`
	Avatar      ClickUpUserGroupAvatarSourceModel   `tfsdk:"avatar"`
}
