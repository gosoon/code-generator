package middleware

import (
	"io"
	"path/filepath"

	clientgentypes "k8s.io/code-generator/cmd/client-gen/types"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/klog"
)

// genMiddleware generates a package for a controller.
type genMiddleware struct {
	generator.DefaultGen
	groups           []clientgentypes.GroupVersions
	groupGoNames     map[clientgentypes.GroupVersion]string
	clientsetPackage string
	inputPackages    []string
	outputPackage    string
	imports          namer.ImportTracker
	serviceGenerated bool

	typeToGenerate *types.Type
	objectMeta     *types.Type
}

var _ generator.Generator = &genMiddleware{}

func (g *genMiddleware) Namers(c *generator.Context) namer.NameSystems {
	return namer.NameSystems{
		"raw": namer.NewRawNamer(g.outputPackage, g.imports),
	}
}

// We only want to call GenerateType() once.
func (g *genMiddleware) Filter(c *generator.Context, t *types.Type) bool {
	ret := !g.serviceGenerated
	g.serviceGenerated = true
	return ret
}

func (g *genMiddleware) Imports(c *generator.Context) (imports []string) {
	imports = append(imports, g.imports.ImportLines()...)
	imports = append(imports, "github.com/spf13/viper")
	imports = append(imports, filepath.Join(g.outputPackage, "server/controller"))
	return
}

func (g *genMiddleware) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	sw := generator.NewSnippetWriter(w, c, "$", "$")

	klog.Infof("processing type %v", t)
	m := map[string]interface{}{}

	sw.Do(authenticateTmpl, m)
	return sw.Error()
}

var authenticateTmpl = `
// AuthenticateMW will create a authenticate middleware
func Authenticate(next http.Handler) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if len(token) != 0 {
            bearerValue := strings.Split(token, " ")[1]
            // token in config
            if bearerValue == viper.GetString(config.token) {
                next.ServeHTTP(w, r)
            }
        }
        controller.Unauthorized(w, r, fmt.Sprintf("Authenticate failed,plz check your token."))
    }
}
`
