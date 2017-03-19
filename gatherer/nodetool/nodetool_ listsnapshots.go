package nodetool

import (
	"strings"
	c "github.com/cloudurable/metricsd/common"
)

func ListSnapshots(nodetoolCommand string) ([]c.Metric, error) {
	output, err := c.ExecCommand(nodetoolCommand, NtFunc_listsnapshots)
	if err != nil {
		return nil, err
	}

	//Snapshot Details:
	//Snapshot Name  Keyspace   Column Family  True Size   Size on Disk
	//<blank line>
	//1387304478196  Keyspace1  Standard1      0 bytes     308.66 MB
	//1387304417755  Keyspace1  Standard1      0 bytes     107.21 MB
	//1387305820866  Keyspace1  Standard2      0 bytes      41.69 MB
	//               Keyspace1  Standard1      0 bytes     308.66 MB
	//<blank line>

	//No snapshots looks like this:
	//
	//Snapshot Details: There are no snapshots

	// uncomment to test
	//output = "" +
	//	"Snapshot Details:\n" +
	//	"Snapshot Name  Keyspace   Column Family  True Size   Size on Disk\n" +
	//	"\n" +
	//	"1387304478196  Keyspace1  Standard1      0 bytes     308.66 MB\n" +
	//	"1387304417755  Keyspace1  Standard1      0 bytes     107.21 MB\n" +
	//	"1387305820866  Keyspace1  Standard2      0 bytes      41.69 MB\n" +
	//	"               Keyspace1  Standard1      0 bytes     308.66 MB\n" +
	//	"\n"

	var metrics = []c.Metric{}

	lines := strings.Split(output, c.NEWLINE)
	var snapCount int64 = 0
	for index, line := range lines {
		if index > 2 && line != c.EMPTY {
			vals := strings.Fields(line)
			if len(vals) == 7 {
				metrics = append(metrics, *c.NewMetricIntString(c.MT_MICROS, vals[1], name("Keyspace", vals[0]), c.PROVIDER_NODETOOL))
				metrics = append(metrics, *c.NewMetricIntString(c.MT_MICROS, vals[2], name("ColumnFamily", vals[0]), c.PROVIDER_NODETOOL))
				metrics = append(metrics, *c.NewMetricIntString(c.ToSizeMetricType(vals[4]), vals[3], name("TrueSize", vals[0]), c.PROVIDER_NODETOOL))
				metrics = append(metrics, *c.NewMetricIntString(c.ToSizeMetricType(vals[6]), vals[5], name("SizeOnDisk", vals[0]), c.PROVIDER_NODETOOL))
			}
		}
	}
	metrics = append(metrics, *c.NewMetricInt(c.MT_COUNT, snapCount, "ntSnapCount", c.PROVIDER_NODETOOL))

	return metrics, nil
}

func name(name string, snapName string) string {
	return "ntSnap" + name + c.COLON + snapName
}

