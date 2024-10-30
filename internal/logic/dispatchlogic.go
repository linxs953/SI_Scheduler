package logic

import (
	"context"

	"scheduler/internal/svc"
	"scheduler/pb/scheduler"

	"github.com/zeromicro/go-zero/core/logx"
)

type DispatchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDispatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DispatchLogic {
	return &DispatchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DispatchLogic) Dispatch(in *scheduler.CreateJobRequest) (*scheduler.CreateJobResponse, error) {
	// todo: add your logic here and delete this line

	return &scheduler.CreateJobResponse{}, nil
}
