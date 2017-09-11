[![CircleCI](https://circleci.com/gh/a4t/packer-post-processor-amazon-launchconfiguration/tree/master.svg?style=svg)](https://circleci.com/gh/a4t/packer-post-processor-amazon-launchconfiguration/tree/master)

# packer-post-processor-amazon-launchconfiguration

Packer post-processor plugin for Launchconfiguration rotate management

## Description

The Amazon Launchconfig Rotate is AWS AutoScaling deploy management tool.
Packer makes only AMI.
This Plugin is create Launchconfiguration and update AutoScaling Group.

## Installing

```
$ git clone git@github.com:a4t/packer-post-processor-amazon-launchconfiguration.git

or

$ wget -P /tmp https://github.com/a4t/packer-post-processor-amazon-launchconfiguration/releases/download/v0.0.2/packer-post-processor-amazon-launchconfiguration-v0.0.2-linux-amd64.zip
$ unzip packer-post-processor-amazon-launchconfiguration-v0.0.2-linux-amd64.zip
$ mv linux-amd64/packer-post-processor-amazon-launchconfiguration ~/.packer.d/plugins or Your PATH
```

## Build

```
$ make deps
$ make

or

$ make cross-build
```

## Setting

```
$ cp -rp example/packer.json{.sample,}
$ vim example/packer.json
```

### Basic

```
{
  "builders": [{
    "type": "amazon-ebs",
    "vpc_id": "vpc-xxxxxxxxxx",
    "subnet_id": "subnet-yyyyyyyyyy",
    "region": "ap-northeast-1",
    "source_ami": "ami-4af5022c",
    "instance_type": "t2.micro",
    "ssh_username": "ec2-user",
    "ami_name": "amazon-launchconfiguration-{{timestamp}}"
  }],
  "post-processors":[{
    "type": "amazon-launchconfiguration",
    "config_name_prefix": "my-service-",
    "instance_type": "c4.large",
    "keep_releases": 3
  }]
}
```

### Update AutoScalingGroup

```
{
  "builders": [{
    "type": "amazon-ebs",
    "vpc_id": "vpc-xxxxxxxxxx",
    "subnet_id": "subnet-yyyyyyyyyy",
    "region": "ap-northeast-1",
    "source_ami": "ami-4af5022c",
    "instance_type": "t2.micro",
    "ssh_username": "ec2-user",
    "ami_name": "amazon-launchconfiguration-{{timestamp}}"
  }],
  "post-processors":[{
    "type": "amazon-launchconfiguration",
    "config_name_prefix": "my-service-",
    "instance_type": "c4.large",
    "keep_releases": 3,
    "auto_scaling_group_names": [
      "hogehoge",
      "mogemoge"
    ]
  }]
}

```

### more options

- Required
  - config_name_prefix
    - LaunchConfiguration prefix name
  - instance_type
  - keep_releases
    - LaunchConfiguration rotate cycle

- Auth
  - profile_name
  - aws_access_key_id
  - aws_secret_access_key

- Setting
  - iam_instance_profile
  - classic_link_vpc_id
  - classic_link_vpc_security_groups
  - associate_public_ip_address
  - ebs_optimized
  - kernel_id
  - key_name
  - placement_tenancy
  - ramdisk_id
  - security_groups
  - spot_price
  - user_data
    - file://example/userdata
    - #!/bin/bash\n\necho mogemoge

- Update autoscaling group
  - auto_scaling_group_name

