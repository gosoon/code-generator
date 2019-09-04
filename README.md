# code-generator
Generate RESTful APIs based on the golang type.

## Getting started

1、config the Golang development environment

2、compile binaries using  `go install`

```
$ go install github.com/gosoon/code-generator/cmd/restfulapi-gen
```

3、define types(the types file path must be under GOPATH),example：github.com/gosoon/code-generator/_examples/types/v1

```
package types

// +genclient

// Namespace xxx
type Namespace struct {
	Name string
}
```

Each type need add the tag for the generated API in comments, current only support use "+genclient" tag, this tag can generator CRUD code for the type. You can define multiple types, each with a corresponding tag.



4、execute the command to generate the code

```
$ restfulapi-gen --input-dirs github.com/gosoon/code-generator/_examples/types/v1 --output-package github.com/gosoon/code-generator/_examples
```


After generating the code, the user modifies the corresponding business logic as needed.



5、Add main.go file,example：github.com/gosoon/code-generator/_examples/main.go

```
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
```



6、start http server

```
$ go run main.go -config=./kubeconfig
```

7、request RESTful API

```
$ curl -s 127.0.0.1:8080/api/v1/namespace/default | jq .
{
  "code": "OK",
  "message": {
    "metadata": {
      "name": "default",
      "selfLink": "/api/v1/namespaces/default",
      "uid": "57a6c84a-85cd-11e9-b1b3-080027cbf73d",
      "resourceVersion": "159",
      "creationTimestamp": "2019-06-03T07:00:59Z"
    },
    "spec": {
      "finalizers": [
        "kubernetes"
      ]
    },
    "status": {
      "phase": "Active"
    }
  }
}
```


Now automatic generation of CRUD code is the most basic feature,more functions please look forward to, welcome your attention.

