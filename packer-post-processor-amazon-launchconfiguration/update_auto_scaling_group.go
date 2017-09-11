package amazonlaunchconfiguration

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

func UpdateAutoScalingGroup(autoscalingClient *autoscaling.AutoScaling, launchConfigurationName string, config *Config) (bool, error) {
	for i := 0; i < len(config.AutoScalingGroupNames); i++ {
		params := &autoscaling.UpdateAutoScalingGroupInput{
			AutoScalingGroupName:    config.AutoScalingGroupNames[i],
			LaunchConfigurationName: aws.String(launchConfigurationName),
		}
		_, err := autoscalingClient.UpdateAutoScalingGroup(params)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	return true, nil
}
