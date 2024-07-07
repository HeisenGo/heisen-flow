package rbac

type Role string

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
	PermissionMoveOwnTask       Permission = "move_own_task"
	PermissionCreateTask     Permission = "create_task"
	PermissionCreateSubtask  Permission = "create_subtask"
	PermissionCommentAnyTask Permission = "comment_any_task"
	PermissionMoveAnyTask       Permission = "move_any_task"
	PermissionManageColumns  Permission = "manage_columns"
	PermissionInviteUsers    Permission = "invite_users"
	PermissionRemoveBoard    Permission = "remove_board"
	// PermissionSetRole TODO
	// PermissionRemoveUser TODO
)

var RolePermissions = map[Role][]Permission{
	RoleViewer: {
		PermissionViewBoard,
		PermissionViewTask,
	},
	RoleEditor: {
		PermissionViewBoard,
		PermissionViewTask,
		PermissionCommentOwnTask,
		PermissionMoveOwnTask,
	},
	RoleMaintainer: {
		PermissionViewBoard,
		PermissionViewTask,
		PermissionCommentOwnTask,
		PermissionMoveOwnTask,
		PermissionMoveAnyTask,
		PermissionCreateTask,
		PermissionCreateSubtask,
		PermissionCommentAnyTask,
		PermissionManageColumns,
	},
	RoleOwner: {
		PermissionViewBoard,
		PermissionViewTask,
		PermissionCommentOwnTask,
		PermissionMoveOwnTask,
		PermissionMoveAnyTask,
		PermissionCreateTask,
		PermissionCreateSubtask,
		PermissionCommentAnyTask,
		PermissionManageColumns,
		PermissionInviteUsers,
	},
}
