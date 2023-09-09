package spaces

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	// "github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/raksul/go-clickup/clickup"
)


var _ datasource.DataSource = &ClickUpSpacesDataSource{}

func NewDataSource() datasource.DataSource {
    return &ClickUpSpacesDataSource{}
}

type ClickUpSpacesDataSource struct {
    client *clickup.Client
}

func (c *ClickUpSpacesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_spaces"
}

func (c *ClickUpSpacesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (c *ClickUpSpacesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ClickUpSpacesDataSourceModel

    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
    if resp.Diagnostics.HasError(){
        return
    }

    spaces, _, err := c.client.Spaces.GetSpaces(ctx, data.TeamId.ValueString(), false)
    if err != nil {
        resp.Diagnostics.AddError(
            "ClickUp Client had issue getting Spaces",
            fmt.Sprintf("Error: %s", err),
        )
        return
    }

    for _, space := range spaces {
        sts := []ClickUpSpaceStatusDataSourceModel{}
        for _, status := range space.Statuses{
            oi, _  := status.Orderindex.Int64()
            s := ClickUpSpaceStatusDataSourceModel{
                Status: types.StringValue(status.Status),
                Type: types.StringValue(status.Type),
                OrderIndex: types.Int64Value(oi),
                Color: types.StringValue(status.Color),
            }
            sts = append(sts, s)
        }

        sp := ClickUpSpaceDataSourceModel{
            Id: types.StringValue(space.ID),
            Name: types.StringValue(space.Name),
            Private: types.BoolValue(space.Private),
            MultipleAssignees: types.BoolValue(space.MultipleAssignees),
            Statuses: sts,
        }
        data.Spaces = append(data.Spaces, sp)
    }
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
