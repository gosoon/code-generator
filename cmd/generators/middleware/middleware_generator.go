package middleware

import (
	"github.com/gosoon/code-generator/pkg/args"
	"k8s.io/gengo/generator"
)

// PackageForMiddleware xxx
func PackageForMiddleware(packagePath string, arguments *args.GeneratorArgs, boilerplate []byte) generator.Package {
	return &generator.DefaultPackage{
		PackageName: "middleware",
		PackagePath: packagePath,
		HeaderText:  boilerplate,
		PackageDocumentation: []byte(
			`// This package has the automatically generated middleware.
`),
		// GeneratorFunc returns a list of generators. Each generator generates a
		// single file.
		GeneratorFunc: func(c *generator.Context) (generators []generator.Generator) {
			generators = []generator.Generator{
				generator.DefaultGen{OptionalName: "doc"},

				&genMiddleware{
					DefaultGen: generator.DefaultGen{
						OptionalName: "auth",
					},
					inputPackages: arguments.InputDirs,
					imports:       generator.NewImportTracker(),
				},
			}
			return generators
		},
	}
}
