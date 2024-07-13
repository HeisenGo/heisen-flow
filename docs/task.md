## Task
`pkg/adapters/storage/task.go`: Represents task entity in the database

This Task entity includes:

- Basic fields like ID, Title, Description, and Status.
- Timestamps for creation, update, and soft delete.
- A relationship to the Board it belongs to.
- A relationship to the Colum for the column it is in to
- A self-referential relationship for subtasks:

- ParentID and Parent for the parent task (null for top-level tasks).
- Subtasks for child tasks.


- Many-to-many relationships for task dependencies:

$\quad$ DependsOn for tasks that this task depends on.
$\quad$ DependentBy for tasks that depend on this task.


A separate TaskDependency struct to represent the many-to-many relationship in the database.

This structure allows for:

- Nested subtasks to any depth.
- Task dependencies between any tasks, regardless of their level in the hierarchy.
Proper tracking of task status.

## Storage

#### Methods:

- CheckCircularDependency: This function uses a depth-first search algorithm to detect if adding a new dependency would create a circular dependency.