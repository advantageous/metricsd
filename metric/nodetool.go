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
	function_tp_stats           = "tpstats"
	function_get_logging_levels = "getlogginglevels"
)

var supportedFunctions = [...]string {
	function_net_stats,
	function_gc_stats,
	function_tp_stats,
	function_get_logging_levels,
}

type NodetoolMetricGatherer struct {
	logger            l.Logger
	debug             bool
	command           string
	nodeFunction      string
}

func nodetoolFunctionSupported(nodeFunction string) bool {
	lower := strings.ToLower(nodeFunction)
	for _,supported := range supportedFunctions {
		if supported == lower {
			return true;
		}
	}
	return false;
}

func NewNodetoolMetricGatherers(logger l.Logger, config *Config) []*NodetoolMetricGatherer {
	gatherers := []*NodetoolMetricGatherer{}

	if (config.NodetoolGather) {
		nodetoolFunctions := strings.Split(config.NodetoolFunctions, SPACE)
		for _, nodeFunction := range nodetoolFunctions {
			if nodetoolFunctionSupported(nodeFunction) {
				gatherers = append(gatherers, newNodetoolMetricGatherer(logger, config, nodeFunction))
			}
		}
	}
	return gatherers
}

func newNodetoolMetricGatherer(logger l.Logger, config *Config, nodeFunction string) *NodetoolMetricGatherer {
	logger = EnsureLogger(logger, config.Debug, PROVIDER_NODETOOL, FLAG_NODE)
	command := readNodetoolConfig(config, logger, nodeFunction)

	return &NodetoolMetricGatherer{
		logger:            logger,
		debug:             config.Debug,
		command:           command,
		nodeFunction:      strings.ToLower(nodeFunction),
	}
}

func readNodetoolConfig(config *Config, logger l.Logger, nodeFunction string) (string) {
	command := "/usr/bin/nodetool"
	label := DEFAULT_LABEL

	if config.NodetoolCommand != EMPTY {
		command = config.NodetoolCommand
		label = CONFIG_LABEL
	}

	if config.Debug {
		logger.Println("Node gatherer initialized by:", label, "as:", command, "function is:", nodeFunction)
	}

	return command
}

func isNodeFunctionRequested(config *Config, inNodeFunction string) (bool) {
	nodetoolFunctions := strings.Split(config.NodetoolFunctions, SPACE)
	for _, nodeFunction := range nodetoolFunctions {
		if nodeFunction == inNodeFunction {
			return true
		}
	}
	return false
}

func (gatherer *NodetoolMetricGatherer) Reload(config *Config) (ReloadResult) {
	if (!config.NodetoolGather || !isNodeFunctionRequested(config, gatherer.nodeFunction)) {
		return RELOAD_EJECT
	}  // eject if not turned on, or the nodeFunction was removed from the list

	gatherer.command = readNodetoolConfig(config, gatherer.logger, gatherer.nodeFunction)
	return RELOAD_SUCCESS
}

func (gatherer *NodetoolMetricGatherer) GetMetrics() ([]Metric, error) {

	var metrics = []Metric{}
	var err error = nil

	switch gatherer.nodeFunction {
	case function_net_stats:			metrics, err = gatherer.netstats()
	case function_gc_stats:				metrics, err = gatherer.gcstats()
	case function_tp_stats:				metrics, err = gatherer.tpstats()
	case function_get_logging_levels:	metrics, err = gatherer.getlogginglevels()
	}

	if err != nil {
		return nil, err
	}

	return metrics, err
}

func (gatherer *NodetoolMetricGatherer) netstats() ([]Metric, error) {
	output, err := ExecCommand(gatherer.command, function_net_stats)
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

	metrics = appendNsPool(metrics, lines[7], "nsPoolLargeMsgs")
	metrics = appendNsPool(metrics, lines[8], "nsPoolSmallMsgs")
	metrics = appendNsPool(metrics, lines[9], "nsPoolGossipMsgs")

	return metrics, nil
}

func appendNsMode(metrics []Metric, line string) []Metric {
	value := value_mode_other
	switch strings.ToLower(FieldByIndex(line, 1)) {
	case "starting":		value = value_mode_starting
	case "normal":			value = value_mode_normal
	case "joining":			value = value_mode_joining
	case "leaving":			value = value_mode_leaving
	case "decommissioned": 	value = value_mode_decommissioned
	case "moving":			value = value_mode_moving
	case "draining":		value = value_mode_draining
	case "drained":			value = value_mode_drained
	}
	return append(metrics, metric{NO_UNIT, MetricValue(value), "nsMode", PROVIDER_NODETOOL})
}

func appendNsReadRepair(metrics []Metric, line string, columnIndex int, name string) []Metric {
	metricValue := MetricValue( ToInt64(FieldByIndex(line, columnIndex), value_error) )
	return append(metrics, metric{COUNT, metricValue, name, PROVIDER_NODETOOL})
}

func appendNsPool(metrics []Metric, line string, prefix string) []Metric {
	valuesOnly := strings.Fields(line)
	metrics = append(metrics, metric{COUNT, numericMetricValue(valuesOnly[2]), prefix + "Active", PROVIDER_NODETOOL})
	metrics = append(metrics, metric{COUNT, numericMetricValue(valuesOnly[3]), prefix + "Pending", PROVIDER_NODETOOL})
	metrics = append(metrics, metric{COUNT, numericMetricValue(valuesOnly[4]), prefix + "Completed", PROVIDER_NODETOOL})
	return append(metrics, metric{COUNT, numericMetricValue(valuesOnly[5]), prefix + "Dropped", PROVIDER_NODETOOL})
}

