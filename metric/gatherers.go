package metric

func LoadGatherers(config *Config) ([]MetricsGatherer) {

	var gatherers = []MetricsGatherer{}

	cpu := NewCPUMetricsGatherer(nil, config)
	if cpu != nil {
		gatherers = append(gatherers, cpu)
	}

	disk := NewDiskMetricsGatherer(nil, config)
	if disk != nil {
		gatherers = append(gatherers, disk)
	}

	free := NewFreeMetricGatherer(nil, config)
	if free != nil {
		gatherers = append(gatherers, free)
	}

	for _, nodetool := range NewNodetoolMetricGatherers(nil, config) {
		gatherers = append(gatherers, nodetool)
	}

	return gatherers
}
