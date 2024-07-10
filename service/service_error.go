package service

import "errors"

var (
	ErrPermissionDenied = errors.New("permission denied: cannot create task")
	ErrNotMember = errors.New("the assignee is not a member of this board")
)
