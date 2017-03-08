package metric

import (
	//"fmt"
	l "github.com/advantageous/go-logback/logging"
	//"os/exec"
	//"runtime"
	//"strings"
	//"fmt"
	"os/exec"
	"runtime"
	"strings"
)

const ALL	= -128
const DEBUG	= 1
const ERROR = 4
const FATAL = 5
const INFO  = 2
const OFF   = 127
const WARN  = 3


type NodeMetricGatherer struct {
	logger            l.Logger
	debug             bool
	command           string
	nodeFunction      string
}

func NewNodeMetricGatherer(logger l.Logger, config *Config, nodeFunction string) *NodeMetricGatherer {

	logger = ensureLogger(logger, config.Debug, "node", "MT_NODE_DEBUG")

	command := "/usr/bin/nodetool"
	label := LINUX_LABEL

	if config.FreeCommand != EMPTY {
		command = config.FreeCommand
		label = CONFIG_LABEL
	} else if runtime.GOOS == GOOS_DARWIN {
		command = "/usr/local/bin/nodetool"
		label = DARWIN_LABEL
	}

	if config.Debug {
		logger.Println("Node gatherer initialized by:", label, "as:", command, "function is:", nodeFunction)
	}

	return &NodeMetricGatherer{
		logger:            logger,
		debug:             config.Debug,
		command:           command,
		nodeFunction:      nodeFunction,
	}
}

func (gatherer *NodeMetricGatherer) GetMetrics() ([]Metric, error) {

	var metrics = []Metric{}
	var err error = nil

	if gatherer.nodeFunction == "getlogginglevels" {
		metrics, err = gatherer.getlogginglevels()
	}

	if err != nil {
		return nil, err
	}

	return metrics, err
}

func (gatherer *NodeMetricGatherer) getlogginglevels() ([]Metric, error) {
	var metrics = []Metric{}
	var output string

	if out, err := exec.Command(gatherer.command, "getlogginglevels" ).Output(); err != nil {
		return nil, err
	} else {
		output = string(out)
	}

	//Logger Name                                        Log Level
	//ROOT                                                    INFO
	//com.thinkaurelius.thrift                               ERROR
	//org.apache.cassandra                                   DEBUG

	lines := strings.Split(output, "\n")
	end := len(lines) - 1
	for i := 2; i < end; i++ {
		gatherer.logger.Println("!!!", lines[i])

		split := strings.Split(lines[i], SPACE)
		name := "cassandra:" + split[0];
		logLevelString := split[len(split)-1]
		logLevel := OFF;
		switch logLevelString {
			case "ALL":		logLevel = ALL
			case "DEBUG":	logLevel = DEBUG
			case"ERROR":	logLevel = ERROR
			case "FATAL":	logLevel = FATAL
			case "INFO":	logLevel = INFO
			case "WARN":	logLevel = WARN
		}
		metrics = append(metrics, metric{CUSTOM_UNIT, MetricValue(logLevel), name, "nodetool", logLevelString})
	}

	return metrics, nil

}
