package comment

import "context"

type Ops struct {
	repo Repo
}

func NewOps(repo Repo) *Ops {
	return &Ops{repo}
}

func (o *Ops) Insert(ctx context.Context, comment *Comment) error {
	return nil
}
