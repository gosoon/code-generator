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

	"github.com/gosoon/test/pkg/types"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

// CreateSecret xxx
// TODO(user): Modify this function to implement your logic.This example use namespace.
func (s *service) CreateSecret(ctx context.Context, secretObj *types.Secret) error {
	clientset := s.opt.KubeClientset
	secret := &apiv1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: secretObj.Name,
		},
	}

	_, err := clientset.CoreV1().Secrets().Create(namespace)
	if err != nil {
		klog.Errorf("create secret failed with:%v", err)
		return err
	}
	return nil
}

// GetSecret xxx
// TODO(user): Modify this function to implement your logic.This example use namespace.
func (s *service) GetSecret(ctx context.Context, name string) (*apiv1.Secret, error) {
	clientset := s.opt.KubeClientset

	secret, err := clientset.CoreV1().Secrets().Get(name, metav1.GetOptions{})
	if err != nil {
		klog.Errorf("get secret %v failed with:%v", name, err)
		return nil, err
	}

	return secret, nil
}

// UpdateSecret xxx
// TODO(user): Modify this function to implement your logic.This example use namespace.
func (s *service) UpdateSecret(ctx context.Context, secretObj *types.Secret) error {
	clientset := s.opt.KubeClientset

	var err error
	secret, err := clientset.CoreV1().Secrets().Get(secretObj.Name, metav1.GetOptions{})
	if err != nil {
		klog.Errorf("get secret %v failed with:%v", secretObj.Name, err)
		return err
	}

	secret, err = clientset.CoreV1().Secrets().Update(secret)
	if err != nil {
		klog.Errorf("update secret failed with:%v", err)
		return err
	}
	return nil
}

// DeleteSecret xxx
// TODO(user): Modify this function to implement your logic.This example use namespace.
func (s *service) DeleteSecret(ctx context.Context, name string) error {
	clientset := s.opt.KubeClientset

	_, err := clientset.CoreV1().Secrets().Get(name, metav1.GetOptions{})
	if err != nil {
		klog.Errorf("get secret %v failed with:%v", name, err)
		return err
	}

	err = clientset.CoreV1().Secrets().Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		klog.Errorf("delete secretObj %v failed with:%v", name, err)
		return err
	}
	return nil
}
