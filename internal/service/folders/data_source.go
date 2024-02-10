package folders

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	// "github.com/hashicorp/terraform-plugin-log/tflog".
	"github.com/raksul/go-clickup/clickup"
)

var _ datasource.DataSource = &ClickUpFoldersDataSource{}

func NewDataSource() datasource.DataSource {
	return &ClickUpFoldersDataSource{}
}

type ClickUpFoldersDataSource struct {
	client *clickup.Client
}

func (c *ClickUpFoldersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_folders"
}

func (c *ClickUpFoldersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (c *ClickUpFoldersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ClickUpFoldersDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	folders, _, err := c.client.Folders.GetFolders(ctx, data.SpaceId.ValueString(), false)
	if err != nil {
		resp.Diagnostics.AddError(
			"ClickUp Client had issue getting Spaces",
			fmt.Sprintf("Error: %s", err),
		)
		return
	}

	for _, folder := range folders {
		oi, _ := folder.Orderindex.Int64()
		f := ClickUpFolderDataSourceModel{
			Id:               types.StringValue(folder.ID),
			Name:             types.StringValue(folder.Name),
			Hidden:           types.BoolValue(folder.Hidden),
			OrderIndex:       types.Int64Value(oi),
			TaskCount:        types.StringValue(folder.TaskCount.String()),
			OverrideStatuses: types.BoolValue(folder.OverrideStatuses),
		}
		data.Folders = append(data.Folders, f)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
