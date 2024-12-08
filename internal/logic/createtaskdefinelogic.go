package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"

	"Scheduler/internal/kube/gvr"
	"Scheduler/internal/svc"
	"Scheduler/scheduler"
)

type CreateTaskDefineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateTaskDefineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTaskDefineLogic {
	return &CreateTaskDefineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateTaskDefineLogic) CreateTaskDefine(in *scheduler.CreateTaskDefineRequest) (*scheduler.CreateTaskDefineResponse, error) {
	// 获取 k8s 配置
	config, err := l.svcCtx.Config.Kubernetes.GetK8sConfig()
	if err != nil {
		return &scheduler.CreateTaskDefineResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	// 创建 dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return &scheduler.CreateTaskDefineResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	// 创建 TaskDefine 客户端
	taskDefineClient := gvr.NewTaskDefineClient(dynamicClient)
	namespace, name := in.Metadata["namespace"], in.Metadata["name"]
	if namespace == "" || name == "" {
		return &scheduler.CreateTaskDefineResponse{
			Code:    500,
			Message: "namespace or name is empty",
		}, nil
	}

	// 构建对象
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": gvr.TaskDefineGVR.Group + "/" + gvr.TaskDefineGVR.Version,
			"kind":       "TaskDefine",
			"metadata": map[string]interface{}{
				"namespace": namespace,
				"name":      name,
			},
			"spec": map[string]interface{}{
				"idlCode":    in.Spec.IdlCode,
				"idlType":    in.Spec.IdlType,
				"idlName":    in.Spec.IdlName,
				"idlVersion": in.Spec.IdlVersion,
				"relatedImage": map[string]interface{}{
					"builder":   in.Spec.RelatedImage.Builder,
					"digest":    in.Spec.RelatedImage.Digest,
					"version":   in.Spec.RelatedImage.Version,
					"namespace": in.Spec.RelatedImage.Namespace,
				},
				"definition": in.Spec.Definition,
			},
		},
	}

	// 创建 TaskDefine 对象
	result, err := taskDefineClient.Create(l.ctx, namespace, obj)
	if err != nil {
		logx.Error(err)
		return &scheduler.CreateTaskDefineResponse{
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
			taskDefine.Spec.Definition = definition
		}
	}

	return &scheduler.CreateTaskDefineResponse{
		Code:    0,
		Message: "success",
		Data:    taskDefine,
	}, nil
}
