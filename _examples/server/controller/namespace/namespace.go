/*
 * Copyright 2019 gosoon.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package namespace

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gosoon/code-generator/_examples/server/controller"
	"github.com/gosoon/code-generator/_examples/types/v1"
	"github.com/gosoon/test/server/middleware"
)

// namespace implements the controller interface.
type namespace struct {
	opt *controller.Options
}

// New is create a namespace object.
func New(opt *controller.Options) controller.Controller {
	return &namespace{opt: opt}
}

// Register is register the routes to router
func (c *namespace) Register(router *mux.Router) {
	router = router.PathPrefix("/api/v1").Subrouter()

	// create
	router.Methods("POST").Path("/namespace").HandlerFunc(
		middleware.Authenticate(http.HandlerFunc((c.createNamespace))))

	// get
	router.Methods("GET").Path("/namespace/{name}").HandlerFunc(
		middleware.Authenticate(http.HandlerFunc((c.getNamespace))))

	// update
	router.Methods("PUT").Path("/namespace").HandlerFunc(
		middleware.Authenticate(http.HandlerFunc((c.updateNamespace))))

	// delete
	router.Methods("DELETE").Path("/namespace").HandlerFunc(
		middleware.Authenticate(http.HandlerFunc((c.deleteNamespace))))
}

// createNamespace
func (c *namespace) createNamespace(w http.ResponseWriter, r *http.Request) {
	namespaceObj := &types.Namespace{}
	err := json.NewDecoder(r.Body).Decode(namespaceObj)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}

	err = c.opt.Service.CreateNamespace(r.Context(), namespaceObj)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	controller.OK(w, r, "success")
}

// getNamespace
func (c *namespace) getNamespace(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	namespaceObj, err := c.opt.Service.GetNamespace(r.Context(), name)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	controller.Response(w, r, http.StatusOK, namespaceObj)
}

// updateNamespace
func (c *namespace) updateNamespace(w http.ResponseWriter, r *http.Request) {
	namespaceObj := &types.Namespace{}
	err := json.NewDecoder(r.Body).Decode(namespaceObj)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}

	err = c.opt.Service.UpdateNamespace(r.Context(), namespaceObj)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	controller.OK(w, r, "success")
}

// deleteNamespace
func (c *namespace) deleteNamespace(w http.ResponseWriter, r *http.Request) {
	namespaceObj := &types.Namespace{}
	err := json.NewDecoder(r.Body).Decode(namespaceObj)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}

	// get object
	namespace, err := c.opt.Service.GetNamespace(r.Context(), namespaceObj.Name)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}

	// delete object
	err = c.opt.Service.DeleteNamespace(r.Context(), namespace.Name)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	controller.OK(w, r, "success")
}
