package generators

import (
	"io"

	clientgentypes "k8s.io/code-generator/cmd/client-gen/types"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
)

// genController generates a package for a controller.
type genController struct {
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

var _ generator.Generator = &genController{}

func (g *genController) Namers(c *generator.Context) namer.NameSystems {
	return namer.NameSystems{
		"raw": namer.NewRawNamer(g.outputPackage, g.imports),
	}
}

// We only want to call GenerateType() once.
func (g *genController) Filter(c *generator.Context, t *types.Type) bool {
	//ret := !g.clientsetGenerated
	//g.clientsetGenerated = true
	return true
}

func (g *genController) Imports(c *generator.Context) (imports []string) {
	//imports = append(imports, g.imports.ImportLines()...)
	//imports = append(imports, "k8s.io/apimachinery/pkg/api/errors")
	return
}

func (g *genController) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	return nil
}
