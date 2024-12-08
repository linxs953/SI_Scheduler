package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"k8s.io/client-go/dynamic"

	"Scheduler/internal/kube/gvr"
	"Scheduler/internal/svc"
	"Scheduler/scheduler"
)

type ListTaskDefinesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListTaskDefinesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTaskDefinesLogic {
	return &ListTaskDefinesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListTaskDefinesLogic) ListTaskDefines(in *scheduler.ListTaskDefinesRequest) (*scheduler.ListTaskDefinesResponse, error) {

	// 获取 k8s 配置
	config, err := l.svcCtx.Config.Kubernetes.GetK8sConfig()
	if err != nil {
		return &scheduler.ListTaskDefinesResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	// 创建 dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return &scheduler.ListTaskDefinesResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	// 创建 TaskDefine 客户端
	taskDefineClient := gvr.NewTaskDefineClient(dynamicClient)

	namespace := in.Metadata["namespace"]
	if namespace == "" {
		return &scheduler.ListTaskDefinesResponse{
			Code:    500,
			Message: "namespace is empty",
		}, nil
	}

	// 获取 TaskDefine 列表
	taskDefines, err := taskDefineClient.List(l.ctx, namespace)
	if err != nil {
		return &scheduler.ListTaskDefinesResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	logx.Infof("TaskDefine list: %v", taskDefines)

	// 将 TaskDefine 转换为 TaskDefine 列表
	tas := make([]*scheduler.TaskDefine, 0, len(taskDefines))
	for _, taskDefine := range taskDefines {
		spec := taskDefine.Object["spec"].(map[string]interface{})
		status := taskDefine.Object["status"].(map[string]interface{})

		taskDefineSpec := &scheduler.TaskDefineSpec{
			IdlName:    getStringFromMap(spec, "idlName"),
			IdlCode:    getStringFromMap(spec, "idlCode"),
			IdlType:    getStringFromMap(spec, "idlType"),
			IdlVersion: getStringFromMap(spec, "idlVersion"),
		}

		if relatedImage, ok := spec["relatedImage"].(map[string]interface{}); ok {
			taskDefineSpec.RelatedImage = &scheduler.TaskDefineSpec_RelatedImage{
				Builder:   getStringFromMap(relatedImage, "builder"),
				Digest:    getStringFromMap(relatedImage, "digest"),
				Version:   getStringFromMap(relatedImage, "version"),
				Namespace: getStringFromMap(relatedImage, "namespace"),
			}
		}

		taskDefineStatus := &scheduler.TaskDefineStatus{
			State:       getStringFromMap(status, "state"),
			Message:     getStringFromMap(status, "message"),
			LastUpdated: getStringFromMap(status, "lastUpdated"),
		}

		tas = append(tas, &scheduler.TaskDefine{
			Spec:   taskDefineSpec,
			Status: taskDefineStatus,
		})
	}

	// 返回 TaskDefine 列表
	return &scheduler.ListTaskDefinesResponse{
		Code:    200,
		Message: "success",
		Data:    tas,
	}, nil
}
