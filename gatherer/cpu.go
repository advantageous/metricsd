package gatherer

import (
	l "github.com/advantageous/go-logback/logging"
	c "github.com/cloudurable/metricsd/common"
	"bufio"
	"os"
	"strings"
)

type CPUMetricsGatherer struct {
	procStatPath string
	lastStats    *CpuStats
	logger       l.Logger
	debug        bool
	reportZeros  bool
}

type CpuTimes struct {
	User      int64
	Nice      int64
	System    int64
	Idle      int64
	IoWait    int64
	Irq       int64
	SoftIrq   int64
	Steal     int64
	Guest     int64
	GuestNice int64
}

type CpuStats struct {
	CpuMap              map[string]CpuTimes
	ContextSwitchCount  int64
	BootTime            int64
	ProcessCount        int64
	ProcessRunningCount int64
	ProcessBlockCount   int64
	InterruptCount      int64
	SoftInterruptCount  int64
}

func NewCPUMetricsGatherer(logger l.Logger, config *c.Config) *CPUMetricsGatherer {

	logger = c.EnsureLogger(logger, config.Debug, c.PROVIDER_CPU, c.FLAG_CPU)
	procStatPath := readCpuConfig(config, logger)

	return &CPUMetricsGatherer{
		procStatPath: procStatPath,
		logger:       logger,
		debug:        config.Debug,
		reportZeros:  config.CpuReportZeros,
	}
}

func readCpuConfig(config *c.Config, logger l.Logger) (string) {
	procStatPath := "/proc/stat"
	label := c.DEFAULT_LABEL
	if config.CpuProcStat != c.EMPTY {
		procStatPath = config.CpuProcStat
		label = c.CONFIG_LABEL
	}
	if config.Debug {
		logger.Println("c.PROVIDER_CPU gatherer initialized by:", label, "as:", procStatPath)
	}

	return procStatPath
}

func (cpu *CPUMetricsGatherer) TestingChangeProcStatPath(inProcStatPath string) {
	cpu.procStatPath = inProcStatPath
}

func (cpu *CPUMetricsGatherer) GetMetrics() ([]c.Metric, error) {

	if cpu.debug {
		cpu.logger.Debug("GetMetrics called")
	}

	var cpuStats *CpuStats
	var err error

	if cpuStats, err = cpu.readProcStat(); err != nil {
		return nil, err
	}

	metrics := cpu.convertToMetrics(cpu.lastStats, cpuStats)
	cpu.lastStats = cpuStats
	if (cpu.debug) {
		cpu.logger.Debug(c.ObjectToString(cpuStats))
	}
	return metrics, nil
}

func (cpu *CPUMetricsGatherer) appendCount(metrics []c.Metric, name string, count int64) []c.Metric {
	if cpu.reportZeros || count > 0 {
		metrics = append(metrics, *c.NewMetricInt(c.MT_COUNT, count, name, c.PROVIDER_CPU))
	}
	return metrics
}

func (cpu *CPUMetricsGatherer) convertToMetrics(lastTimeStats *CpuStats, nowStats *CpuStats) []c.Metric {
	var metrics = []c.Metric{}

	if lastTimeStats != nil {

		metrics = cpu.appendCount(metrics, "softIrqCnt", int64(nowStats.SoftInterruptCount - lastTimeStats.SoftInterruptCount))
		metrics = cpu.appendCount(metrics, "intrCnt", int64(nowStats.InterruptCount - lastTimeStats.InterruptCount))
		metrics = cpu.appendCount(metrics, "ctxtCnt", int64(nowStats.ContextSwitchCount - lastTimeStats.ContextSwitchCount))
		metrics = cpu.appendCount(metrics, "processesStrtCnt", int64(nowStats.ProcessCount - lastTimeStats.ProcessCount))

		for cpuName, nowCpuTimes := range nowStats.CpuMap {
			lastCpuTimes, found := lastTimeStats.CpuMap[cpuName]
			if !found {
				lastCpuTimes = CpuTimes{}
			}
			metrics = cpu.appendCount(metrics, "GuestJif", int64(nowCpuTimes.Guest - lastCpuTimes.Guest))
			metrics = cpu.appendCount(metrics, "UsrJif", int64(nowCpuTimes.User - lastCpuTimes.User))
			metrics = cpu.appendCount(metrics, "IdleJif", int64(nowCpuTimes.Idle - lastCpuTimes.Idle))
			metrics = cpu.appendCount(metrics, "IowaitJif", int64(nowCpuTimes.IoWait - lastCpuTimes.IoWait))
			metrics = cpu.appendCount(metrics, "IrqJif", int64(nowCpuTimes.Irq - lastCpuTimes.Irq))
			metrics = cpu.appendCount(metrics, "GuestniceJif", int64(nowCpuTimes.GuestNice - lastCpuTimes.GuestNice))
			metrics = cpu.appendCount(metrics, "StealJif", int64(nowCpuTimes.Steal - lastCpuTimes.Steal))
			metrics = cpu.appendCount(metrics, "NiceJif", int64(nowCpuTimes.Nice - lastCpuTimes.Nice))
			metrics = cpu.appendCount(metrics, "SysJif", int64(nowCpuTimes.System - lastCpuTimes.System))
			metrics = cpu.appendCount(metrics, "SoftIrqJif", int64(nowCpuTimes.SoftIrq - lastCpuTimes.SoftIrq))
		}
	}

	metrics = append(metrics, *c.NewMetricInt(c.MT_COUNT, int64(nowStats.ProcessRunningCount), "procsRunning", c.PROVIDER_CPU))
	metrics = append(metrics, *c.NewMetricInt(c.MT_COUNT, int64(nowStats.ProcessBlockCount), "procsBlocked", c.PROVIDER_CPU))

	return metrics
}

