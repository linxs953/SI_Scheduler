package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"k8s.io/client-go/dynamic"

	"Scheduler/internal/kube/gvr"
	"Scheduler/internal/svc"
	"Scheduler/scheduler"
)

type GetTaskDefineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTaskDefineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskDefineLogic {
	return &GetTaskDefineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTaskDefineLogic) GetTaskDefine(in *scheduler.GetTaskDefineRequest) (*scheduler.GetTaskDefineResponse, error) {
	// 获取 k8s 配置
	config, err := l.svcCtx.Config.Kubernetes.GetK8sConfig()
	if err != nil {
		return &scheduler.GetTaskDefineResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	// 创建 dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return &scheduler.GetTaskDefineResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	// 创建 TaskDefine 客户端
	taskDefineClient := gvr.NewTaskDefineClient(dynamicClient)
	namespace, name := in.Metadata["namespace"], in.Metadata["name"]
	if namespace == "" || name == "" {
		return &scheduler.GetTaskDefineResponse{
			Code:    500,
			Message: fmt.Sprintf("namespace or name is empty, namespace: %s, name: %s", namespace, name),
		}, nil
	}

	// 获取 TaskDefine
	obj, err := taskDefineClient.Get(l.ctx, namespace, name)
	if err != nil {
		return &scheduler.GetTaskDefineResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	// 打印整个对象用于调试
	objJSON, _ := obj.MarshalJSON()
	l.Logger.Infof("Retrieved object: %s", string(objJSON))

	// 构造返回结果
	taskDefine := &scheduler.TaskDefine{
		Spec:   &scheduler.TaskDefineSpec{},
		Status: &scheduler.TaskDefineStatus{},
	}

	// 获取 spec
	if spec, ok := obj.Object["spec"].(map[string]interface{}); ok {
		taskDefine.Spec.IdlCode = getStringFromMap(spec, "idlCode")
		taskDefine.Spec.IdlType = getStringFromMap(spec, "idlType")
		taskDefine.Spec.IdlName = getStringFromMap(spec, "idlName")
		taskDefine.Spec.IdlVersion = getStringFromMap(spec, "idlVersion")

		// 获取 relatedImage
		if relatedImage, ok := spec["relatedImage"].(map[string]interface{}); ok {
			taskDefine.Spec.RelatedImage = &scheduler.TaskDefineSpec_RelatedImage{
				Builder:   getStringFromMap(relatedImage, "builder"),
				Digest:    getStringFromMap(relatedImage, "digest"),
				Version:   getStringFromMap(relatedImage, "version"),
				Namespace: getStringFromMap(relatedImage, "namespace"),
			}
		}

		// 获取 definition
		if definition, ok := spec["definition"].(string); ok {
			taskDefine.Spec.Definition = definition
		}
	}

	// 获取 status
	if status, ok := obj.Object["status"].(map[string]interface{}); ok {
		taskDefine.Status.State = getStringFromMap(status, "state")
		taskDefine.Status.Message = getStringFromMap(status, "message")
		taskDefine.Status.LastUpdated = getStringFromMap(status, "lastUpdated")
	}

	// 打印 status 用于调试
	if status, ok := obj.Object["status"].(map[string]interface{}); ok {
		statusJSON, _ := json.Marshal(status)
		l.Logger.Infof("Status map: %s", string(statusJSON))
	} else {
		l.Logger.Errorf("Failed to get status as map[string]interface{}, actual type: %T", obj.Object["status"])
	}

	statusJSON, _ := json.Marshal(taskDefine.Status)
	l.Logger.Infof("Final status: %s", string(statusJSON))

	return &scheduler.GetTaskDefineResponse{
		Code:    0,
		Message: "success",
		Data:    taskDefine,
	}, nil
}

// getStringFromMap 安全地从 map 中获取字符串值
func getStringFromMap(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	if v, ok := m[key]; ok && v != nil {
		return fmt.Sprint(v)
	}
	return ""
}

// getStringFromSpec 从 spec 中获取字符串值
func getStringFromSpec(spec map[string]interface{}, key string) string {
	if spec == nil {
		return ""
	}
	if v, ok := spec[key]; ok && v != nil {
		return fmt.Sprint(v)
	}
	return ""
}
