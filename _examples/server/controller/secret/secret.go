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
package secret

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gosoon/code-generator/server/controller"
	"github.com/gosoon/test/pkg/types"
)

// secret implements the controller interface.
type secret struct {
	opt *controller.Options
}

// New is create a secret object.
func New(opt *controller.Options) controller.Controller {
	return &secret{opt: opt}
}

// Register is register the routes to router
func (c *secret) Register(router *mux.Router) {
	router = router.PathPrefix("/api/v1").Subrouter()

	// create
	router.Methods("POST").Path("/secret").HandlerFunc(
		(c.createSecret))

	// get
	router.Methods("GET").Path("/secret/{name}").HandlerFunc(
		(c.getSecret))

	// update
	router.Methods("PUT").Path("/secret").HandlerFunc(
		(c.updateSecret))

	// delete
	router.Methods("DELETE").Path("/secret").HandlerFunc(
		(c.deleteSecret))
}

// createSecret
func (c *secret) createSecret(w http.ResponseWriter, r *http.Request) {
	secretObj := &types.Secret{}
	err := json.NewDecoder(r.Body).Decode(secretObj)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}

	err = c.opt.Service.CreateSecret(r.Context(), secretObj)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	controller.OK(w, r, "success")
}

// getSecret
func (c *secret) getSecret(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	secretObj, err := c.opt.Service.GetSecret(r.Context(), name)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	controller.Response(w, r, http.StatusOK, secretObj)
}

// updateSecret
func (c *secret) updateSecret(w http.ResponseWriter, r *http.Request) {
	secretObj := &types.Secret{}
	err := json.NewDecoder(r.Body).Decode(secretObj)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}

	err = c.opt.Service.UpdateSecret(r.Context(), secretObj)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	controller.OK(w, r, "success")
}

// deleteSecret
func (c *secret) deleteSecret(w http.ResponseWriter, r *http.Request) {
	secretObj := &types.Secret{}
	err := json.NewDecoder(r.Body).Decode(secretObj)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}

	// get object
	secret, err := c.opt.Service.GetSecret(r.Context(), secretObj.Name)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}

	// delete object
	err = c.opt.Service.DeleteSecret(r.Context(), secret.Name)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	controller.OK(w, r, "success")
}
