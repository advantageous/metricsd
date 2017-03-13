package repeater

import (
	c "github.com/cloudurable/metricsd/common"
)

func LoadRepeaters(config *c.Config) ([]c.MetricsRepeater) {

	var repeaters = []c.MetricsRepeater{}

	for _,provider := range config.Repeaters {
		switch provider {
		case c.REPEATER_AWS:
			repeater := NewAwsCloudMetricRepeater(config)
			if repeater != nil {
				repeaters = append(repeaters, repeater)
			}

		case c.REPEATER_LOGGER:
			repeater := NewLogMetricsRepeater()
			if repeater != nil {
				repeaters = append(repeaters, repeater)
			}
		}
	}

	return repeaters
}
