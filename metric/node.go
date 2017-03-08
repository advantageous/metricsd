package metric

import (
	l "github.com/advantageous/go-logback/logging"
	"strings"
)

const (
	value_level_all   = 0
	value_level_debug = 1
	value_level_info  = 2
	value_level_warn  = 3
	value_level_error = 4
	value_level_fatal = 5
	value_level_off   = 127

	value_na    = -125
	value_nan   = -126
	value_error = -127

	value_mode_starting = 0
	value_mode_normal = 1
	value_mode_joining = 2
	value_mode_leaving = 3
	value_mode_decommissioned = 4
	value_mode_moving = 5
	value_mode_draining = 6
	value_mode_drained = 7
	value_mode_other = 99
)

const (
	in_value_na = "n/a"
	in_value_nan = "NaN"
)

const (
	function_net_stats          = "netstats"
	function_gc_stats           = "gcstats"
	function_get_logging_levels = "getlogginglevels"
)

var SUPPORTED = [...]string {
	function_gc_stats,
	function_get_logging_levels,
	function_net_stats,
}

type NodeMetricGatherer struct {
	logger            l.Logger
	debug             bool
	command           string
	nodeFunction      string
}

func NodeFunctionSupported(nodeFunction string) bool {
	lower := strings.ToLower(nodeFunction)
	for _,supported := range SUPPORTED {
		if supported == lower {
			return true;
		}
	}
	return false;
}

func NewNodeMetricGatherer(logger l.Logger, config *Config, nodeFunction string) *NodeMetricGatherer {

	logger = ensureLogger(logger, config.Debug, PROVIDER_NODE, FLAG_NODE)

	command := "/usr/bin/nodetool"
	label := LINUX_LABEL

	if config.NodetoolCommand != EMPTY {
		command = config.NodetoolCommand
		label = CONFIG_LABEL
	}

	if config.Debug {
		logger.Println("Node gatherer initialized by:", label, "as:", command, "function is:", nodeFunction)
	}

	return &NodeMetricGatherer{
		logger:            logger,
		debug:             config.Debug,
		command:           command,
		nodeFunction:      strings.ToLower(nodeFunction),
	}
}

func (gatherer *NodeMetricGatherer) GetMetrics() ([]Metric, error) {

	var metrics = []Metric{}
	var err error = nil

	switch gatherer.nodeFunction {
	case function_net_stats:			metrics, err = gatherer.netstats()
	case function_gc_stats:				metrics, err = gatherer.gcstats()
	case function_get_logging_levels:	metrics, err = gatherer.getlogginglevels()
	}

	if err != nil {
		return nil, err
	}

	return metrics, err
}

func (gatherer *NodeMetricGatherer) netstats() ([]Metric, error) {
	output, err := execCommand(gatherer.command, function_net_stats)
	if err != nil {
		return nil, err
	}

	// -- sample netstats output --
	// [0] // Mode: NORMAL
	// [1] // Not sending any streams.
	// [2] // Read Repair Statistics:
	// [3] // Attempted: 0
	// [4] // Mismatch (Blocking): 0
	// [5] // Mismatch (Background): 0
	// [6] // Pool Name                    Active   Pending      Completed   Dropped
	// [7] // Large messages                  n/a         0              0         0
	// [8] // Small messages                  n/a         0              2         0
	// [9] // Gossip messages                 n/a         0              0         0

	lines := strings.Split(output, NEWLINE)

	var metrics = appendNsMode([]Metric{}, lines[0])

	metrics = appendNsReadRepair(metrics, lines[3], 1, "nsRrAttempted")
	metrics = appendNsReadRepair(metrics, lines[4], 2, "nsRrBlocking")
	metrics = appendNsReadRepair(metrics, lines[4], 2, "nsRrBackground")

	metrics = appendNsPool(metrics, lines[7], "nsLargeMsgs")
	metrics = appendNsPool(metrics, lines[8], "nsSmallMsgs")
	metrics = appendNsPool(metrics, lines[9], "nsGossipMsgs")

	return metrics, nil
}

func appendNsMode(metrics []Metric, line string) []Metric {
	value := value_mode_other
	switch strings.ToLower(parseForColumn(line, 1)) {
	case "starting":		value = value_mode_starting
	case "normal":			value = value_mode_normal
	case "joining":			value = value_mode_joining
	case "leaving":			value = value_mode_leaving
	case "decommissioned": 	value = value_mode_decommissioned
	case "moving":			value = value_mode_moving
	case "draining":		value = value_mode_draining
	case "drained":			value = value_mode_drained
	}
	return append(metrics, metric{NO_UNIT, MetricValue(value), "nsMode", PROVIDER_NODE})
}

