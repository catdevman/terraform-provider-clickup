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
    Seats ClickUpTeamSeatsSourceModel `tfsdk:"seats"`
    Plan ClickUpTeamPlanSourceModel `tfsdk:"plan"`
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

type ClickUpTeamSeatsSourceModel struct {
    Members ClickUpSeatMembersSourceModel `tfsdk:"members" json:"members"`
    Guests ClickUpSeatGuestsSourceModel `tfsdk:"guests" json:"guests"`
}

type ClickUpSeatMembersSourceModel struct {
    FilledSeats types.Int64 `tfsdk:"filled_members_seats" json:"filled_members_seats"`
    TotalSeats types.Int64 `tfsdk:"total_member_seats" json:"total_member_seats"`
    EmptySeats types.Int64 `tfsdk:"empty_member_seats" json:"empty_member_seats"`
}


type ClickUpSeatGuestsSourceModel struct {
    FilledSeats types.Int64 `tfsdk:"filled_guest_seats" json:"filled_guest_seats"`
    TotalSeats types.Int64 `tfsdk:"total_guest_seats" json:"total_guest_seats"` 
    EmptySeats types.Int64 `tfsdk:"empty_guest_seats" json:"empty_guest_seats"`
}

type ClickUpTeamPlanSourceModel struct {
    Id types.Int64 `tfsdk:"id" json:"plan_id"`
    Name types.String `tfsdk:"name" json:"plan_name"`
}
