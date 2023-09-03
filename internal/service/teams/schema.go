package teams

import (
	"context"

	"github.com/catdevman/terraform-provider-clickup/internal/consts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (c *ClickUpTeamsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse){
    resp.Schema = schema.Schema{
        MarkdownDescription: "Use this data source to get the current authenticated users Teams",
        Attributes: map[string]schema.Attribute{
            "teams": schema.ListNestedAttribute{
                Computed: true,
                NestedObject: schema.NestedAttributeObject{
                    Attributes: map[string]schema.Attribute{
                        consts.IDSchemaKey: schema.StringAttribute{
                            Description: consts.IDSchemaDescription,
                            Computed: true,
                        },
                        "name": schema.StringAttribute{
                            Computed: true,
                        },
                        "color": schema.StringAttribute{
                            Computed: true,
                        },
                        "members": schema.ListNestedAttribute{
                            Computed: true,
                            NestedObject: schema.NestedAttributeObject{
                                Attributes: map[string]schema.Attribute{
                                    "user": schema.SingleNestedAttribute{
                                        Computed: true,
                                        Optional: true,
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
                                            "profile_picture": schema.StringAttribute{
                                                Computed: true,
                                            },
                                            "initials": schema.StringAttribute{
                                                Computed: true,
                                            },
                                            "role": schema.Int64Attribute{
                                                Computed: true,
                                            },
                                            "last_active": schema.StringAttribute{
                                                Computed: true,
                                            },
                                            "date_joined": schema.StringAttribute{
                                                Computed: true,
                                            },
                                            "date_invited": schema.StringAttribute{
                                                Computed: true,
                                            },
                                        },
                                    },
                                    "invited_by": schema.SingleNestedAttribute{
                                        Computed: true,
                                        Optional: true,
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
                                            "profile_picture": schema.StringAttribute{
                                                Computed: true,
                                            }, 
                                            "initials": schema.StringAttribute{
                                                Computed: true,
                                            }, 
                                        },
                                    },
                                },
                            },
                        },
                        "seats": schema.SingleNestedAttribute{
                            Computed: true,
                            Attributes: map[string]schema.Attribute{
                                "members": schema.SingleNestedAttribute{
                                    Computed: true,
                                    Attributes: map[string]schema.Attribute{
                                        "filled_members_seats": schema.Int64Attribute{
                                            Computed: true,
                                        },
                                        "total_member_seats": schema.Int64Attribute{
                                            Computed: true,
                                        },
                                        "empty_member_seats": schema.Int64Attribute{
                                            Computed: true,
                                        },
                                    },
                                },
                                "guests": schema.SingleNestedAttribute{
                                    Computed: true,
                                    Attributes: map[string]schema.Attribute{
                                        "filled_guest_seats": schema.Int64Attribute{
                                            Computed: true,
                                        },
                                        "total_guest_seats": schema.Int64Attribute{
                                            Computed: true,
                                        },
                                        "empty_guest_seats": schema.Int64Attribute{
                                            Computed: true,
                                        },
                                    },
                                },
                            },
                        },
                        "plan": schema.SingleNestedAttribute{
                            Computed: true,
                            Attributes: map[string]schema.Attribute{
                                "id": schema.Int64Attribute{
                                    Computed: true,
                                },
                                "name": schema.StringAttribute{
                                    Computed: true,
                                },
                            },
                        },
                    },
                },
            },
        },
    }
}
