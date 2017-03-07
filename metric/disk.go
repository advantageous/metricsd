package metric

import (
	"fmt"
	l "github.com/advantageous/go-logback/logging"
	"os/exec"
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
	debug  bool
	diskCommand string
	diskArgs    string
}

func NewDiskMetricsGatherer(logger l.Logger, config *Config) *DiskMetricsGatherer {

	if logger == nil {
		if config.Debug {
			logger = l.NewSimpleDebugLogger("disk")
		} else {
			logger = l.GetSimpleLogger("MT_DISK_DEBUG", "disk")
		}
	}

	return &DiskMetricsGatherer{
		logger: logger,
		debug:  config.Debug,
		diskCommand: config.DiskCommand,
		diskArgs:    config.DiskArgs,
	}
}

func (disk *DiskMetricsGatherer) GetMetrics() ([]Metric, error) {
	var metrics = []Metric{}

	var output string

	command :=  "/usr/bin/df"
	args := []string{"-B", "512"}
	label := "Linux"
	argText := "/usr/bin/df -B 512"

	if disk.diskCommand != "" {
		command = disk.diskCommand
		args = strings.Split(disk.diskArgs, " ")
		label = "Config"
		argText = disk.diskArgs
	} else if runtime.GOOS == "darwin" {
		command = "/bin/df"
		args = []string{"-b", "-l"}
		label = "Darwin"
		argText = "/bin/df -b -l"
	}

	if disk.debug {
		disk.logger.Println(label, "Disk gatherer initialized by:", label, "as:", command, argText)
	}

	if out, err := exec.Command(command, args...).Output(); err != nil {
		return nil, err
	} else {
		output = string(out)
	}

	for _, line := range strings.Split(output, "\n") {
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
				provider:   "disk",
			})

		}
	}

	return metrics, nil

}
