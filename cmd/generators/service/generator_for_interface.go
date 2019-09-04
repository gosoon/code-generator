package service

import (
	"io"

	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/klog"
)

// genServiceInterface generates a package for a controller.
type genServiceInterface struct {
	generator.DefaultGen
	clientsetPackage string
	inputPackages    []string
	outputPackage    string
	imports          namer.ImportTracker
	serviceGenerated bool
	typesToGenerate  []*types.Type
}

var _ generator.Generator = &genServiceInterface{}

func (g *genServiceInterface) Namers(c *generator.Context) namer.NameSystems {
	return namer.NameSystems{
		"raw": namer.NewRawNamer(g.outputPackage, g.imports),
	}
}

// We only want to call GenerateType() once.
func (g *genServiceInterface) Filter(c *generator.Context, t *types.Type) bool {
	ret := !g.serviceGenerated
	g.serviceGenerated = true
	return ret
}

func (g *genServiceInterface) Imports(c *generator.Context) (imports []string) {
	imports = append(imports, g.imports.ImportLines()...)
	imports = append(imports, "k8s.io/client-go/kubernetes")
	imports = append(imports, "apiv1 \"k8s.io/api/core/v1\"")

	for _, pkg := range g.inputPackages {
		imports = append(imports, pkg)
	}
	return
}

func (g *genServiceInterface) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	sw := generator.NewSnippetWriter(w, c, "$", "$")

	klog.Infof("processing type %v", t)
	m := map[string]interface{}{
		"types": g.typesToGenerate,
	}

	sw.Do(typeOptionsStruct, m)
	sw.Do(typeServiceStruct, m)
	sw.Do(newServiceTmpl, m)
	sw.Do(serviceInterfaceTmpl, m)
	return sw.Error()
}

var typeOptionsStruct = `
// Options contains the config by service
type Options struct {
	KubeClientset kubernetes.Interface
}
`

var typeServiceStruct = `
// service implements the Service interface.
type service struct {
	opt *Options
}
`

var newServiceTmpl = `
// New is create a service object.
func New(opt *Options) Interface {
	return &service{opt: opt}
}
`

var serviceInterfaceTmpl = `
// Interface is definition service all method.
type Interface interface {
$range .types$
	Create$.|public$(ctx context.Context, $.|private$Obj *types.$.|public$) error
	Get$.|public$(ctx context.Context, name string) (*apiv1.$.|public$, error)
	Update$.|public$(ctx context.Context, $.|private$Obj *types.$.|public$) error
	Delete$.|public$(ctx context.Context, name string) error
$end$
}
`
