package nodetool

import (
	c "github.com/cloudurable/metricsd/common"
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
	NodetoolFunction_netstats          = "netstats"
	NodetoolFunction_gcstats           = "gcstats"
	NodetoolFunction_tpstats           = "tpstats"
	NodetoolFunction__getlogginglevels = "getlogginglevels"
)

var NodetoolAllSupportedFunctions = [...]string {
	NodetoolFunction_netstats,
	NodetoolFunction_gcstats,
	NodetoolFunction_tpstats,
	NodetoolFunction__getlogginglevels,
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
