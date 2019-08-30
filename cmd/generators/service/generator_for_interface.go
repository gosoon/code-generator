package service

import (
	"io"
	"path/filepath"

	clientgentypes "k8s.io/code-generator/cmd/client-gen/types"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/klog"
)

// genServiceInterface generates a package for a controller.
type genServiceInterface struct {
	generator.DefaultGen
	groups           []clientgentypes.GroupVersions
	groupGoNames     map[clientgentypes.GroupVersion]string
	clientsetPackage string
	outputPackage    string
	imports          namer.ImportTracker
	serviceGenerated bool

	typeToGenerate *types.Type
	objectMeta     *types.Type
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
	imports = append(imports, filepath.Join(g.outputPackage, "pkg/types"))
	return
}

func (g *genServiceInterface) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	sw := generator.NewSnippetWriter(w, c, "$", "$")

	klog.Infof("processing type %v", t)
	m := map[string]interface{}{
		"type": t,
	}

	sw.Do(typeOptionsStruct, m)
	sw.Do(typeServiceStruct, m)
	sw.Do(newServiceTmpl, m)
	sw.Do(serviceInterfaceTmpl, m)
	return sw.Error()
}

var typeOptionsStruct = `
type Options struct {
	//KubeClientset              kubernetes.Interface
}
`

var typeServiceStruct = `
type service struct {
	opt *Options
}
`

var newServiceTmpl = `
func New(opt *Options) Interface {
	return &service{opt: opt}
}
`

var serviceInterfaceTmpl = `
type Interface interface {
	// cluster
	CreateCluster(ctx context.Context, region string, namespace string, name string, clusterInfo *types.EcsClient) error
	DeleteCluster(ctx context.Context, region string, namespace string, name string, clusterInfo *types.EcsClient) error
}
`
