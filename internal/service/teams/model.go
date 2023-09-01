package teams

import "github.com/hashicorp/terraform-plugin-framework/types"

type ClickUpTeamsDataSourceModel struct {
    Teams []ClickUpTeamDataSourceModel `tfsdk:"teams"`
}

type ClickUpTeamDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    Color types.String `tfsdk:"color"`
    Members []ClickUpTeamMemberDataSourceModel `tfsdk:"members"`
}

type ClickUpTeamMemberDataSourceModel struct {
    User ClickUpTeamUserDataSourceModel `tfsdk:"user"`
    InvitedBy ClickUpTeamInvitedByDataSourceModel `tfsdk:"invited_by"`
}

type ClickUpTeamUserDataSourceModel struct {
        Id types.Int64 `tfsdk:"id"`
        Username types.String `tfsdk:"username"`
        Email types.String `tfsdk:"email"`
        Color types.String `tfsdk:"color"`
        ProfilePicture types.String `tfsdk:"profile_picture"`
        Initials types.String `tfsdk:"initials"`
        Role types.Int64 `tfsdk:"role"`
        LastActive types.String `tfsdk:"last_active"`
        DateJoined types.String `tfsdk:"date_joined"`
        DateInvited types.String `tfsdk:"date_invited"`
}

type ClickUpTeamInvitedByDataSourceModel struct{
        Id types.Int64 `tfsdk:"id"`
        Username types.String `tfsdk:"username"`
        Email types.String `tfsdk:"email"`
        Color types.String `tfsdk:"color"`
        ProfilePicture types.String `tfsdk:"profile_picture"`
        Initials types.String `tfsdk:"initials"`
}
