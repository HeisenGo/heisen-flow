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

func (o *Ops) CreateNotification(ctx context.Context, notif *Notification) error {
	err := o.repo.CreateNotification(ctx, notif)
	return err
}

func (o *Ops) GetUserUnseenNotifications(ctx context.Context, userID uuid.UUID) ([]Notification, error) {
	notification, err := o.repo.GetUserUnseenNotifications(ctx, userID)
	if err != nil {
		return nil, err
	}
	return notification, nil
}

func (o *Ops) MarkNotificationAsSeen(ctx context.Context, notificationID uuid.UUID) (*Notification, error) {
	notif, err := o.repo.MarkNotificationAsSeen(ctx, notificationID)
	if err != nil {
		return nil, err
	}
	return notif, nil
}


func (o *Ops)GetNotificationByID(ctx context.Context, notificationID uuid.UUID)(*Notification, error){
	return o.repo.GetNotificationByID(ctx, notificationID)
}