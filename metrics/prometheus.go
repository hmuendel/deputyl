package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metric struct {
	PatchVersions float64
	MinorVersions float64
	MajorVersions float64
	Repository    string
	Labels        map[string]string
}

var patchOpts = prometheus.GaugeOpts{
	Namespace:   "deputyl",
	Subsystem:   "versions",
	Name:        "patch",
	Help:        "Number of newer patch semvers found in the upstream repository",
	ConstLabels: nil,
}

var minorOpts = prometheus.GaugeOpts{
	Namespace:   "deputyl",
	Subsystem:   "versions",
	Name:        "minor",
	Help:        "Number of newer minor semvers found in the upstream repository",
	ConstLabels: nil,
}

var majorOpts = prometheus.GaugeOpts{
	Namespace:   "deputyl",
	Subsystem:   "versions",
	Name:        "major",
	Help:        "Number of newer major semvers found in the upstream repository",
	ConstLabels: nil,
}

type PrometheusEmitter struct {
	patchVersions *prometheus.GaugeVec
	minorVersions *prometheus.GaugeVec
	majorVersions *prometheus.GaugeVec
}

func NewPrometheusEmitter() PrometheusEmitter {
	return PrometheusEmitter{
		patchVersions: promauto.NewGaugeVec(patchOpts, []string{"repository", "pod", "namespace"}),
		minorVersions: promauto.NewGaugeVec(minorOpts, []string{"repository", "pod", "namespace"}),
		majorVersions: promauto.NewGaugeVec(majorOpts, []string{"repository", "pod", "namespace"}),
	}
}

func (pe PrometheusEmitter) Emit(metric Metric) {
	labels := prometheus.Labels(metric.Labels)
	pe.patchVersions.MustCurryWith(labels).WithLabelValues(metric.Repository).Set(metric.PatchVersions)
	pe.minorVersions.MustCurryWith(labels).WithLabelValues(metric.Repository).Set(metric.MinorVersions)
	pe.majorVersions.MustCurryWith(labels).WithLabelValues(metric.Repository).Set(metric.MajorVersions)
}
