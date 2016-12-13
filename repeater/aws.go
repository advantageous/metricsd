package repeater

import (
	m "github.com/advantageous/metricsd/metric"
	lg "github.com/advantageous/metricsd/logger"
	"github.com/advantageous/metricsd/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"time"
	"strings"
)

type AwsCloudMetricRepeater struct {
	logger lg.Logger
	conn   *cloudwatch.CloudWatch
	config *m.Config
}

func (cw AwsCloudMetricRepeater)ProcessMetrics(metrics []m.Metric) error {

	timestamp := aws.Time(time.Now())

	aDatum := func(name string) *cloudwatch.MetricDatum {
		return &cloudwatch.MetricDatum{
			MetricName: aws.String(name),
			Timestamp:  timestamp,
		}
	}

	data := []*cloudwatch.MetricDatum{}

	var err error

	for index, d := range metrics {

		if cw.config.Debug {
			cw.logger.Printf("%s %d %d", d.GetName(), d.GetType(), d.GetValue())
		}

		switch(d.GetType()) {
		case m.COUNT:
			value := float64(d.GetValue())
			datum := aDatum(d.GetName())
			if strings.HasSuffix(d.GetName(), "Per") {
				datum.Unit = aws.String(cloudwatch.StandardUnitCount)
			} else {
				datum.Unit = aws.String(cloudwatch.StandardUnitPercent)
			}
			datum.Value = aws.Float64(float64(value))
			data = append(data, datum)
		case m.LEVEL:
			value := float64(d.GetValue())
			datum := aDatum(d.GetName())
			if strings.HasSuffix(d.GetName(), "Per") {
				datum.Unit = aws.String(cloudwatch.StandardUnitKilobytes)
			} else {
				datum.Unit = aws.String(cloudwatch.StandardUnitPercent)
			}
			datum.Value = aws.Float64(float64(value))
			data = append(data, datum)
		case m.TIMING:
			value := float64(d.GetValue())
			datum := aDatum(d.GetName())
			datum.Unit = aws.String(cloudwatch.StandardUnitMilliseconds)
			datum.Value = aws.Float64(float64(value))
			data = append(data, datum)

		}

		if index % 20 == 0  && index != 0{
			data = []*cloudwatch.MetricDatum{}

			if  len (data) > 0 {
				request := &cloudwatch.PutMetricDataInput{
					Namespace:  aws.String(cw.config.MetricPrefix),
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


	if  len (data) > 0 {
		request := &cloudwatch.PutMetricDataInput{
			Namespace:  aws.String(cw.config.MetricPrefix),
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