package gatherer

import (
	l "github.com/advantageous/go-logback/logging"
	c "github.com/cloudurable/metricsd/common"
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type CpuTimeType byte
type CpuTime uint64
type CpuCount uint64

type CPUMetricsGatherer struct {
	procStatPath string
	lastTime     *CpuStats
	logger       l.Logger
	debug        bool
	reportZeros  bool
}

type CpuStats struct {
	CpuTimeList         []CpuTimes
	ContextSwitchCount  CpuCount
	BootTime            CpuTime
	ProcessCount        CpuCount
	ProcessRunningCount CpuCount
	ProcessBlockCount   CpuCount
	InterruptCount      CpuCount
	SoftInterruptCount  CpuCount
}

type CpuTimes struct {
	Name      string
	User      CpuTime
	Nice      CpuTime
	System    CpuTime
	Idle      CpuTime
	IoWait    CpuTime
	Irq       CpuTime
	SoftIrq   CpuTime
	Steal     CpuTime
	Guest     CpuTime
	GuestNice CpuTime
}

//cpu  13761633 12805 3121528 171746421 158717 0 7708 0 0 0
//cpu0 3446824 3018 816772 42845024 25964 0 712 0 0 0
//cpu1 3461972 4316 817690 42760371 93670 0 2989 0 0 0
//cpu2 3395036 2706 758030 43069668 14793 0 3312 0 0 0
//cpu3 3457800 2764 729034 43071357 24288 0 693 0 0 0
//intr 676492475 20 15 0 0 0 0 0 0 1 76 0 0 329 0 0 0 66768 0 0 0 0 0 0 1211524 0 0 1932502 2739244 14 17870042 661 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
//ctxt 1663680895
//btime 1489105326
//processes 3326310
//procs_running 2
//procs_blocked 0
//softirq 370470175 1215735 135209363 21947 3127466 1690346 0 174465 124692639 0 104338214

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

	if cpuStats, err = cpu.readCpuStats(); err != nil {
		return nil, err
	}

	metrics := cpu.convertToMetrics(cpu.lastTime, cpuStats)
	cpu.lastTime = cpuStats
	if (cpu.debug) {
		cpu.logger.Debugf("%+v \n", cpuStats)
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

		for index, cput := range nowStats.CpuTimeList {
			metrics = cpu.appendCount(metrics, "GuestJif", int64(cput.Guest - lastTimeStats.CpuTimeList[index].Guest))
			metrics = cpu.appendCount(metrics, "UsrJif", int64(cput.User - lastTimeStats.CpuTimeList[index].User))
			metrics = cpu.appendCount(metrics, "IdleJif", int64(cput.Idle - lastTimeStats.CpuTimeList[index].Idle))
			metrics = cpu.appendCount(metrics, "IowaitJif", int64(cput.IoWait - lastTimeStats.CpuTimeList[index].IoWait))
			metrics = cpu.appendCount(metrics, "IrqJif", int64(cput.Irq - lastTimeStats.CpuTimeList[index].Irq))
			metrics = cpu.appendCount(metrics, "GuestniceJif", int64(cput.GuestNice - lastTimeStats.CpuTimeList[index].GuestNice))
			metrics = cpu.appendCount(metrics, "StealJif", int64(cput.Steal - lastTimeStats.CpuTimeList[index].Steal))
			metrics = cpu.appendCount(metrics, "NiceJif", int64(cput.Nice - lastTimeStats.CpuTimeList[index].Nice))
			metrics = cpu.appendCount(metrics, "SysJif", int64(cput.System - lastTimeStats.CpuTimeList[index].System))
			metrics = cpu.appendCount(metrics, "SoftIrqJif", int64(cput.SoftIrq - lastTimeStats.CpuTimeList[index].SoftIrq))
		}
	}

	metrics = append(metrics, *c.NewMetricInt(c.MT_COUNT, int64(nowStats.ProcessRunningCount), "procsRunning", c.PROVIDER_CPU))
	metrics = append(metrics, *c.NewMetricInt(c.MT_COUNT, int64(nowStats.ProcessBlockCount), "procsBlocked", c.PROVIDER_CPU))

	return metrics
}

