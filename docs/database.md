# Task

### Contents

1. [Composite Types](#composite-types)
    - [task_creation_t](#task_creation_t)
    - [task_update_t](#task_update_t)

2. [Views](#views)
    - [completed_tasks](#completed_tasks)
    - [archived_tasks](#archived_tasks)
    - [trashed_tasks](#trashed_tasks)

3. [Routines](#routines)
    - [Procedures](#procedures)
        - [assert_task_exists](#assert_task_exists)
    - [Functions](#functions)
      - [make_task](#make_task)
      - [duplicate_task](#duplicate_task)
      - [fetch_task_by_id](#fetch_task_by_id)
      - [fetch_tasks](#fetch_tasks)
      - [fetch_tasks_from_today_list](#fetch_tasks_from_today_list)
      - [fetch_tasks_from_tomorrow_list](#fetch_tasks_from_tomorrow_list)
      - [fetch_tasks_from_deferred_list](#fetch_tasks_from_deferred_list)
      - [update_task](#update_task)
      - [reorder_task_in_list](#reorder_task_in_list)
      - [set_task_reminder_date](#set_task_reminder_date)
      - [set_task_priority](#set_task_priority)
      - [set_task_due_date](#set_task_due_date)
      - [set_task_as_completed](#set_task_as_completed)
      - [set_task_as_uncompleted](#set_task_as_uncompleted)
      - [pin_task](#pin_task)
      - [defer_tasks_in_today_list](#defer_tasks_in_today_list)
      - [move_tasks_from_tomorrow_to_today_list](#move_tasks_from_tomorrow_to_today_list)
      - [move_task_from_list](#move_task_from_list)
      - [move_task_to_today_list](#move_task_to_today_list)
      - [move_task_to_tomorrow_list](#move_task_to_tomorrow_list)
      - [move_task_to_deferred_list](#move_task_to_deferred_list)
      - [trash_task](#trash_task)
      - [restore_task_from_trash](#restore_task_from_trash)
      - [delete_task](#make_task)

## Composite Types

### `task_creation_t`

This composite type represents the specifications for creating a new task.

**Fields:**
1. `title` (VARCHAR(100)): The title of the task.
2. `headline` (VARCHAR): A headline for the task (optional).
3. `description` (TEXT): The detailed description of the task (optional).
4. `priority` (task_priority_t): The priority level of the task.
5. `due_date` (TIMESTAMP): The due date of the task (optional).
6. `remind_at` (TIMESTAMP): The date to set a reminder for the task (optional).

### `task_update_t`

This composite type represents the specifications for updating a task.

**Fields:**
1. `title` (VARCHAR(100)): The new title of the task (optional).
2. `headline` (VARCHAR): The new headline for the task (optional).
3. `description` (TEXT): The new detailed description of the task (optional).

## Views

### `completed_tasks`

This view displays a list of tasks that have been marked as completed.

### `archived_tasks`

This view displays a list of archived tasks.

### `trashed_tasks`

This view displays a list of tasks that have been moved to the trash.

## Routines

### Procedures

### `assert_task_exists`

Asserts the existence of a task that belongs to a user.

**Parameters**

1. `IN owner_id UUID`: The owner of the task.
2. `IN list_uuid UUID`: The list the task belongs to.
3. `IN task_id UUID`: The task to assert existence.

### Functions

### `make_task`

Creates a new task that'll belong to `list_uuid`. `list_uuid` can be the ID of any user-defined or special list except for deferred list, in that case it throws a list not found exception. If `list_uuid` is `NULL` the task will belong to the today list of `owner_id`.

**Parameters**

1. `IN owner_id UUID`: The owner of the task.
2. `IN list_uuid? UUID`: The list the task belongs to.
3. `IN creation task_creation_t`: The creation object.

**Returns** `UUID`

The ID of the newly created task.

### `duplicate_task`

Duplicates a task and its content, into the same list.

**Parameters**

1. `IN owner_id UUID`: The owner of the task.
2. `IN task_id UUID`: The task to duplicate.

**Returns** `UUID`

The ID of the replica task.

### `fetch_task_by_id`

Retrieves a task by its ID.

**Parameters**

1. `IN owner_id UUID`: The owner of the task.
2. `IN list_uuid UUID`: The list the task belongs to.
3. `IN task_id UUID`: The task to duplicate.

**Returns** `SETOF "task"`

### `fetch_tasks`

Retrieves the tasks that belong to `list_uuid`.

**Parameters**

1. `IN owner_id UUID`: The owner of the task.
2. `IN list_uuid UUID`: The list the task belongs to.

**Returns** `SETOF "task"`

### `fetch_tasks_from_today_list`

Retrieves the tasks that belong to the today list of `owner_id`.

**Parameters**

1. `IN owner_id UUID`: The owner of the tasks.

**Returns** `SETOF "task"`

### `fetch_tasks_from_tomorrow_list`

Retrieves the tasks that belong to the tomorrow list of `owner_id`.

**Parameters**

1. `IN owner_id UUID`: The owner of the tasks.

**Returns** `SETOF "task"`

### `fetch_tasks_from_deferred_list`

Retrieves the tasks that belong to the deferred list of `owner_id`.

**Parameters**

1. `IN owner_id UUID`: The owner of the tasks.

**Returns** `SETOF "task"`

### `update_task`

Updates the title, headline or description of a task individually or all of them at once. For updating other fields, see the next functions.

**Parameters**

1. `IN owner_id UUID`: The owner of the task.
2. `IN list_uuid UUID`: The list the task belongs to.
3. `IN update task_update_t`: The update object.

**Returns** `BOOLEAN`

True if any of the task fields was changed. False otherwise.

### `reorder_task_in_list`

Changes the order of the task inside the list.

**Parameters**

1. `IN owner_id UUID`: The owner of the task.
2. `IN list_uuid UUID`: The list the task belongs to.
3. `IN new_position pos_t`: The new position of the task.

**Returns** `BOOLEAN`

### `set_task_reminder_date`

Sets the date to remind for this task.

**Parameters**

1. `IN owner_id UUID`: The owner of the task.
2. `IN list_uuid UUID`: The list the task belongs to.
3. `IN remind_at TIMESTAMPTZ`: The date to remind.

**Returns** `BOOLEAN`

### `set_task_priority`

Sets the priority of the task.

**Parameters**

1. `IN owner_id UUID`: The owner of the task.
2. `IN list_uuid UUID`: The list the task belongs to.
3. `IN new_priority task_priority_t`: The new priority.

**Returns** `BOOLEAN`

### `set_task_due_date`

Sets the due date of the task.

**Parameters**

1. `IN owner_id UUID`: The owner of the task.
2. `IN list_uuid UUID`: The list the task belongs to.
3. `IN due_date TIMESTAMPTZ`: The due date.

**Returns** `BOOLEAN`

### `set_task_as_completed`

Sets the status of the task as completed.

**Parameters**

1. `IN owner_id UUID`: The owner of the task.
2. `IN list_uuid UUID`: The list the task belongs to.

**Returns** `BOOLEAN`

### `set_task_as_uncompleted`

Sets the status of the task as uncompleted.

**Parameters**

1. `IN owner_id UUID`: The owner of the task.
2. `IN list_uuid UUID`: The list the task belongs to.

**Returns** `BOOLEAN`

### `pin_task`

Pins task. Pinned tasks should always be retrieved first.

**Parameters**

1. `IN owner_id UUID`: The owner of the task.
2. `IN list_uuid UUID`: The list the task belongs to.
3. `IN task_id UUID`: The task to pin.

**Returns** `BOOLEAN`

### `defer_tasks_in_today_list`

Moves all the tasks inside the today list to the deferred list. 

**Parameters**

1. `IN owner_id UUID`: The owner of the lists.

**Returns** `BOOLEAN`

### `move_tasks_from_tomorrow_to_today_list`

Moves all the tasks inside the tomorrow list to the today list.

**Parameters**

1. `IN owner_id UUID`: The owner of the lists.

**Returns** `BOOLEAN`

### `move_task_from_list`

Moves a task from one list to another one.

**Parameters**

1. `IN owner_id UUID`: The owner of the task.
2. `IN task_id UUID`: The task to move.
3. `IN dst_list UUID`: The list to move the task to.

**Returns** `BOOLEAN`

### `move_task_to_today_list`

Moves a task to the today list.

**Parameters**

1. `IN owner_id UUID`: The owner of the task.
2. `IN task_id UUID`: The task to move.

**Returns** `BOOLEAN`

### `move_task_to_tomorrow_list`

Moves a task to the tomorrow list.

**Parameters**

1. `IN owner_id UUID`: The owner of the task.
2. `IN task_id UUID`: The task to move.

**Returns** `BOOLEAN`

### `move_task_to_deferred_list`

Moves a task to the deferred list.

**Parameters**

1. `IN owner_id UUID`: The owner of the task.
2. `IN task_id UUID`: The task to move.

**Returns** `BOOLEAN`

### `trash_task`

Throws a task to the trash. Moves a task to the `trashed_task` table.

**Parameters**

1. `IN owner_id UUID`: The owner of the task.
2. `IN list_uuid UUID`: The list the task belongs to.
3. `IN task_id UUID`: The task to trash.

**Returns** `BOOLEAN`

### `restore_task_from_trash`

Recovers a task from trash. Moves a task back to the `task` table.

**Parameters**

1. `IN owner_id UUID`: The owner of the task.
2. `IN list_uuid UUID`: The list the task belongs to.
3. `IN task_id UUID`: The task to recover from trash.

**Returns** `BOOLEAN`

### `delete_task`

Permanently deletes a task.

**Parameters**

1. `IN owner_id UUID`: The owner of the task.
2. `IN list_uuid UUID`: The list the task belongs to.
3. `IN task_id UUID`: The task to delete.

**Returns** `BOOLEAN`
