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

package controller

import (
	"path/filepath"
	"strings"

	"github.com/gosoon/code-generator/pkg/args"

	//"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
)

func PackageForControllerMeta(packagePath string, arguments *args.GeneratorArgs, boilerplate []byte) generator.Package {
	return &generator.DefaultPackage{
		PackageName: "controller",
		PackagePath: packagePath, // output path, "pkg/server/controller/{type}"
		HeaderText:  boilerplate, // Licensed
		PackageDocumentation: []byte(
			`// This package has the automatically generated type controller.
`),
		GeneratorFunc: func(c *generator.Context) (generators []generator.Generator) {
			generators = []generator.Generator{
				generator.DefaultGen{OptionalName: "doc"},
				&genControllerMeta{
					DefaultGen: generator.DefaultGen{
						OptionalName: "controller", // filename : kubernetescluster.go
					},
					outputPackage: arguments.OutputPackagePath, //github.com/gosoon/code-generator
					//typeToGenerate: t,                           // github.com/gosoon/test/pkg/apis/ecs/v1.KubernetesCluster
					imports: generator.NewImportTracker(),
				},
				&genControllerUtils{
					DefaultGen: generator.DefaultGen{
						OptionalName: "utils", // filename : kubernetescluster.go
					},
					outputPackage: arguments.OutputPackagePath, //github.com/gosoon/code-generator
					//typeToGenerate: t,                           // github.com/gosoon/test/pkg/apis/ecs/v1.KubernetesCluster
					imports: generator.NewImportTracker(),
				},
			}
			return generators
		},
	}
}

func PackageForTypesController(packagePath string, arguments *args.GeneratorArgs, t *types.Type, boilerplate []byte) generator.Package {
	packageName := strings.ToLower(t.Name.Name)
	return &generator.DefaultPackage{
		PackageName: packageName,
		PackagePath: filepath.Join(packagePath, packageName), // output path, "pkg/server/controller/{type}"
		HeaderText:  boilerplate,                             // Licensed
		PackageDocumentation: []byte(
			`// This package has the automatically generated type controller.
`),
		GeneratorFunc: func(c *generator.Context) (generators []generator.Generator) {
			generators = []generator.Generator{

				generator.DefaultGen{OptionalName: "doc"},
				&genTypesController{
					DefaultGen: generator.DefaultGen{
						OptionalName: packageName,
					},
					inputPackages:  arguments.InputDirs,
					outputPackage:  arguments.OutputPackagePath,
					typeToGenerate: t,
					imports:        generator.NewImportTracker(),
				},
			}
			return generators
		},
	}
}
