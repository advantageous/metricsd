package gatherer

import (
	l "github.com/advantageous/go-logback/logging"
	c "github.com/cloudurable/metricsd/common"
	nt "github.com/cloudurable/metricsd/gatherer/nodetool"
	"strings"
)

type NodetoolMetricGatherer struct {
	logger            l.Logger
	debug             bool
	command           string
	nodeFunction      string
}

func nodetoolFunctionSupported(nodeFunction string) bool {
	lower := strings.ToLower(nodeFunction)
	for _,supported := range nt.NodetoolAllSupportedFunctions {
		if supported == lower {
			return true
		}
	}
	return false
}

func NewNodetoolMetricGatherers(logger l.Logger, config *c.Config) []*NodetoolMetricGatherer {

	if config.NodetoolFunctions == nil || len(config.NodetoolFunctions) == 0 {
		return nil
	}

	gatherers := []*NodetoolMetricGatherer{}
	for _, nodeFunction := range config.NodetoolFunctions {
		if nodetoolFunctionSupported(nodeFunction) {
			gatherers = append(gatherers, newNodetoolMetricGatherer(logger, config, nodeFunction))
		} else {
			logger.Warn("Unsupported or unknown Nodetool function", nodeFunction)
		}
	}

	return gatherers
}

func newNodetoolMetricGatherer(logger l.Logger, config *c.Config, nodeFunction string) *NodetoolMetricGatherer {
	logger = c.EnsureLogger(logger, config.Debug, c.PROVIDER_NODETOOL, c.FLAG_NODETOOL)
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
	case nt.NtFunc_netstats:		    metrics, err = nt.Netstats(gatherer.command)
	case nt.NtFunc_gcstats:			    metrics, err = nt.Gcstats(gatherer.command)
	case nt.NtFunc_tpstats:			    metrics, err = nt.Tpstats(gatherer.command)
	case nt.NtFunc_getlogginglevels:    metrics, err = nt.Getlogginglevels(gatherer.command)
	case nt.NtFunc_gettimeout:	        metrics, err = nt.Gettimeout(gatherer.command)
	case nt.NtFunc_cfstats:	            metrics, err = nt.Cfstats(gatherer.command)
	case nt.NtFunc_proxyhistograms:     metrics, err = nt.ProxyHistograms(gatherer.command)
	case nt.NtFunc_listsnapshots:       metrics, err = nt.ListSnapshots(gatherer.command)
	case nt.NtFunc_statuses:            metrics, err = nt.Statuses(gatherer.command)
	case nt.NtFunc_getstreamthroughput: metrics, err = nt.GetStreamThroughput(gatherer.command)
	}

	if err != nil {
		return nil, err
	}

	return metrics, err
}
