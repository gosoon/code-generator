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
	"io"
	"path/filepath"

	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
)

// genControllerMeta generates a package for a controller.
type genControllerMeta struct {
	generator.DefaultGen
	clientsetPackage    string
	outputPackage       string
	imports             namer.ImportTracker
	controllerGenerated bool

	typeToGenerate *types.Type
	objectMeta     *types.Type
}

var _ generator.Generator = &genControllerMeta{}

func (g *genControllerMeta) Namers(c *generator.Context) namer.NameSystems {
	return namer.NameSystems{
		"raw": namer.NewRawNamer(g.outputPackage, g.imports),
	}
}

// We only want to call GenerateType() once.
func (g *genControllerMeta) Filter(c *generator.Context, t *types.Type) bool {
	ret := !g.controllerGenerated
	g.controllerGenerated = true
	return ret
}

func (g *genControllerMeta) Imports(c *generator.Context) (imports []string) {
	imports = append(imports, g.imports.ImportLines()...)
	imports = append(imports, filepath.Join(g.outputPackage, "server/service"))
	imports = append(imports, "github.com/gorilla/mux")
	imports = append(imports, "k8s.io/client-go/kubernetes")
	return
}

func (g *genControllerMeta) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	sw := generator.NewSnippetWriter(w, c, "$", "$")

	m := map[string]interface{}{}

	sw.Do(typeOptionsStruct, m)
	sw.Do(typeControllerInterface, m)
	return sw.Error()
}

var typeOptionsStruct = `
// Options contains the config by controller
type Options struct {
    KubeClientset              kubernetes.Interface
    Service                    service.Interface
}
`
var typeControllerInterface = `
// Controller helps register to router. 
type Controller interface {
    Register(router *mux.Router)
}
`
