provider "aws" {
  region = "us-east-1"
}

terraform {
  required_providers {
    aws = "5.94.1"
  }
}

resource "aws_ebs_volume" "foo" {
    availability_zone = "us-east-1a"
    size              = 10

    tags = {
        Name = "Foo Volume"
    }
}
