package promtest_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mdlayher/promtest"
	"github.com/prometheus/client_golang/prometheus"
)

func TestCollect(t *testing.T) {
	tests := []struct {
		name string
		c    prometheus.Collector
		s    string
	}{
		{
			name: "OK",
			c: prometheus.NewGaugeFunc(
				prometheus.GaugeOpts{
					Name: "promtest_value",
					Help: "A metric for promtest testing.",
				},
				func() float64 { return 1.0 },
			),
			s: `# HELP promtest_value A metric for promtest testing.
# TYPE promtest_value gauge
promtest_value 1
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := string(promtest.Collect(t, tt.c))
			if diff := cmp.Diff(tt.s, s); diff != "" {
				t.Fatalf("unexpected body (-want +got):\n%s", diff)
			}
		})
	}
}

func TestLint(t *testing.T) {
	tests := []struct {
		name string
		s    string
		ok   bool
	}{
		{
			name: "bad",
			s: `# HELP
# TYPE promtest_value counter
promtest_value 1
`,
		},
		{
			name: "OK",
			s: `# HELP promtest_value A metric for promtest testing.
# TYPE promtest_value gauge
promtest_value 1
`,
			ok: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok := promtest.Lint(t, []byte(tt.s))
			if diff := cmp.Diff(tt.ok, ok); diff != "" {
				t.Fatalf("unexpected OK value (-want +got):\n%s", diff)
			}
		})
	}
}

func TestMatch(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		metrics []string
		ok      bool
	}{
		{
			name: "bad",
			s: `# HELP promtest_value A metric for promtest testing.
# TYPE promtest_value gauge
promtest_value 1
`,
			metrics: []string{`promtest_value 100`},
		},
		{
			name: "OK",
			s: `# HELP promtest_value A metric for promtest testing.
# TYPE promtest_value gauge
promtest_value 1
# HELP promtest_labels A metric for promtest testing.
# TYPE promtest_labels gauge
promtest_value{foo="bar"} 2
`,
			metrics: []string{
				`promtest_value 1`,
				`promtest_value{foo="bar"} 2`,
			},
			ok: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok := promtest.Match(t, []byte(tt.s), tt.metrics)
			if diff := cmp.Diff(tt.ok, ok); diff != "" {
				t.Fatalf("unexpected OK value (-want +got):\n%s", diff)
			}
		})
	}
}
