package logic

import (
	"context"

	"Scheduler/internal/svc"
	"Scheduler/scheduler"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListImageBuildsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListImageBuildsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListImageBuildsLogic {
	return &ListImageBuildsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListImageBuildsLogic) ListImageBuilds(in *scheduler.ListBuildsRequest) (*scheduler.ListBuildsResponse, error) {
	// todo: add your logic here and delete this line

	return &scheduler.ListBuildsResponse{}, nil
}
