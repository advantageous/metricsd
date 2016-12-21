package metric

import (
	l "github.com/advantageous/go-logback/logging"
	"github.com/hashicorp/hcl"
	"io/ioutil"
	"time"
)

type Config struct {
	AWSRegion         string        `hcl:"aws_region"`
	IpAddress         string        `hcl:"ip_address"`
	EC2InstanceId     string        `hcl:"ec2_instance_id"`
	Debug             bool          `hcl:"debug"`
	Local             bool          `hcl:"local"`
	NameSpace         string        `hcl:"namespace"`
	Env               string        `hcl:"env"`
	TimePeriodSeconds time.Duration `hcl:"interval_seconds"`
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

	return config, nil

}
