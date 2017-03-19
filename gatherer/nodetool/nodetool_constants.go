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
	NtFunc_netstats            = "netstats"
	NtFunc_gcstats             = "gcstats"
	NtFunc_tpstats             = "tpstats"
	NtFunc_getlogginglevels    = "getlogginglevels"
	NtFunc_gettimeout          = "gettimeout"
	NtFunc_cfstats             = "cfstats"
	NtFunc_proxyhistograms     = "proxyhistograms"
	NtFunc_listsnapshots       = "listsnapshots"
	NtFunc_statuses            = "statuses"
	NtFunc_getstreamthroughput = "getstreamthroughput"
)

var NodetoolAllSupportedFunctions = []string {
	NtFunc_netstats,
	NtFunc_gcstats,
	NtFunc_tpstats,
	NtFunc_getlogginglevels,
	NtFunc_gettimeout,
	NtFunc_cfstats,
	NtFunc_proxyhistograms,
	NtFunc_listsnapshots,
	NtFunc_statuses,
	NtFunc_getstreamthroughput,
}