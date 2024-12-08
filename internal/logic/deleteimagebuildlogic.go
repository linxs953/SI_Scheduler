package logic

import (
	"context"

	"Scheduler/internal/svc"
	"Scheduler/scheduler"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteImageBuildLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteImageBuildLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteImageBuildLogic {
	return &DeleteImageBuildLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteImageBuildLogic) DeleteImageBuild(in *scheduler.DeleteBuildRequest) (*scheduler.DeleteBuildResponse, error) {
	// todo: add your logic here and delete this line

	return &scheduler.DeleteBuildResponse{}, nil
}
