package repeater

import (
	lg "github.com/advantageous/go-logback/logging"
	m "github.com/advantageous/metricsd/metric"
	"github.com/advantageous/metricsd/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"time"
)

type AwsCloudMetricRepeater struct {
	logger lg.Logger
	conn   *cloudwatch.CloudWatch
	config *m.Config
}

func (cw AwsCloudMetricRepeater) ProcessMetrics(context m.MetricContext, metrics []m.Metric) error {

	timestamp := aws.Time(time.Now())

	createDatum := func(name string, provider string) *cloudwatch.MetricDatum {

		dimensions := make([]*cloudwatch.Dimension, 0, 3)

		if context.SendId() {
			instanceIdDim := &cloudwatch.Dimension{
				Name:  aws.String("InstanceId"),
				Value: aws.String(cw.config.EC2InstanceId),
			}
			dimensions = append(dimensions, instanceIdDim)

			if cw.config.IpAddress != "" {
				ipDim := &cloudwatch.Dimension{
					Name:  aws.String("IpAddress"),
					Value: aws.String(cw.config.IpAddress),
				}
				dimensions = append(dimensions, ipDim)
			}

			if cw.config.EC2InstanceNameTag != "" {
				dim := &cloudwatch.Dimension{
					Name:  aws.String("InstanceName"),
					Value: aws.String(cw.config.EC2InstanceNameTag),
				}
				dimensions = append(dimensions, dim)
			}
		}
		if context.GetEnv() != "" {
			dim := &cloudwatch.Dimension{
				Name:  aws.String("Environment"),
				Value: aws.String(context.GetEnv()),
			}
			dimensions = append(dimensions, dim)
		}


		if context.GetRole() != "" {
			dim := &cloudwatch.Dimension{
				Name:  aws.String("Role"),
				Value: aws.String(context.GetRole()),
			}
			dimensions = append(dimensions, dim)
		}

		if provider != "" {
			dim := &cloudwatch.Dimension{
				Name:  aws.String("Provider"),
				Value: aws.String(provider),
			}
			dimensions = append(dimensions, dim)
		}

		return &cloudwatch.MetricDatum{
			MetricName: aws.String(name),
			Timestamp:  timestamp,
			Dimensions: dimensions,
		}
	}

	data := []*cloudwatch.MetricDatum{}

	var err error

	for index, d := range metrics {

		if cw.config.Debug {
			cw.logger.Printf("%s %d %d %s", d.GetName(), d.GetType(), d.GetValue(), d.GetProvider())
		}
		value := float64(d.GetValue())
		datum := createDatum(d.GetName(), d.GetProvider())
		datum.Value = aws.Float64(float64(value))

		switch d.GetType() {
		case m.COUNT:
			datum.Unit = aws.String(cloudwatch.StandardUnitCount)
		case m.LEVEL:
			datum.Unit = aws.String(cloudwatch.StandardUnitKilobytes)
		case m.LEVEL_PERCENT:
			datum.Unit = aws.String(cloudwatch.StandardUnitPercent)
		case m.TIMING:
			datum.Unit = aws.String(cloudwatch.StandardUnitMilliseconds)
		default:
			datum.Unit = aws.String(cloudwatch.StandardUnitCount)
		}

		data = append(data, datum)


		if index%20 == 0 && index != 0 {
			data = []*cloudwatch.MetricDatum{}

			if len(data) > 0 {
				request := &cloudwatch.PutMetricDataInput{
					Namespace:  aws.String(context.GetNameSpace()),
					MetricData: data,
				}
				_, err = cw.conn.PutMetricData(request)
				if err != nil {
					cw.logger.PrintError("Error writing metrics", err)
					cw.logger.Error("Error writing metrics", err, index)
				} else {
					if cw.config.Debug {
						cw.logger.Info("SENT..........................")
					}
				}
			}

		}

	}

	if len(data) > 0 {
		request := &cloudwatch.PutMetricDataInput{
			Namespace:  aws.String(context.GetNameSpace()),
			MetricData: data,
		}
		_, err = cw.conn.PutMetricData(request)

	}
	return err
}

func NewAwsCloudMetricRepeater(config *m.Config) AwsCloudMetricRepeater {
	session := util.NewAWSSession(config)
	logger := lg.NewSimpleLogger("log-repeater")
	return AwsCloudMetricRepeater{logger, cloudwatch.New(session), config}
}
