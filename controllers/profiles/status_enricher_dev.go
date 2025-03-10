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

package profiles

import (
	"context"
	"fmt"

	openshiftv1 "github.com/openshift/api/route/v1"
	"knative.dev/pkg/apis"

	"github.com/kiegroup/kogito-serverless-operator/controllers/workflowdef"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	operatorapi "github.com/kiegroup/kogito-serverless-operator/api/v1alpha08"
)

func defaultDevStatusEnricher(ctx context.Context, c client.Client, workflow *operatorapi.KogitoServerlessWorkflow) (client.Object, error) {
	//If the workflow Status hasn't got a NodePort Endpoint, we are ensuring it will be set
	// If we aren't on OpenShift we will enrich the status with 2 info:
	// - Address the service can be reached
	// - Node port used
	service := &v1.Service{}
	err := c.Get(ctx, types.NamespacedName{Namespace: workflow.Namespace, Name: workflow.Name}, service)
	if err != nil {
		return nil, err
	}
	//If the service has got a Port that is a nodePort we have to use it to create the workflow's NodePort Endpoint
	if service.Spec.Ports != nil && len(service.Spec.Ports) > 0 {
		if port := findNodePortFromPorts(service.Spec.Ports); port > 0 {
			labels := workflowdef.GetDefaultLabels(workflow)

			podList := &v1.PodList{}
			opts := []client.ListOption{
				client.InNamespace(workflow.Namespace),
				client.MatchingLabels{workflowdef.LabelApp: labels[workflowdef.LabelApp]},
			}
			err := c.List(ctx, podList, opts...)
			if err != nil {
				return nil, err
			}
			var ipaddr string
			for _, p := range podList.Items {
				ipaddr = p.Status.HostIP
				break
			}

			url, err := apis.ParseURL("http://" + ipaddr + ":" + fmt.Sprint(port) + "/" + workflow.Name)
			if err != nil {
				return nil, err
			}
			workflow.Status.Endpoint = url
		}
	}

	return workflow, nil
}

func devStatusEnricherForOpenShift(ctx context.Context, client client.Client, workflow *operatorapi.KogitoServerlessWorkflow) (client.Object, error) {
	// On OpenShift we need to retrieve the Route to have the URL the service is available to
	route := &openshiftv1.Route{}
	err := client.Get(ctx, types.NamespacedName{Namespace: workflow.Namespace, Name: workflow.Name}, route)
	if err != nil {
		return nil, err
	}
	var url *apis.URL
	if route.Spec.TLS != nil {
		url = apis.HTTPS(route.Spec.Host)
	} else {
		url = apis.HTTP(route.Spec.Host)
	}

	workflow.Status.Endpoint = url

	return workflow, nil
}