func (cpu *CPUMetricsGatherer) readCpuStats() (*CpuStats, error) {
	org, err := os.Open(cpu.procStatPath)

	fd := bufio.NewReader(org)

	if err != nil {
		cpu.logger.Emergencyf("Error reading file %v", err)
	}

	stats := CpuStats{}
	stats.CpuTimeList = make([]CpuTimes, 0)

	for {
		var name string
		var value uint64
		var line string

		if bytes, _, err := fd.ReadLine(); err == nil {
			line = string(bytes)
		} else if err.Error() == "EOF" {
			//Error EOF is ok
			cpu.logger.Debug("EOF while reading /proc/stat file")
			break
		} else {
			//Others errors are not ok
			cpu.logger.PrintError("Error reading line from /proc/stat", err)
			return nil, err
		}

		if count, err := fmt.Sscanf(line, "%s %d ", &name, &value); err != nil {
			cpu.logger.PrintError("Error scanning name / value from a line from /proc/stat", err)
			return nil, err
		} else if count == 0 {
			cpu.logger.Debug("Count was 0 when scanning /proc/stat line")
			break

		}

		if err = cpu.parseLine(name, value, line, &stats); err != nil {
			return nil, err
		}

		cpu.logger.Debugf("%s %d", name, value)
	}
	cpu.logger.Debugf("%+v \n", stats)
	return &stats, nil
}

func (cpu *CPUMetricsGatherer) parseLine(name string, value uint64, line string, stats *CpuStats) error {

	//if cpu.debug {
	//	cpu.logger.Println("LINE", line)
	//	cpu.logger.Println(name, value)
	//}

	switch name {
	case "ctxt":
		stats.ContextSwitchCount = CpuCount(value)
	case "btime":
		stats.BootTime = CpuTime(value)
	case "processes":
		stats.ProcessCount = CpuCount(value)
	case "procs_running":
		stats.ProcessBlockCount = CpuCount(value)
	case "procs_blocked":
		stats.ProcessBlockCount = CpuCount(value)
	case "intr":
		stats.InterruptCount = CpuCount(value)
	case "softirq":
		stats.SoftInterruptCount = CpuCount(value)
	default:
		if strings.HasPrefix(name, c.PROVIDER_CPU) {
			t := CpuTimes{}
			t.Name = name
			t.User = CpuTime(value)
			count, err := fmt.Sscanf(line, "%s %d %d %d %d %d %d %d %d %d %d", &t.Name,
				&t.User, &t.Nice, &t.System,
				&t.Idle, &t.IoWait, &t.Irq,
				&t.SoftIrq, &t.Steal, &t.Guest,
				&t.GuestNice)

			if cpu.debug {
				cpu.logger.Printf("Name = %s,\t User=%d,\t Nice=%d,\t System=%d \n"+
					"Idle = %d,\t IoWait = %d,\t Irq = %d,\tSftIrq=%d \n"+
					"Steal = %d,\t Guest=%d,\t GuestNice=%d",
					t.Name, t.User, t.Nice, t.System,
					&t.Idle, &t.IoWait, &t.Irq, &t.SoftIrq,
					&t.Steal, &t.Guest, &t.GuestNice)
			}

			if err != nil {
				cpu.logger.PrintError("Failure parsing cpu stats", err)
				return err
			}

			if count != 11 {
				cpu.logger.Errorf("cpu scan amount is off expected 11, but got %d", count)
				return errors.New("Unable to scan cpu times")
			}
			stats.CpuTimeList = append(stats.CpuTimeList, t)
		} else {
			return fmt.Errorf("Not sure what this is %s", name)
		}
	}
	return nil
}
