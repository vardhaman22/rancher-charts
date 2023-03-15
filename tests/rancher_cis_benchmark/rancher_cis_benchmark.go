package rancher_cis_benchmark

import (
	"github.com/aiyengar2/hull/pkg/chart"
	"github.com/aiyengar2/hull/pkg/checker"
	"github.com/aiyengar2/hull/pkg/test"
	"github.com/aiyengar2/hull/pkg/utils"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ChartPath = utils.MustGetLatestChartVersionPathFromIndex("../index.yaml", "rancher-cis-benchmark", true)

var (
	DefaultReleaseName = "rancher-cis-benchmark"
	DefaultNamespace   = "cis-operator-system"
)

var defaultTolerations = []corev1.Toleration{
	{
		Key:      "cattle.io/os",
		Operator: corev1.TolerationOpEqual,
		Value:    "linux",
		Effect:   corev1.TaintEffectNoSchedule,
	},
}

var testTolerations = []corev1.Toleration{
	{
		Key:      "test",
		Operator: corev1.TolerationOpEqual,
		Value:    "test",
		Effect:   corev1.TaintEffectNoSchedule,
	},
}

var suite = test.Suite{
	ChartPath: ChartPath,

	Cases: []test.Case{
		{
			Name: "Using Defaults",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace),
		},
		{
			Name: "Set Test Tolerations",

			TemplateOptions: chart.NewTemplateOptions(DefaultReleaseName, DefaultNamespace).Set("tolerations", testTolerations),
		},
	},

	NamedChecks: []test.NamedCheck{
		{
			Name: "Check Pod Tolerations Values",
			Covers: []string{
				".Values.tolerations",
			},

			Checks: test.Checks{
				checker.PerWorkload(func(tc *checker.TestContext, obj metav1.Object, podTemplateSpec corev1.PodTemplateSpec) {
					if obj.GetNamespace() != checker.MustRenderValue[string](tc, ".Release.Namespace") {
						return
					}

					tolerationsAddedByValues := checker.MustRenderValue[[]corev1.Toleration](tc, ".Values.tolerations")
					expectedArgs := append(defaultTolerations, tolerationsAddedByValues...)
					if len(expectedArgs) == 0 {
						assert.Nil(tc.T, podTemplateSpec.Spec.Tolerations,
							"expected pod %s in %T %s to have no args",
							podTemplateSpec.Name, obj, checker.Key(obj),
						)
					} else {
						assert.Equal(tc.T,
							expectedArgs, podTemplateSpec.Spec.Tolerations,
							"container %s in %T %s does not have correct args",
							podTemplateSpec.Name, obj, checker.Key(obj),
						)
					}
				}),
			},
		},
	},
}
