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
package e2e

import (
	"testing"

	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"

	analtyicsapi "github.com/iter8-tools/iter8-controller/pkg/analytics/api/v1alpha2"
	iter8v1alpha2 "github.com/iter8-tools/iter8-controller/pkg/apis/iter8/v1alpha2"
	"github.com/iter8-tools/iter8-controller/pkg/controller/experiment/routing/router/istio"
	"github.com/iter8-tools/iter8-controller/pkg/controller/experiment/util"
	"github.com/iter8-tools/iter8-controller/test"
)

const (
	//Bookinfo sample constant values
	ReviewsV1Image = "istio/examples-bookinfo-reviews-v1:1.11.0"
	ReviewsV2Image = "istio/examples-bookinfo-reviews-v2:1.11.0"
	ReviewsV3Image = "istio/examples-bookinfo-reviews-v3:1.11.0"

	ReviewsPort = 9080
)

// TestKubernetesExperiment tests various experiment scenarios on Kubernetes platform
func TestExperiment(t *testing.T) {
	service := test.StartAnalytics()
	defer service.Close()
	testCases := map[string]testCase{
		"rolltowinner": func(name string, exp *iter8v1alpha2.Experiment) testCase {
			return testCase{
				mocks: map[string]analtyicsapi.Response{
					name: test.GetRollToWinnerMockResponse(exp, 0),
				},
				initObjects: []runtime.Object{
					getReviewsService(),
					getReviewsDeployment("v1"),
					getReviewsDeployment("v2"),
					getReviewsDeployment("v3"),
				},
				object: exp,
				wantState: test.WantAllStates(
					test.CheckExperimentCompleted,
				),
				wantResults: []runtime.Object{
					getDestinationRule("reviews", name,
						[]string{istio.SubsetBaseline, istio.CandidateSubsetName(0), istio.CandidateSubsetName(1)},
						[]runtime.Object{getReviewsDeployment("v1"), getReviewsDeployment("v2"), getReviewsDeployment("v3")},
					),
					getVirtualServiceForDeployments("reviews", name,
						[]string{istio.SubsetBaseline, istio.CandidateSubsetName(0), istio.CandidateSubsetName(1)},
						[]int32{0, 100, 0},
					),
				},
			}
		}("rolltowinner", getFastKubernetesExperiment("rolltowinner", "reviews", "reviews-v1", service.GetURL(), []string{"reviews-v2", "reviews-v3"})),
		"rollbackward": func(name string, exp *iter8v1alpha2.Experiment) testCase {
			return testCase{
				mocks: map[string]analtyicsapi.Response{
					name: test.GetRollbackMockResponse(exp),
				},
				initObjects: []runtime.Object{
					getReviewsService(),
					getReviewsDeployment("v1"),
					getReviewsDeployment("v2"),
					getReviewsDeployment("v3"),
				},
				object: exp,
				wantState: test.WantAllStates(
					test.CheckExperimentCompleted,
				),
				wantResults: []runtime.Object{
					getDestinationRule("reviews", name,
						[]string{istio.SubsetBaseline, istio.CandidateSubsetName(0), istio.CandidateSubsetName(1)},
						[]runtime.Object{getReviewsDeployment("v1"), getReviewsDeployment("v2"), getReviewsDeployment("v3")},
					),
					getVirtualServiceForDeployments("reviews", name,
						[]string{istio.SubsetBaseline, istio.CandidateSubsetName(0), istio.CandidateSubsetName(1)},
						[]int32{100, 0, 0},
					),
				},
			}
		}("rollbackward", getFastKubernetesExperiment("rollbackward", "reviews", "reviews-v1", service.GetURL(), []string{"reviews-v2", "reviews-v3"})),
		"ongoingdelete": func(name string, exp *iter8v1alpha2.Experiment) testCase {
			return testCase{
				mocks: map[string]analtyicsapi.Response{
					name: test.GetRollToWinnerMockResponse(exp, 0),
				},
				initObjects: []runtime.Object{
					getReviewsService(),
					getReviewsDeployment("v1"),
					getReviewsDeployment("v2"),
					getReviewsDeployment("v3"),
				},
				object:    exp,
				wantState: test.CheckServiceFound,
				wantResults: []runtime.Object{
					getDestinationRule("reviews", name,
						[]string{istio.SubsetBaseline, istio.CandidateSubsetName(0), istio.CandidateSubsetName(1)},
						[]runtime.Object{getReviewsDeployment("v1"), getReviewsDeployment("v2"), getReviewsDeployment("v3")},
					),
					getVirtualServiceForDeployments("reviews", name,
						[]string{istio.SubsetBaseline, istio.CandidateSubsetName(0), istio.CandidateSubsetName(1)},
						[]int32{100, 0, 0},
					),
				},
				postHook: test.DeleteExperiment("ongoingdelete", Flags.Namespace),
			}
		}("ongoingdelete", getSlowKubernetesExperiment("ongoingdelete", "reviews", "reviews-v1", service.GetURL(), []string{"reviews-v2", "reviews-v3"})),
		"completedelete": func(name string, exp *iter8v1alpha2.Experiment) testCase {
			return testCase{
				mocks: map[string]analtyicsapi.Response{
					name: test.GetRollToWinnerMockResponse(exp, 0),
				},
				initObjects: []runtime.Object{
					getReviewsService(),
					getReviewsDeployment("v1"),
					getReviewsDeployment("v2"),
					getReviewsDeployment("v3"),
				},
				object:    exp,
				wantState: test.CheckExperimentCompleted,
				frozenObjects: []runtime.Object{
					getDestinationRule("reviews", name,
						[]string{istio.SubsetBaseline, istio.CandidateSubsetName(0), istio.CandidateSubsetName(1)},
						[]runtime.Object{getReviewsDeployment("v1"), getReviewsDeployment("v2"), getReviewsDeployment("v3")},
					),
					getVirtualServiceForDeployments("reviews", name,
						[]string{istio.SubsetBaseline, istio.CandidateSubsetName(0), istio.CandidateSubsetName(1)},
						[]int32{100, 0, 0},
					),
				},
				postHook: test.DeleteExperiment(name, Flags.Namespace),
			}
		}("completedelete", getFastKubernetesExperiment("completedelete", "reviews", "reviews-v1", service.GetURL(), []string{"reviews-v2", "reviews-v3"})),
		"abortexperiment": func(name string, exp *iter8v1alpha2.Experiment) testCase {
			return testCase{
				mocks: map[string]analtyicsapi.Response{
					name: test.GetAbortExperimentResponse(exp),
				},
				initObjects: []runtime.Object{
					getReviewsService(),
					getReviewsDeployment("v1"),
					getReviewsDeployment("v2"),
					getReviewsDeployment("v3"),
				},
				object: exp,
				wantState: test.WantAllStates(
					test.CheckExperimentCompleted,
				),
				wantResults: []runtime.Object{
					getDestinationRule("reviews", name,
						[]string{istio.SubsetBaseline, istio.CandidateSubsetName(0), istio.CandidateSubsetName(1)},
						[]runtime.Object{getReviewsDeployment("v1"), getReviewsDeployment("v2"), getReviewsDeployment("v3")},
					),
					getVirtualServiceForDeployments("reviews", name,
						[]string{istio.SubsetBaseline, istio.CandidateSubsetName(0), istio.CandidateSubsetName(1)},
						[]int32{100, 0, 0},
					),
				},
			}
		}("abortexperiment", getSlowKubernetesExperiment("abortexperiment", "reviews", "reviews-v1", service.GetURL(), []string{"reviews-v2", "reviews-v3"})),
		"emptycriterion": func(name string, exp *iter8v1alpha2.Experiment) testCase {
			return testCase{
				initObjects: []runtime.Object{
					getReviewsService(),
					getReviewsDeployment("v1"),
					getReviewsDeployment("v2"),
					getReviewsDeployment("v3"),
				},
				object: exp,
				wantState: test.WantAllStates(
					test.CheckExperimentCompleted,
				),
				wantResults: []runtime.Object{
					getDestinationRule("reviews", name,
						[]string{istio.SubsetBaseline, istio.CandidateSubsetName(0), istio.CandidateSubsetName(1)},
						[]runtime.Object{getReviewsDeployment("v1"), getReviewsDeployment("v2"), getReviewsDeployment("v3")},
					),
					getVirtualServiceForDeployments("reviews", name,
						[]string{istio.SubsetBaseline, istio.CandidateSubsetName(0), istio.CandidateSubsetName(1)},
						[]int32{100, 0, 0},
					),
				},
			}
		}("emptycriterion", getDefaultKubernetesExperiment("emptycriterion", "reviews", "reviews-v1", []string{"reviews-v2", "reviews-v3"})),
		"cleanupdelete": func(name string, exp *iter8v1alpha2.Experiment) testCase {
			return testCase{
				initObjects: []runtime.Object{
					getReviewsService(),
					getReviewsDeployment("v1"),
					getReviewsDeployment("v2"),
					getReviewsDeployment("v3"),
				},
				object: exp,
				wantState: test.WantAllStates(
					test.CheckExperimentCompleted,
				),
				postHook: test.CheckObjectDeleted(getReviewsDeployment("v2"), getReviewsDeployment("v3")),
			}
		}("cleanupdelete", getCleanUpDeleteExperiment("cleanupdelete", "reviews", "reviews-v1", []string{"reviews-v2", "reviews-v3"})),
		"attach-gateway": func(name string, exp *iter8v1alpha2.Experiment) testCase {
			return testCase{
				mocks: map[string]analtyicsapi.Response{
					name: test.GetRollToWinnerMockResponse(exp, 0),
				},
				initObjects: []runtime.Object{
					getReviewsService(),
					getReviewsDeployment("v1"),
					getReviewsDeployment("v2"),
					getReviewsDeployment("v3"),
				},
				object: exp,
				wantState: test.WantAllStates(
					test.CheckExperimentCompleted,
				),
				wantResults: []runtime.Object{
					getDestinationRule("reviews", name,
						[]string{istio.SubsetBaseline, istio.CandidateSubsetName(0), istio.CandidateSubsetName(1)},
						[]runtime.Object{getReviewsDeployment("v1"), getReviewsDeployment("v2"), getReviewsDeployment("v3")},
					),
					getVirtualServiceWithGateway("reviews", name,
						[]string{istio.SubsetBaseline, istio.CandidateSubsetName(0), istio.CandidateSubsetName(1)},
						[]int32{0, 100, 0}, "reviews.com", "gateway-testing",
					),
				},
			}
		}("attach-gateway", getExperimentWithGateway("attach-gateway", "reviews", "reviews-v1", service.GetURL(),
			[]string{"reviews-v2", "reviews-v3"}, "reviews.com", "gateway-testing")),
		// "pauseresume": func(name string, exp *iter8v1alpha2.Experiment) testCase {
		// 	return testCase{
		// 		mocks: map[string]analtyicsapi.Response{
		// 			name: test.GetRollToWinnerMockResponse(exp, 0),
		// 		},
		// 		initObjects: []runtime.Object{
		// 			getReviewsService(),
		// 			getReviewsDeployment("v1"),
		// 			getReviewsDeployment("v2"),
		// 			getReviewsDeployment("v3"),
		// 		},
		// 		object: exp,
		// 		wantState: test.WantAllStates(
		// 			test.CheckExperimentPause,
		// 		),
		// 		postHook: test.ResumeExperiment(exp),
		// 		wantResults: []runtime.Object{
		// 			getStableDestinationRule("reviews", name, getReviewsDeployment("v2")),
		// 			getStableVirtualService("reviews", name),
		// 		},
		// 	}
		// }("pauseresume", getPauseExperiment("pauseresume", "reviews", "reviews-v1", service.GetURL(), []string{"reviews-v2", "reviews-v3"})),
		// "deletebaseline": func(name string) testCase {
		// 	return testCase{
		// 		mocks: map[string]analtyicsapi.Response{
		// 			name: test.GetSuccessMockResponse(),
		// 		},
		// 		initObjects: []runtime.Object{
		// 			getReviewsService(),
		// 			getReviewsDeployment("v1"),
		// 			getReviewsDeployment("v2"),
		// 		},
		// 		object:    getSlowKubernetesExperiment(name, "reviews", "reviews-v1", "reviews-v2", service.GetURL()),
		// 		wantState: test.CheckServiceFound,
		// 		postHook:  test.DeleteObject(getReviewsDeployment("v1")),
		// 		wantResults: []runtime.Object{
		// 			getStableDestinationRule("reviews", name, getReviewsDeployment("v2")),
		// 			getStableVirtualService("reviews", name),
		// 		},
		// 	}
		// }("deletebaseline"),
		// "duplicate-service": func(name string) testCase {
		// 	return testCase{
		// 		mocks: map[string]analtyicsapi.Response{
		// 			name: test.GetSuccessMockResponse(),
		// 		},
		// 		initObjects: []runtime.Object{
		// 			getReviewsService(),
		// 			getReviewsDeployment("v1"),
		// 			getReviewsDeployment("v2"),
		// 		},
		// 		preHook:   []test.Hook{test.CreateObject(getDefaultKubernetesExperiment(name, "reviews", "reviews-v1", "reviews-v2"))},
		// 		object:    getDefaultKubernetesExperiment(name+"duplicate", "reviews", "reviews-v1", "reviews-v2"),
		// 		wantState: test.CheckServiceNotFound("TargetsNotFound"),
		// 		wantResults: []runtime.Object{
		// 			getStableDestinationRule("reviews", name, getReviewsDeployment("v2")),
		// 			getStableVirtualService("reviews", name),
		// 		},
		// 		finalizers: []test.Hook{
		// 			test.DeleteObject(getDefaultKubernetesExperiment(name+"duplicate", "reviews", "reviews-v1", "reviews-v2")),
		// 		},
		// 	}
		// }("duplicate-service"),
	}

	runTestCases(t, service, testCases)
}

