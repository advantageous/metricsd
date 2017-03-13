package gatherer

import (
	l "github.com/advantageous/go-logback/logging"
	c "github.com/cloudurable/metricsd/common"
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

func NewNodetoolMetricGatherers(logger l.Logger, config *c.Config) []*NodetoolMetricGatherer {

	if config.NodetoolFunctions == nil || len(config.NodetoolFunctions) == 0 {
		return nil
	}

	gatherers := []*NodetoolMetricGatherer{}
	for _, nodeFunction := range config.NodetoolFunctions {
		if nodetoolFunctionSupported(nodeFunction) {
			gatherers = append(gatherers, newNodetoolMetricGatherer(logger, config, nodeFunction))
		}
	}

	return gatherers
}

func newNodetoolMetricGatherer(logger l.Logger, config *c.Config, nodeFunction string) *NodetoolMetricGatherer {
	logger = c.EnsureLogger(logger, config.Debug, c.PROVIDER_NODETOOL, c.FLAG_NODE)
	command := readNodetoolConfig(config, logger, nodeFunction)

	return &NodetoolMetricGatherer{
		logger:            logger,
		debug:             config.Debug,
		command:           command,
		nodeFunction:      strings.ToLower(nodeFunction),
	}
}

func readNodetoolConfig(config *c.Config, logger l.Logger, nodeFunction string) (string) {
	command := "/usr/bin/nodetool"
	label := c.DEFAULT_LABEL

	if config.NodetoolCommand != c.EMPTY {
		command = config.NodetoolCommand
		label = c.CONFIG_LABEL
	}

	if config.Debug {
		logger.Println("Node gatherer initialized by:", label, "as:", command, "function is:", nodeFunction)
	}

	return command
}

func (gatherer *NodetoolMetricGatherer) GetMetrics() ([]c.Metric, error) {

	var metrics = []c.Metric{}
	var err error = nil

	switch gatherer.nodeFunction {
	case function_net_stats:			metrics, err = netstats(gatherer.command)
	case function_gc_stats:				metrics, err = gcstats(gatherer.command)
	case function_tp_stats:				metrics, err = tpstats(gatherer.command)
	case function_get_logging_levels:	metrics, err = getlogginglevels(gatherer.command)
	}

	if err != nil {
		return nil, err
	}

	return metrics, err
}

func numericMetricValue(value string) c.MetricValue {
	if value == in_value_na {
		return c.MetricValue(value_na)
	}

	if value == in_value_nan {
		return c.MetricValue(value_nan)
	}

	return c.MetricValue(c.ToInt64(value, value_error))
}
