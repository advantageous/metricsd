package metric

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

func NewCPUMetricsGatherer(logger l.Logger, config *c.Config) *CPUMetricsGatherer {

	logger = c.EnsureLogger(logger, config.Debug, c.PROVIDER_CPU, c.FLAG_CPU)
	procStatPath := readCpuConfig(config, logger)

	return &CPUMetricsGatherer{
		procStatPath: procStatPath,
		logger:       logger,
		debug:        config.Debug,
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

func (cpu *CPUMetricsGatherer) GetMetrics() ([]c.Metric, error) {

	if cpu.debug {
		cpu.logger.Debug("GetMetrics called")
	}

	var cpuStats *CpuStats
	var err error

	if cpuStats, err = cpu.readCpuStats(); err != nil {
		return nil, err
	}

	metrics := convertToMetrics(cpu.lastTime, cpuStats)
	cpu.lastTime = cpuStats
	if (cpu.debug) {
		cpu.logger.Debugf("%+v \n", cpuStats)
	}
	return metrics, nil

}

func convertToMetrics(lastTimeStats *CpuStats, nowStats *CpuStats) []c.Metric {
	var metrics = []c.Metric{}

	if lastTimeStats != nil {

		softInterruptCount := nowStats.SoftInterruptCount - lastTimeStats.SoftInterruptCount
		if softInterruptCount > 0 {
			metrics = append(metrics, c.Metric{c.COUNT, c.MetricValue(softInterruptCount), "softIrqCnt", c.PROVIDER_CPU})
		}

		interruptCount := nowStats.InterruptCount - lastTimeStats.InterruptCount
		if interruptCount > 0 {
			metrics = append(metrics, c.Metric{c.COUNT, c.MetricValue(interruptCount), "intrCnt", c.PROVIDER_CPU})
		}

		contextSwitchCount := nowStats.ContextSwitchCount - lastTimeStats.ContextSwitchCount
		if contextSwitchCount > 0 {
			metrics = append(metrics, c.Metric{c.COUNT, c.MetricValue(contextSwitchCount), "ctxtCnt", c.PROVIDER_CPU})
		}

		processCount := nowStats.ProcessCount - lastTimeStats.ProcessCount
		if processCount > 0 {
			metrics = append(metrics, c.Metric{c.COUNT, c.MetricValue(processCount), "processesStrtCnt", c.PROVIDER_CPU})
		}

		for index, cput := range nowStats.CpuTimeList {

			guest := cput.Guest - lastTimeStats.CpuTimeList[index].Guest
			if guest > 0 {
				metrics = append(metrics, c.Metric{c.COUNT, c.MetricValue(guest), cput.Name + "GuestJif", c.PROVIDER_CPU})
			}

			user := cput.User - lastTimeStats.CpuTimeList[index].User
			if user > 0 {
				metrics = append(metrics, c.Metric{c.COUNT, c.MetricValue(user), cput.Name + "UsrJif", c.PROVIDER_CPU})
			}

			idle := cput.Idle - lastTimeStats.CpuTimeList[index].Idle
			if idle > 0 {
				metrics = append(metrics, c.Metric{c.COUNT, c.MetricValue(idle), cput.Name + "IdleJif", c.PROVIDER_CPU})
			}

			IoWait := cput.IoWait - lastTimeStats.CpuTimeList[index].IoWait
			if IoWait > 0 {
				metrics = append(metrics, c.Metric{c.COUNT, c.MetricValue(IoWait), cput.Name + "IowaitJif", c.PROVIDER_CPU})
			}

			Irq := cput.Irq - lastTimeStats.CpuTimeList[index].Irq
			if Irq > 0 {
				metrics = append(metrics, c.Metric{c.COUNT, c.MetricValue(Irq), cput.Name + "IrqJif", c.PROVIDER_CPU})
			}

			GuestNice := cput.GuestNice - lastTimeStats.CpuTimeList[index].GuestNice
			if GuestNice > 0 {
				metrics = append(metrics, c.Metric{c.COUNT, c.MetricValue(GuestNice), cput.Name + "GuestniceJif", c.PROVIDER_CPU})
			}

			Steal := cput.Steal - lastTimeStats.CpuTimeList[index].Steal
			if Steal > 0 {
				metrics = append(metrics, c.Metric{c.COUNT, c.MetricValue(Steal), cput.Name + "StealJif", c.PROVIDER_CPU})
			}

			Nice := cput.Nice - lastTimeStats.CpuTimeList[index].Nice
			if Nice > 0 {
				metrics = append(metrics, c.Metric{c.COUNT, c.MetricValue(Nice), cput.Name + "NiceJif", c.PROVIDER_CPU})
			}

			System := cput.System - lastTimeStats.CpuTimeList[index].System
			if System > 0 {
				metrics = append(metrics, c.Metric{c.COUNT, c.MetricValue(System), cput.Name + "SysJif", c.PROVIDER_CPU})
			}

			SoftIrq := cput.SoftIrq - lastTimeStats.CpuTimeList[index].SoftIrq
			if SoftIrq > 0 {
				metrics = append(metrics, c.Metric{c.COUNT, c.MetricValue(SoftIrq), cput.Name + "SoftIrqJif", c.PROVIDER_CPU})
			}

		}

	}

	metrics = append(metrics, c.Metric{c.COUNT, c.MetricValue(nowStats.ProcessRunningCount), "procsRunning", c.PROVIDER_CPU})
	metrics = append(metrics, c.Metric{c.COUNT, c.MetricValue(nowStats.ProcessBlockCount), "procsBlocked", c.PROVIDER_CPU})

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
