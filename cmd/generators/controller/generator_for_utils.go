package controller

import (
	"io"

	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/klog"
)

// genControllerUtils generates a package for a controller.
type genControllerUtils struct {
	generator.DefaultGen
	clientsetPackage    string
	outputPackage       string
	imports             namer.ImportTracker
	controllerGenerated bool

	typeToGenerate *types.Type
	objectMeta     *types.Type
}

var _ generator.Generator = &genControllerUtils{}

func (g *genControllerUtils) Namers(c *generator.Context) namer.NameSystems {
	return namer.NameSystems{
		"raw": namer.NewRawNamer(g.outputPackage, g.imports),
	}
}

// We only want to call GenerateType() once.
func (g *genControllerUtils) Filter(c *generator.Context, t *types.Type) bool {
	ret := !g.controllerGenerated
	g.controllerGenerated = true
	return ret
}

func (g *genControllerUtils) Imports(c *generator.Context) (imports []string) {
	imports = append(imports, g.imports.ImportLines()...)
	imports = append(imports, "k8s.io/klog")
	return
}

func (g *genControllerUtils) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	sw := generator.NewSnippetWriter(w, c, "$", "$")

	klog.Infof("processing type %v", t)
	m := map[string]interface{}{
		"type": t,
	}

	sw.Do(typeCommRespStruct, m)
	sw.Do(respDefine, m)
	return sw.Error()
}

var typeCommRespStruct = `
type commResp struct {
    Code    string` + "    `json:\"code\"`" + `
    Message interface{}` + "    `json:\"message\"`" + `
}
`

var respDefine = `
// OK reply
func OK(w http.ResponseWriter, r *http.Request, message string) {
	Response(w, r, http.StatusOK, message)
}

// ResourceNotFound will return an error message indicating that the resource is not exist
func ResourceNotFound(w http.ResponseWriter, r *http.Request, message string) {
	Response(w, r, http.StatusNotFound, message)
}

// BadRequest will return an error message indicating that the request is invalid
func BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	Response(w, r, http.StatusBadRequest, err.Error())
}

// Forbidden will block user access the resource, not authorized
func Forbidden(w http.ResponseWriter, r *http.Request, err error) {
	Response(w, r, http.StatusForbidden, err.Error())
}

// Unauthorized will block user access the api, not login
func Unauthorized(w http.ResponseWriter, r *http.Request, err error) {
	Response(w, r, http.StatusUnauthorized, err.Error())
}

// InternalError will return an error message indicating that the something is error inside the controller
func InternalError(w http.ResponseWriter, r *http.Request, err error) {
	Response(w, r, http.StatusInternalServerError, err.Error())
}

// ServiceUnavailable will return an error message indicating that the service is not available now
func ServiceUnavailable(w http.ResponseWriter, r *http.Request, err error) {
	Response(w, r, http.StatusServiceUnavailable, err.Error())
}

// Conflict xxx
func Conflict(w http.ResponseWriter, r *http.Request, err error) {
	Response(w, r, http.StatusConflict, err.Error())
}

// NotAcceptable xxx
func NotAcceptable(w http.ResponseWriter, r *http.Request, err error) {
	Response(w, r, http.StatusNotAcceptable, err.Error())
}

// Response : http response func (no return http code)
func Response(w http.ResponseWriter, r *http.Request, httpCode int, message interface{}) {
	resp := commResp{
		Code:    http.StatusText(httpCode),
		Message: message,
	}

	jsonByte, err := json.Marshal(resp)
	if err != nil {
		klog.Errorf("marshal [%v] failed with err [%v]", resp, err)
	}
	_, err = r.Cookie("WriteHeader")
	// if no WriteHeader
	if err != nil {
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpCode)
		w.Write(jsonByte)
	}
}
`
