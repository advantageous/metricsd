package nodetool

import (
	"strings"
	c "github.com/cloudurable/metricsd/common"
)

func Getlogginglevels(nodetoolCommand string) ([]c.Metric, error) {
	output, err := c.ExecCommand(nodetoolCommand, NtFunc_getlogginglevels)
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
			name := "ntLl:" + split[0]
			logLevelString := strings.ToLower(split[len(split)-1])
			code := value_level_off
			switch logLevelString {
			case "all":		code = value_level_all
			case "debug":	code = value_level_debug
			case "error":	code = value_level_error
			case "fatal":	code = value_level_fatal
			case "info":	code = value_level_info
			case "warn":	code = value_level_warn
			}
			metrics = append(metrics, *c.NewMetricStringCode(c.MT_NONE, logLevelString, code, name, c.PROVIDER_NODETOOL))
		}
	}

	return metrics, nil
}