func appendNsReadRepair(metrics []Metric, line string, columnIndex int, name string) []Metric {
	metricValue := MetricValue( toInt64(parseForColumn(line, columnIndex), value_error) )
	return append(metrics, metric{COUNT, metricValue, name, PROVIDER_NODE})
}

func appendNsPool(metrics []Metric, line string, prefix string) []Metric {
	parsed := splitValuesOnly(line)
	metrics = append(metrics, metric{COUNT, nsPoolMetricValue(parsed, 2), prefix + "Active", PROVIDER_NODE})
	metrics = append(metrics, metric{COUNT, nsPoolMetricValue(parsed, 3), prefix + "Pending", PROVIDER_NODE})
	metrics = append(metrics, metric{COUNT, nsPoolMetricValue(parsed, 4), prefix + "Completed", PROVIDER_NODE})
	metrics = append(metrics, metric{COUNT, nsPoolMetricValue(parsed, 5), prefix + "Dropped", PROVIDER_NODE})
	return metrics
}

func nsPoolMetricValue(parsed []string, columnIndex int) MetricValue {
	value := parsed[columnIndex]
	if value == in_value_na {
		return MetricValue(value_na)
	}
	return MetricValue(toInt64(value, value_error))
}

func (gatherer *NodeMetricGatherer) gcstats() ([]Metric, error) {
	output, err := execCommand(gatherer.command, function_gc_stats)
	if err != nil {
		return nil, err
	}

	// -- sample gcstats output --
	//Interval (ms) Max GC Elapsed (ms)Total GC Elapsed (ms)Stdev GC Elapsed (ms)   GC Reclaimed (MB)         Collections      Direct Memory Bytes
	//3491665                   0                   0                 NaN                   0                   0                       -1

	lines := strings.Split(output, NEWLINE)
	values := splitValuesOnly(lines[1])

	var metrics = []Metric{}
	metrics = append(metrics, metric{TIMING_MS, gcstatsMetricValue(values[0]), "gcInterval", PROVIDER_NODE})
	metrics = append(metrics, metric{TIMING_MS, gcstatsMetricValue(values[1]), "gcMaxElapsed", PROVIDER_NODE})
	metrics = append(metrics, metric{TIMING_MS, gcstatsMetricValue(values[2]), "gcTotalElapsed", PROVIDER_NODE})
	metrics = append(metrics, metric{TIMING_MS, gcstatsMetricValue(values[3]), "gcStdevElapsed", PROVIDER_NODE})
	metrics = append(metrics, metric{SIZE_MB, gcstatsMetricValue(values[4]), "gcReclaimed", PROVIDER_NODE})
	metrics = append(metrics, metric{COUNT, gcstatsMetricValue(values[5]), "gcCollections", PROVIDER_NODE})
	metrics = append(metrics, metric{SIZE_B, gcstatsMetricValue(values[6]), "gcDirectMemoryBytes", PROVIDER_NODE})

	return metrics, nil
}

func gcstatsMetricValue(value string) MetricValue {
	if value == in_value_nan {
		return MetricValue(value_nan)
	}
	return MetricValue(toInt64(value, value_error))
}

func (gatherer *NodeMetricGatherer) getlogginglevels() ([]Metric, error) {
	output, err := execCommand(gatherer.command, function_get_logging_levels)
	if err != nil {
		return nil, err
	}

	// -- sample getlogginglevels output --
	// <blank line>
	// Logger Name                                        Log Level
	// ROOT                                                    INFO
	// com.thinkaurelius.thrift                               ERROR
	// org.apache.cassandra                                   DEBUG

	var metrics = []Metric{}
	lines := strings.Split(output, NEWLINE)
	end := len(lines) - 1
	for i := 0; i < end; i++ {
		line := lines[i]

		if line != "" && !strings.Contains(line, "Logger Name") {
			split := strings.Split(line, SPACE)
			name := "logging.level." + split[0];
			logLevelString := strings.ToLower(split[len(split)-1])
			value := value_level_off;
			switch logLevelString {
			case "all":		value = value_level_all
			case "debug":	value = value_level_debug
			case "error":	value = value_level_error
			case "fatal":	value = value_level_fatal
			case "info":	value = value_level_info
			case "warn":	value = value_level_warn
			}
			metrics = append(metrics, metric{NO_UNIT, MetricValue(value), name, PROVIDER_NODE})
		}
	}

	return metrics, nil
}
