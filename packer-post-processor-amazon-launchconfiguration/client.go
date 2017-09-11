package amazonlaunchconfiguration

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

func getClientSession(region string, config *Config) *session.Session {
	if len(config.ProfileName) != 0 {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			Config:  aws.Config{Region: aws.String(region)},
			Profile: config.ProfileName,
		}))
		return sess
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	return sess
}

func autoscalingClient(region string, config *Config) *autoscaling.AutoScaling {
	if len(config.AwsAccessKeyId) != 0 {
		creds := credentials.NewStaticCredentials(config.AwsAccessKeyId, config.AwsSecretAccessKey, "")
		svc := autoscaling.New(
			getClientSession(region, config),
			aws.NewConfig().WithRegion(region).WithCredentials(creds),
		)
		return svc
	}

	svc := autoscaling.New(getClientSession(region, config))
	return svc
}
