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
	"context"

	"github.com/gosoon/code-generator/_examples/types/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

// CreateNamespace xxx
// TODO(user): Modify this function to implement your logic.This example use namespace.
func (s *service) CreateNamespace(ctx context.Context, namespaceObj *types.Namespace) error {
	clientset := s.opt.KubeClientset
	namespace := &apiv1.Namespace{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Namespace",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: namespaceObj.Name,
		},
	}

	_, err := clientset.CoreV1().Namespaces().Create(namespace)
	if err != nil {
		klog.Errorf("create namespace failed with:%v", err)
		return err
	}
	return nil
}

// GetNamespace xxx
// TODO(user): Modify this function to implement your logic.This example use namespace.
func (s *service) GetNamespace(ctx context.Context, name string) (*apiv1.Namespace, error) {
	clientset := s.opt.KubeClientset

	namespace, err := clientset.CoreV1().Namespaces().Get(name, metav1.GetOptions{})
	if err != nil {
		klog.Errorf("get namespace %v failed with:%v", name, err)
		return nil, err
	}

	return namespace, nil
}

// UpdateNamespace xxx
// TODO(user): Modify this function to implement your logic.This example use namespace.
func (s *service) UpdateNamespace(ctx context.Context, namespaceObj *types.Namespace) error {
	clientset := s.opt.KubeClientset

	var err error
	namespace, err := clientset.CoreV1().Namespaces().Get(namespaceObj.Name, metav1.GetOptions{})
	if err != nil {
		klog.Errorf("get namespace %v failed with:%v", namespaceObj.Name, err)
		return err
	}

	namespace, err = clientset.CoreV1().Namespaces().Update(namespace)
	if err != nil {
		klog.Errorf("update namespace failed with:%v", err)
		return err
	}
	return nil
}

// DeleteNamespace xxx
// TODO(user): Modify this function to implement your logic.This example use namespace.
func (s *service) DeleteNamespace(ctx context.Context, name string) error {
	clientset := s.opt.KubeClientset

	_, err := clientset.CoreV1().Namespaces().Get(name, metav1.GetOptions{})
	if err != nil {
		klog.Errorf("get namespace %v failed with:%v", name, err)
		return err
	}

	err = clientset.CoreV1().Namespaces().Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		klog.Errorf("delete namespaceObj %v failed with:%v", name, err)
		return err
	}
	return nil
}
