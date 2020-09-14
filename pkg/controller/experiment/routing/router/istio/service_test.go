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

import (
	"fmt"
	"testing"

	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"

	iter8v1alpha2 "github.com/iter8-tools/iter8/pkg/apis/iter8/v1alpha2"
	"github.com/iter8-tools/iter8/pkg/controller/experiment/util"
)

// Logger creates logger for the test
func Logger(t *testing.T) logr.Logger {
	log.SetLogger(log.ZapLogger(false))
	return log.Log.WithName(t.Name())
}

func TestServiceHandlerValidation(t *testing.T) {
	logger := Logger(t)

	drl := &v1alpha3.DestinationRuleList{
		Items: []*networkingv1alpha3.DestinationRule{

		}
	}
}

func TestServiceHandlerBuildDestination(t *testing.T) {
	logger := Logger(t)
	testPort := 9080
	testWeight := 50

	opts := destinationOptions{
		name: "test-name",
		subset: "test-subset",
		weight: testWeight,
		port: &port,
	}

	exp := &iter8v1alpha2.Experiment{
		Namespace: "test-namespace",
	}

	destination := &networkingv1alpha3.HTTPRouteDestination{
		Weight: 50,
		Destination: &networkingv1alpha3.Destination{
			Host: "test-name.test-namespace.svc.cluster.local"
			Port: &networkingv1alpha3.PortSelector{
				Number: testPort,
			}
		}
	}
}