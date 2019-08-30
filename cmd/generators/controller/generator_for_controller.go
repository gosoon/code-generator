package controller

import (
	"io"
	"path/filepath"

	clientgentypes "k8s.io/code-generator/cmd/client-gen/types"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/klog"
)

// genTypesController generates a package for a controller.
type genTypesController struct {
	generator.DefaultGen
	groups             []clientgentypes.GroupVersions
	groupGoNames       map[clientgentypes.GroupVersion]string
	clientsetPackage   string
	outputPackage      string
	imports            namer.ImportTracker
	clientsetGenerated bool

	typeToGenerate *types.Type
	objectMeta     *types.Type
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
type $.type|private$ struct {
	opt *controller.Options
}
`

var newObject = `
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
        (c.create$.type|public$))
	
	// update 
    router.Methods("PUT").Path("/$.type|lowercaseSingular$").HandlerFunc(
        (c.create$.type|public$))
	
	// delete
    router.Methods("DELETE").Path("/$.type|lowercaseSingular$").HandlerFunc(
        (c.create$.type|public$))
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
    name := mux.Vars(r)["name"]

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
	controller.OK(w, r, $.type|private$Obj)
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
	$.type|private$Obj,err = c.opt.Service.Get$.type|public$(r.Context(), $.type|private$Obj.Name)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}

	// delete object
	err = c.opt.Service.Delete$.type|public$(r.Context(), $.type|private$Obj.ID)
	if err != nil {
		controller.BadRequest(w, r, err)
		return
	}
	controller.OK(w, r, "success")
}
`
