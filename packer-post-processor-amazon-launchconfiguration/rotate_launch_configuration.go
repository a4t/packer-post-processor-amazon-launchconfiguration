package amazonlaunchconfiguration

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

func RotateLaunchConfiguration(autoscalingClient *autoscaling.AutoScaling, config *Config) []string {
	res := describeLaunchConfiguration(autoscalingClient)
	filteringLaunchConfigurationNames := filteringLaunchConfiguration(res, config)
	deleteLaunchConfigurations := deleteLaunchConfiguration(autoscalingClient, filteringLaunchConfigurationNames, config)
	return deleteLaunchConfigurations
}

func describeLaunchConfiguration(autoscalingClient *autoscaling.AutoScaling) *autoscaling.DescribeLaunchConfigurationsOutput {
	params := &autoscaling.DescribeLaunchConfigurationsInput{
		MaxRecords: aws.Int64(100), // maximum value is 100.
	}
	res, err := autoscalingClient.DescribeLaunchConfigurations(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return res
}

func filteringLaunchConfiguration(res *autoscaling.DescribeLaunchConfigurationsOutput, config *Config) []string {
	var tempNames = make([]string, 1000)
	j := 0
	for i := 0; i < len(res.LaunchConfigurations); i++ {
		if strings.Index(*res.LaunchConfigurations[i].LaunchConfigurationName, config.ConfigNamePrefix) == 0 {
			tempNames[j] = *res.LaunchConfigurations[i].LaunchConfigurationName
			j++
		}
	}
	filteringLaunchConfigurationNames := make([]string, j)
	copy(filteringLaunchConfigurationNames, tempNames)
	return filteringLaunchConfigurationNames
}

func deleteLaunchConfiguration(autoscalingClient *autoscaling.AutoScaling, names []string, config *Config) []string {
	var deleteLaunchConfigurations []string
	sort.Sort(sort.Reverse(sort.StringSlice(names)))
	for i := config.KeepReleases; i < len(names); i++ {
		params := &autoscaling.DeleteLaunchConfigurationInput{
			LaunchConfigurationName: aws.String(names[i]),
		}
		_, err := autoscalingClient.DeleteLaunchConfiguration(params)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		deleteLaunchConfigurations = append(deleteLaunchConfigurations, names[i])
	}
	return deleteLaunchConfigurations
}
