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
					outputPackage: arguments.OutputPackagePath,
					inputPackages: arguments.InputDirs,
					imports:       generator.NewImportTracker(),
				},
			}
			return generators
		},
	}
}
