package metric

import (
	"io/ioutil"
	"github.com/hashicorp/hcl"
	l "github.com/advantageous/metricsd/logger"
)

type Config struct {
	AWSRegion      string 	`hcl:"aws_region"`
	EC2InstanceId  string 	`hcl:"ec2_instance_id"`
	Debug          bool    	`hcl:"debug"`
	Local          bool    	`hcl:"local"`
	MetricPrefix  string 	`hcl:"metric_prefix"`
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

	return config, nil

}