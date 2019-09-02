package generators

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gosoon/code-generator/cmd/generators/controller"
	"github.com/gosoon/code-generator/cmd/generators/service"
	"github.com/gosoon/code-generator/pkg/args"
	"github.com/gosoon/code-generator/pkg/tags"

	clientgentypes "k8s.io/code-generator/cmd/client-gen/types"

	//"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/klog"
)

// NameSystems returns the name system used by the generators in this package.
func NameSystems() namer.NameSystems {
	return namer.NameSystems{
		"public":             namer.NewPublicNamer(0),
		"private":            namer.NewPrivateNamer(0),
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

func packageForServer(serverPackagePath string, boilerplate []byte) generator.Package {
	return &generator.DefaultPackage{
		PackageName: "server",
		PackagePath: serverPackagePath,
		HeaderText:  boilerplate,
		PackageDocumentation: []byte(
			`// This package has the automatically generated server.
`),
		// GeneratorFunc returns a list of generators. Each generator generates a
		// single file.
		GeneratorFunc: func(c *generator.Context) (generators []generator.Generator) {
			generators = []generator.Generator{
				generator.DefaultGen{OptionalName: "doc"},
				&genServer{
					DefaultGen: generator.DefaultGen{
						OptionalName: "server",
					},
					imports: generator.NewImportTracker(),
				},
			}
			return generators
		},
	}
}

// Packages makes the client package definition.
func Packages(context *generator.Context, arguments *args.GeneratorArgs) generator.Packages {
	fmt.Println("in package")
	// 1. 加载 license
	boilerplate, err := arguments.LoadGoBoilerplate()
	if err != nil {
		klog.Fatalf("Failed loading boilerplate: %v", err)
	}

	var packageList []generator.Package
	for _, inputDir := range arguments.InputDirs {
		// 2.Package returns the Package for the given path.
		// package save all types and tags
		p := context.Universe.Package(inputDir)
		fmt.Println("universe package success")
		fmt.Printf("p :%+v \n", p)

		var gv clientgentypes.GroupVersion
		var internalGVPkg string

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
			fmt.Printf("----------> t:%+v\n", t)
			// 7.  t.SecondClosestCommentLines is tags, t.CommentLines is comments
			fmt.Printf("tags:%+v  %+v\n", t.SecondClosestCommentLines, t.CommentLines)
			// panic
			tags := tags.MustParseClientGenTags(append(t.SecondClosestCommentLines, t.CommentLines...))
			fmt.Printf("filter tags:%+v \n", tags)
			if !tags.GenerateClient {
				continue
			}
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
		serverPackagePath := filepath.Join(arguments.OutputPackagePath, "server")
		servicePackagePath := filepath.Join(arguments.OutputPackagePath, "server/service")

		packageList = append(packageList, packageForServer(serverPackagePath, boilerplate))
		packageList = append(packageList, controller.PackageForControllerMeta(packagePath, arguments, boilerplate))
		packageList = append(packageList, service.PackageForService(servicePackagePath, arguments, boilerplate))

		// 为每个 types 生成一个目录以及对应的 CRUD 方法
		for _, t := range typesToGenerate {
			packageList = append(packageList, controller.PackageForTypesController(packagePath,
				arguments, t, boilerplate))
			packageList = append(packageList, service.PackageForTypes(servicePackagePath, arguments, t, boilerplate))
		}
	}
	return generator.Packages(packageList)
}

// isInternal returns true if the tags for a member do not contain a json tag
func isInternal(m types.Member) bool {
	return !strings.Contains(m.Tags, "json")
}
