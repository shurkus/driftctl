provider "aws" {
  region = "us-east-1"
}

terraform {
  required_providers {
    aws = "5.94.1"
  }
}

resource "aws_ami" "test-ami" {
    name                = "test"
    virtualization_type = "hvm"
    root_device_name    = "/dev/xvda"

    ebs_block_device {
        device_name = "/dev/xvda"
        snapshot_id = aws_ebs_snapshot.test-ebs-snapshot.id
        volume_size = 10
        iops = 0
        encrypted = false
    }

    timeouts {
        create = "20m"
    }
}

resource "aws_ebs_volume" "test-ebs-volume" {
    availability_zone = "us-east-1a"
    size              = 10

    tags = {
        Name = "HelloWorld"
    }
}

resource "aws_ebs_snapshot" "test-ebs-snapshot" {
    volume_id = aws_ebs_volume.test-ebs-volume.id

    tags = {
        Name = "HelloWorld_snap"
    }
}