func getReviewsService() runtime.Object {
	return test.NewKubernetesService("reviews", Flags.Namespace).
		WithSelector(map[string]string{"app": "reviews"}).
		WithPorts(map[string]int{"http": ReviewsPort}).
		Build()
}

func getReviewsDeployment(version string) runtime.Object {
	labels := map[string]string{
		"app":     "reviews",
		"version": version,
	}

	image := ""

	switch version {
	case "v1":
		image = ReviewsV1Image
	case "v2":
		image = ReviewsV2Image
	case "v3":
		image = ReviewsV3Image
	default:
		image = ReviewsV1Image
	}

	return test.NewKubernetesDeployment("reviews-"+version, Flags.Namespace).
		WithLabels(labels).
		WithContainer("reviews", image, ReviewsPort).
		Build()
}

func getPauseExperiment(name, serviceName, baseline, analyticsHost string, candidates []string) *iter8v1alpha2.Experiment {
	exp := getFastKubernetesExperiment(name, serviceName, baseline, analyticsHost, candidates)
	exp.Spec.ManualOverride = &iter8v1alpha2.ManualOverride{
		Action: iter8v1alpha2.ActionPause,
	}
	return exp
}

func getCleanUpDeleteExperiment(name, serviceName, baseline string, candidates []string) *iter8v1alpha2.Experiment {
	exp := getDefaultKubernetesExperiment(name, serviceName, baseline, candidates)
	cleanup := true
	exp.Spec.Cleanup = &cleanup
	return exp
}

