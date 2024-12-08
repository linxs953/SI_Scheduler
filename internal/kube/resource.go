package kube

import (
	"k8s.io/client-go/dynamic"
)

// ResourceClient 封装资源操作的客户端
type ResourceClient struct {
	dynamicClient dynamic.Interface
}

// NewResourceClient 创建新的资源客户端
func NewResourceClient(dynamicClient dynamic.Interface) *ResourceClient {
	return &ResourceClient{
		dynamicClient: dynamicClient,
	}
}