func (gatherer *NodetoolMetricGatherer) tpstats() ([]Metric, error) {
	output, err := ExecCommand(gatherer.command, function_tp_stats)
	if err != nil {
		return nil, err
	}

	//Pool Name                         Active   Pending      Completed   Blocked  All time blocked
	//ReadStage                              0         0              3         0                 0
	//MiscStage                              0         0              0         0                 0
	//CompactionExecutor                     0         0          51133         0                 0
	//MutationStage                          0         0              1         0                 0
	//...                                    0         0              1         0                 0
	//
	//Message type           Dropped
	//READ                         0
	//RANGE_SLICE                  0
	//_TRACE                       0
	//HINT                         0
	//MUTATION                     0
	//COUNTER_MUTATION             0
	//BATCH_STORE                  0
	//BATCH_REMOVE                 0
	//REQUEST_RESPONSE             0
	//PAGED_RANGE                  0
	//READ_REPAIR                  0

	var metrics = []Metric{}

	lines := strings.Split(output, NEWLINE)
	state := 0
	for _,line := range lines {
		if (state == 0 || state == 2) {
			state++ // skip the line
		} else if state == 1 {
			if line != EMPTY {
				metrics = appendTpPool(metrics, line)
			} else {
				state = 2
			}
		} else if state == 3 {
			if line != EMPTY {
				metrics = appendTpMessageType(metrics, line)
			}
		}
	}

	return metrics, nil
}

func appendTpMessageType(metrics []Metric, line string) []Metric {
	valuesOnly := strings.Fields(line)
	parts := strings.Split(valuesOnly[0], UNDER)
	name := "tpMsgType"
	for _,part := range parts {
		if part != EMPTY {
			name = name + part[0:1] + strings.ToLower(part[1:])
		}
	}
	return append(metrics, metric{COUNT, numericMetricValue(valuesOnly[1]), name, PROVIDER_NODETOOL})
}

func appendTpPool(metrics []Metric, line string) []Metric {
	valuesOnly := strings.Fields(line)
	prefix := "tpPool" + valuesOnly[0]
	metrics = append(metrics, metric{COUNT, numericMetricValue(valuesOnly[1]), prefix + "Active", PROVIDER_NODETOOL})
	metrics = append(metrics, metric{COUNT, numericMetricValue(valuesOnly[2]), prefix + "Pending", PROVIDER_NODETOOL})
	metrics = append(metrics, metric{COUNT, numericMetricValue(valuesOnly[3]), prefix + "Completed", PROVIDER_NODETOOL})
	metrics = append(metrics, metric{COUNT, numericMetricValue(valuesOnly[4]), prefix + "Blocked", PROVIDER_NODETOOL})
	return append(metrics, metric{COUNT, numericMetricValue(valuesOnly[5]), prefix + "AllTimeBlocked", PROVIDER_NODETOOL})
}

func (gatherer *NodetoolMetricGatherer) gcstats() ([]Metric, error) {
	output, err := ExecCommand(gatherer.command, function_gc_stats)
	if err != nil {
		return nil, err
	}

	// -- sample gcstats output --
	//Interval (ms) Max GC Elapsed (ms)Total GC Elapsed (ms)Stdev GC Elapsed (ms)   GC Reclaimed (MB)         Collections      Direct Memory Bytes
	//3491665                   0                   0                 NaN                   0                   0                       -1

	lines := strings.Split(output, NEWLINE)
	values := strings.Fields(lines[1])

	var metrics = []Metric{}
	metrics = append(metrics, metric{TIMING_MS, numericMetricValue(values[0]), "gcInterval", PROVIDER_NODETOOL})
	metrics = append(metrics, metric{TIMING_MS, numericMetricValue(values[1]), "gcMaxElapsed", PROVIDER_NODETOOL})
	metrics = append(metrics, metric{TIMING_MS, numericMetricValue(values[2]), "gcTotalElapsed", PROVIDER_NODETOOL})
	metrics = append(metrics, metric{TIMING_MS, numericMetricValue(values[3]), "gcStdevElapsed", PROVIDER_NODETOOL})
	metrics = append(metrics, metric{SIZE_MB, numericMetricValue(values[4]), "gcReclaimed", PROVIDER_NODETOOL})
	metrics = append(metrics, metric{COUNT, numericMetricValue(values[5]), "gcCollections", PROVIDER_NODETOOL})
	metrics = append(metrics, metric{SIZE_B, numericMetricValue(values[6]), "gcDirectMemoryBytes", PROVIDER_NODETOOL})

	return metrics, nil
}

func (gatherer *NodetoolMetricGatherer) getlogginglevels() ([]Metric, error) {
	output, err := ExecCommand(gatherer.command, function_get_logging_levels)
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
			name := "loggingLevel:" + split[0];
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
			metrics = append(metrics, metric{NO_UNIT, MetricValue(value), name, PROVIDER_NODETOOL})
		}
	}

	return metrics, nil
}

func numericMetricValue(value string) MetricValue {
	if value == in_value_na {
		return MetricValue(value_na)
	}

	if value == in_value_nan {
		return MetricValue(value_nan)
	}

	return MetricValue(ToInt64(value, value_error))
}
