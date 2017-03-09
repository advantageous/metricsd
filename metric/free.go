package metric

import (
	"fmt"
	l "github.com/advantageous/go-logback/logging"
	"runtime"
	"strings"
)

type FreeMetricGatherer struct {
	logger l.Logger
	debug  bool
	command string
}

func NewFreeMetricGatherer(logger l.Logger, config *Config) *FreeMetricGatherer {

	logger = ensureLogger(logger, config.Debug, PROVIDER_FREE, FLAG_FREE)

	command := "/usr/bin/free"
	label := DEFAULT_LABEL

	if config.FreeCommand != EMPTY {
		command = config.FreeCommand
		label = CONFIG_LABEL
	}

	if config.Debug {
		logger.Println("Free gatherer initialized by:", label, "as:", command)
	}

	return &FreeMetricGatherer{
		logger: logger,
		debug:  config.Debug,
		command: command,
	}
}

func (gatherer *FreeMetricGatherer) GetMetrics() ([]Metric, error) {
	output, err := execCommand(gatherer.command)
	if err != nil {
		return nil, err
	}

	var metrics = []Metric{}

	lines := strings.Split(output, NEWLINE)
	line1 := lines[1]
	line2 := lines[2]

	var total uint64
	var free uint64
	var used uint64
	var shared uint64
	var buffer uint64
	var available uint64
	var mem string

	fmt.Sscanf(line1, "%s %d %d %d %d %d %d", &mem, &total, &used, &free, &shared, &buffer, &available)

	if gatherer.debug {
		gatherer.logger.Printf("name %s total %d, used %d, free %d,"+
			" shared %d , buffer %d, available %d\n", mem, total, used, free, shared, buffer, available)
	}

	metrics = append(metrics, metric{LEVEL, MetricValue(free), "mFreeLvl", PROVIDER_RAM})
	metrics = append(metrics, metric{LEVEL, MetricValue(used), "mUsedLvl", PROVIDER_RAM})
	metrics = append(metrics, metric{LEVEL, MetricValue(shared), "mSharedLvl", PROVIDER_RAM})
	metrics = append(metrics, metric{LEVEL, MetricValue(buffer), "mBufLvl", PROVIDER_RAM})
	metrics = append(metrics, metric{LEVEL, MetricValue(available), "mAvailableLvl", PROVIDER_RAM})

	totalF := float64(total)

	freePercent := (float64(free) / totalF) * 100.0
	metrics = append(metrics, metric{LEVEL_PERCENT, MetricValue(int64(freePercent)), "mFreePer", PROVIDER_RAM})

	usedPercent := (float64(used) / totalF) * 100.0
	metrics = append(metrics, metric{LEVEL_PERCENT, MetricValue(int64(usedPercent)), "mUsedPer", PROVIDER_RAM})

	fmt.Sscanf(line2, "%s %d %d %d", &mem, &total, &used, &free)

	if free == 0 && used == 0 && total == 0 {
		// do nothing
	} else {
		metrics = append(metrics, metric{LEVEL, MetricValue(free), "mSwpFreeLvl", PROVIDER_RAM})
		metrics = append(metrics, metric{LEVEL, MetricValue(used), "mSwpUsedLvl", PROVIDER_RAM})
		metrics = append(metrics, metric{LEVEL, MetricValue(shared), "mSwpSharedLvl", PROVIDER_RAM})

		totalF = float64(total)
		freePercent = (float64(free) / totalF) * 100.0
		metrics = append(metrics, metric{LEVEL_PERCENT, MetricValue(int64(freePercent)), "mSwpFreePer", PROVIDER_RAM})
		usedPercent = (float64(used) / totalF) * 100.0
		metrics = append(metrics, metric{LEVEL_PERCENT, MetricValue(int64(usedPercent)), "mSwpUsedPer", PROVIDER_RAM})
	}

	return metrics, nil

}
