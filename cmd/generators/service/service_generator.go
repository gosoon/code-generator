package service

import (
	"k8s.io/gengo/generator"
)

func PackageForService(servicePackage string, boilerplate []byte) generator.Package {
	return &generator.DefaultPackage{
		PackageName: "service",
		PackagePath: servicePackage,
		HeaderText:  boilerplate,
		PackageDocumentation: []byte(
			`// This package has the automatically generated service.
`),
		// GeneratorFunc returns a list of generators. Each generator generates a
		// single file.
		GeneratorFunc: func(c *generator.Context) (generators []generator.Generator) {
			generators = []generator.Generator{
				generator.DefaultGen{OptionalName: "doc"},

				&genServiceInterface{
					DefaultGen: generator.DefaultGen{
						OptionalName: "interface",
					},
					//groups:           customArgs.Groups,
					//groupGoNames:     groupGoNames,
					//clientsetPackage: clientsetPackage,
					//outputPackage:    customArgs.ClientsetName,
					imports: generator.NewImportTracker(),
				},
			}
			return generators
		},
	}
}
