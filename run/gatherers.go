package run

import (
	c "github.com/cloudurable/metricsd/common"
	g "github.com/cloudurable/metricsd/gatherer"
)

func LoadGatherers(config *c.Config) ([]c.MetricsGatherer) {

	var gatherers = []c.MetricsGatherer{}

	for _,provider := range config.Gatherers {
		switch provider {
		case c.PROVIDER_CPU:
			cpu := g.NewCPUMetricsGatherer(nil, config)
			if cpu != nil {
				gatherers = append(gatherers, cpu)
			}

		case c.PROVIDER_DISK:
			disk := g.NewDiskMetricsGatherer(nil, config)
			if disk != nil {
				gatherers = append(gatherers, disk)
			}

		case c.PROVIDER_FREE:
			free := g.NewFreeMetricGatherer(nil, config)
			if free != nil {
				gatherers = append(gatherers, free)
			}

		case c.PROVIDER_NODETOOL:
			for _, nodetool := range g.NewNodetoolMetricGatherers(nil, config) {
				gatherers = append(gatherers, nodetool)
			}
		}
	}

	return gatherers
}
