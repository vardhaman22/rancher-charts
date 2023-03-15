package rancher_cis_benchmark

import (
	"testing"

	"github.com/aiyengar2/hull/pkg/test"
)

func TestChart(t *testing.T) {
	opts := test.GetRancherOptions()
	// opts.Coverage.IncludeSubcharts = true
	opts.Coverage.Disabled = false
	opts.YAMLLint.Enabled = true
	suite.Run(t, opts)
}
