package service

import (
	"io"

	clientgentypes "k8s.io/code-generator/cmd/client-gen/types"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/klog"
)

// genTypesService generates a package for a controller.
type genTypesService struct {
	generator.DefaultGen
	groups           []clientgentypes.GroupVersions
	groupGoNames     map[clientgentypes.GroupVersion]string
	clientsetPackage string
	inputPackages    []string
	outputPackage    string
	imports          namer.ImportTracker
	serviceGenerated bool

	typeToGenerate *types.Type
	objectMeta     *types.Type
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
	return
}

func (g *genTypesService) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	sw := generator.NewSnippetWriter(w, c, "$", "$")

	klog.Infof("processing type %v", t)
	m := map[string]interface{}{
		"type": t,
	}

	sw.Do(createObjectService, m)
	sw.Do(getObjectService, m)
	sw.Do(updateObjectService, m)
	sw.Do(deleteObjectService, m)
	return sw.Error()
}

var createObjectService = `
// Create$.type|public$ xxx
func (s *service) Create$.type|public$(tenant *types.Tenant) error {
    clientset := s.opt.KubeClientset
    namespace := &apiv1.Namespace{
        TypeMeta: metav1.TypeMeta{
            APIVersion: "v1",
            Kind:       "Namespace",
        },
        ObjectMeta: metav1.ObjectMeta{
            Name: tenant.Name,
        },
    }

   /*
    _, err := clientset.CoreV1().Namespaces().Create(namespace)
    if err != nil {
        clog.Errorf("create tenant failed with:%v", err)
        return err
    } */
    return nil
}
`
var getObjectService = `
// Get$.type|public$ xxx
func (s *service) Get$.type|public$(name string) (*apiv1.Namespace, error) {
    clientset := s.opt.KubeClient

    namespace, err := clientset.CoreV1().Namespaces().Get(name, metav1.GetOptions{})
    if err != nil {
        clog.Errorf("get tenant %v failed with:%v", name, err)
        return nil, err
    }

    if _, existed := namespace.Annotations[enum.TenantAnnotation]; !existed {
        return nil, fmt.Errorf("not found tenant %v", name)
    }
    return namespace, nil
}
`

var updateObjectService = `
func (s *service) Update$.type|public$(tenant *types.Tenant) (*apiv1.Namespace, error) {
    clientset := s.opt.KubeClient

    namespace, err := clientset.CoreV1().Namespaces().Get(tenant.Name, metav1.GetOptions{})
    if err != nil {
        clog.Errorf("get tenant %v failed with:%v", name, err)
        return err
    }

    curNamespace, err := clientset.CoreV1().Namespaces().Update(namespace)
    if err != nil {
        clog.Errorf("update tenant failed with:%v", err)
        return nil, err
    }
    return curNamespace, nil
}
`

var deleteObjectService = `
func (s *service) Delete$.type|public$(name string) error {
    clientset := s.opt.KubeClient

    namespace, err := clientset.CoreV1().Namespaces().Get(name, metav1.GetOptions{})
    if err != nil {
        clog.Errorf("get tenant %v failed with:%v", name, err)
        return err
    }

    // Delete(name string, options *metav1.DeleteOptions) error
    err = clientset.CoreV1().Namespaces().Delete(name, &metav1.DeleteOptions{})
    if err != nil {
        clog.Errorf("delete tenant %v failed with:%v", name, err)
        return err
    }
    return nil
}
`
