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

package processor

import (
	"istio.io/client-go/pkg/apis/networking/v1alpha3"

	iter8v1alpha2 "github.com/iter8-tools/iter8-controller/pkg/apis/iter8/v1alpha2"
	"github.com/iter8-tools/iter8-controller/pkg/controller/experiment/targets"
	"github.com/iter8-tools/iter8-controller/pkg/controller/experiment/util"
)

// Deployment implements processor Interface and manipulates routing rules based on deployment kind of targets
type Deployment struct {
	destinationRule *v1alpha3.DestinationRule
	virtualService  *v1alpha3.VirtualService
}

func New(dr *v1alpha3.DestinationRule, vs *v1alpha3.VirtualService) *Deployment {
	return &Deployment{
		destinationRule: dr,
		virtualService:  vs,
	}
}

func (p *Deployment) ValidateAndInit(drl *v1alpha3.DestinationRuleList, vsl *v1alpha3.VirtualServiceList, instance *iter8v1alpha1.Experiment) error {
	svcNamespace := instance.ServiceNamespace()
	host := util.GetHost(instance)
	expFullName := util.FullExperimentName(instance)
	if len(drl.Items) == 0 && len(vsl.Items) == 0 {
		// init rule
		p.destinationRule = NewDestinationRule(host, expFullName, svcNamespace).
			WithInitLabel().
			Build()
		p.virtualService = NewVirtualService(host, expFullName, svcNamespace).
			WithInitLabel().
			Build()
	} else if len(drl.Items) == 1 && len(vsl.Items) == 1 {
		drrole, drok := drl.Items[0].GetLabels()[experimentRole]
		vsrole, vsok := vsl.Items[0].GetLabels()[experimentRole]
		if drok && vsok {
			if drrole == Stable && vsrole == Stable {
				// Valid stable rules detected
				p.destinationRule = drl.Items[0].DeepCopy()
				p.virtualService = vsl.Items[0].DeepCopy()
			} else {
				drLabel, drok := drl.Items[0].GetLabels()[experimentLabel]
				vsLabel, vsok := vsl.Items[0].GetLabels()[experimentLabel]
				if drok && vsok {
					if drLabel == expFullName && vsLabel == expFullName {
						// valid progressing rules found
						p.destinationRule = drl.Items[0].DeepCopy()
						p.virtualService = vsl.Items[0].DeepCopy()
					} else {
						return fmt.Errorf("Progressing rules being involved in other experiments")
					}
				} else {
					return fmt.Errorf("Experiment label missing in dr or vs")
				}
			}
		} else {
			return fmt.Errorf("experiment role label missing in dr or vs")
		}
	} else if len(drl.Items) == 0 && len(vsl.Items) == 1 {
		vsrole, vsok := vsl.Items[0].GetLabels()[experimentRole]
		if vsok && vsrole == Stable {
			// Valid stable rules detected
			p.virtualService = vsl.Items[0].DeepCopy()
			p.destinationRule = NewDestinationRule(host, expFullName, svcNamespace).
				WithInitLabel().
				Build()
		} else {
			return fmt.Errorf("0 dr and 1 unstable vs found")
		}
	} else {
		return fmt.Errorf("%d dr and %d vs detected", len(drl.Items), len(vsl.Items))
	}

	return nil
}
