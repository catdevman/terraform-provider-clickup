// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/catdevman/terraform-provider-clickup/internal/consts"
	"github.com/catdevman/terraform-provider-clickup/internal/service/folder"
	folderlesslists "github.com/catdevman/terraform-provider-clickup/internal/service/folderless_lists"
	"github.com/catdevman/terraform-provider-clickup/internal/service/folders"
	"github.com/catdevman/terraform-provider-clickup/internal/service/lists"
	"github.com/catdevman/terraform-provider-clickup/internal/service/space"
	"github.com/catdevman/terraform-provider-clickup/internal/service/spaces"
	"github.com/catdevman/terraform-provider-clickup/internal/service/teams"
	"github.com/catdevman/terraform-provider-clickup/internal/service/usergroups"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/raksul/go-clickup/clickup"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &ClickUpProvider{}

// ScaffoldingProvider defines the provider implementation.
type ClickUpProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ClickUpProviderModel describes the provider data model.
type ClickUpProviderModel struct {
	APIToken types.String `tfsdk:"api_token"`
}

func (p *ClickUpProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "clickup"
	resp.Version = p.version
}

func (p *ClickUpProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			consts.APITokenSchemaKey: schema.StringAttribute{
				MarkdownDescription: "ClickUp API Token - needed to talk to ClickUp API",
                Sensitive: true, 
				Required: true,
			},
		},
	}
}

func (p *ClickUpProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ClickUpProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

    client := clickup.NewClient(nil, data.APIToken.ValueString())
    _, _, err := client.Authorization.GetAuthorizedUser(ctx)
    if err != nil {
        resp.Diagnostics.Append(
            diag.NewErrorDiagnostic("Unable to create ClickUp client", "ClickUp client requires authorization to function"),
        )
        return
    }

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *ClickUpProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
	}
}

func (p *ClickUpProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
        teams.NewDataSource,
        usergroups.NewDataSource,
        spaces.NewDataSource,
        space.NewDataSource,
        folders.NewDataSource,
        folder.NewDataSource,
        lists.NewDataSource,
        folderlesslists.NewDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ClickUpProvider{
			version: version,
		}
	}
}
