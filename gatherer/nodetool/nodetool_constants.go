package nodetool

const (
	value_level_all int64 = 0
	value_level_debug int64 = 1
	value_level_info  int64 = 2
	value_level_warn  int64 = 3
	value_level_error int64 = 4
	value_level_fatal int64 = 5
	value_level_off   int64 = 127
)

const (
	value_mode_starting int64 = 0
	value_mode_normal int64 = 1
	value_mode_joining int64 = 2
	value_mode_leaving int64 = 3
	value_mode_decommissioned int64 = 4
	value_mode_moving int64 = 5
	value_mode_draining int64 = 6
	value_mode_drained int64 = 7
	value_mode_other int64 = 99
)

const (
	NodetoolFunction_netstats         = "netstats"
	NodetoolFunction_gcstats          = "gcstats"
	NodetoolFunction_tpstats          = "tpstats"
	NodetoolFunction_getlogginglevels = "getlogginglevels"
	NodetoolFunction_gettimeout       = "gettimeout"
	NodetoolFunction_cfstats          = "cfstats"
)

var NodetoolAllSupportedFunctions = []string {
	NodetoolFunction_netstats,
	NodetoolFunction_gcstats,
	NodetoolFunction_tpstats,
	NodetoolFunction_getlogginglevels,
	NodetoolFunction_gettimeout,
	NodetoolFunction_cfstats,
}