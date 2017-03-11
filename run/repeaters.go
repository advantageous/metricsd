package run

import (
	c "github.com/cloudurable/metricsd/common"
    r "github.com/cloudurable/metricsd/repeater"
)

func LoadRepeaters(config *c.Config) ([]c.MetricsRepeater) {

	var repeaters = []c.MetricsRepeater{}

	for _,provider := range config.Repeaters {
		switch provider {
		case c.REPEATER_AWS:
			aws := r.NewAwsCloudMetricRepeater(config);
			if aws != nil {
				repeaters = append(repeaters, aws)
			}

		case c.REPEATER_LOGGER:
			lgr := r.NewLogMetricsRepeater();
			if lgr != nil {
				repeaters = append(repeaters, lgr)
			}
		}
	}

	return repeaters
}
