package teams

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/raksul/go-clickup/clickup"
)

var _ datasource.DataSource = &ClickUpTeamsDataSource{}

func NewDataSource() datasource.DataSource {
	return &ClickUpTeamsDataSource{}
}

type ClickUpTeamsDataSource struct {
	client *clickup.Client
}

func (c *ClickUpTeamsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_teams"
}

func (c *ClickUpTeamsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (c *ClickUpTeamsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ClickUpTeamsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	teams, _, err := c.client.Teams.GetTeams(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to make call to ClickUp API",
			fmt.Sprintf("err: %s", err),
		)
	}

	var team ClickUpTeamDataSourceModel

	for _, t := range teams {
		seats, err := getTeamSeats(ctx, c.client, t.ID)
		if err != nil {
			resp.Diagnostics.Append(diag.NewWarningDiagnostic("failed to get Team seats", fmt.Sprintf("Failed to get Team seat for team (workspace) with id: %s and error: %s", t.ID, err)))
		}

		plan, err := getTeamPlan(ctx, c.client, t.ID)
		if err != nil {
			resp.Diagnostics.Append(diag.NewWarningDiagnostic("failled to get Team plan", fmt.Sprintf("Failed to get Team plan for team (workspace) with id: %s and error: %s", t.ID, err)))
		}
		team = ClickUpTeamDataSourceModel{
			Id:      types.StringValue(t.ID),
			Name:    types.StringValue(t.Name),
			Color:   types.StringValue(t.Color),
			Members: convertTeamMembers(t.Members),
			Seats:   seats,
			Plan:    plan,
		}
		data.Teams = append(data.Teams, team)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func convertTeamMembers(members []clickup.TeamMember) []ClickUpTeamMemberDataSourceModel {
	if len(members) == 0 {
		return []ClickUpTeamMemberDataSourceModel{}
	}
	result := make([]ClickUpTeamMemberDataSourceModel, len(members))

	for i, member := range members {
		mem := ClickUpTeamMemberDataSourceModel{
			User: ClickUpTeamUserDataSourceModel{
				Id:             types.Int64Value(int64(member.User.ID)),
				Username:       types.StringValue(member.User.Username),
				Email:          types.StringValue(member.User.Email),
				Color:          types.StringValue(member.User.Color),
				ProfilePicture: types.StringValue(member.User.ProfilePicture),
				Role:           types.Int64Value(int64(member.User.Role)),
				LastActive:     types.StringValue(member.User.LastActive),
				DateJoined:     types.StringValue(member.User.DateJoined),
				DateInvited:    types.StringValue(member.User.DateInvited),
			},
			InvitedBy: ClickUpTeamInvitedByDataSourceModel{
				Id:             types.Int64Value(int64(member.InvitedBy.ID)),
				Username:       types.StringValue(member.InvitedBy.Username),
				Email:          types.StringValue(member.InvitedBy.Email),
				Color:          types.StringValue(member.InvitedBy.Color),
				ProfilePicture: types.StringValue(member.InvitedBy.ProfilePicture),
				Initials:       types.StringValue(member.InvitedBy.Initials),
			},
		}
		result[i] = mem
	}
	return result
}

// FIXME: Figure out how to get ints into types.Int64Value so we don't need this struct to translate.
func getTeamSeats(ctx context.Context, client *clickup.Client, teamId string) (ClickUpTeamSeatsSourceModel, error) {
	s := ClickUpTeamSeatsSourceModel{}
	req, err := client.NewRequest(http.MethodGet, fmt.Sprintf("team/%s/seats", teamId), nil)
	if err != nil {
		return ClickUpTeamSeatsSourceModel{}, err
	}

	ty := struct {
		Members struct {
			FilledSeats int64 `json:"filled_members_seats"`
			TotalSeats  int64 `json:"total_member_seats"`
			EmptySeats  int64 `json:"empty_member_seats"`
		} `json:"members"`
		Guests struct {
			FilledSeats int64 `json:"filled_guest_seats"`
			TotalSeats  int64 `json:"total_guest_seats"`
			EmptySeats  int64 `json:"empty_guest_seats"`
		} `json:"guests"`
	}{}

	resp, err := client.Do(context.Background(), req, &ty)
	if err != nil {
		return ClickUpTeamSeatsSourceModel{}, err
	}
	if resp.StatusCode != http.StatusOK {
		tflog.Warn(ctx, "Odd that status code was not 200")
	}
	s.Members = ClickUpSeatMembersSourceModel{
		FilledSeats: types.Int64Value(ty.Members.FilledSeats),
		TotalSeats:  types.Int64Value(ty.Members.TotalSeats),
		EmptySeats:  types.Int64Value(ty.Members.EmptySeats),
	}
	s.Guests = ClickUpSeatGuestsSourceModel{
		FilledSeats: types.Int64Value(ty.Guests.FilledSeats),
		TotalSeats:  types.Int64Value(ty.Guests.TotalSeats),
		EmptySeats:  types.Int64Value(ty.Guests.EmptySeats),
	}
	return s, nil
}

func getTeamPlan(ctx context.Context, client *clickup.Client, teamId string) (ClickUpTeamPlanSourceModel, error) {
	p := ClickUpTeamPlanSourceModel{}
	req, err := client.NewRequest(http.MethodGet, fmt.Sprintf("team/%s/plan", teamId), nil)
	if err != nil {
		return ClickUpTeamPlanSourceModel{}, err
	}

	ty := struct {
		Id   int64  `json:"plan_id"`
		Name string `json:"plan_name"`
	}{}

	resp, err := client.Do(context.Background(), req, &ty)
	if err != nil {
		return ClickUpTeamPlanSourceModel{}, err
	}
	if resp.StatusCode != http.StatusOK {
		tflog.Warn(ctx, "Odd that status code was not 200")
	}
	p.Id = types.Int64Value(ty.Id)
	p.Name = types.StringValue(ty.Name)
	tflog.Info(ctx, fmt.Sprintf("plan: %+v", p))
	return p, nil
}
