package folderlesslists

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (c *ClickUpListsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"space_id": schema.StringAttribute{
				Required: true,
			},
			"lists": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"orderindex": schema.Int64Attribute{
							Computed: true,
						},
						"content": schema.StringAttribute{
							Computed: true,
						},
						"status": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"status": schema.StringAttribute{
									Computed: true,
								},
								"color": schema.StringAttribute{
									Computed: true,
								},
								"hide_label": schema.BoolAttribute{
									Computed: true,
								},
							},
						},
						"priority": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"priority": schema.StringAttribute{
									Computed: true,
								},
								"color": schema.StringAttribute{
									Computed: true,
								},
							},
						},
						"assignee": schema.StringAttribute{
							Computed: true,
						},
						"task_count": schema.StringAttribute{
							Computed: true,
						},
						"due_date": schema.StringAttribute{
							Computed: true,
						},
						"start_date": schema.StringAttribute{
							Computed: true,
						},
						"folder": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed: true,
								},
								"name": schema.StringAttribute{
									Computed: true,
								},
								"hidden": schema.BoolAttribute{
									Computed: true,
								},
								"access": schema.BoolAttribute{
									Computed: true,
								},
							},
						},
						"space": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed: true,
								},
								"name": schema.StringAttribute{
									Computed: true,
								},
								"access": schema.BoolAttribute{
									Computed: true,
								},
							},
						},
						"archived": schema.BoolAttribute{
							Computed: true,
						},
						"override_statuses": schema.BoolAttribute{
							Computed: true,
						},
						"permission_level": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}
