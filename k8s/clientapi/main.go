package main

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 在 .kube config 中使用当前上下文
	// path-to-kube config -- 例如 /root/.kube/config
	config, _ := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	// 创建 clientSet
	clientSet, _ := kubernetes.NewForConfig(config)
	// 访问 API 以列出 Pod
	pods, _ := clientSet.CoreV1().Pods("default").List(context.TODO(), v1.ListOptions{})
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
}
