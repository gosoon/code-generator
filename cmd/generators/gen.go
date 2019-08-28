package generators

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"k8s.io/code-generator/cmd/client-gen/generators/util"
	clientgentypes "k8s.io/code-generator/cmd/client-gen/types"
	"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/klog"
)

// NameSystems returns the name system used by the generators in this package.
func NameSystems() namer.NameSystems {
	return namer.NameSystems{
		"public":             namer.NewPublicNamer(0),    // types
		"private":            namer.NewPrivateNamer(0),   // Types
		"raw":                namer.NewRawNamer("", nil), // package.Types
		"publicPlural":       namer.NewPublicPluralNamer(nil),
		"allLowercasePlural": namer.NewAllLowercasePluralNamer(nil),
		"lowercaseSingular":  &lowercaseSingularNamer{},
	}
}

// lowercaseSingularNamer implements Namer
type lowercaseSingularNamer struct{}

// Name returns t's name in all lowercase.
func (n *lowercaseSingularNamer) Name(t *types.Type) string {
	return strings.ToLower(t.Name.Name)
}

// DefaultNameSystem returns the default name system for ordering the types to be
// processed by the generators in this package.
func DefaultNameSystem() string {
	return "public"
}

//func packageForController(customArgs *clientgenargs.CustomArgs, clientsetPackage string, groupGoNames map[clientgentypes.GroupVersion]string, boilerplate []byte) generator.Package {
//return &generator.DefaultPackage{
//PackageName: customArgs.ClientsetName,
//PackagePath: clientsetPackage,
//HeaderText:  boilerplate,
//PackageDocumentation: []byte(
//`// This package has the automatically generated clientset.
//`),
//// GeneratorFunc returns a list of generators. Each generator generates a
//// single file.
//GeneratorFunc: func(c *generator.Context) (generators []generator.Generator) {
//generators = []generator.Generator{
//// Always generate a "doc.go" file.
//generator.DefaultGen{OptionalName: "doc"},

//&genController{
//DefaultGen: generator.DefaultGen{
//OptionalName: "controller",
//},
//groups:           customArgs.Groups,
//groupGoNames:     groupGoNames,
//clientsetPackage: clientsetPackage,
//outputPackage:    customArgs.ClientsetName,
//imports:          generator.NewImportTracker(),
//},
//}
//return generators
//},
//}
//}

// Packages makes the client package definition.
func Packages(context *generator.Context, arguments *args.GeneratorArgs) generator.Packages {
	// 1. 加载 license
	boilerplate, err := arguments.LoadGoBoilerplate()
	if err != nil {
		klog.Fatalf("Failed loading boilerplate: %v", err)
	}

	var packageList generator.Packages
	for _, inputDir := range arguments.InputDirs {
		// 2.Package returns the Package for the given path.
		// package save all types and tags
		p := context.Universe.Package(inputDir)
		fmt.Printf("p.types:%+v\n", p.Types)
		fmt.Printf("p.path:%+v\n", p.Path)
		fmt.Printf("p.comments:%+v\n", p.Comments)

		// 3. objectMeta is have genclient's types
		objectMeta, internal, err := objectMetaForPackage(p)
		if err != nil {
			klog.Fatal(err)
		}
		if objectMeta == nil {
			// no types in this package had genclient
			continue
		}

		var gv clientgentypes.GroupVersion
		var internalGVPkg string

		// 4. if objectMeta no json tag and is internal
		if internal {
			lastSlash := strings.LastIndex(p.Path, "/")
			if lastSlash == -1 {
				klog.Fatalf("error constructing internal group version for package %q", p.Path)
			}
			gv.Group = clientgentypes.Group(p.Path[lastSlash+1:])
			internalGVPkg = p.Path
		} else {
			parts := strings.Split(p.Path, "/")
			gv.Group = clientgentypes.Group(parts[len(parts)-2])
			gv.Version = clientgentypes.Version(parts[len(parts)-1])

			internalGVPkg = strings.Join(parts[0:len(parts)-1], "/")
		}
		//groupPackageName := strings.ToLower(gv.Group.NonEmpty())

		// If there's a comment of the form "// +groupName=somegroup" or
		// "// +groupName=somegroup.foo.bar.io", use the first field (somegroup) as the name of the
		// group when generating.
		if override := types.ExtractCommentTags("+", p.Comments)["groupName"]; override != nil {
			gv.Group = clientgentypes.Group(strings.SplitN(override[0], ".", 2)[0])
		}

		// 5. gv is  "Group:pks,Version:v1"
		fmt.Printf("gv:%+v \n", gv, internalGVPkg)

		// 6.filter have GenTags types
		var typesToGenerate []*types.Type
		for _, t := range p.Types {
			// 7.  t.SecondClosestCommentLines is tags, t.CommentLines is comments
			fmt.Printf("tags:%+v  %+v\n", t.SecondClosestCommentLines, t.CommentLines)
			tags := util.MustParseClientGenTags(append(t.SecondClosestCommentLines, t.CommentLines...))
			if !tags.GenerateClient || !tags.HasVerb("list") || !tags.HasVerb("get") {
				continue
			}
			fmt.Printf("filter tags:%+v \n", tags)
			typesToGenerate = append(typesToGenerate, t)
		}
		if len(typesToGenerate) == 0 {
			continue
		}
		orderer := namer.Orderer{Namer: namer.NewPrivateNamer(0)}
		typesToGenerate = orderer.OrderTypes(typesToGenerate)

		//  8.github.com/gosoon/test/pkg/apis/ecs/v1.KubernetesCluster
		fmt.Printf("typesToGenerate : %+v \n", typesToGenerate)

		// 9. packagePath is "listers/ecs/v1"
		// github.com/gosoon/code-generator ecs v1
		//packagePath := filepath.Join(arguments.OutputPackagePath, groupPackageName, strings.ToLower(gv.Version.NonEmpty()))
		packagePath := filepath.Join(arguments.OutputPackagePath, "server/controller")

		// 为每个 types 生成一个目录以及对应的 CRUD 方法
		for _, t := range typesToGenerate {
			packageName := strings.ToLower(t.Name.Name)

			defaultPkg := &generator.DefaultPackage{
				PackageName: packageName,
				PackagePath: filepath.Join(packagePath, packageName), // output path, "pkg/server/controller"
				HeaderText:  boilerplate,                             // Licensed
				GeneratorFunc: func(c *generator.Context) (generators []generator.Generator) {
					fmt.Printf("----------strings.ToLower(t.Name.Name):%+v \n", strings.ToLower(t.Name.Name))
					fmt.Printf("arguments.OutputPackagePath:%+v\n", arguments.OutputPackagePath)
					fmt.Printf("gv:%+v \n", gv)
					fmt.Printf("internalGVPkg:%+v\n", internalGVPkg)
					fmt.Printf("t:%+v\n", t)
					fmt.Printf("generator.NewImportTracker():%+v\n", generator.NewImportTracker())
					fmt.Printf("objectMeta:%+v\n", objectMeta)

					generators = append(generators, &genController{
						DefaultGen: generator.DefaultGen{
							OptionalName: packageName, // kubernetescluster
						},
						outputPackage: arguments.OutputPackagePath, //github.com/gosoon/code-generator
						//groupVersion:   gv,                          // ecs  v1
						//internalGVPkg:  internalGVPkg,               // github.com/gosoon/test/pkg/apis/ecs
						typeToGenerate: t, // github.com/gosoon/test/pkg/apis/ecs/v1.KubernetesCluster
						imports:        generator.NewImportTracker(),
						objectMeta:     objectMeta, // k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta
					})
					return generators
				},
			}
			packageList = append(packageList, defaultPkg)
		}

	}
	return packageList
}

