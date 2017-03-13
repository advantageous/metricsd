package gatherer

import (
	c "github.com/cloudurable/metricsd/common"
)

func LoadGatherers(config *c.Config) ([]c.MetricsGatherer) {

	var gatherers = []c.MetricsGatherer{}

	for _,provider := range config.Gatherers {
		switch provider {
		case c.PROVIDER_CPU:
			cpu := NewCPUMetricsGatherer(nil, config)
			if cpu != nil {
				gatherers = append(gatherers, cpu)
			}

		case c.PROVIDER_DISK:
			disk := NewDiskMetricsGatherer(nil, config)
			if disk != nil {
				gatherers = append(gatherers, disk)
			}

		case c.PROVIDER_FREE:
			free := NewFreeMetricGatherer(nil, config)
			if free != nil {
				gatherers = append(gatherers, free)
			}

		case c.PROVIDER_NODETOOL:
			tools := NewNodetoolMetricGatherers(nil, config)
			if tools != nil {
				for _, tool := range tools {
					gatherers = append(gatherers, tool)
				}
			}
		}
	}

	return gatherers
}
