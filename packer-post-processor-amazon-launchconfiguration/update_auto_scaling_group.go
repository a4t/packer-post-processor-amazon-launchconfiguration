package amazonlaunchconfiguration

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

func UpdateAutoScalingGroup(autoscalingClient *autoscaling.AutoScaling, launchConfigurationName string, config *Config) *autoscaling.UpdateAutoScalingGroupOutput {
	params := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName:    aws.String(config.AutoScalingGroupName),
		LaunchConfigurationName: aws.String(launchConfigurationName),
	}
	res, err := autoscalingClient.UpdateAutoScalingGroup(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return res
}
