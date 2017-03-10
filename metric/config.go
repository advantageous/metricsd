package metric

import (
	l "github.com/advantageous/go-logback/logging"
	"github.com/hashicorp/hcl"
	"io/ioutil"
	"time"
)

type Config struct {
	AWSRegion			string			`hcl:"aws_region"`
	ServerRole			string			`hcl:"server_role"`
	IpAddress			string			`hcl:"ip_address"`
	EC2InstanceId		string			`hcl:"ec2_instance_id"`
	EC2InstanceNameTag	string			`hcl:"ec2_instance_name"`
	Debug				bool			`hcl:"debug"`
	Local             	bool			`hcl:"local"`
	NameSpace         	string			`hcl:"namespace"`
	Env               	string			`hcl:"env"`

	TimePeriodSeconds 	time.Duration	`hcl:"interval_seconds"`
	ReadConfigSeconds 	time.Duration	`hcl:"interval_read_config_seconds"`

	DiskGather        	bool			`hcl:"disk_gather"`
	DiskCommand       	string			`hcl:"disk_command"`
	DiskIncludes		string			`hcl:"disk_includes"`

	CpuGather         	bool			`hcl:"cpu_gather"`
	CpuProcStat       	string			`hcl:"cpu_proc_stat"`

	FreeGather        	bool			`hcl:"free_gather"`
	FreeCommand       	string			`hcl:"free_command"`

	NodetoolGather    	bool			`hcl:"nodetool_gather"`
	NodetoolCommand   	string			`hcl:"nodetool_command"`
	NodetoolFunctions 	string			`hcl:"nodetool_functions"`
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
	return c1.AWSRegion      == c2.AWSRegion &&
		c1.ServerRole            == c2.ServerRole &&
		c1.IpAddress             == c2.IpAddress &&
		c1.EC2InstanceId         == c2.EC2InstanceId &&
		c1.EC2InstanceNameTag    == c2.EC2InstanceNameTag &&
		c1.Debug                 == c2.Debug &&
		c1.Local                 == c2.Local &&
		c1.NameSpace             == c2.NameSpace &&
		c1.Env                   == c2.Env &&
		c1.TimePeriodSeconds     == c2.TimePeriodSeconds &&
		c1.ReadConfigSeconds     == c2.ReadConfigSeconds &&
		c1.DiskGather            == c2.DiskGather &&
		c1.DiskCommand           == c2.DiskCommand &&
		c1.DiskIncludes          == c2.DiskIncludes &&
		c1.CpuGather             == c2.CpuGather &&
		c1.CpuProcStat           == c2.CpuProcStat &&
		c1.FreeGather            == c2.FreeGather &&
		c1.FreeCommand           == c2.FreeCommand &&
		c1.NodetoolGather        == c2.NodetoolGather &&
		c1.NodetoolCommand       == c2.NodetoolCommand &&
		c1.NodetoolFunctions     == c2.NodetoolFunctions
}

func ConfigJsonString(cfg *Config) (string) {
	return "{" +
		Jstr("AWSRegion", cfg.AWSRegion, false) +
		Jstr("ServerRole", cfg.ServerRole, false) +
		Jstr("IpAddress", cfg.IpAddress, false) +
		Jstr("EC2InstanceId", cfg.EC2InstanceId, false) +
		Jstr("EC2InstanceNameTag", cfg.EC2InstanceNameTag, false) +
		Jbool("Debug", cfg.Debug, false) +
		Jbool("Local", cfg.Local, false) +
		Jstr("NameSpace", cfg.NameSpace, false) +
		Jstr("Env", cfg.Env, false) +
		Jdur("TimePeriodSeconds", cfg.TimePeriodSeconds, false) +
		Jdur("ReadConfigSeconds", cfg.ReadConfigSeconds, false) +
		Jbool("DiskGather", cfg.DiskGather, false) +
		Jstr("DiskCommand", cfg.DiskCommand, false) +
		Jstr("DiskIncludes", cfg.DiskIncludes, false) +
		Jbool("CpuGather", cfg.CpuGather, false) +
		Jstr("CpuProcStat", cfg.CpuProcStat, false) +
		Jbool("FreeGather", cfg.FreeGather, false) +
		Jstr("FreeCommand", cfg.FreeCommand, false) +
		Jbool("NodetoolGather", cfg.NodetoolGather, false) +
		Jstr("NodetoolCommand", cfg.NodetoolCommand, false) +
		Jstr("NodetoolFunctions", cfg.NodetoolFunctions, true) +
		"}"
}
