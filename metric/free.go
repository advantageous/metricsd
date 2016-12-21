package metric

import (
	"fmt"
	l "github.com/advantageous/go-logback/logging"
	"os/exec"
	"runtime"
	"strings"
)

type FreeMetricGatherer struct {
	logger l.Logger
	debug  bool
}

func NewFreeMetricGatherer(logger l.Logger, config *Config) *FreeMetricGatherer {

	if logger == nil {
		if config.Debug {
			logger = l.NewSimpleDebugLogger("free")
		} else {
			logger = l.GetSimpleLogger("MT_FREE_DEBUG", "free")
		}

	}
	return &FreeMetricGatherer{
		logger: logger,
		debug: config.Debug,
	}
}

func (disk *FreeMetricGatherer) GetMetrics() ([]Metric, error) {
	var metrics = []Metric{}

	var output string

	var command string
	if runtime.GOOS == "linux" {
		command = "/usr/bin/free"
	} else if runtime.GOOS == "darwin" {
		command = "/usr/local/bin/free"
	}
	if out, err := exec.Command(command).Output(); err != nil {
		return nil, err
	} else {
		output = string(out)
	}

	lines := strings.Split(output, "\n")
	line1 := lines[1]
	line2 := lines[2]

	var total uint64
	var free uint64
	var used uint64
	var shared uint64
	var buffer uint64
	var available uint64
	var mem string

	fmt.Sscanf(line1, "%s %d %d %d %d %d %d", &mem, &total, &free, &used, &shared, &buffer, &available)

	metrics = append(metrics, metric{LEVEL, MetricValue(free), "mFree", "ram"})
	metrics = append(metrics, metric{LEVEL, MetricValue(used), "mUsed", "ram"})
	metrics = append(metrics, metric{LEVEL, MetricValue(shared), "mShared","ram"})
	metrics = append(metrics, metric{LEVEL, MetricValue(buffer), "mBuf","ram"})
	metrics = append(metrics, metric{LEVEL, MetricValue(available), "mAvailable","ram"})

	totalF := float64(total)

	freePercent := (float64(free) / totalF) * 100.0
	metrics = append(metrics, metric{LEVEL, MetricValue(int64(freePercent)), "mFreePer","ram"})

	usedPercent := (float64(used) / totalF) * 100.0
	metrics = append(metrics, metric{LEVEL, MetricValue(int64(usedPercent)), "mUsedPer","ram"})

	fmt.Sscanf(line2, "%s %d %d %d", &mem, &total, &free, &used)

	if free == 0 && used == 0 && total == 0 {

	} else {
		metrics = append(metrics, metric{LEVEL, MetricValue(free), "mSwpFree","ram"})
		metrics = append(metrics, metric{LEVEL, MetricValue(used), "mSwpUsed","ram"})
		metrics = append(metrics, metric{LEVEL, MetricValue(shared), "mSwpShared","ram"})

		totalF = float64(total)
		freePercent = (float64(free) / totalF) * 100.0
		metrics = append(metrics, metric{LEVEL, MetricValue(int64(freePercent)), "mSwpFreePer","ram"})
		usedPercent = (float64(used) / totalF) * 100.0
		metrics = append(metrics, metric{LEVEL, MetricValue(int64(usedPercent)), "mSwpUsedPer","ram"})
	}

	return metrics, nil

}
