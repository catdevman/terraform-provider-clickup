package lists

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	//	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raksul/go-clickup/clickup"
)


var _ datasource.DataSource = &ClickUpListsDataSource{}

func NewDataSource() datasource.DataSource {
    return &ClickUpListsDataSource{}
}

type ClickUpListsDataSource struct {
    client *clickup.Client
}

func (c *ClickUpListsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_lists"
}

func (c *ClickUpListsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (c *ClickUpListsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ClickUpListsWrapperDataSourceModel

    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
    if resp.Diagnostics.HasError(){
        return
    }

    lists, _, err := c.client.Lists.GetLists(ctx, data.FolderId.ValueString(), false)
    if err != nil {
        resp.Diagnostics.AddError(
            "ClickUp Client had issue getting Spaces",
            fmt.Sprintf("Error: %s", err),
        )
        return
    }

    for _, l := range lists{
        st := ClickUpListStatusDataSourceModel{
            Status: types.StringValue(l.Status.Status),
            Color: types.StringValue(l.Status.Color),
            HideLabel: types.BoolValue(l.Status.HideLabel),
        }
        p := ClickUpListPriorityDataSourceModel{
            Priority: types.StringValue(l.Priority.Priority),
            Color: types.StringValue(l.Priority.Color),
        }
        
        oi, _ := l.Orderindex.Int64()
        list := ClickUpListDataSourceModel{
            Id: types.StringValue(l.ID),
            Name: types.StringValue(l.Name),
            OrderIndex: types.Int64Value(oi),
            Content: types.StringValue(l.Content),
            Status: st,
            Priority: p,
            Assignee: types.StringValue(l.Assignee.Username),
            TaskCount: types.StringValue(l.TaskCount.String()),
            DueDate: types.StringValue(l.DueDate),
            StartDate: types.StringValue(l.StartDate),
        }
        data.Lists = append(data.Lists, list)
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

