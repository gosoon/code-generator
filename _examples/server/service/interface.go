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
	"k8s.io/client-go/kubernetes"
)

// Options contains the config by service
type Options struct {
	KubeClientset kubernetes.Interface
}

// service implements the Service interface.
type service struct {
	opt *Options
}

// New is create a service object.
func New(opt *Options) Interface {
	return &service{opt: opt}
}

// Interface is definition service all method.
type Interface interface {
	CreateNamespace(ctx context.Context, namespaceObj *types.Namespace) error
	GetNamespace(ctx context.Context, name string) (*apiv1.Namespace, error)
	UpdateNamespace(ctx context.Context, namespaceObj *types.Namespace) error
	DeleteNamespace(ctx context.Context, name string) error

	CreateSecret(ctx context.Context, secretObj *types.Secret) error
	GetSecret(ctx context.Context, name string) (*apiv1.Secret, error)
	UpdateSecret(ctx context.Context, secretObj *types.Secret) error
	DeleteSecret(ctx context.Context, name string) error
}
