package main

import (
	"context"
	"flag"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", home+"/.kube/config", "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully created clientset")
	pods, err := clientSet.CoreV1().Pods("default").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("List of pods:")
	fmt.Println("============")
	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
	}

	deployments, err := clientSet.AppsV1().Deployments("default").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("List of deployments:")
	fmt.Println("====================")
	for _, deployment := range deployments.Items {
		fmt.Println(deployment.Name)
	}
}