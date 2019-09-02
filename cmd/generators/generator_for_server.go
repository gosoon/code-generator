package generators

import (
	"fmt"
	"io"
	"path/filepath"

	clientgentypes "k8s.io/code-generator/cmd/client-gen/types"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/klog"
)

// genServer generates a package for a controller.
type genServer struct {
	generator.DefaultGen
	groups              []clientgentypes.GroupVersions
	groupGoNames        map[clientgentypes.GroupVersion]string
	clientsetPackage    string
	outputPackage       string
	imports             namer.ImportTracker
	controllerGenerated bool
	typeToGenerate      *types.Type
	objectMeta          *types.Type
}

var _ generator.Generator = &genServer{}

func (g *genServer) Namers(c *generator.Context) namer.NameSystems {
	return namer.NameSystems{
		"raw": namer.NewRawNamer(g.outputPackage, g.imports),
	}
}

// We only want to call GenerateType() once.
func (g *genServer) Filter(c *generator.Context, t *types.Type) bool {
	ret := !g.controllerGenerated
	g.controllerGenerated = true
	return ret
}

func (g *genServer) Imports(c *generator.Context) (imports []string) {
	imports = append(imports, g.imports.ImportLines()...)
	imports = append(imports, filepath.Join(g.outputPackage, "server/service"))
	imports = append(imports, fmt.Sprintf("ctrl \"%v\"", filepath.Join(g.outputPackage, "server/controller")))
	imports = append(imports, "github.com/gorilla/mux")
	return
}

func (g *genServer) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	sw := generator.NewSnippetWriter(w, c, "$", "$")

	klog.Infof("processing type %v", t)
	m := map[string]interface{}{
		"type": t,
	}

	sw.Do(typeServerInterface, m)
	sw.Do(typeOptionsStruct, m)
	sw.Do(typeServerStruct, m)
	sw.Do(serverNewFunc, m)
	sw.Do(serveHTTPFunc, m)
	sw.Do(listenAndServeFunc, m)
	return sw.Error()
}

var typeServerInterface = `
// Server helps start a http server.
type Server interface {
	http.Handler
	ListenAndServe() error
}
`

var typeOptionsStruct = `
// Options contains the config required by server
type Options struct {
	CtrlOptions *ctrl.Options
	ListenAddr  string
}
`

var typeServerStruct = `
// server implements the Server interface.
type server struct {
	opt    Options
	router *mux.Router
}
`

var serverNewFunc = `
// New is create a server object.
func New(opt Options) Server {
	// init service
	options := &service.Options{
		//KubernetesClusterClientset: opt.CtrlOptions.KubernetesClusterClientset,
		KubeClientset:              opt.CtrlOptions.KubeClientset,
	}

	opt.CtrlOptions.Service = service.New(options)

	router := mux.NewRouter().StrictSlash(true)
	//$.type|private$.New(opt.CtrlOptions).Register(router)

	return &server{
		opt:    opt,
		router: router,
	}
}
`

var serveHTTPFunc = `
// ServeHTTP dispatches the handler registered in the matched route.
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
`

var listenAndServeFunc = `
// ListenAndServe start a http server.
func (s *server) ListenAndServe() error {
	server := &http.Server{
		Handler: s.router,
		Addr:    s.opt.ListenAddr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout:   15 * time.Second,
		ReadTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
`
