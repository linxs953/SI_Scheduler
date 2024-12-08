package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"k8s.io/client-go/dynamic"

	"Scheduler/internal/kube/gvr"
	"Scheduler/internal/svc"
	"Scheduler/scheduler"
)

type GetImageBuildLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetImageBuildLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetImageBuildLogic {
	return &GetImageBuildLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetImageBuildLogic) GetImageBuild(in *scheduler.GetBuildRequest) (*scheduler.GetBuildResponse, error) {
	logx.Error("开始执行 获取 ImageBuild")

	// 获取 k8s 配置
	config, err := l.svcCtx.Config.Kubernetes.GetK8sConfig()
	if err != nil {
		return nil, err
	}

	// 创建 dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// 创建 TaskDefine 客户端
	taskDefineClient := gvr.NewTaskDefineClient(dynamicClient)

	// 获取 TaskDefine spec
	specJSON, err := taskDefineClient.GetSpecJSON(l.ctx, "default", "taskdefine-sample")
	if err != nil {
		return nil, err
	}

	logx.Infof("Serialized TaskDefine spec: %s", specJSON)

	return &scheduler.GetBuildResponse{}, nil
}
