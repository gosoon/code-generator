package service

import (
	"io"

	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/klog"
)

// genTypesService generates a package for a controller.
type genTypesService struct {
	generator.DefaultGen
	clientsetPackage string
	inputPackages    []string
	outputPackage    string
	imports          namer.ImportTracker
	serviceGenerated bool
	typeToGenerate   *types.Type
}

var _ generator.Generator = &genTypesService{}

func (g *genTypesService) Namers(c *generator.Context) namer.NameSystems {
	return namer.NameSystems{
		"raw": namer.NewRawNamer(g.outputPackage, g.imports),
	}
}

// We only want to call GenerateType() once.
func (g *genTypesService) Filter(c *generator.Context, t *types.Type) bool {
	ret := !g.serviceGenerated
	g.serviceGenerated = true
	return ret
}

func (g *genTypesService) Imports(c *generator.Context) (imports []string) {
	imports = append(imports, g.imports.ImportLines()...)
	imports = append(imports, "k8s.io/klog")
	imports = append(imports, "apiv1 \"k8s.io/api/core/v1\"")
	imports = append(imports, "metav1 \"k8s.io/apimachinery/pkg/apis/meta/v1\"")

	for _, pkg := range g.inputPackages {
		imports = append(imports, pkg)
	}
	return
}

func (g *genTypesService) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	sw := generator.NewSnippetWriter(w, c, "$", "$")

	klog.Infof("processing type %v", g.typeToGenerate)
	m := map[string]interface{}{
		"type": g.typeToGenerate,
	}

	sw.Do(createObjectService, m)
	sw.Do(getObjectService, m)
	sw.Do(updateObjectService, m)
	sw.Do(deleteObjectService, m)
	return sw.Error()
}

var createObjectService = `
// Create$.type|public$ xxx
// TODO(user): Modify this function to implement your logic.This example use namespace.
func (s *service) Create$.type|public$(ctx context.Context, $.type|private$Obj *types.$.type|public$) error {
    clientset := s.opt.KubeClientset
    $.type|private$ := &apiv1.$.type|public${
        TypeMeta: metav1.TypeMeta{
            APIVersion: "v1",
            Kind:       "$.type|public$",
        },
        ObjectMeta: metav1.ObjectMeta{
            Name: $.type|private$Obj.Name,
        },
    }

    _, err := clientset.CoreV1().$.type|publicPlural$().Create(namespace)
    if err != nil {
        klog.Errorf("create $.type|private$ failed with:%v", err)
        return err
    } 
    return nil
}
`
var getObjectService = `
// Get$.type|public$ xxx
// TODO(user): Modify this function to implement your logic.This example use namespace.
func (s *service) Get$.type|public$(ctx context.Context, name string) (*apiv1.$.type|public$, error) {
    clientset := s.opt.KubeClientset

    $.type|private$, err := clientset.CoreV1().$.type|publicPlural$().Get(name, metav1.GetOptions{})
    if err != nil {
        klog.Errorf("get $.type|private$ %v failed with:%v", name, err)
        return nil, err
    }

    return $.type|private$, nil
}
`

var updateObjectService = `
// Update$.type|public$ xxx
// TODO(user): Modify this function to implement your logic.This example use namespace. 
func (s *service) Update$.type|public$(ctx context.Context, $.type|private$Obj *types.$.type|public$) error {
    clientset := s.opt.KubeClientset

	var err error
	$.type|private$, err := clientset.CoreV1().$.type|publicPlural$().Get($.type|private$Obj.Name, metav1.GetOptions{})
    if err != nil {
        klog.Errorf("get $.type|private$ %v failed with:%v", $.type|private$Obj.Name, err)
        return err
    }

    $.type|private$, err = clientset.CoreV1().$.type|publicPlural$().Update($.type|private$)
    if err != nil {
        klog.Errorf("update $.type|private$ failed with:%v", err)
        return err
    }
    return nil
}
`

var deleteObjectService = `
// Delete$.type|public$ xxx
// TODO(user): Modify this function to implement your logic.This example use namespace.
func (s *service) Delete$.type|public$(ctx context.Context, name string) error {
    clientset := s.opt.KubeClientset

    _, err := clientset.CoreV1().$.type|publicPlural$().Get(name, metav1.GetOptions{})
    if err != nil {
        klog.Errorf("get $.type|private$ %v failed with:%v", name, err)
        return err
    }

    err = clientset.CoreV1().$.type|publicPlural$().Delete(name, &metav1.DeleteOptions{})
    if err != nil {
        klog.Errorf("delete $.type|private$Obj %v failed with:%v", name, err)
        return err
    }
    return nil
}
`