func getDefaultKubernetesExperiment(name, serviceName, baseline string, candidates []string) *iter8v1alpha2.Experiment {
	exp := test.NewExperiment(name, Flags.Namespace).
		WithKubernetesTargetService(serviceName, baseline, candidates).
		Build()

	onesec := "1s"
	one := int32(1)
	exp.Spec.Duration = &iter8v1alpha2.Duration{
		Interval:      &onesec,
		MaxIterations: &one,
	}

	return exp
}

func getFastKubernetesExperiment(name, serviceName, baseline, analyticsHost string, candidates []string) *iter8v1alpha2.Experiment {
	experiment := test.NewExperiment(name, Flags.Namespace).
		WithKubernetesTargetService(serviceName, baseline, candidates).
		WithAnalyticsEndpoint(analyticsHost).
		WithDummyCriterion().
		Build()

	onesec := "1s"
	one := int32(1)
	experiment.Spec.Duration = &iter8v1alpha2.Duration{
		Interval:      &onesec,
		MaxIterations: &one,
	}

	return experiment
}

func getExperimentWithGateway(name, serviceName, baseline, analyticsHost string, candidates []string, host, gw string) *iter8v1alpha2.Experiment {
	experiment := test.NewExperiment(name, Flags.Namespace).
		WithKubernetesTargetService(serviceName, baseline, candidates).
		WithHostInTargetService(host, gw).
		WithAnalyticsEndpoint(analyticsHost).
		WithDummyCriterion().
		Build()

	onesec := "1s"
	one := int32(1)
	experiment.Spec.Duration = &iter8v1alpha2.Duration{
		Interval:      &onesec,
		MaxIterations: &one,
	}

	return experiment
}

