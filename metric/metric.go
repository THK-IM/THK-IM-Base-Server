package metric

import "github.com/prometheus/client_golang/prometheus"

type Metric struct {
	MetricCollector prometheus.Collector
	ID              string
	Name            string
	Description     string
	Type            string
	Args            []string
}
