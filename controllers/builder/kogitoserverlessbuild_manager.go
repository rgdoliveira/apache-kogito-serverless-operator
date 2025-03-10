// Copyright 2023 Red Hat, Inc. and/or its affiliates
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

package builder

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/kiegroup/kogito-serverless-operator/controllers/platform"

	operatorapi "github.com/kiegroup/kogito-serverless-operator/api/v1alpha08"
)

var _ KogitoServerlessBuildManager = &kogitoServerlessBuildManager{}

type kogitoServerlessBuildManager struct {
	client client.Client
	ctx    context.Context
}

func (k *kogitoServerlessBuildManager) MarkToRestart(build *operatorapi.KogitoServerlessBuild) error {
	build.Status.BuildPhase = operatorapi.BuildPhaseNone
	return k.client.Status().Update(k.ctx, build)
}

func (k *kogitoServerlessBuildManager) GetOrCreateBuild(workflow *operatorapi.KogitoServerlessWorkflow) (*operatorapi.KogitoServerlessBuild, error) {
	buildInstance := &operatorapi.KogitoServerlessBuild{}
	buildInstance.ObjectMeta.Namespace = workflow.Namespace
	buildInstance.ObjectMeta.Name = workflow.Name

	if err := k.client.Get(k.ctx, client.ObjectKeyFromObject(workflow), buildInstance); err != nil {
		if errors.IsNotFound(err) {
			plat := &operatorapi.KogitoServerlessPlatform{}
			if plat, err = platform.GetActivePlatform(k.ctx, k.client, workflow.Namespace); err != nil {
				return nil, err
			}
			buildInstance.Spec.BuildTemplate = plat.Spec.BuildTemplate
			if err = controllerutil.SetControllerReference(workflow, buildInstance, k.client.Scheme()); err != nil {
				return nil, err
			}
			if err = k.client.Create(k.ctx, buildInstance); err != nil {
				return nil, err
			}
			return buildInstance, nil
		}
		return nil, err
	}

	return buildInstance, nil
}

type KogitoServerlessBuildManager interface {
	// GetOrCreateBuild gets or creates a new instance of KogitoServerlessBuild for the given KogitoServerlessWorkflow.
	//
	// Only one build is allowed per workflow instance
	GetOrCreateBuild(workflow *operatorapi.KogitoServerlessWorkflow) (*operatorapi.KogitoServerlessBuild, error)
	// MarkToRestart tell the controller to restart this build in the next iteration
	MarkToRestart(build *operatorapi.KogitoServerlessBuild) error
}

// NewKogitoServerlessBuildManager entry point to manage KogitoServerlessBuild instances.
// Won't start a build, but once it creates a new instance, the controller will take place and start the build in the cluster context.
func NewKogitoServerlessBuildManager(ctx context.Context, client client.Client) KogitoServerlessBuildManager {
	return &kogitoServerlessBuildManager{
		client: client,
		ctx:    ctx,
	}
}
