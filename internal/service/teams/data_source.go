package teams

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/raksul/go-clickup/clickup"
)

var _ datasource.DataSource = &ClickUpTeamsDataSource{}

func NewDataSource() datasource.DataSource {
    return &ClickUpTeamsDataSource{}
}

type ClickUpTeamsDataSource struct {
    client *clickup.Client
}

func (c *ClickUpTeamsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse){
    resp.TypeName = req.ProviderTypeName + "_teams"
}

func (c *ClickUpTeamsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse){
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

func (c *ClickUpTeamsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse){
    var data ClickUpTeamsDataSourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
    if resp.Diagnostics.HasError(){
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
        team = ClickUpTeamDataSourceModel{
            Id: types.StringValue(t.ID),
            Name: types.StringValue(t.Name),
            Color: types.StringValue(t.Color),
            Members: convertTeamMembers(t.Members),
        }
        data.Teams = append(data.Teams, team)
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func convertTeamMembers(members []clickup.TeamMember) []ClickUpTeamMemberDataSourceModel{
    if len(members) == 0 {
        return []ClickUpTeamMemberDataSourceModel{}
    }
    result := make([]ClickUpTeamMemberDataSourceModel, len(members), len(members))

    for _, member := range members{
        mem := ClickUpTeamMemberDataSourceModel {
            User: ClickUpTeamUserDataSourceModel {
                Id: types.Int64Value(int64(member.User.ID)),
                Username: types.StringValue(member.User.Username),
                Email: types.StringValue(member.User.Email),
                Color: types.StringValue(member.User.Color),
                ProfilePicture: types.StringValue(member.User.ProfilePicture),
                Role: types.Int64Value(int64(member.User.Role)),
                LastActive: types.StringValue(member.User.LastActive),
                DateJoined: types.StringValue(member.User.DateJoined),
                DateInvited: types.StringValue(member.User.DateInvited),
            },
            InvitedBy: ClickUpTeamInvitedByDataSourceModel {
                Id: types.Int64Value(int64(member.InvitedBy.ID)),
                Username: types.StringValue(member.InvitedBy.Username),
                Email: types.StringValue(member.InvitedBy.Email),
                Color: types.StringValue(member.InvitedBy.Color),
                ProfilePicture: types.StringValue(member.InvitedBy.ProfilePicture),
                Initials: types.StringValue(member.InvitedBy.Initials),
            },
        }
        result = append(result, mem)
    }
    return result
}
