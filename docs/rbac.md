# Role-Based Access Control (RBAC)

[Role-Based Access Control (RBAC)](https://en.wikipedia.org/wiki/Role-based_access_control) is a method of regulating access to resources based on the roles of individual users within an organization. In RBAC, permissions are associated with roles, and users are assigned to appropriate roles, thereby acquiring the permissions of those roles.

## Roles and Permissions

### Roles

- **Viewer**: Can only view boards, their info, and tasks.
- **Editor**: Can comment on their own tasks and move them between columns (e.g., from "in progress" to "done").
- **Maintainer**: Has editors permissions and also Can create tasks and subtasks, comment on them, change their columns, create new columns, remove a column, or reorder them.
- **Owner**: Has full control over the board. Can do everything a maintainer can and also invite people to the board and specify their roles.

## Package Structure

### RBAC Definitions

#### File: `pkg/rbac/rbac.go`

This file defines the basic structure of our RBAC system. It includes:

- **Permission type**: Represents individual permissions in the system.
- **Role type**: Represents user roles.
- **RolePermissions map**: Maps each role to its allowed permissions.

This structure allows for easy extension of roles and permissions in the future. New permissions or roles by simply can be added adding new constants and updating the RolePermissions map.

```go
const (
    RoleViewer     Role = "viewer"
    RoleEditor     Role = "editor"
    RoleMaintainer Role = "maintainer"
    RoleOwner      Role = "owner"
)

type Permission string

const (
    PermissionViewBoard      Permission = "view_board"
    PermissionViewTask       Permission = "view_task"
    PermissionCommentOwnTask Permission = "comment_own_task"
    PermissionMoveOwnTask    Permission = "move_own_task"
    PermissionCreateTask     Permission = "create_task"
    PermissionCreateSubtask  Permission = "create_subtask"
    PermissionCommentAnyTask Permission = "comment_any_task"
    PermissionMoveAnyTask    Permission = "move_any_task"
    PermissionManageColumns  Permission = "manage_columns"
    PermissionInviteUsers    Permission = "invite_users"
    PermissionRemoveBoard    Permission = "remove_board"
)
```

## Permission Checking
**File: `pkg/rbac/checker.go`**
This file provides utility functions for checking permissions:

- **HasPermission**: Checks if a given role has a specific permission.
- **HasAllPermissions**: Checks if a given role has all of the specified permissions.
- **HasAnyPermission**: Checks if a given role has any of the specified permissions.
- **IsAPossibleRole**: Check if a given role string is one of the defined roles in RBAC system.

These functions will be useful when implementing permission checks throughout application. They provide flexibility in how to check permissions - whether to check for a single permission, all of a set of permissions, or any of a set of permissions.

## User-Board Role Management
- **File: `pkg/adapters/storage/entities/user_board_role.go`**:
 ***UserBoardRole struct***: Represents the relationship between a user, a board, and the user's role on that board.
This UserBoardRole entity includes:

- A unique ID using UUID.
- UserID and BoardID to link users to boards.
- A Role field to store the user's role for the specific board.
- Timestamps for creation, update, and soft delete.
- Relationships to the User and Board entities.

This structure allows for:

- Assigning different roles to users for different boards.
- Easily querying user roles for specific boards.
- Maintaining the relationship between users and boards with associated roles.

- **File: `pkg/adapters/storage/user_board_role.go`**
This file provides methods for managing user-board roles in the storage layer:

- **GetUserBoardRole**: Retrieves the role of a user for a specific board.
- **SetUserBoardRole**: Sets the role of a user for a specific board.
- **RemoveUserBoardRole**: Removes the role of a user for a specific board.
- **GetUserBoardRoleObj**: Retrieves the user-board-role record of a user for a specific board.

## Permission Checks in Service Layer
When RBAC is needed in a service, each method first retrieves the user's role for the specific board and then checks if that role has the required permission for the action. If the permission check fails, it returns an error. If the check passes, it proceeds with the action.

Example: Invite User to Board
Only the owner can invite users to the board.

```go
inviterRole, err := s.userBoardRoleOps.GetUserBoardRole(ctx, inviterID, userBoardRole.BoardID)
if err != nil {
    return ErrPermissionDeniedToInvite
}

if !rbac.HasPermission(inviterRole, rbac.PermissionInviteUsers) {
    return ErrPermissionDeniedToInvite
}
```

## Advantages of checking permissions in the service layer:

- **Separation of concerns**: The service layer is responsible for business logic, which includes enforcing access control rules. By keeping permission checks in the service layer, you maintain a clear separation between HTTP-specific concerns (handled in handlers) and core business logic.
- **Reusability**: Your service methods might be called from different entry points (e.g., HTTP handlers, gRPC handlers, or even internal processes). By including permission checks in the service layer, you ensure that these rules are enforced regardless of how the service is invoked.
- **Granularity**: Some operations might require complex permission checks that depend on the specific data being accessed or modified. The service layer has access to this context, allowing for more fine-grained access control.
- **Testing**: It's easier to unit test permission logic when it's part of the service layer, as you can test it independently of the HTTP layer.
- **Security**: Implementing permission checks at the service layer provides an additional layer of security. Even if a middleware fails or is accidentally omitted, the core business logic still enforces access control.

