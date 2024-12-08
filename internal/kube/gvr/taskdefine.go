package gvr

import (
	"context"
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

const (
	// TaskDefine GVR constants
	TaskDefineGroup    = "lct.kube.inspect"
	TaskDefineVersion  = "v1"
	TaskDefineResource = "taskdefines"
)

var (
	// TaskDefineGVR is the GroupVersionResource for TaskDefine
	TaskDefineGVR = schema.GroupVersionResource{
		Group:    TaskDefineGroup,
		Version:  TaskDefineVersion,
		Resource: TaskDefineResource,
	}
)

// TaskDefineClient 封装 TaskDefine 资源操作的客户端
type TaskDefineClient struct {
	dynamicClient dynamic.Interface
}

// NewTaskDefineClient 创建新的 TaskDefine 客户端
func NewTaskDefineClient(dynamicClient dynamic.Interface) *TaskDefineClient {
	return &TaskDefineClient{
		dynamicClient: dynamicClient,
	}
}

// Get 获取指定的 TaskDefine 资源
func (c *TaskDefineClient) Get(ctx context.Context, namespace, name string) (*unstructured.Unstructured, error) {
	return c.dynamicClient.Resource(TaskDefineGVR).
		Namespace(namespace).
		Get(ctx, name, metav1.GetOptions{})
}

func (c *TaskDefineClient) List(ctx context.Context, namespace string) ([]unstructured.Unstructured, error) {
	list, err := c.dynamicClient.Resource(TaskDefineGVR).
		Namespace(namespace).
		List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	items := make([]unstructured.Unstructured, 0, len(list.Items))
	items = append(items, list.Items...)
	return items, nil
}

// GetSpec 获取 TaskDefine 的 spec
func (c *TaskDefineClient) GetSpec(ctx context.Context, namespace, name string) (map[string]interface{}, error) {
	obj, err := c.Get(ctx, namespace, name)
	if err != nil {
		return nil, err
	}

	spec, ok := obj.Object["spec"].(map[string]interface{})
	if !ok {
		return nil, nil
	}

	return spec, nil
}

// GetStatus 获取 TaskDefine 的 status
func (c *TaskDefineClient) GetStatus(ctx context.Context, namespace, name string) (map[string]interface{}, error) {
	obj, err := c.Get(ctx, namespace, name)
	if err != nil {
		return nil, err
	}

	status, ok := obj.Object["status"].(map[string]interface{})
	if !ok {
		return nil, nil
	}

	return status, nil
}

// GetSpecJSON 获取 TaskDefine 的 spec 并序列化为 JSON
func (c *TaskDefineClient) GetSpecJSON(ctx context.Context, namespace, name string) (string, error) {
	spec, err := c.GetSpec(ctx, namespace, name)
	if err != nil {
		return "", err
	}

	specJSON, err := json.Marshal(spec)
	if err != nil {
		return "", err
	}

	return string(specJSON), nil
}

// GetStatusJSON 获取 TaskDefine 的 status 并序列化为 JSON
func (c *TaskDefineClient) GetStatusJSON(ctx context.Context, namespace, name string) (string, error) {
	status, err := c.GetStatus(ctx, namespace, name)
	if err != nil {
		return "", err
	}

	statusJSON, err := json.Marshal(status)
	if err != nil {
		return "", err
	}

	return string(statusJSON), nil
}

// Create 创建 TaskDefine 资源
func (c *TaskDefineClient) Create(ctx context.Context, namespace string, obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	return c.dynamicClient.Resource(TaskDefineGVR).
		Namespace(namespace).
		Create(ctx, obj, metav1.CreateOptions{})
}

// CreateFromSpec 从 spec 创建 TaskDefine 资源
func (c *TaskDefineClient) CreateFromSpec(ctx context.Context, namespace, name string, spec map[string]interface{}) (*unstructured.Unstructured, error) {
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": TaskDefineGroup + "/" + TaskDefineVersion,
			"kind":       "TaskDefine",
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": namespace,
			},
			"spec": spec,
		},
	}
	return c.Create(ctx, namespace, obj)
}

// Update 更新 TaskDefine 资源
func (c *TaskDefineClient) Update(ctx context.Context, namespace string, obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	return c.dynamicClient.Resource(TaskDefineGVR).
		Namespace(namespace).
		Update(ctx, obj, metav1.UpdateOptions{})
}

// UpdateByName 通过 name 更新 TaskDefine 资源
func (c *TaskDefineClient) UpdateByName(ctx context.Context, namespace, name string, obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	// 确保 name 匹配
	obj.SetName(name)
	return c.Update(ctx, namespace, obj)
}

// UpdateFromSpec 从 spec 更新 TaskDefine 资源
func (c *TaskDefineClient) UpdateFromSpec(ctx context.Context, namespace, name string, spec map[string]interface{}) (*unstructured.Unstructured, error) {
	// 首先获取现有对象
	existing, err := c.Get(ctx, namespace, name)
	if err != nil {
		return nil, err
	}

	// 更新 spec
	existing.Object["spec"] = spec

	return c.Update(ctx, namespace, existing)
}

// Delete 删除 TaskDefine 资源
func (c *TaskDefineClient) Delete(ctx context.Context, namespace, name string) error {
	return c.dynamicClient.Resource(TaskDefineGVR).
		Namespace(namespace).
		Delete(ctx, name, metav1.DeleteOptions{})
}

// ResourceClient 封装资源操作的客户端
type ResourceClient struct {
	dynamicClient dynamic.Interface
}

// GetTaskDefine 获取指定的 TaskDefine 资源
func (r *ResourceClient) GetTaskDefine(ctx context.Context, namespace, name string) (*unstructured.Unstructured, error) {
	return r.dynamicClient.Resource(TaskDefineGVR).
		Namespace(namespace).
		Get(ctx, name, metav1.GetOptions{})
}

// GetTaskDefineSpec 获取 TaskDefine 的 spec
func (r *ResourceClient) GetTaskDefineSpec(ctx context.Context, namespace, name string) (map[string]interface{}, error) {
	obj, err := r.GetTaskDefine(ctx, namespace, name)
	if err != nil {
		return nil, err
	}

	spec, ok := obj.Object["spec"].(map[string]interface{})
	if !ok {
		return nil, nil
	}

	return spec, nil
}

// GetTaskDefineStatus 获取 TaskDefine 的 status
func (r *ResourceClient) GetTaskDefineStatus(ctx context.Context, namespace, name string) (map[string]interface{}, error) {
	obj, err := r.GetTaskDefine(ctx, namespace, name)
	if err != nil {
		return nil, err
	}

	status, ok := obj.Object["status"].(map[string]interface{})
	if !ok {
		return nil, nil
	}

	return status, nil
}

// GetTaskDefineSpecJSON 获取 TaskDefine 的 spec 并序列化为 JSON
func (r *ResourceClient) GetTaskDefineSpecJSON(ctx context.Context, namespace, name string) (string, error) {
	spec, err := r.GetTaskDefineSpec(ctx, namespace, name)
	if err != nil {
		return "", err
	}

	specJSON, err := json.Marshal(spec)
	if err != nil {
		return "", err
	}

	return string(specJSON), nil
}

// GetTaskDefineStatusJSON 获取 TaskDefine 的 status 并序列化为 JSON
func (r *ResourceClient) GetTaskDefineStatusJSON(ctx context.Context, namespace, name string) (string, error) {
	status, err := r.GetTaskDefineStatus(ctx, namespace, name)
	if err != nil {
		return "", err
	}

	statusJSON, err := json.Marshal(status)
	if err != nil {
		return "", err
	}

	return string(statusJSON), nil
}
