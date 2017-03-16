package common

import (
	l "github.com/advantageous/go-logback/logging"
	"github.com/hashicorp/hcl"
	"io/ioutil"
	"time"
)

type Config struct {
	Env               	string			`hcl:"env"`
	Local             	bool			`hcl:"local"`
	Debug				bool			`hcl:"debug"`
	TimePeriodSeconds 	time.Duration	`hcl:"interval_seconds"`
	ReadConfigSeconds 	time.Duration	`hcl:"interval_read_config_seconds"`

	Repeaters           []string		`hcl:"repeaters"`
	Gatherers           []string		`hcl:"gatherers"`

	AWSRegion			string			`hcl:"aws_region"`
	ServerRole			string			`hcl:"server_role"`
	IpAddress			string			`hcl:"ip_address"`
	EC2InstanceId		string			`hcl:"ec2_instance_id"`
	EC2InstanceNameTag	string			`hcl:"ec2_instance_name"`
	NameSpace         	string			`hcl:"namespace"`

	DiskCommand       	string			`hcl:"disk_command"`
	DiskFileSystems     []string		`hcl:"disk_file_systems"`
	DiskFields          []string		`hcl:"disk_fields"`

	CpuProcStat       	string			`hcl:"cpu_proc_stat"`
	CpuReportZeros     	bool			`hcl:"cpu_report_zeros"`

	FreeCommand       	string			`hcl:"free_command"`

	NodetoolCommand   	string			`hcl:"nodetool_command"`
	NodetoolFunctions 	[]string		`hcl:"nodetool_functions"`
}

func LoadConfig(filename string, logger l.Logger) (*Config, error) {

	if logger == nil {
		logger = l.NewSimpleLogger("config")
	}

	logger.Printf("Loading config %s", filename)

	var err error

	configBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return LoadConfigFromString(string(configBytes), logger)
}

func LoadConfigFromString(data string, logger l.Logger) (*Config, error) {

	if logger == nil {
		logger = l.NewSimpleLogger("config")
	}

	config := &Config{}
	logger.Println("Loading log...")

	err := hcl.Decode(&config, data)
	if err != nil {
		return nil, err
	}

	if config.TimePeriodSeconds <= 0 {
		config.TimePeriodSeconds = 30
	}

	if config.ReadConfigSeconds <= 0 {
		config.ReadConfigSeconds = 60
	}

	if config.NameSpace == "" {
		config.NameSpace = "Linux System"
	}

	return config, nil

}

func (config *Config) GetEnv() string {
	return config.Env
}

func (config *Config) GetNameSpace() string {
	return config.NameSpace
}

func (config *Config) GetRole() string {
	return config.ServerRole
}

func (config *Config) SendId() bool {
	return true
}

func (config *Config) GetNoIdContext() MetricContext {
	return context{
		env:       config.Env,
		namespace: config.NameSpace,
		role:      config.ServerRole,
	}
}

type context struct {
	env       string
	namespace string
	role      string
}

func (c context) GetEnv() string {
	return c.env
}

func (c context) GetNameSpace() string {
	return c.namespace
}

func (c context) GetRole() string {
	return c.role
}

func (c context) SendId() bool {
	return false
}

func ConfigEquals(c1 *Config, c2 *Config) (bool) {
	return c1.Env    == c2.Env &&
		c1.Local     == c2.Local &&
		c1.Debug     == c2.Debug &&
		c1.NameSpace == c2.NameSpace &&

		c1.TimePeriodSeconds == c2.TimePeriodSeconds &&
		c1.ReadConfigSeconds == c2.ReadConfigSeconds &&

		StringArraysEqual(c1.Repeaters, c2.Repeaters) &&
		StringArraysEqual(c1.Gatherers, c2.Gatherers) &&

		c1.AWSRegion          == c2.AWSRegion &&
		c1.ServerRole         == c2.ServerRole &&
		c1.IpAddress          == c2.IpAddress &&
		c1.EC2InstanceId      == c2.EC2InstanceId &&
		c1.EC2InstanceNameTag == c2.EC2InstanceNameTag &&

		c1.DiskCommand           == c2.DiskCommand &&
		c1.CpuProcStat           == c2.CpuProcStat &&
		c1.CpuReportZeros        == c2.CpuReportZeros &&
		c1.FreeCommand           == c2.FreeCommand &&
		c1.NodetoolCommand       == c2.NodetoolCommand &&

		StringArraysEqual(c1.DiskFileSystems, c2.DiskFileSystems) &&
		StringArraysEqual(c1.DiskFields, c2.DiskFields) &&
		StringArraysEqual(c1.NodetoolFunctions, c2.NodetoolFunctions)
}