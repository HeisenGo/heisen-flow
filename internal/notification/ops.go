package notification

import (
	"context"
	"github.com/google/uuid"
)

type Ops struct {
	repo Repo
}

func NewOps(repo Repo) *Ops {
	return &Ops{repo}
}

func (o *Ops) CreateNotification(ctx context.Context, userID, boardID uuid.UUID) error {
	err := o.repo.CreateNotification(ctx, userID,boardID)
	return err
}

func (o *Ops) DisplyNotification(ctx context.Context, userID, boardID uuid.UUID) (*Notification,error){
	notification , err := o.repo.DisplyNotification(ctx,userID,boardID)
	if err != nil{
		return nil , err
	}
	return notification, nil
}

func (o *Ops) DeleteNotification(ctx context.Context, notif *Notification) error{
	err := o.repo.DeleteNotification(ctx,notif)
	if err != nil {
		return err
	}
	return nil
}