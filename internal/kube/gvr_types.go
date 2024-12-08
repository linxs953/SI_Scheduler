package kube

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
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