func (cpu *CPUMetricsGatherer) readProcStat() (*CpuStats, error) {
	org, err := os.Open(cpu.procStatPath)
	fd := bufio.NewReader(org)
	if err != nil {
		cpu.logger.Emergencyf("Error reading file %v", err)
	}

	stats := CpuStats{}
	stats.CpuMap = make(map[string]CpuTimes)

	for {
		theLine := c.EMPTY
		bytes, _, err := fd.ReadLine()
		if err != nil {
			if err.Error() != "EOF" { // EOF error is ok, other errors are not ok
				cpu.logger.PrintError("Error reading line from /proc/stat", err)
				return nil, err
			}
			break // EOF, leave the for loop
		}
		if len(bytes) == 0 {
			break
		}

		theLine = string(bytes)

		valuesOnly := strings.Fields(theLine)
		lineName := valuesOnly[0]
		value := c.ToInt64(valuesOnly[1], 0)

		switch lineName {
		case "ctxt":          stats.ContextSwitchCount = value
		case "btime":         stats.BootTime = value
		case "processes":     stats.ProcessCount = value
		case "procs_running": stats.ProcessRunningCount = value
		case "procs_blocked": stats.ProcessBlockCount = value
		case "intr":          stats.InterruptCount = value
		case "softirq":       stats.SoftInterruptCount = value
		default:
			if strings.HasPrefix(lineName, "cpu") {
				cpuTimes := CpuTimes{}
				for i := 1; i < len(valuesOnly); i++ {
					value = c.ToInt64(valuesOnly[i], 0)
					switch i {
					case  1: cpuTimes.User = value
					case  2: cpuTimes.Nice = value
					case  3: cpuTimes.System = value
					case  4: cpuTimes.Idle = value
					case  5: cpuTimes.IoWait = value
					case  6: cpuTimes.Irq = value
					case  7: cpuTimes.SoftIrq = value
					case  8: cpuTimes.Steal = value
					case  9: cpuTimes.Guest = value
					case 10: cpuTimes.GuestNice = value
					default:
						if (cpu.debug) {
							cpu.logger.Debug("Unknown cpu time column, index:", i, "found in", theLine)
						}
					}
				}
				stats.CpuMap[lineName] = cpuTimes
			} else {
				if (cpu.debug) {
					cpu.logger.Debug("Unknown Data", theLine)
				}
			}
		}
	}

	return &stats, nil
}

	/*
	cpu  5017 1 3356 7561462 1674 2 53 3 4 5
	cpu0 1105 0 1113 1890502 345 0 35 0 0 0
	cpu1 1291 0 792 1889318 496 0 6 0 0 0
	cpu2 1251 0 713 1890968 482 0 5 0 0 0
	cpu3 1370 0 738 1890674 351 0 7 0 0 0
	intr 1488221 27 0 0 0 348 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 159 0 12047 0 22568 1 0 3143 0 76 0 91494 83893 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
	ctxt 2710734
	btime 1480893191
	processes 9277
	procs_running 3
	procs_blocked 0
	softirq 655105 0 221348 21 39766 0 0 1 215075 0 178894
	*/
