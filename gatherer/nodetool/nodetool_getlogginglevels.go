package nodetool

import (
	"strings"
	c "github.com/cloudurable/metricsd/common"
)

func Getlogginglevels(nodetoolCommand string) ([]c.Metric, error) {
	output, err := c.ExecCommand(nodetoolCommand, NodetoolFunction__getlogginglevels)
	if err != nil {
		return nil, err
	}

	// -- sample getlogginglevels output --
	// <blank line>
	// Logger Name                                        Log Level
	// ROOT                                                    INFO
	// com.thinkaurelius.thrift                               ERROR
	// org.apache.cassandra                                   DEBUG

	var metrics = []c.Metric{}
	lines := strings.Split(output, c.NEWLINE)
	end := len(lines) - 1
	for i := 0; i < end; i++ {
		line := lines[i]

		if line != "" && !strings.Contains(line, "Logger Name") {
			split := strings.Split(line, c.SPACE)
			name := "loggingLevel:" + split[0]
			logLevelString := strings.ToLower(split[len(split)-1])
			value := value_level_off
			switch logLevelString {
			case "all":		value = value_level_all
			case "debug":	value = value_level_debug
			case "error":	value = value_level_error
			case "fatal":	value = value_level_fatal
			case "info":	value = value_level_info
			case "warn":	value = value_level_warn
			}
			metrics = append(metrics, c.Metric{c.NO_UNIT, c.MetricValue(value), name, c.PROVIDER_NODETOOL})
		}
	}

	return metrics, nil
}
