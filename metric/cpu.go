package metric

import (
	"bufio"
	"errors"
	"fmt"
	l "github.com/advantageous/go-logback/logging"
	"os"
	"runtime"
	"strings"
)

type CpuTimeType byte
type CpuTime uint64
type CpuCount uint64

type CPUMetricsGatherer struct {
	path string
	lastTime *CpuStats
	logger l.Logger
	debug bool
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

func NewCPUMetricsGathererWithPath(path string, logger l.Logger, debug bool) *CPUMetricsGatherer {

	logger = ensureLogger(logger, debug, "cpu", "MT_CPU_DEBUG")

	return &CPUMetricsGatherer{
		path:   path,
		logger: logger,
		debug:  debug,
	}
}

func NewCPUMetricsGatherer(logger l.Logger, config *Config) *CPUMetricsGatherer {

	logger = ensureLogger(logger, config.Debug, "cpu", "MT_CPU_DEBUG")

	statFile := "/proc/stat"
	label := LINUX_LABEL

	if config.CpuProcStat != EMPTY {
		statFile = config.CpuProcStat
		label = CONFIG_LABEL
	} else if runtime.GOOS == GOOS_DARWIN {
		dir, _ := os.Getwd()
		statFile = dir + "/metric/test-data/proc/stat"
		label = DARWIN_LABEL
	}

	if config.Debug {
		logger.Println("CPU gatherer initialized by:", label, "as:", statFile)
	}

	return NewCPUMetricsGathererWithPath(statFile, logger, config.Debug)
}

func (cpu *CPUMetricsGatherer) SetPath(path string) {
	cpu.path = path
}

func (cpu *CPUMetricsGatherer) GetMetrics() ([]Metric, error) {

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
	cpu.logger.Debugf("%+v \n", cpuStats)
	return metrics, nil

}

func convertToMetrics(lastTimeStats *CpuStats, nowStats *CpuStats) []Metric {
	var metrics = []Metric{}

	if lastTimeStats != nil {

		softInterruptCount := nowStats.SoftInterruptCount - lastTimeStats.SoftInterruptCount
		if softInterruptCount > 0 {
			metrics = append(metrics, metric{COUNT, MetricValue(softInterruptCount), "softIrqCnt", "cpu", EMPTY})
		}

		interruptCount := nowStats.InterruptCount - lastTimeStats.InterruptCount
		if interruptCount > 0 {
			metrics = append(metrics, metric{COUNT, MetricValue(interruptCount), "intrCnt", "cpu", EMPTY})
		}

		contextSwitchCount := nowStats.ContextSwitchCount - lastTimeStats.ContextSwitchCount
		if contextSwitchCount > 0 {
			metrics = append(metrics, metric{COUNT, MetricValue(contextSwitchCount), "ctxtCnt", "cpu", EMPTY})
		}

		processCount := nowStats.ProcessCount - lastTimeStats.ProcessCount
		if processCount > 0 {
			metrics = append(metrics, metric{COUNT, MetricValue(processCount), "processesStrtCnt", "cpu", EMPTY})
		}

		for index, c := range nowStats.CpuTimeList {

			guest := c.Guest - lastTimeStats.CpuTimeList[index].Guest
			if guest > 0 {
				metrics = append(metrics, metric{COUNT, MetricValue(guest), c.Name + "GuestJif", "cpu", EMPTY})
			}

			user := c.User - lastTimeStats.CpuTimeList[index].User
			if user > 0 {
				metrics = append(metrics, metric{COUNT, MetricValue(user), c.Name + "UsrJif", "cpu", EMPTY})
			}

			idle := c.Idle - lastTimeStats.CpuTimeList[index].Idle
			if idle > 0 {
				metrics = append(metrics, metric{COUNT, MetricValue(idle), c.Name + "IdleJif", "cpu", EMPTY})
			}

			IoWait := c.IoWait - lastTimeStats.CpuTimeList[index].IoWait
			if IoWait > 0 {
				metrics = append(metrics, metric{COUNT, MetricValue(IoWait), c.Name + "IowaitJif", "cpu", EMPTY})
			}

			Irq := c.Irq - lastTimeStats.CpuTimeList[index].Irq
			if Irq > 0 {
				metrics = append(metrics, metric{COUNT, MetricValue(Irq), c.Name + "IrqJif", "cpu", EMPTY})
			}

			GuestNice := c.GuestNice - lastTimeStats.CpuTimeList[index].GuestNice
			if GuestNice > 0 {
				metrics = append(metrics, metric{COUNT, MetricValue(GuestNice), c.Name + "GuestniceJif", "cpu", EMPTY})
			}

			Steal := c.Steal - lastTimeStats.CpuTimeList[index].Steal
			if Steal > 0 {
				metrics = append(metrics, metric{COUNT, MetricValue(Steal), c.Name + "StealJif", "cpu", EMPTY})
			}

			Nice := c.Nice - lastTimeStats.CpuTimeList[index].Nice
			if Nice > 0 {
				metrics = append(metrics, metric{COUNT, MetricValue(Nice), c.Name + "NiceJif", "cpu", EMPTY})
			}

			System := c.System - lastTimeStats.CpuTimeList[index].System
			if System > 0 {
				metrics = append(metrics, metric{COUNT, MetricValue(System), c.Name + "SysJif", "cpu", EMPTY})
			}

			SoftIrq := c.SoftIrq - lastTimeStats.CpuTimeList[index].SoftIrq
			if SoftIrq > 0 {
				metrics = append(metrics, metric{COUNT, MetricValue(SoftIrq), c.Name + "SoftIrqJif", "cpu", EMPTY})
			}

		}

	}

	metrics = append(metrics, metric{COUNT, MetricValue(nowStats.ProcessRunningCount), "procsRunning", "cpu", EMPTY})
	metrics = append(metrics, metric{COUNT, MetricValue(nowStats.ProcessBlockCount), "procsBlocked", "cpu", EMPTY})

	return metrics
}

func (cpu *CPUMetricsGatherer) readCpuStats() (*CpuStats, error) {
	org, err := os.Open(cpu.path)

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

	if cpu.debug {
		cpu.logger.Println("LINE", line)
		cpu.logger.Println(name, value)
	}

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
		if strings.HasPrefix(name, "cpu") {
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
