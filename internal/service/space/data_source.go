package space

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/raksul/go-clickup/clickup"
)

var _ datasource.DataSource = &ClickUpSpaceDataSource{}

func NewDataSource() datasource.DataSource {
	return &ClickUpSpaceDataSource{}
}

type ClickUpSpaceDataSource struct {
	client *clickup.Client
}

func (c *ClickUpSpaceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_space"
}

func (c *ClickUpSpaceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (c *ClickUpSpaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ClickUpSpaceWrapperDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	space, _, err := c.client.Spaces.GetSpace(ctx, data.SpaceId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"ClickUp Client had issue getting Space",
			fmt.Sprintf("Error: %s", err),
		)
		return
	}

	sts := []ClickUpSpaceStatusDataSourceModel{}
	for _, status := range space.Statuses {
		oi, _ := status.Orderindex.Int64()
		s := ClickUpSpaceStatusDataSourceModel{
			Status:     types.StringValue(status.Status),
			Type:       types.StringValue(status.Type),
			OrderIndex: types.Int64Value(oi),
			Color:      types.StringValue(status.Color),
		}
		sts = append(sts, s)
	}

	data.Space = &ClickUpSpaceDataSourceModel{
		Id:                types.StringValue(space.ID),
		Name:              types.StringValue(space.Name),
		Private:           types.BoolValue(space.Private),
		MultipleAssignees: types.BoolValue(space.MultipleAssignees),
		Statuses:          sts,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