// objectMetaForPackage returns the type of ObjectMeta used by package p.
func objectMetaForPackage(p *types.Package) (*types.Type, bool, error) {
	generatingForPackage := false
	for _, t := range p.Types {
		// filter out types which dont have genclient.
		if !util.MustParseClientGenTags(append(t.SecondClosestCommentLines, t.CommentLines...)).GenerateClient {
			continue
		}
		generatingForPackage = true
		for _, member := range t.Members {
			if member.Name == "ObjectMeta" {
				return member.Type, isInternal(member), nil
			}
		}
	}
	if generatingForPackage {
		return nil, false, fmt.Errorf("unable to find ObjectMeta for any types in package %s", p.Path)
	}
	return nil, false, nil
}

// isInternal returns true if the tags for a member do not contain a json tag
func isInternal(m types.Member) bool {
	return !strings.Contains(m.Tags, "json")
}

// listerGenerator produces a file of listers for a given GroupVersion and
// type.
type listerGenerator struct {
	generator.DefaultGen
	outputPackage  string
	groupVersion   clientgentypes.GroupVersion
	internalGVPkg  string
	typeToGenerate *types.Type
	imports        namer.ImportTracker
	objectMeta     *types.Type
}

var _ generator.Generator = &listerGenerator{}

func (g *listerGenerator) Filter(c *generator.Context, t *types.Type) bool {
	return t == g.typeToGenerate
}

func (g *listerGenerator) Namers(c *generator.Context) namer.NameSystems {
	return namer.NameSystems{
		"raw": namer.NewRawNamer(g.outputPackage, g.imports),
	}
}

func (g *listerGenerator) Imports(c *generator.Context) (imports []string) {
	imports = append(imports, g.imports.ImportLines()...)
	imports = append(imports, "k8s.io/apimachinery/pkg/api/errors")
	imports = append(imports, "k8s.io/apimachinery/pkg/labels")
	// for Indexer
	imports = append(imports, "k8s.io/client-go/tools/cache")
	return
}

func (g *listerGenerator) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	sw := generator.NewSnippetWriter(w, c, "$", "$")

	klog.Infof("processing type %v", t)
	m := map[string]interface{}{
		"Resource":   c.Universe.Function(types.Name{Package: t.Name.Package, Name: "Resource"}),
		"type":       t,
		"objectMeta": g.objectMeta,
	}

	//tags, err := util.ParseClientGenTags(append(t.SecondClosestCommentLines, t.CommentLines...))
	//if err != nil {
	//return err
	//}

	//if tags.NonNamespaced {
	//sw.Do(typeListerInterface_NonNamespaced, m)
	//} else {
	sw.Do(typeListerInterface, m)
	//}

	return sw.Error()
}

var typeListerInterface = `
// $.type|public$Lister helps list $.type|publicPlural$.
type $.type|public$Lister interface {
    // List lists all $.type|publicPlural$ in the indexer.
    List(selector labels.Selector) (ret []*$.type|raw$, err error)
    // $.type|publicPlural$ returns an object that can list and get $.type|publicPlural$.
    $.type|publicPlural$(namespace string) $.type|public$NamespaceLister
    $.type|public$ListerExpansion
}
`
