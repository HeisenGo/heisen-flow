package service

import "errors"

var (
	ErrPermissionDenied = errors.New("permission denied")
	ErrNotMember        = errors.New("the assignee is not a member of this board")
	ErrCantAssigned     = errors.New("assignee is a viewer")
)
