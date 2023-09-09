package usergroups

import (
	"context"

	"github.com/catdevman/terraform-provider-clickup/internal/consts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (c *ClickUpUserGroupsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse){
    resp.Schema = schema.Schema{
        MarkdownDescription: "Use this data source to get UserGroups for Team (Workspace)",
        Attributes: map[string]schema.Attribute{
            "team_id": schema.StringAttribute{
                Description: "",
                Optional: true,
            },
            "groups": schema.ListNestedAttribute{
                Computed: true,
                NestedObject: schema.NestedAttributeObject{
                    Attributes: map[string]schema.Attribute{
                        consts.IDSchemaKey: schema.StringAttribute{
                            Description: consts.IDSchemaDescription,
                            Computed: true,
                        },
                        "userid": schema.StringAttribute{
                            Computed: true,
                        },
                        "name": schema.StringAttribute{
                            Computed: true,
                        },
                        "handle": schema.StringAttribute{
                            Computed: true,
                        },
                        "date_created": schema.StringAttribute{
                            Computed: true,
                        },
                        "initials": schema.StringAttribute{
                            Computed: true,
                        },
                        "members": schema.ListNestedAttribute{
                            Computed: true,
                            NestedObject: schema.NestedAttributeObject{
                                Attributes: map[string]schema.Attribute{
                                    "id": schema.Int64Attribute{
                                        Computed: true,
                                    },
                                    "username": schema.StringAttribute{
                                        Computed: true,
                                    },
                                    "email": schema.StringAttribute{
                                        Computed: true,
                                    },
                                    "color": schema.StringAttribute{
                                        Computed: true,
                                    },
                                    "initials": schema.StringAttribute{
                                        Computed: true,
                                    },
                                    "profile_picture": schema.StringAttribute{
                                        Computed: true,
                                    },
                                },
                            },
                        },
                        "avatar": schema.SingleNestedAttribute{
                            Computed: true,
                            Attributes: map[string]schema.Attribute{
                                "attachment_id": schema.StringAttribute{
                                    Computed: true,
                                },
                                "color": schema.StringAttribute{
                                    Computed: true,
                                },
                                "source": schema.StringAttribute{
                                    Computed: true,
                                },
                                "icon": schema.StringAttribute{
                                    Computed: true                             },
                            },
                        },
                    },
                },
            },
        },
    }
}
