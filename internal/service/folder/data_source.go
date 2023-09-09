package folder

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	// "github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/raksul/go-clickup/clickup"
)


var _ datasource.DataSource = &ClickUpFolderDataSource{}

func NewDataSource() datasource.DataSource {
    return &ClickUpFolderDataSource{}
}

type ClickUpFolderDataSource struct {
    client *clickup.Client
}

func (c *ClickUpFolderDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_folder"
}

func (c *ClickUpFolderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (c *ClickUpFolderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ClickUpFolderWrapperDataSourceModel

    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
    if resp.Diagnostics.HasError(){
        return
    }

    folder, _, err := c.client.Folders.GetFolder(ctx, data.FolderId.ValueString())
    if err != nil {
        resp.Diagnostics.AddError(
            "ClickUp Client had issue getting Spaces",
            fmt.Sprintf("Error: %s", err),
        )
        return
    }

    oi, _ := folder.Orderindex.Int64()
    data.Folder = &ClickUpFolderDataSourceModel{
        Id: types.StringValue(folder.ID),
        Name: types.StringValue(folder.Name),
        Hidden: types.BoolValue(folder.Hidden),
        OrderIndex: types.Int64Value(oi),
        TaskCount: types.StringValue(folder.TaskCount.String()),
        OverrideStatuses: types.BoolValue(folder.OverrideStatuses),
    }
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
