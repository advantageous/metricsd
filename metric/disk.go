package metric

import (
	"fmt"
	l "github.com/advantageous/go-logback/logging"
	"runtime"
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
	args []string
}

func NewDiskMetricsGatherer(logger l.Logger, config *Config) *DiskMetricsGatherer {

	logger = ensureLogger(logger, config.Debug, PROVIDER_DISK, FLAG_DISK)

	command :=  "/usr/bin/df"
	args := []string{"-B", "512"}
	label := LINUX_LABEL
	argText := "-B 512"

	if config.DiskCommand != EMPTY {
		command = config.DiskCommand
		if (config.DiskArgs != EMPTY) {
			args = strings.Split(config.DiskArgs, SPACE)
		}
		label = CONFIG_LABEL
		argText = config.DiskArgs
	} else if runtime.GOOS == GOOS_DARWIN {
		command = "/bin/df"
		args = []string{"-b", "-l"}
		label = DARWIN_LABEL
		argText = "-b -l"
	}

	if config.Debug {
		logger.Println("Disk gatherer initialized by:", label, "as:", command, argText)
	}

	return &DiskMetricsGatherer{
		logger: logger,
		debug: config.Debug,
		command: command,
		args: args,
	}
}

func (disk *DiskMetricsGatherer) GetMetrics() ([]Metric, error) {

	output, err := execCommand(disk.command, disk.args...)
	if err != nil {
		return nil, err
	}

	var metrics = []Metric{}

	// TODO read config disk_filesystems

	for _, line := range strings.Split(output, NEWLINE) {
		if strings.HasPrefix(line, "/dev/") {

			var name string
			var total, used, available uint64
			fmt.Sscanf(line, "%s %d %d %d", &name, &total, &used, &available)
			var totalF, availableF float64

			totalF = float64(total)
			availableF = float64(available)

			var calc = availableF / totalF * 100.0

			if disk.debug {
				disk.logger.Printf("name %s total %d used %d available %d calc %2.2f\n",
					name, total, used, available, calc)
			}

			metrics = append(metrics, metric{
				metricType: LEVEL_PERCENT,
				name:       "dUVol" + name[5:] + "AvailPer",
				value:      MetricValue(calc),
				provider:   PROVIDER_DISK,
			})

		}
	}

	return metrics, nil

}
