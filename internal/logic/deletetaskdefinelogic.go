package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"k8s.io/client-go/dynamic"

	"Scheduler/internal/kube/gvr"
	"Scheduler/internal/svc"
	"Scheduler/scheduler"
)

type DeleteTaskDefineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteTaskDefineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTaskDefineLogic {
	return &DeleteTaskDefineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteTaskDefineLogic) DeleteTaskDefine(in *scheduler.DeleteTaskDefineRequest) (*scheduler.DeleteTaskDefineResponse, error) {
	// 获取 k8s 配置
	config, err := l.svcCtx.Config.Kubernetes.GetK8sConfig()
	if err != nil {
		return &scheduler.DeleteTaskDefineResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	// 创建 dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return &scheduler.DeleteTaskDefineResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	// 创建 TaskDefine 客户端
	taskDefineClient := gvr.NewTaskDefineClient(dynamicClient)
	namespace, name := in.Metadata["namespace"], in.Metadata["name"]
	if namespace == "" || name == "" {
		return &scheduler.DeleteTaskDefineResponse{
			Code:    500,
			Message: "namespace or name is empty",
		}, nil
	}

	// 检查资源是否存在
	_, err = taskDefineClient.Get(l.ctx, namespace, name)
	if err != nil {
		return &scheduler.DeleteTaskDefineResponse{
			Code:    404,
			Message: "TaskDefine not found",
		}, nil
	}
	err = taskDefineClient.Delete(l.ctx, namespace, name)
	if err != nil {
		return &scheduler.DeleteTaskDefineResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	return &scheduler.DeleteTaskDefineResponse{
		Code:    200,
		Message: "success",
	}, nil
}
