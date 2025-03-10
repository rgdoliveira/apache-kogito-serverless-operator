/*
 * Copyright 2022 Red Hat, Inc. and/or its affiliates.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"os"
	"time"

	builder "github.com/kiegroup/kogito-serverless-operator/container-builder/builder/kubernetes"

	v1 "k8s.io/api/core/v1"
	resource2 "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kiegroup/kogito-serverless-operator/container-builder/api"
	"github.com/kiegroup/kogito-serverless-operator/container-builder/client"
)

/*
Usage example. Please note that you must have a valid Kubernetes environment up and running.
*/

func main() {
	cli, err := client.NewOutOfClusterClient("")

	dockerFile, err := os.ReadFile("examples/dockerfiles/Kogito.dockerfile")
	if err != nil {
		panic("Can't read dockerfile")
	}
	source, err := os.ReadFile("examples/sources/kogitogreetings.sw.json")
	if err != nil {
		panic("Can't read source file")
	}

	if err != nil {
		fmt.Println("Failed to create client")
		fmt.Println(err.Error())
	}
	platform := api.PlatformContainerBuild{
		ObjectReference: api.ObjectReference{
			Namespace: "kogito-builder",
			Name:      "testPlatform",
		},
		Spec: api.PlatformContainerBuildSpec{
			BuildStrategy:   api.ContainerBuildStrategyPod,
			PublishStrategy: api.PlatformBuildPublishStrategyKaniko,
			Registry: api.ContainerRegistrySpec{
				Insecure: true,
			},
			Timeout: &metav1.Duration{
				Duration: 5 * time.Minute,
			},
		},
	}

	cpuQty, _ := resource2.ParseQuantity("1")
	memQty, _ := resource2.ParseQuantity("4Gi")

	build, err := builder.NewBuild(builder.ContainerBuilderInfo{FinalImageName: "greetings:latest", BuildUniqueName: "kogito-test", Platform: platform}).
		WithResource("Dockerfile", dockerFile).WithResource("greetings.sw.json", source).
		WithAdditionalArgs([]string{"--build-arg=QUARKUS_PACKAGE_TYPE=mutable-jar", "--build-arg=QUARKUS_LAUNCH_DEVMODE=true", "--build-arg=SCRIPT_DEBUG=false"}).
		WithResourceRequirements(v1.ResourceRequirements{
			Limits: v1.ResourceList{
				v1.ResourceCPU:    cpuQty,
				v1.ResourceMemory: memQty,
			},
			Requests: v1.ResourceList{
				v1.ResourceCPU:    cpuQty,
				v1.ResourceMemory: memQty,
			},
		}).
		WithClient(cli).
		Schedule()
	if err != nil {
		fmt.Println(err.Error())
		panic("Can't create build")
	}

	// from now the Reconcile method can be called until the build is finished
	for build.Status.Phase != api.ContainerBuildPhaseSucceeded &&
		build.Status.Phase != api.ContainerBuildPhaseError &&
		build.Status.Phase != api.ContainerBuildPhaseFailed {
		fmt.Printf("\nBuild status is %s", build.Status.Phase)
		build, err = builder.FromBuild(build).WithClient(cli).Reconcile()
		if err != nil {
			fmt.Println("Failed to run test")
			panic(fmt.Errorf("build %v just failed", build))
		}
		time.Sleep(10 * time.Second)
	}

}
