/*
 * Copyright 2019 gosoon.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"strings"

	"github.com/gosoon/code-generator/pkg/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
)

// PackageForService xxx
func PackageForService(packageName string, arguments *args.GeneratorArgs, types []*types.Type, boilerplate []byte) generator.Package {
	return &generator.DefaultPackage{
		PackageName: "service",
		PackagePath: packageName,
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
					typesToGenerate: types,
					inputPackages:   arguments.InputDirs,
					imports:         generator.NewImportTracker(),
				},
			}
			return generators
		},
	}
}

// PackageForTypes xxx
func PackageForTypes(packagePath string, arguments *args.GeneratorArgs, t *types.Type, boilerplate []byte) generator.Package {
	packageName := strings.ToLower(t.Name.Name)
	return &generator.DefaultPackage{
		PackageName: "service",
		PackagePath: packagePath,
		HeaderText:  boilerplate,
		PackageDocumentation: []byte(
			`// This package has the automatically generated service.
`),
		// GeneratorFunc returns a list of generators. Each generator generates a
		// single file.
		GeneratorFunc: func(c *generator.Context) (generators []generator.Generator) {
			generators = []generator.Generator{
				generator.DefaultGen{OptionalName: "doc"},

				&genTypesService{
					DefaultGen: generator.DefaultGen{
						OptionalName: packageName,
					},
					typeToGenerate: t,
					inputPackages:  arguments.InputDirs,
					imports:        generator.NewImportTracker(),
				},
			}
			return generators
		},
	}
}
