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
	"strings"
	"testing"

	iter8v1alpha2 "github.com/iter8-tools/iter8/pkg/apis/iter8/v1alpha2"
)

func TestExperimentKey(t *testing.T) {
	exp := &iter8v1alpha2.Experiment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-name",
			Namespace: "test-namespace",
		},
	}
	expected := "test-namespace/test-name"
	res := experimentKey(exp)

	if res != expected {
		t.Errorf("unexpected resule, want %s; got %s", expected, res)
	}
}

func TestTargetKey(t *testing.T) {
	ns, name := "test-namespace", "test-name"
	expected := "test-namespace/test-name"
	res := targetKey(name, namespace)

	if res != expected {
		t.Errorf("unexpected resule, want %s; got %s", expected, res)
	}
}

// TestResolveExperimentKey tests function resolveExperimentKey
func TestResolveExperimentKey(t *testing.T) {
	ns, name := "test-namespace", "test-name"
	key := "test-namespace/test-name"
	res := resolveExperimentKey(key)

	if res[0] != ns || res[1] != name {
		t.Errorf("unexpected resule, want %s, %s; got %s, %s", ns, name, res[0], res[1])
	}
}
