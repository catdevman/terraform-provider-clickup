package usergroups

import "github.com/hashicorp/terraform-plugin-framework/types"

type ClickUpUserGroupsDataSourceModel struct {
    TeamId types.String `tfsdk:"team_id"`
    Groups []ClickUpUserGroupDataSourceModel `tfsdk:"groups"`
}

type ClickUpUserGroupDataSourceModel struct {
    Id types.String `tfsdk:"id"`
}

