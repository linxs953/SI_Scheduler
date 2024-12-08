package logic

import (
	"context"

	"Scheduler/internal/svc"
	"Scheduler/scheduler"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateImageBuildLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateImageBuildLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateImageBuildLogic {
	return &CreateImageBuildLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ImageBuild 相关操作
func (l *CreateImageBuildLogic) CreateImageBuild(in *scheduler.CreateImageBuildRequest) (*scheduler.CreateImageBuildResponse, error) {
	// todo: add your logic here and delete this line

	return &scheduler.CreateImageBuildResponse{}, nil
}
