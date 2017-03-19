package nodetool

import (
	"strings"
	c "github.com/cloudurable/metricsd/common"
)

func Cfstats(nodetoolCommand string) ([]c.Metric, error) {
	output, err := c.ExecCommand(nodetoolCommand, NtFunc_cfstats)
	if err != nil {
		return nil, err
	}

	var metrics = []c.Metric{}

	currentKeyspace := c.EMPTY
	currentTable := c.EMPTY

	lines := strings.Split(output, c.NEWLINE)
	for _,line := range lines {
		if line != c.EMPTY  && !strings.HasPrefix(line, "---") {
			if strings.HasPrefix(line, "Total number") {
				value := c.SplitGetLastField(line)
				metrics = append(metrics, *c.NewMetricIntString(c.MT_COUNT, value, "ntCfTotalTables", c.PROVIDER_NODETOOL))
			} else {
				fields := fields(line)
				if strings.HasPrefix(fields[0], "Keyspace") {
					currentKeyspace = fields[1]
					currentTable = c.EMPTY
				} else if strings.HasPrefix(fields[0], "Table") {
					currentTable = fields[1]
				} else if strings.HasPrefix(fields[0], "SStables") {
					// ignore
				} else {
					lastIndex := c.GetLastIndex(fields)
					if lastIndex != -1 {
						mt := c.MT_COUNT
						temp := fields[lastIndex]
						if temp == "ms." || temp == "ms" {
							mt = c.MT_MILLIS
							lastIndex -= 1
						}
						value := fields[lastIndex]
						name := cfName(fields, lastIndex, currentKeyspace, currentTable)
						var metric c.Metric
						if strings.Contains(value, c.DOT) {
							metric = *c.NewMetricFloatString(mt, value, name, c.PROVIDER_NODETOOL)
						} else {
							metric = *c.NewMetricIntString(mt, value, name, c.PROVIDER_NODETOOL)
						}
						metrics = append(metrics, metric)
					}
				}
			}
		}
	}

	return metrics, nil
}
func cfName(fields []string, stopIndex int, currentKeyspace string, currentTable string) string {
	name := "ntCf:" + currentKeyspace
	if currentTable != c.EMPTY {
		name = name + ":" + currentTable
	}
	name = name + ":"
	for x := 0; x < stopIndex; x++ {
		name = name + c.UpFirst(fields[x])
	}
	return name
}

func fields(s string) []string {
	fields := []string{}

	current := c.EMPTY

	for x := 0; x < len(s); x++ {
		bite := s[x]
		if (bite >= 'A' && bite <= 'Z') || (bite >= 'a' && bite <= 'z') || (bite >= '0' && bite <= '9') || bite == '.' || bite == '-' || bite == '_' {
			current = current + s[x:x+1]

		} else if current != c.EMPTY {
			fields = append(fields, current)
			current = c.EMPTY
		}
	}

	if current != c.EMPTY {
		fields = append(fields, current)
		current = c.EMPTY
	}
	return fields
}
