provider "aws" {
  region = "us-east-1"
}

terraform {
  required_providers {
    aws = "5.94.1"
  }
}

resource "aws_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.vpc.id

  tags = {
    Name = "main"
  }
}

resource "aws_default_route_table" "default" {
  default_route_table_id = aws_vpc.vpc.default_route_table_id

  route {
    cidr_block = "10.1.1.0/24"
    gateway_id = aws_internet_gateway.main.id
  }

  route {
    cidr_block = "10.1.2.0/24"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = {
    Name = "default_table"
  }
}

resource "aws_route_table" "table2" {
  vpc_id = aws_vpc.vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  route {
    ipv6_cidr_block = "::/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = {
    Name = "table2"
  }
}

resource "aws_route_table" "table1" {
  vpc_id = aws_vpc.vpc.id

  tags = {
    Name = "table1"
  }
}

resource "aws_route" "route1" {
  route_table_id = aws_route_table.table1.id
  gateway_id = aws_internet_gateway.main.id
  destination_cidr_block = "1.1.1.1/32"
}

resource "aws_route" "route_v6" {
  route_table_id = aws_route_table.table1.id
  gateway_id = aws_internet_gateway.main.id
  destination_ipv6_cidr_block = "::/0"
}

resource "aws_route_table" "table3" {
  vpc_id = aws_vpc.vpc.id

  tags = {
    Name = "table3"
  }
}
