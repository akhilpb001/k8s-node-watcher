package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	log.Println("[INFO] Shared Informer app started")
	// kubeconfig := os.Getenv("KUBECONFIG")
	userHomeDir, err := os.UserHomeDir()

	if err != nil {
		log.Panic(fmt.Errorf("[ERROR] Error getting user home dir. Error: %s", err.Error()))
	}

	kubeconfigDefaultPath := fmt.Sprintf("%s/.kube/config", userHomeDir)
	kubeconfig := flag.String("kubeconfig", kubeconfigDefaultPath, "")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	if err != nil {
		log.Println(fmt.Errorf("[ERROR] Building config from flags. Error: %s", err.Error()))
		config, err = rest.InClusterConfig()

		if err != nil {
			// log.Panic(err.Error())
			log.Panic(fmt.Errorf("[ERROR] Error getting in cluster config. Error: %s", err.Error()))
		}
	} else {
		log.Println("[INFO] Building config from flags successful")
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		// log.Panic(err.Error())
		log.Panic(fmt.Errorf("[ERROR] Error creating new clientset. Error: %s", err.Error()))
	} else {
		log.Println("[INFO] Creating new clientset successful")
	}

	factory := informers.NewSharedInformerFactory(clientset, 0)
	informer := factory.Core().V1().Nodes().Informer()
	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAdd,
		UpdateFunc: onUpdate,
		DeleteFunc: onDelete,
	})

	go informer.Run(stopper)

	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("[ERROR] Timed out waiting for caches to sync"))
		return
	}

	<-stopper
}

func onAdd(obj interface{}) {
	// Cast the obj as node
	node := obj.(*corev1.Node)
	log.Printf("[INFO] New node added - node.Name: %v", node.Name)
}

func onUpdate(oldObj, newObj interface{}) {
	// Cast the obj as node
	node := newObj.(*corev1.Node)
	log.Printf("[INFO] Node updated - node.Name: %v", node.Name)
}

func onDelete(obj interface{}) {
	// Cast the obj as node
	node := obj.(*corev1.Node)
	log.Printf("[INFO] Node deleted - node.Name: %v", node)
}
