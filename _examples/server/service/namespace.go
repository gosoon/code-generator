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
	"fmt"

	"git.pacloud.io/pks/tenant-service/pkg/enum"
	"git.pacloud.io/pks/tenant-service/pkg/types"
)

// CreateNamespace xxx
func (s *service) CreateNamespace(tenant *types.Tenant) error {
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

// GetNamespace xxx
func (s *service) GetNamespace(name string) (*apiv1.Namespace, error) {
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

func (s *service) UpdateNamespace(tenant *types.Tenant) (*apiv1.Namespace, error) {
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

func (s *service) DeleteNamespace(name string) error {
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
