package config

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetK8sConfig 根据配置获取 Kubernetes 配置
func (c *KubernetesConfig) GetK8sConfig() (*rest.Config, error) {
	if c.InCluster {
		// 使用集群内配置
		return rest.InClusterConfig()
	}

	// 使用本地 kubeconfig
	kubeconfigPath := c.KubeConfig
	if kubeconfigPath == "" {
		// 如果未指定 kubeconfig 路径，使用默认路径
		home := os.Getenv("HOME")
		kubeconfigPath = filepath.Join(home, ".kube", "config")
	}

	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}
