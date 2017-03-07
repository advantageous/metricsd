package util

import (
	l "github.com/advantageous/go-logback/logging"
	m "github.com/cloudurable/metricsd/metric"
	"github.com/aws/aws-sdk-go/aws"
	awsCredentials "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"os"
)

var awsLogger = l.NewSimpleLogger("aws")

func NewAWSSession(cfg *m.Config) *awsSession.Session {

	metaDataClient, session := getClient(cfg)
	credentials := getCredentials(metaDataClient)

	if credentials != nil {

		credentials := getCredentials(metaDataClient)
		readMeta(metaDataClient, cfg, session)

		awsConfig := &aws.Config{
			Credentials: credentials,
			Region:      aws.String(cfg.AWSRegion),
			MaxRetries:  aws.Int(3),
		}
		return awsSession.New(awsConfig)
	} else {
		readMeta(metaDataClient, cfg, session)

		return awsSession.New(&aws.Config{
			Region:     aws.String(cfg.AWSRegion),
			MaxRetries: aws.Int(3),
		})
	}

}

func getClient(config *m.Config) (*ec2metadata.EC2Metadata, *awsSession.Session) {
	if !config.Local {
		awsLogger.Debug("Config NOT set to local using meta-data client to find local")
		var session = awsSession.New(&aws.Config{})
		return ec2metadata.New(session), session
	} else {
		awsLogger.Println("Config set to local")
		return nil, nil
	}
}

func readMeta(client *ec2metadata.EC2Metadata, config *m.Config, session *awsSession.Session) {

	if client == nil {
		awsLogger.Info("Client missing using config to set region")
		if config.AWSRegion == "" {
			awsLogger.Info("AWSRegion missing using default region us-west-2")
			config.AWSRegion = "us-west-2"
		}
	} else {
		region, err := client.Region()
		if err != nil {
			awsLogger.Error("Unable to get region from aws meta client : %s %v", err.Error(), err)
			os.Exit(3)
		}

		config.AWSRegion = region
		config.IpAddress = findLocalIp(client)
		config.EC2InstanceId, err = client.GetMetadata("instance-id")

		config.EC2InstanceNameTag = findInstanceName(config.EC2InstanceId, config.AWSRegion, session)
		if err != nil {
			awsLogger.Error("Unable to get instance id from aws meta client : %s %v", err.Error(), err)
			os.Exit(4)
		}
	}

}
func findLocalIp(metaClient *ec2metadata.EC2Metadata) string {
	ip, err := metaClient.GetMetadata("local-ipv4")

	if err != nil {
		awsLogger.Error("Unable to get private ip address from aws meta client : %s %v", err.Error(), err)
		os.Exit(6)
	}

	return ip

}

func getCredentials(client *ec2metadata.EC2Metadata) *awsCredentials.Credentials {

	if client == nil {
		awsLogger.Info("Client missing credentials not looked up")
		return nil
	} else {
		return awsCredentials.NewChainCredentials([]awsCredentials.Provider{
			&awsCredentials.EnvProvider{},
			&ec2rolecreds.EC2RoleProvider{
				Client: client,
			},
		})
	}

}

func findAZ(metaClient *ec2metadata.EC2Metadata) string {

	az, err := metaClient.GetMetadata("placement/availability-zone")

	if err != nil {
		awsLogger.Errorf("Unable to get az from aws meta client : %s %v", err.Error(), err)
		os.Exit(5)
	}

	return az
}

func findInstanceName(instanceId string, region string, session *awsSession.Session) string {

	var name = "NO_NAME"
	var err error

	ec2Service := ec2.New(session, aws.NewConfig().WithRegion(region))

	params := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId), // Required
			// More values...
		},
	}

	resp, err := ec2Service.DescribeInstances(params)

	if err != nil {
		awsLogger.Errorf("Unable to get instance name tag DescribeInstances failed : %s %v", err.Error(), err)
		return name
	}

	if len(resp.Reservations) > 0 && len(resp.Reservations[0].Instances) > 0 {
		var instance = resp.Reservations[0].Instances[0]
		if len(instance.Tags) > 0 {

			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					return *tag.Value
				}
			}
		}
		awsLogger.Errorf("Unable to get find name tag ")
		return name

	} else {
		awsLogger.Errorf("Unable to get find name tag ")
		return name
	}
}
