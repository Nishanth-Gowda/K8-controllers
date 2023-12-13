package main

import (
	"context"
	"flag"
	"fmt"
	"log"

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
		log.Fatalf("error %s building config from flags", err.Error())
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("error %s building config from flags", err.Error())
	}

	fmt.Println("Successfully created clientset")
	pods, err := clientSet.CoreV1().Pods("default").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Fatalf("error %s building config from flags", err.Error())
	}


	if len(pods.Items) == 0 {
		fmt.Println("No pods found in default namespace.")
	} else {
		for _, pod := range pods.Items {
			// Print the name of each pod
			fmt.Println("List of pods:")
			fmt.Println("============")
			fmt.Println(pod.Name)
		}
	}

	deployments, err := clientSet.AppsV1().Deployments("default").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Fatalf("error %s building config from flags", err.Error())
	}

	if len(deployments.Items) == 0 {
		fmt.Println("No deployments found in default namespace.")
	} else {
		for _, deployment := range deployments.Items {
			// Print the name of each deployment
			fmt.Println("List of deployments:")
			fmt.Println("====================")
			fmt.Println(deployment.Name)
		}
	
	}
}