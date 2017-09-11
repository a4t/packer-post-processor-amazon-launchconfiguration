package amazonlaunchconfiguration

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

func CreateLaunchConfiguration(autoscalingClient *autoscaling.AutoScaling, amiId string, config *Config) *string {
	params := &autoscaling.CreateLaunchConfigurationInput{
		LaunchConfigurationName: aws.String(config.ConfigNamePrefix + time.Now().Format("20060102150405")),
		ImageId:                 aws.String(amiId),
		InstanceType:            aws.String(config.InstanceType),
	}

	params = setNonRequireParams(params, config)
	_, err := autoscalingClient.CreateLaunchConfiguration(params)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	log.Println("Create: LaunchConfigName " + *params.LaunchConfigurationName)
	return params.LaunchConfigurationName
}

func setNonRequireParams(params *autoscaling.CreateLaunchConfigurationInput, config *Config) *autoscaling.CreateLaunchConfigurationInput {
	params = setIamInstanceProfile(params, config)
	params = setAssociatePublicIpAddress(params, config)
	params = setClassicLinkVPCId(params, config)
	params = setClassicLinkVPCSecurityGroups(params, config)
	params = setEbsOptimized(params, config)
	params = setInstanceMonitoring(params, config)
	params = setKernelId(params, config)
	params = setKeyName(params, config)
	params = setPlacementTenancy(params, config)
	params = setRamdiskId(params, config)
	params = setSecurityGroups(params, config)
	params = setSpotPrice(params, config)
	params = setUserData(params, config)

	return params
}

func setIamInstanceProfile(params *autoscaling.CreateLaunchConfigurationInput, config *Config) *autoscaling.CreateLaunchConfigurationInput {
	if len(config.IamInstanceProfile) != 0 {
		params.IamInstanceProfile = aws.String(config.IamInstanceProfile)
	}
	return params
}

func setAssociatePublicIpAddress(params *autoscaling.CreateLaunchConfigurationInput, config *Config) *autoscaling.CreateLaunchConfigurationInput {
	if !config.AssociatePublicIpAddress {
		params.AssociatePublicIpAddress = aws.Bool(config.AssociatePublicIpAddress)
	}
	return params
}

func setClassicLinkVPCId(params *autoscaling.CreateLaunchConfigurationInput, config *Config) *autoscaling.CreateLaunchConfigurationInput {
	if len(config.ClassicLinkVPCId) != 0 {
		params.ClassicLinkVPCId = aws.String(config.ClassicLinkVPCId)
	}
	return params
}

func setClassicLinkVPCSecurityGroups(params *autoscaling.CreateLaunchConfigurationInput, config *Config) *autoscaling.CreateLaunchConfigurationInput {
	if len(config.ClassicLinkVPCSecurityGroups) != 0 {
		params.ClassicLinkVPCSecurityGroups = config.ClassicLinkVPCSecurityGroups
	}
	return params
}

func setEbsOptimized(params *autoscaling.CreateLaunchConfigurationInput, config *Config) *autoscaling.CreateLaunchConfigurationInput {
	if !config.EbsOptimized {
		params.EbsOptimized = aws.Bool(config.EbsOptimized)
	}
	return params
}

func setInstanceMonitoring(params *autoscaling.CreateLaunchConfigurationInput, config *Config) *autoscaling.CreateLaunchConfigurationInput {
	if config.InstanceMonitoring != nil {
		params.InstanceMonitoring = config.InstanceMonitoring
	}
	return params
}

func setKernelId(params *autoscaling.CreateLaunchConfigurationInput, config *Config) *autoscaling.CreateLaunchConfigurationInput {
	if len(config.KernelId) != 0 {
		params.KernelId = aws.String(config.KernelId)
	}
	return params
}

func setKeyName(params *autoscaling.CreateLaunchConfigurationInput, config *Config) *autoscaling.CreateLaunchConfigurationInput {
	if len(config.KeyName) != 0 {
		params.KeyName = aws.String(config.KeyName)
	}
	return params
}

func setPlacementTenancy(params *autoscaling.CreateLaunchConfigurationInput, config *Config) *autoscaling.CreateLaunchConfigurationInput {
	if len(config.PlacementTenancy) != 0 {
		params.PlacementTenancy = aws.String(config.PlacementTenancy)
	}
	return params
}

func setRamdiskId(params *autoscaling.CreateLaunchConfigurationInput, config *Config) *autoscaling.CreateLaunchConfigurationInput {
	if len(config.RamdiskId) != 0 {
		params.RamdiskId = aws.String(config.RamdiskId)
	}
	return params
}

func setSecurityGroups(params *autoscaling.CreateLaunchConfigurationInput, config *Config) *autoscaling.CreateLaunchConfigurationInput {
	if len(config.SecurityGroups) != 0 {
		params.SecurityGroups = config.SecurityGroups
	}
	return params
}

func setSpotPrice(params *autoscaling.CreateLaunchConfigurationInput, config *Config) *autoscaling.CreateLaunchConfigurationInput {
	if len(config.SpotPrice) != 0 {
		params.SpotPrice = aws.String(config.SpotPrice)
	}
	return params
}

func setUserData(params *autoscaling.CreateLaunchConfigurationInput, config *Config) *autoscaling.CreateLaunchConfigurationInput {
	if len(config.UserData) != 0 {
		userDataBase64 := getBase64UserData(config.UserData)
		params.UserData = aws.String(userDataBase64)
	}
	return params
}

func getBase64UserData(userData string) string {
	if strings.Index(userData, "file://") == 0 {
		filepath := userData[7:]
		_, err := os.Stat(filepath)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		file, _ := os.Open(filepath)
		defer file.Close()

		fi, _ := file.Stat()
		size := fi.Size()

		data := make([]byte, size)
		file.Read(data)
		return base64.StdEncoding.EncodeToString(data)
	} else {
		return base64.StdEncoding.EncodeToString([]byte(userData))
	}
}
