package folders

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (c *ClickUpFoldersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"space_id": schema.StringAttribute{
				Required: true,
			},
			"folders": schema.ListNestedAttribute{
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
						"hidden": schema.BoolAttribute{
							Computed: true,
						},
						"override_statuses": schema.BoolAttribute{
							Computed: true,
						},
						// "space": schema.SingleNestedAttribute{
						//     Computed: true,
						//     Attributes: map[string]schema.Attribute{
						//         "id": schema.StringAttribute{
						//             Computed: true,
						//         },
						//         "name": schema.StringAttribute{
						//             Computed: true,
						//         },
						//         "access": schema.BoolAttribute{
						//             Computed: true,
						//         },
						//     },
						// },
						"task_count": schema.StringAttribute{
							Computed: true,
						},
						// "lists": schema.ListNestedAttribute{
						//     Computed: true,
						//     NestedObject: schema.NestedAttributeObject{
						//         Attributes: map[string]schema.Attribute{
						//             "status": schema.StringAttribute{
						//                 Computed: true,
						//             },
						//             "type": schema.StringAttribute{
						//                 Computed: true,
						//             },
						//             "order_index": schema.Int64Attribute{
						//                 Computed: true,
						//             },
						//             "color": schema.StringAttribute{
						//                 Computed: true,
						//             },
						//         },
						//     },
						// },
					},
				},
			},
		},
	}
}
