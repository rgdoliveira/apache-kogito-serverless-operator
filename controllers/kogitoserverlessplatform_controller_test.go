// Copyright 2023 Red Hat, Inc. and/or its affiliates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controllers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/kiegroup/kogito-serverless-operator/test"

	"github.com/kiegroup/kogito-serverless-operator/api/v1alpha08"
)

func TestKogitoServerlessPlatformController(t *testing.T) {
	t.Run("verify that a basic reconcile is performed without error", func(t *testing.T) {
		// Create a KogitoServerlessPlatform object with metadata and spec.
		ksp := test.GetKogitoServerlessPlatform("../config/samples/sw.kogito_v1alpha08_kogitoserverlessplatform.yaml")

		// Create a fake client to mock API calls.
		cl := test.NewKogitoClientBuilder().WithRuntimeObjects(ksp).Build()
		// Create a KogitoServerlessPlatformReconciler object with the scheme and fake client.
		r := &KogitoServerlessPlatformReconciler{cl, cl, cl.Scheme(), &rest.Config{}, &record.FakeRecorder{}}

		// Mock request to simulate Reconcile() being called on an event for a
		// watched resource .
		req := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      ksp.Name,
				Namespace: ksp.Namespace,
			},
		}
		_, err := r.Reconcile(context.TODO(), req)
		if err != nil {
			t.Fatalf("reconcile: (%v)", err)
		}

		assert.NoError(t, cl.Get(context.TODO(), types.NamespacedName{Name: ksp.Name, Namespace: ksp.Namespace}, ksp))

		// Perform some checks on the created CR
		assert.Equal(t, "quay.io/kiegroup", ksp.Spec.BuildPlatform.Registry.Address)
		assert.Equal(t, "regcred", ksp.Spec.BuildPlatform.Registry.Secret)
		assert.Equal(t, v1alpha08.OperatorBuildStrategy, ksp.Spec.BuildPlatform.BuildStrategy)
		assert.Equal(t, v1alpha08.PlatformClusterKubernetes, ksp.Status.Cluster)

		assert.Equal(t, v1alpha08.PlatformPhaseCreating, ksp.Status.Phase)
	})
}
