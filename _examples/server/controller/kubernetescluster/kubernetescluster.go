/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by main. DO NOT EDIT.

package kubernetescluster

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gosoon/code-generator/server/controller"
)

type kubernetesCluster struct {
	opt *controller.Options
}

func New(opt *controller.Options) controller.Controller {
	return &kubernetesCluster{opt: opt}
}

// Register is register the routes to router
func (c *kubernetesCluster) Register(router *mux.Router) {
	router = router.PathPrefix("/api/v1").Subrouter()

	// create
	router.Methods("POST").Path("/kubernetescluster").HandlerFunc(
		(c.createKubernetesCluster))

	// get
	router.Methods("GET").Path("/kubernetescluster/{name}").HandlerFunc(
		(c.createKubernetesCluster))

	// update
	router.Methods("PUT").Path("/kubernetescluster").HandlerFunc(
		(c.createKubernetesCluster))

	// delete
	router.Methods("DELETE").Path("/kubernetescluster").HandlerFunc(
		(c.createKubernetesCluster))
}

// createKubernetesCluster
func (c *kubernetesCluster) createKubernetesCluster(w http.ResponseWriter, r *http.Request) {
	kubernetesClusterObj := &types.KubernetesCluster{}
	err := json.NewDecoder(r.Body).Decode(kubernetesClusterObj)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	name := mux.Vars(r)["name"]

	err = c.opt.Service.CreateKubernetesCluster(r.Context(), kubernetesClusterObj)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	controller.OK(w, r, "success")
}

// getKubernetesCluster
func (c *kubernetesCluster) getKubernetesCluster(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	kubernetesClusterObj, err := c.opt.Service.GetKubernetesCluster(r.Context(), name)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	controller.OK(w, r, kubernetesClusterObj)
}

// updateKubernetesCluster
func (c *kubernetesCluster) updateKubernetesCluster(w http.ResponseWriter, r *http.Request) {
	kubernetesClusterObj := &types.KubernetesCluster{}
	err := json.NewDecoder(r.Body).Decode(kubernetesClusterObj)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}

	err = c.opt.Service.UpdateKubernetesCluster(r.Context(), kubernetesClusterObj)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	controller.OK(w, r, "success")
}

// deleteKubernetesCluster
func (c *kubernetesCluster) deleteKubernetesCluster(w http.ResponseWriter, r *http.Request) {
	kubernetesClusterObj := &types.KubernetesCluster{}
	err := json.NewDecoder(r.Body).Decode(kubernetesClusterObj)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}

	// get object
	kubernetesClusterObj, err = c.opt.Service.GetKubernetesCluster(r.Context(), kubernetesClusterObj.Name)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}

	// delete object
	err = c.opt.Service.DeleteKubernetesCluster(r.Context(), kubernetesClusterObj.ID)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	controller.OK(w, r, "success")
}