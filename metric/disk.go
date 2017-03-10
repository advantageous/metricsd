package metric

import (
	"fmt"
	l "github.com/advantageous/go-logback/logging"
	"strings"
)

//
//type DiskStatCount uint64
//type TimeSpent uint64
//
//type DiskStats struct {
//	NumReads              DiskStatCount //Field  1 -- # of reads completed.     This is the total number of reads completed successfully.
//	ReadsMerged           DiskStatCount //Field  2 -- # of reads merged, field 6 -- # of writes merged
//	SectorsRead           DiskStatCount //Field  3 -- # of sectors read --  This is the total number of sectors read successfully.
//	TimeReading           TimeSpent     //Field  4 -- # of milliseconds spent reading . This is the total number of milliseconds spent by all reads (as measured from __make_request() to end_that_request_last()).
//	WriteMerged           DiskStatCount //Field  6 -- # of writes merged. See the description of field 2.
//	SectorsWritten        DiskStatCount //Field  7 -- # of sectors written. This is the total number of sectors written successfully.
//	TimeWriting           TimeSpent     //Field  8 -- # of milliseconds spent writing. This is the total number of milliseconds spent by all writes (as measured from __make_request() to end_that_request_last()).
//	IOCurrentlyInProgress DiskStatCount //Field  9 -- # of I/Os currently in progress. The only field that should go to zero. Incremented as requests are given to appropriate struct request_queue and decremented as they finish.
//	TimeSpendDoingIO      TimeSpent     //Field 10 -- # of milliseconds spent doing I/Os. This field increases so long as field 9 is nonzero.
//}
/**
https://www.kernel.org/doc/Documentation/iostats.txt
Field 11 -- weighted # of milliseconds spent doing I/Os
    This field is incremented at each I/O start, I/O completion, I/O
    merge, or read of these stats by the number of I/Os in progress
    (field 9) times the number of milliseconds spent doing I/O since the
    last update of this field.  This can provide an easy measure of both
    I/O completion time and the backlog that may be accumulating.

*/

type DiskMetricsGatherer struct {
	logger l.Logger
	debug bool
	command string
	includes []diskInclude
}

type diskInclude struct {
	starts bool
	value string
}

func NewDiskMetricsGatherer(logger l.Logger, config *Config) *DiskMetricsGatherer {

	if (!config.DiskGather) { return nil } // don't return anything if not turned on

	logger = EnsureLogger(logger, config.Debug, PROVIDER_DISK, FLAG_DISK)
	command, includes := readDiskConfig(config, logger)

	return &DiskMetricsGatherer{
		logger: logger,
		debug: config.Debug,
		command: command,
		includes: includes,
	}
}

func readDiskConfig(config *Config, logger l.Logger) (string, []diskInclude) {
	command :=  "/usr/bin/df"
	label := DEFAULT_LABEL

	if config.DiskCommand != EMPTY {
		command = config.DiskCommand
		label = CONFIG_LABEL
	}

	var includes = []diskInclude{}
	var includesLabel string
	var includesString string
	if config.DiskIncludes != EMPTY {
		includesLabel = CONFIG_LABEL
		includesString = config.DiskIncludes
		for _,inc := range strings.Split(includesString, SPACE) {
			if strings.HasSuffix(inc, "*") {
				includes = append(includes, diskInclude{true, inc[:len(inc)-1]})
			} else {
				includes = append(includes, diskInclude{false, inc})
			}
		}
	} else {
		includesLabel = DEFAULT_LABEL
		includesString = "/dev/*"
		includes = append(includes, diskInclude{true, "/dev/"})
	}

	if config.Debug {
		logger.Println("Disk gatherer initialized by:", label, "as:", command, "with includes by:", includesLabel, "of:", includesString)
	}

	return command, includes;
}

func (cpu *DiskMetricsGatherer) Reload(config *Config) (ReloadResult) {
	if (!config.DiskGather) { return RELOAD_EJECT }  // eject if not turned on

	cpu.command, cpu.includes = readDiskConfig(config, cpu.logger);
	return RELOAD_SUCCESS
}

func (disk *DiskMetricsGatherer) GetMetrics() ([]Metric, error) {

	output, err := ExecCommand(disk.command, "-k", "-l") // k for 1K, l for local only
	if err != nil {
		return nil, err
	}

	// Filesystem     1K-blocks    Used Available Use% Mounted on
	// udev             4019524       0   4019524   0% /dev
	// tmpfs             808148    9700    798448   2% /run
	// /dev/sda5       88339720 9280388  74548868  12% /
	// tmpfs            4040720  122236   3918484   4% /dev/shm
	// tmpfs               5120       4      5116   1% /run/lock
	// tmpfs            4040720       0   4040720   0% /sys/fs/cgroup
	// tmpfs             808148     120    808028   1% /run/user/1000

	var metrics = []Metric{}
	first := true // skip first line
	for _, line := range strings.Split(output, NEWLINE) {
		if first {
			first = false
		} else if disk.includeDisk(line) {
			metrics = disk.appendDu(metrics, line)
		}
	}

	return metrics, nil

}

func (disk *DiskMetricsGatherer) includeDisk(line string) bool {
	fsname := FieldByIndex(line, 0);
	for _,include := range disk.includes {
		if include.starts {
			if strings.HasPrefix(fsname, include.value) {
				return true;
			}
		} else {
			if fsname == include.value {
				return true;
			}
		}
	}
	return false
}

func (disk *DiskMetricsGatherer) appendDu(metrics []Metric, line string) []Metric {
	var name string
	var total, used, available uint64
	fmt.Sscanf(line, "%s %d %d %d", &name, &total, &used, &available)

	var totalF = float64(total)
	var availableF = float64(available)

	var calc = availableF / totalF * 100.0

	if disk.debug {
		disk.logger.Printf("name %s total %d used %d available %d calc %2.2f\n", name, total, used, available, calc)
	}

	metrics = append(metrics, metric{LEVEL_PERCENT, MetricValue(calc),  "dUVolAvailPct:" + name, PROVIDER_DISK})

	return metrics
}
