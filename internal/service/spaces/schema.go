package spaces

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (c *ClickUpSpacesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse){
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "team_id": schema.StringAttribute{
                Required: true,
            },
            "archived": schema.BoolAttribute{
                Optional: true,   
            },
            "spaces": schema.ListNestedAttribute{
                Computed: true,
                NestedObject: schema.NestedAttributeObject{
                    Attributes: map[string]schema.Attribute{
                        "id": schema.StringAttribute{
                            Computed: true,
                        },
                        "name": schema.StringAttribute{
                            Computed: true,
                        },
                        "private": schema.BoolAttribute{
                            Computed: true,
                        },
                        "statuses": schema.ListNestedAttribute{
                            Computed: true,
                            NestedObject: schema.NestedAttributeObject{
                                Attributes: map[string]schema.Attribute{
                                    "status": schema.StringAttribute{
                                        Computed: true,
                                    },
                                    "type": schema.StringAttribute{
                                        Computed: true,
                                    },
                                    "order_index": schema.Int64Attribute{
                                        Computed: true,
                                    },
                                    "color": schema.StringAttribute{
                                        Computed: true,
                                    },
                                },
                            },
                        },
                        "multiple_assignees": schema.BoolAttribute{
                            Computed: true,
                        },
                    },
                },
            },
        },
    }
}