func getSlowKubernetesExperiment(name, serviceName, baseline, analyticsHost string, candidates []string) *iter8v1alpha2.Experiment {
	experiment := test.NewExperiment(name, Flags.Namespace).
		WithKubernetesTargetService(serviceName, baseline, candidates).
		WithAnalyticsEndpoint(analyticsHost).
		WithDummyCriterion().
		Build()

	tensecs := "10s"
	two := int32(2)
	experiment.Spec.Duration = &iter8v1alpha2.Duration{
		Interval:      &tensecs,
		MaxIterations: &two,
	}

	return experiment
}

func getDestinationRule(serviceName, name string, subsets []string, objs []runtime.Object) runtime.Object {
	drb := istio.NewDestinationRule(util.ServiceToFullHostName(serviceName, Flags.Namespace), name, Flags.Namespace)
	for i, subset := range subsets {
		drb.WithSubset(objs[i].(*appsv1.Deployment), subset)
	}
	return drb.Build()
}

func getVirtualServiceForDeployments(serviceName, name string, subsets []string, weights []int32) runtime.Object {
	host := util.ServiceToFullHostName(serviceName, Flags.Namespace)
	vsb := istio.NewVirtualService(host, name, Flags.Namespace)
	rb := istio.NewEmptyHTTPRoute()
	for i, subset := range subsets {
		destination := istio.NewHTTPRouteDestination().
			WithHost(host).
			WithSubset(subset).
			WithWeight(weights[i]).Build()
		rb = rb.WithDestination(destination)
	}

	return vsb.WithHTTPRoute(rb.Build()).WithMeshGateway().WithHosts([]string{host}).Build()
}

func getVirtualServiceForServices(serviceName, name string, destinations []string, weights []int32) runtime.Object {
	host := util.ServiceToFullHostName(serviceName, Flags.Namespace)
	vsb := istio.NewVirtualService(host, name, Flags.Namespace)
	rb := istio.NewEmptyHTTPRoute()
	for i, name := range destinations {
		destination := istio.NewHTTPRouteDestination().
			WithHost(name).
			WithWeight(weights[i]).Build()
		rb = rb.WithDestination(destination)
	}

	return vsb.WithHTTPRoute(rb.Build()).WithMeshGateway().WithHosts([]string{host}).Build()
}

func getVirtualServiceWithGateway(serviceName, name string, subsets []string, weights []int32, host, gw string) runtime.Object {
	return istio.NewVirtualServiceBuilder(getVirtualServiceForDeployments(serviceName, name, subsets, weights).(*v1alpha3.VirtualService)).
		WithGateways([]string{gw}).
		WithHosts([]string{host}).
		Build()
}
