provider "aws" {
  region = "us-east-1"
}

terraform {
  required_providers {
    aws = "5.94.1"
  }
}

resource "aws_rds_cluster" "postgresql" {
  cluster_identifier      = "aurora-cluster-demo"
  engine                  = "aurora-postgresql"
  availability_zones      = ["us-east-1a", "us-east-1b", "us-east-1d"]
  database_name           = "mydb"
  master_username         = "foo"
  master_password         = "bar12345678"
  backup_retention_period = 5
  preferred_backup_window = "07:00-09:00"
  skip_final_snapshot     = true
}
