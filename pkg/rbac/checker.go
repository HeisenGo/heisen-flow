/*
implementing the permission checker
This file provides utility functions for checking permissions:

HasPermission: Checks if a given role has a specific permission.
HasAllPermissions: Checks if a given role has all of the specified permissions.
HasAnyPermission: Checks if a given role has any of the specified permissions.
*/

package rbac

func HasPermission(role Role, permission Permission) bool {
	permissions, exists := RolePermissions[role]
	if !exists {
		return false
	}
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

func HasAllPermissions(role Role, permissions ...Permission) bool {
	for _, permission := range permissions {
		if !HasPermission(role, permission) {
			return false
		}
	}
	return true
}

func HasAnyPermission(role Role, permissions ...Permission) bool {
	for _, permission := range permissions {
		if HasPermission(role, permission) {
			return true
		}
	}
	return false
}
