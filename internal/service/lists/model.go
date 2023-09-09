package lists

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)


type ClickUpListsWrapperDataSourceModel struct {
    FolderId types.String `tfsdk:"folder_id"`
    Lists []ClickUpListDataSourceModel `tfsdk:"lists"`
}
type ClickUpFolderlessListsWrapperDataSourceModel struct {
    FolderId types.String `tfsdk:"folder_id"`
    Lists []ClickUpListDataSourceModel `tfsdk:"lists"`
}

type ClickUpListDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    OrderIndex types.Int64 `tfsdk:"orderindex"`
    Content types.String `tfsdk:"content"`
    Status ClickUpListStatusDataSourceModel `tfsdk:"status"`
    Priority ClickUpListPriorityDataSourceModel `tfsdk:"priority"`
    Assignee types.String `tfsdk:"assignee"`
    TaskCount types.String `tfsdk:"task_count"`
    DueDate types.String `tfsdk:"due_date"`
    StartDate types.String `tfsdk:"start_date"`
    Folder ClickUpListFolderDataSourceModel `tfsdk:"folder"`
    Space ClickUpListSpaceDataSourceModel `tfsdk:"space"`
    Archived types.Bool `tfsdk:"archived"`
    OverrideStatuses types.Bool `tfsdk:"override_statuses"`
    PermissionLevel types.String `tfsdk:"permission_level"`
}

type ClickUpListStatusDataSourceModel struct {
    Status types.String `tfsdk:"status"`
    Color types.String `tfsdk:"color"`
    HideLabel types.Bool `tfsdk:"hide_label"`
}

type ClickUpListPriorityDataSourceModel struct {
    Priority types.String `tfsdk:"priority"`
    Color types.String `tfsdk:"color"`
}

type ClickUpListFolderDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    Hidden types.Bool `tfsdk:"hidden"`
    Access types.Bool `tfsdk:"access"`
}

type ClickUpListSpaceDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    Access types.Bool `tfsdk:"access"`
}
