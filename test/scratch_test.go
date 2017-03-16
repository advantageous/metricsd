package test

import (
	"testing"
	c "github.com/cloudurable/metricsd/common"
	"fmt"
)

func TestScratch(test *testing.T) {

	logger := GetTestLogger(test, "scratch")
	config := GetTestConfig(logger)

	fmt.Println(c.ObjectToString(config))

	config = &c.Config{
		Debug: false,
		DiskCommand: "df",
		DiskFileSystems: []string{"/dev/*", "udev"},
		DiskFields: []string{"total", "used", "available", "usedpct", "availablepct", "capacitypct", "mount"},
	}

	fmt.Println(c.ObjectToString(config))
}
