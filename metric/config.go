package metric

import (
	l "github.com/advantageous/go-logback/logging"
	"github.com/hashicorp/hcl"
	"io/ioutil"
	"time"
)

type Config struct {
	AWSRegion          string        `hcl:"aws_region"`
	ServerRole         string        `hcl:"server_role"`
	IpAddress          string        `hcl:"ip_address"`
	EC2InstanceId      string        `hcl:"ec2_instance_id"`
	EC2InstanceNameTag string        `hcl:"ec2_instance_name"`
	Debug              bool          `hcl:"debug"`
	Local              bool          `hcl:"local"`
	NameSpace          string        `hcl:"namespace"`
	Env                string        `hcl:"env"`
	TimePeriodSeconds  time.Duration `hcl:"interval_seconds"`
	ReadConfigSeconds  time.Duration `hcl:"interval_read_config_seconds"`
}

func LoadConfig(filename string, logger l.Logger) (*Config, error) {

	if logger == nil {
		logger = l.NewSimpleLogger("config")
	}

	logger.Printf("Loading config %s", filename)

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

	if config.TimePeriodSeconds == 0 {
		config.TimePeriodSeconds = 30
	}

	if config.ReadConfigSeconds == 0 {
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
