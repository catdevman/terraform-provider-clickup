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

func (c *ClickUpUserGroupsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse){
    resp.TypeName = req.ProviderTypeName + "_teams"
}

func (c *ClickUpUserGroupsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse){
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

func (c *ClickUpUserGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse){
    var data ClickUpUserGroupsDataSourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
    if resp.Diagnostics.HasError(){
        return
    }

    groups, _, err := c.client.UserGroups.GetUserGroups(ctx, &clickup.GetUserGroupsOptions{
        TeamID: data.TeamId.ValueString(),
    })
    if err != nil {
        resp.Diagnostics.AddError(
            "failed to make call to ClickUp API",
            fmt.Sprintf("err: %s", err),
        )
    }


    var group ClickUpUserGroupDataSourceModel

    for _, g := range groups {
        group = ClickUpUserGroupDataSourceModel{
            Id: types.StringValue(g.ID),
        }
        data.Groups = append(data.Groups, group)
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
