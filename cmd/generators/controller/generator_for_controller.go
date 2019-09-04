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

package controller

import (
	"io"
	"path/filepath"

	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/klog"
)

// genTypesController generates a package for a controller.
type genTypesController struct {
	generator.DefaultGen
	clientsetPackage   string
	inputPackages      []string
	outputPackage      string
	imports            namer.ImportTracker
	clientsetGenerated bool
	typeToGenerate     *types.Type
}

var _ generator.Generator = &genTypesController{}

func (g *genTypesController) Namers(c *generator.Context) namer.NameSystems {
	return namer.NameSystems{
		"raw": namer.NewRawNamer(g.outputPackage, g.imports),
	}
}

// We only want to call GenerateType() once.
func (g *genTypesController) Filter(c *generator.Context, t *types.Type) bool {
	return t == g.typeToGenerate
}

func (g *genTypesController) Imports(c *generator.Context) (imports []string) {
	imports = append(imports, g.imports.ImportLines()...)
	imports = append(imports, filepath.Join(g.outputPackage, "server/controller"))
	// add input types
	for _, pkg := range g.inputPackages {
		imports = append(imports, pkg)
	}
	imports = append(imports, "github.com/gorilla/mux")
	return
}

func (g *genTypesController) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	sw := generator.NewSnippetWriter(w, c, "$", "$")

	klog.Infof("processing type %v", t)
	m := map[string]interface{}{
		"type": t,
	}

	sw.Do(typeObjectStruct, m)
	sw.Do(newObject, m)
	sw.Do(packRegister, m)
	sw.Do(createObjectHandler, m)
	sw.Do(getObjectHandler, m)
	sw.Do(updateObjectHandler, m)
	sw.Do(deleteObjectHandler, m)

	return sw.Error()
}

var typeObjectStruct = `
// $.type|private$ implements the controller interface.
type $.type|private$ struct {
	opt *controller.Options
}
`

var newObject = `
// New is create a $.type|private$ object.
func New(opt *controller.Options) controller.Controller {
	return &$.type|private${opt: opt}
}
`

var packRegister = `
// Register is register the routes to router
func (c *$.type|private$) Register(router *mux.Router) {
    router = router.PathPrefix("/api/v1").Subrouter()

    // create
    router.Methods("POST").Path("/$.type|lowercaseSingular$").HandlerFunc(
        (c.create$.type|public$))
    
	// get 
    router.Methods("GET").Path("/$.type|lowercaseSingular$/{name}").HandlerFunc(
        (c.get$.type|public$))
	
	// update 
    router.Methods("PUT").Path("/$.type|lowercaseSingular$").HandlerFunc(
        (c.update$.type|public$))
	
	// delete
    router.Methods("DELETE").Path("/$.type|lowercaseSingular$").HandlerFunc(
        (c.delete$.type|public$))
}
`

var createObjectHandler = `
// create$.type|public$
func (c *$.type|private$) create$.type|public$(w http.ResponseWriter, r *http.Request) {
    $.type|private$Obj := &types.$.type|public${}
    err := json.NewDecoder(r.Body).Decode($.type|private$Obj)
    if err != nil {
        controller.BadRequest(w, r, err)
        return
    }

    err = c.opt.Service.Create$.type|public$(r.Context(), $.type|private$Obj)
    if err != nil {
        controller.BadRequest(w, r, err)
        return
    }
    controller.OK(w, r, "success")
}
`

var getObjectHandler = `
// get$.type|public$
func (c *$.type|private$) get$.type|public$(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	$.type|private$Obj,err := c.opt.Service.Get$.type|public$(r.Context(), name)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	controller.Response(w, r, http.StatusOK, $.type|private$Obj)
}
`

var updateObjectHandler = `
// update$.type|public$
func (c *$.type|private$) update$.type|public$(w http.ResponseWriter, r *http.Request) {
    $.type|private$Obj := &types.$.type|public${}
	err := json.NewDecoder(r.Body).Decode($.type|private$Obj)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}

	err = c.opt.Service.Update$.type|public$(r.Context(), $.type|private$Obj)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	controller.OK(w, r, "success")
}
`

var deleteObjectHandler = `
// delete$.type|public$
func (c *$.type|private$) delete$.type|public$(w http.ResponseWriter, r *http.Request) {
    $.type|private$Obj := &types.$.type|public${}
	err := json.NewDecoder(r.Body).Decode($.type|private$Obj)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}

	// get object
	$.type|private$,err := c.opt.Service.Get$.type|public$(r.Context(), $.type|private$Obj.Name)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}

	// delete object
	err = c.opt.Service.Delete$.type|public$(r.Context(), $.type|private$.Name)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	controller.OK(w, r, "success")
}
`
