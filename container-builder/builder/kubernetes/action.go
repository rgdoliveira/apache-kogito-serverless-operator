// Copyright 2022 Red Hat, Inc. and/or its affiliates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubernetes

import (
	"context"

	"github.com/kiegroup/kogito-serverless-operator/container-builder/api"
	"github.com/kiegroup/kogito-serverless-operator/container-builder/client"
	"github.com/kiegroup/kogito-serverless-operator/container-builder/util/log"
)

type Action interface {
	log.Injectable
	client.Injectable
	// Name returns user-friendly name for the action
	Name() string

	// CanHandle returns true if the action can handle the build
	CanHandle(build *api.ContainerBuild) bool

	// Handle executes the handling function
	Handle(ctx context.Context, build *api.ContainerBuild) (*api.ContainerBuild, error)
}

type baseAction struct {
	client client.Client
	L      log.Logger
}

// TODO: implement our client wrapper

func (action *baseAction) InjectClient(client client.Client) {
	action.client = client
}

func (action *baseAction) InjectLogger(log log.Logger) {
	action.L = log
}
