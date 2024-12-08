package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"k8s.io/client-go/dynamic"

	"Scheduler/internal/kube/gvr"
	"Scheduler/internal/svc"
	"Scheduler/scheduler"
)

type UpdateTaskDefineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTaskDefineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTaskDefineLogic {
	return &UpdateTaskDefineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateTaskDefineLogic) UpdateTaskDefine(in *scheduler.UpdateTaskDefineRequest) (*scheduler.UpdateTaskDefineResponse, error) {
	// 获取 k8s 配置
	config, err := l.svcCtx.Config.Kubernetes.GetK8sConfig()
	if err != nil {
		return &scheduler.UpdateTaskDefineResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	// 创建 dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return &scheduler.UpdateTaskDefineResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	// 创建 TaskDefine 客户端
	taskDefineClient := gvr.NewTaskDefineClient(dynamicClient)
	namespace, name := in.Metadata["namespace"], in.Metadata["name"]
	if namespace == "" || name == "" {
		return &scheduler.UpdateTaskDefineResponse{
			Code:    500,
			Message: "namespace or name is empty",
		}, nil
	}

	// 获取现有对象
	existingObj, err := taskDefineClient.Get(l.ctx, namespace, name)
	if err != nil {
		return &scheduler.UpdateTaskDefineResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	// 更新 existingObj 的字段
	if in.Spec != nil {
		spec := existingObj.Object["spec"].(map[string]interface{})

		if in.Spec.IdlCode != "" {
			spec["idlCode"] = in.Spec.IdlCode
		}
		if in.Spec.IdlType != "" {
			spec["idlType"] = in.Spec.IdlType
		}
		if in.Spec.IdlName != "" {
			spec["idlName"] = in.Spec.IdlName
		}
		if in.Spec.IdlVersion != "" {
			spec["idlVersion"] = in.Spec.IdlVersion
		}
		if in.Spec.RelatedImage != nil {
			relatedImage := make(map[string]interface{})
			if in.Spec.RelatedImage.Builder != "" {
				relatedImage["builder"] = in.Spec.RelatedImage.Builder
			}
			if in.Spec.RelatedImage.Digest != "" {
				relatedImage["digest"] = in.Spec.RelatedImage.Digest
			}
			if in.Spec.RelatedImage.Version != "" {
				relatedImage["version"] = in.Spec.RelatedImage.Version
			}
			if in.Spec.RelatedImage.Namespace != "" {
				relatedImage["namespace"] = in.Spec.RelatedImage.Namespace
			}
			spec["relatedImage"] = relatedImage
		}
		existingObj.Object["spec"] = spec
	}

	// 准备更新对象
	obj := existingObj

	// 更新 TaskDefine 对象
	result, err := taskDefineClient.UpdateByName(l.ctx, namespace, name, obj)
	if err != nil {
		return &scheduler.UpdateTaskDefineResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	// 从返回的对象中提取数据
	taskDefine := &scheduler.TaskDefine{
		Spec: &scheduler.TaskDefineSpec{},
	}

	if result.Object["spec"] != nil {
		spec := result.Object["spec"].(map[string]interface{})
		if idlCode, ok := spec["idlCode"].(string); ok {
			taskDefine.Spec.IdlCode = idlCode
		}
		if idlType, ok := spec["idlType"].(string); ok {
			taskDefine.Spec.IdlType = idlType
		}
		if idlName, ok := spec["idlName"].(string); ok {
			taskDefine.Spec.IdlName = idlName
		}
		if idlVersion, ok := spec["idlVersion"].(string); ok {
			taskDefine.Spec.IdlVersion = idlVersion
		}
		if relatedImage, ok := spec["relatedImage"].(map[string]interface{}); ok {
			taskDefine.Spec.RelatedImage = &scheduler.TaskDefineSpec_RelatedImage{}
			if builder, ok := relatedImage["builder"].(string); ok {
				taskDefine.Spec.RelatedImage.Builder = builder
			}
			if digest, ok := relatedImage["digest"].(string); ok {
				taskDefine.Spec.RelatedImage.Digest = digest
			}
			if version, ok := relatedImage["version"].(string); ok {
				taskDefine.Spec.RelatedImage.Version = version
			}
			if namespace, ok := relatedImage["namespace"].(string); ok {
				taskDefine.Spec.RelatedImage.Namespace = namespace
			}
		}

		// 处理 definition
		if definition, ok := spec["definition"].(string); ok {
			// 初始化 definition map
			taskDefine.Spec.Definition = definition

		}
	}

	return &scheduler.UpdateTaskDefineResponse{
		Code:    200,
		Message: "success",
		Data:    taskDefine,
	}, nil
}
