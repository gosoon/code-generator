package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/gosoon/code-generator/_examples/server"
	ctrl "github.com/gosoon/code-generator/_examples/server/controller"

	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

var (
	kubeconfig string
	masterURL  string
)

func main() {
	defer runtime.HandleCrash()

	var cfg *rest.Config
	var err error
	if kubeconfig == "" {
		cfg, err = rest.InClusterConfig()
	} else {
		cfg, err = clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	}
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	go func() {
		opt := &ctrl.Options{KubeClientset: kubeClient}
		server := server.New(server.Options{CtrlOptions: opt, ListenAddr: ":8080"})
		if err := server.ListenAndServe(); err != nil {
			klog.Fatalf("Failed to listen and serve admission webhook server: %v", err)
		}
	}()

	// listening OS shutdown singal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	klog.Infof("Got OS shutdown signal, shutting down webhook server gracefully...")
}

func init() {
	flag.StringVar(&kubeconfig, "config", "", "config file")
	flag.Parse()
}
