package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Kubernetes KubernetesConfig `json:"kubernetes"`
}

type KubernetesConfig struct {
	InCluster  bool   `json:"inCluster"`
	KubeConfig string `json:"kubeConfig,omitempty"`
}
