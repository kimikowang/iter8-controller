/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package util

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	iter8v1alpha2 "github.com/iter8-tools/iter8/pkg/apis/iter8/v1alpha2"
)

// TestServiceToFullHostName tests function ServiceToFullHostName
func TestServiceToFullHostName(t *testing.T) {
	svcName, namespace := "testName", "testNamespace"
	fullHostName := "testName.testNamespace.svc.cluster.local"

	res := ServiceToFullHostName(svcName, namespace)
	if res != fullHostName {
		t.Errorf("unexpected reponse, want %s, got %s", fullHostName, res)
	}
}

// TestFullExperimentName tests function FullExperimentName
func TestFullExperimentName(t *testing.T) {
	expected := "test-name.test-namespace"
	exp := &iter8v1alpha2.Experiment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-name",
			Namespace: "test-namespace",
		},
	}

	res := FullExperimentName(exp)
	if res != expected {
		t.Errorf("unexpected reponse, want %s, got %s", expected, res)
	}
}

// TestGetDefaultHost tests function GetDefaultHost
func TestGetDefaultHost(t *testing.T) {
	exp := &iter8v1alpha2.Experiment{
		Spec: iter8v1alpha2.ExperimentSpec{
			Service: iter8v1alpha2.Service{
				ObjectReference: &corev1.ObjectReference{
					Name:      "test-service",
					Namespace: "test-service-namespace",
				},
			},
		},
	}

	expected := "test-service.test-service-namespace.svc.cluster.local"
	res := GetDefaultHost(exp)
	if res != expected {
		t.Errorf("unexpected reponse, want %s, got %s", expected, res)
	}

	exp = &iter8v1alpha2.Experiment{
		Spec: iter8v1alpha2.ExperimentSpec{
			Networking: &iter8v1alpha2.Networking{
				Hosts: []iter8v1alpha2.Host{
					{
						Name:    "test-host",
						Gateway: "test-gateway",
					},
				},
			},
		},
	}

	expected = "test-host"
	res = GetDefaultHost(exp)
	if res != expected {
		t.Errorf("unexpected reponse, want %s, got %s", expected, res)
	}

	exp = &iter8v1alpha2.Experiment{
		Spec: iter8v1alpha2.ExperimentSpec{
			Service: iter8v1alpha2.Service{
				ObjectReference: &corev1.ObjectReference{
					Name:      "test-service",
					Namespace: "test-service-namespace",
				},
			},
			Networking: &iter8v1alpha2.Networking{
				Hosts: []iter8v1alpha2.Host{
					{
						Name:    "test-host",
						Gateway: "test-gateway",
					},
				},
			},
		},
	}

	expected = "test-service.test-service-namespace.svc.cluster.local"
	res = GetDefaultHost(exp)
	if res != expected {
		t.Errorf("unexpected reponse, want %s, got %s", expected, res)
	}
}
