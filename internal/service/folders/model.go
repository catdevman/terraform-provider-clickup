package folders

import "github.com/hashicorp/terraform-plugin-framework/types"

type ClickUpFoldersDataSourceModel struct {
	SpaceId types.String                   `tfsdk:"space_id"`
	Folders []ClickUpFolderDataSourceModel `tfsdk:"folders"`
}

type ClickUpFolderDataSourceModel struct {
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	OrderIndex       types.Int64  `tfsdk:"orderindex"`
	OverrideStatuses types.Bool   `tfsdk:"override_statuses"`
	Hidden           types.Bool   `tfsdk:"hidden"`
	TaskCount        types.String `tfsdk:"task_count"`
	//Space
	//Lists
}
