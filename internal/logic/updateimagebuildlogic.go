package logic

import (
	"context"

	"Scheduler/internal/svc"
	"Scheduler/scheduler"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateImageBuildLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateImageBuildLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateImageBuildLogic {
	return &UpdateImageBuildLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateImageBuildLogic) UpdateImageBuild(in *scheduler.UpdateBuildRequest) (*scheduler.UpdateBuildResponse, error) {
	// todo: add your logic here and delete this line

	return &scheduler.UpdateBuildResponse{}, nil
}
