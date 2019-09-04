package generators

import (
	"path/filepath"
	"strings"

	"github.com/gosoon/code-generator/cmd/generators/controller"
	"github.com/gosoon/code-generator/cmd/generators/middleware"
	"github.com/gosoon/code-generator/cmd/generators/service"
	"github.com/gosoon/code-generator/pkg/args"
	"github.com/gosoon/code-generator/pkg/tags"

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

func packageForServer(serverPackagePath string, arguments *args.GeneratorArgs, types []*types.Type, boilerplate []byte) generator.Package {
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
					imports:         generator.NewImportTracker(),
					outputPackage:   arguments.OutputPackagePath,
					typesToGenerate: types,
				},
			}
			return generators
		},
	}
}

// Packages makes the client package definition.
func Packages(context *generator.Context, arguments *args.GeneratorArgs) generator.Packages {
	// load license
	boilerplate, err := arguments.LoadGoBoilerplate()
	if err != nil {
		klog.Fatalf("Failed loading boilerplate: %v", err)
	}

	var packageList []generator.Package
	for _, inputDir := range arguments.InputDirs {
		// Package returns the Package for the given path.
		// package save all types and tags
		p := context.Universe.Package(inputDir)

		// filter have GenTags types
		var typesToGenerate []*types.Type
		for _, t := range p.Types {
			tags := tags.MustParseClientGenTags(append(t.SecondClosestCommentLines, t.CommentLines...))
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

		packagePath := filepath.Join(arguments.OutputPackagePath, "server/controller")
		serverPackagePath := filepath.Join(arguments.OutputPackagePath, "server")
		servicePackagePath := filepath.Join(arguments.OutputPackagePath, "server/service")
		middlewarePackagePath := filepath.Join(arguments.OutputPackagePath, "server/middleware")

		packageList = append(packageList, packageForServer(serverPackagePath, arguments, typesToGenerate, boilerplate))
		packageList = append(packageList, controller.PackageForControllerMeta(packagePath, arguments, boilerplate))
		packageList = append(packageList, service.PackageForService(servicePackagePath, arguments, typesToGenerate, boilerplate))

		// middleware
		packageList = append(packageList, middleware.PackageForMiddleware(middlewarePackagePath, arguments, boilerplate))
		// 为每个 types 生成一个目录以及对应的 CRUD 方法
		for _, t := range typesToGenerate {
			packageList = append(packageList, controller.PackageForTypesController(packagePath,
				arguments, t, boilerplate))
			packageList = append(packageList, service.PackageForTypes(servicePackagePath, arguments, t, boilerplate))
		}
	}
	return generator.Packages(packageList)
}
