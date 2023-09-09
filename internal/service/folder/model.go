package folder

import "github.com/hashicorp/terraform-plugin-framework/types"


type ClickUpFolderWrapperDataSourceModel struct {
    FolderId types.String `tfsdk:"folder_id"`
    Folder *ClickUpFolderDataSourceModel `tfsdk:"folder"`
}

type ClickUpFolderDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    OrderIndex types.Int64 `tfsdk:"orderindex"`
    OverrideStatuses types.Bool `tfsdk:"override_statuses"`
    Hidden types.Bool `tfsdk:"hidden"`
    TaskCount types.String `tfsdk:"task_count"`
    //Space
    //Lists
}
