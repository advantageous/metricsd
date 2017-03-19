package gatherer

import (
	l "github.com/advantageous/go-logback/logging"
	c "github.com/cloudurable/metricsd/common"
	"strings"
)

const (
	DiskField_totalk = "totalk"
	DiskField_usedk = "usedk"
	DiskField_availablek = "availablek"
	DiskField_usedpct = "usedpct"
	DiskField_availablepct = "availablepct"
	DiskField_capacitypct = "capacitypct"
	DiskField_mount = "mount"
)

type DiskMetricsGatherer struct {
	logger l.Logger
	debug bool
	command string
	includes []diskInclude
	fields []string
}

type diskInclude struct {
	starts bool
	value string
}

func NewDiskMetricsGatherer(logger l.Logger, config *c.Config) *DiskMetricsGatherer {

	logger = c.EnsureLogger(logger, config.Debug, c.PROVIDER_DISK, c.FLAG_DISK)
	command, includes, fields := readDiskConfig(config, logger)

	return &DiskMetricsGatherer{
		logger: logger,
		debug: config.Debug,
		command: command,
		includes: includes,
		fields: fields,
	}
}

func readDiskConfig(config *c.Config, logger l.Logger) (string, []diskInclude, []string) {
	command :=  "/usr/bin/df"
	label := c.DEFAULT_LABEL

	// DiskCommand
	if config.DiskCommand != c.EMPTY {
		command = config.DiskCommand
		label = c.CONFIG_LABEL
	}

	// DiskFileSystem
	var includes = []diskInclude{}
	var includesLabel string
	var includesString string
	if config.DiskFileSystems != nil && len(config.DiskFileSystems) > 0 {
		includesLabel = c.CONFIG_LABEL
		includesString = c.ArrayToString(config.DiskFileSystems)
		for _, dfs := range config.DiskFileSystems {
			if strings.HasSuffix(dfs, "*") {
				includes = append(includes, diskInclude{true, dfs[:len(dfs)-1]})
			} else {
				includes = append(includes, diskInclude{false, dfs})
			}
		}
	} else {
		includesLabel = c.DEFAULT_LABEL
		includesString = "/dev/*"
		includes = append(includes, diskInclude{true, "/dev/"})
	}

	// DiskFields
	var fields []string
	var fieldsLabel string
	var fieldsString string
	if config.DiskFields != nil && len(config.DiskFields) > 0 {
		fieldsLabel = c.CONFIG_LABEL
		fields = config.DiskFields
	} else {
		fieldsLabel = c.DEFAULT_LABEL
		fields = []string{DiskField_availablepct}
	}
	fieldsString = c.ArrayToString(fields)

	if config.Debug {
		logger.Println("Disk gatherer initialized by:", label, "as:", command, "with includes by:", includesLabel, "of", includesString, "with fields by:", fieldsLabel, "of", fieldsString)
	}

	return command, includes, fields;
}

func (disk *DiskMetricsGatherer) GetMetrics() ([]c.Metric, error) {

	output, err := c.ExecCommand(disk.command, "-P", "-k", "-l") // P for posix compatibility output, k for 1K blocks, l for local only
	if err != nil {
		return nil, err
	}

	var metrics = []c.Metric{}
	first := true // skip first line
	for _, line := range strings.Split(output, c.NEWLINE) {
		if first {
			first = false
		} else if disk.shouldReportDisk(line) {
			metrics = disk.appendDf(metrics, line)
		}
	}

	return metrics, nil

}

func (disk *DiskMetricsGatherer) shouldReportDisk(line string) bool {
	fsname := c.SplitGetFieldByIndex(line, 0)
	for _,include := range disk.includes {
		if include.starts {
			if strings.HasPrefix(fsname, include.value) {
				return true
			}
		} else {
			if fsname == include.value {
				return true
			}
		}
	}
	return false
}

func (disk *DiskMetricsGatherer) appendDf(metrics []c.Metric, line string) []c.Metric {

	// Filesystem     1024-blocks    Used Available Capacity Mounted on
	// udev               4019524       0   4019524       0% /dev
	// tmpfs               808140    9648    798492       2% /run
	// /dev/sda5         88339720 9322112  74507144      12% /
	// tmpfs              4040700  119244   3921456       3% /dev/shm

	valuesOnly := strings.Fields(line)
	name := valuesOnly[0]
	total := c.ToInt64(valuesOnly[1], 0)
	used := c.ToInt64(valuesOnly[2], 0)
	available := c.ToInt64(valuesOnly[3], 0)
	capacity := c.ToInt64( valuesOnly[4][0:len(valuesOnly[4])-1], 0)
	mount := valuesOnly[5]

	var totalF = float64(total)

	var upct = c.Percent(float64(used), totalF)
	var apct = c.Percent(float64(available), totalF)
	var urnd = c.Round(upct)
	var arnd = c.Round(apct)

	if disk.debug {
		disk.logger.Printf("name %s, total %d, used %d, available %d, usedpct %2.2f (%d), availablepct %2.2f (%d), capacity %d, mount %s",
			                name,    total,    used,    available,    upct, urnd,         apct, arnd,              capacity,    mount)
	}

	for _,field := range disk.fields {
		switch field {
		case DiskField_totalk:
			metrics = append(metrics, *c.NewMetricInt(c.MT_SIZE_KB, total, "diskTotalK:" + name, c.PROVIDER_DISK))
		case DiskField_usedk:
			metrics = append(metrics, *c.NewMetricInt(c.MT_SIZE_KB, used, "diskUsedK:" + name, c.PROVIDER_DISK))
		case DiskField_availablek:
			metrics = append(metrics, *c.NewMetricInt(c.MT_SIZE_KB, available, "diskAvailableK:" + name, c.PROVIDER_DISK))
		case DiskField_usedpct:
			metrics = append(metrics, *c.NewMetricInt(c.MT_PERCENT, urnd, "diskUsedPct:" + name, c.PROVIDER_DISK))
		case DiskField_availablepct:
			metrics = append(metrics, *c.NewMetricInt(c.MT_PERCENT, arnd, "diskAvailPct:" + name, c.PROVIDER_DISK))
		case DiskField_capacitypct:
			metrics = append(metrics, *c.NewMetricInt(c.MT_PERCENT, capacity, "diskCapacityPct:" + name, c.PROVIDER_DISK))
		case DiskField_mount:
			metrics = append(metrics, *c.NewMetricString(mount, "diskAvailMount:" + name, c.PROVIDER_DISK))
		}
	}

	return metrics
}
