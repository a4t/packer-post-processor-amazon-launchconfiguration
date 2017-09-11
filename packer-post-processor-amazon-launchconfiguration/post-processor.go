package amazonlaunchconfiguration

import (
	"errors"
	"regexp"

	"github.com/aws/aws-sdk-go/service/autoscaling"
	awscommon "github.com/hashicorp/packer/builder/amazon/common"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
)

type Config struct {
	common.PackerConfig          `mapstructure:",squash"`
	awscommon.AccessConfig       `mapstructure:",squash"`
	ProfileName                  string                          `mapstructure:"profile_name"`
	AwsAccessKeyId               string                          `mapstructure:"aws_access_key_id"`
	AwsSecretAccessKey           string                          `mapstructure:"aws_secret_access_key"`
	ConfigNamePrefix             string                          `mapstructure:"config_name_prefix"`
	AutoScalingGroupName         string                          `mapstructure:"auto_scaling_group_name"`
	InstanceType                 string                          `mapstructure:"instance_type"`
	KeepReleases                 int                             `mapstructure:"keep_releases"`
	IamInstanceProfile           string                          `mapstructure:"iam_instance_profile"`
	ClassicLinkVPCId             string                          `mapstructure:"classic_link_vpc_id"`
	ClassicLinkVPCSecurityGroups []*string                       `mapstructure:"classic_link_vpc_security_groups"`
	AssociatePublicIpAddress     bool                            `mapstructure:"associate_public_ip_address"`
	EbsOptimized                 bool                            `mapstructure:"ebs_optimized"`
	KernelId                     string                          `mapstructure:"kernel_id"`
	KeyName                      string                          `mapstructure:"key_name"`
	PlacementTenancy             string                          `mapstructure:"placement_tenancy"`
	RamdiskId                    string                          `mapstructure:"ramdisk_id"`
	SecurityGroups               []*string                       `mapstructure:"security_groups"`
	SpotPrice                    string                          `mapstructure:"spot_price"`
	UserData                     string                          `mapstructure:"user_data"`
	InstanceMonitoring           *autoscaling.InstanceMonitoring `mapstructure:"instance_monitoring"`
	ctx                          interpolate.Context
}

type PostProcessor struct {
	config Config
}

func amazonLaunchconfigRotate(ui packer.Ui, artifact packer.Artifact, p *PostProcessor) {
	amiId := p.GetImageId(artifact)
	region := p.GetRegion(artifact)
	autoscalingClient := autoscalingClient(region, &p.config)

	createLaunchConfigurationName := CreateLaunchConfiguration(autoscalingClient, amiId, &p.config)
	ui.Say(*createLaunchConfigurationName)

	deleteLaunchConfigurations := RotateLaunchConfiguration(autoscalingClient, &p.config)
	for _, value := range deleteLaunchConfigurations {
		ui.Say("Delete: LaunchConfigName " + value)
	}

	if p.config.AutoScalingGroupName != "" {
		UpdateAutoScalingGroup(autoscalingClient, *createLaunchConfigurationName, &p.config)
		ui.Say("Update: AutoScalingGroupName " + p.config.AutoScalingGroupName)
	}
}

func (p *PostProcessor) Configure(raws ...interface{}) error {
	p.config.ctx.Funcs = awscommon.TemplateFuncs
	err := config.Decode(&p.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &p.config.ctx,
	}, raws...)
	if err != nil {
		return err
	}

	var errs *packer.MultiError
	if p.config.ConfigNamePrefix == "" {
		errs = packer.MultiErrorAppend(errs,
			errors.New("`config_name_prefix' must be specified."))
	}

	if p.config.InstanceType == "" {
		errs = packer.MultiErrorAppend(errs,
			errors.New("`instance_type' must be specified."))
	}

	if p.config.KeepReleases == 0 {
		errs = packer.MultiErrorAppend(errs,
			errors.New("`keep_releases' must be specified."))
	}

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}
	return nil
}

func (p *PostProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, error) {
	amazonLaunchconfigRotate(ui, artifact, p)
	return artifact, true, nil
}

func (p *PostProcessor) GetImageId(artifact packer.Artifact) string {
	r, _ := regexp.Compile("ami-[a-z0-9]+")
	amiId := r.FindString(artifact.Id())
	return amiId
}

func (p *PostProcessor) GetRegion(artifact packer.Artifact) string {
	r, _ := regexp.Compile(`^[a-z0-9\-]+`)
	region := r.FindString(artifact.Id())
	return region
}
