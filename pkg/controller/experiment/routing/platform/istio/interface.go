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

package istio

// Processor processes routing rules based on target runtime objects

import (
	iter8v1alpha2 "github.com/iter8-tools/iter8-controller/pkg/apis/iter8/v1alpha2"
	"github.com/iter8-tools/iter8-controller/pkg/controller/experiment/targets"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
)

const (
	// subset names used in destianation rule
	subsetBaseline  = "iter8-baseline"
	subsetCandidate = "iter8-candidate"
	subsetStable    = "iter8-stable"

	// labels used in routing rules
	experimentInit  = "iter8-tools/init"
	experimentRole  = "iter8-tools/role"
	experimentLabel = "iter8-tools/experiment"
	experimentHost  = "iter8-tools/host"

	// values for experimentRole
	roleInitializing = "initializing"
	roleStable       = "stable"
	roleProgressing  = "progressing"
)

// Handler interface defines functions that a handler should implement
type Handler interface {
	// Validate detected routing lists from cluster to see if they are valid for next experiment
	// Init routing rules after if validation passes
	ValidateAndInit(drl *v1alpha3.DestinationRuleList, vsl *v1alpha3.VirtualServiceList, instance *iter8v1alpha1.Experiment) error
}
