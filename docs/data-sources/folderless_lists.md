---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "clickup_folderless_lists Data Source - terraform-provider-clickup"
subcategory: ""
description: |-
  
---

# clickup_folderless_lists (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `space_id` (String)

### Read-Only

- `lists` (Attributes List) (see [below for nested schema](#nestedatt--lists))

<a id="nestedatt--lists"></a>
### Nested Schema for `lists`

Read-Only:

- `archived` (Boolean)
- `assignee` (String)
- `content` (String)
- `due_date` (String)
- `folder` (Attributes) (see [below for nested schema](#nestedatt--lists--folder))
- `id` (String)
- `name` (String)
- `orderindex` (Number)
- `override_statuses` (Boolean)
- `permission_level` (String)
- `priority` (Attributes) (see [below for nested schema](#nestedatt--lists--priority))
- `space` (Attributes) (see [below for nested schema](#nestedatt--lists--space))
- `start_date` (String)
- `status` (Attributes) (see [below for nested schema](#nestedatt--lists--status))
- `task_count` (String)

<a id="nestedatt--lists--folder"></a>
### Nested Schema for `lists.folder`

Read-Only:

- `access` (Boolean)
- `hidden` (Boolean)
- `id` (String)
- `name` (String)


<a id="nestedatt--lists--priority"></a>
### Nested Schema for `lists.priority`

Read-Only:

- `color` (String)
- `priority` (String)


<a id="nestedatt--lists--space"></a>
### Nested Schema for `lists.space`

Read-Only:

- `access` (Boolean)
- `id` (String)
- `name` (String)


<a id="nestedatt--lists--status"></a>
### Nested Schema for `lists.status`

Read-Only:

- `color` (String)
- `hide_label` (Boolean)
- `status` (String)
