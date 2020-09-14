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

package adapter

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	iter8v1alpha2 "github.com/iter8-tools/iter8/pkg/apis/iter8/v1alpha2"
)

func TestAdapter(t *testing.T) {
	ctx := context.Background()
	testAdapter := New(nil)

	exp := &iter8v1alpha2.Experiment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-name",
			Namespace: "test-namespace",
		},
		Spec: iter8v1alpha2.ExperimentSpec{
			Service: iter8v1alpha2.Service{
				Name:     "test-service",
				Baseline: "test-baseline",
				Candidates: []string{
					"test-candidate-0",
					"test-candidate-1",
				},
			},
		},
	}

	ctx, err := testAdapter.RegisterExperiment(ctx, exp)
}
