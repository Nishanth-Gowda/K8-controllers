package main

import (
	"flag"
	"fmt"
	"time"

	"log"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
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

	informerfactory := informers.NewSharedInformerFactory(clientSet, 30 * time.Second)

	podinformer := informerfactory.Core().V1().Pods()
	podinformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(new interface{}) {
			fmt.Println("Add was called")
		},
		UpdateFunc: func(old, new interface{}) {
			fmt.Println("Update was called")
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("Delete was called")
		},
	})
	informerfactory.Start(wait.NeverStop)
	informerfactory.WaitForCacheSync(wait.NeverStop)
	pods, err := podinformer.Lister().Pods("default").Get("default")
	if err != nil {
		log.Fatalf("error %s listing pods", err.Error())
	}
	fmt.Println(pods)

}