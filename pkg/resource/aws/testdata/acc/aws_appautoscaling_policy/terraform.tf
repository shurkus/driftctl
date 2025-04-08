provider "aws" {
    region = "us-east-1"
}

terraform {
    required_providers {
        aws = "5.94.1"
    }
}

resource "aws_dynamodb_table" "basic-dynamodb-table" {
    name           = "GameScores"
    billing_mode   = "PROVISIONED"
    read_capacity  = 20
    write_capacity = 20
    hash_key       = "UserId"
    range_key      = "GameTitle"

    attribute {
        name = "UserId"
        type = "S"
    }

    attribute {
        name = "GameTitle"
        type = "S"
    }

    attribute {
        name = "TopScore"
        type = "N"
    }

    ttl {
        attribute_name = "TimeToExist"
        enabled        = false
    }

    global_secondary_index {
        name               = "GameTitleIndex"
        hash_key           = "GameTitle"
        range_key          = "TopScore"
        write_capacity     = 10
        read_capacity      = 10
        projection_type    = "INCLUDE"
        non_key_attributes = ["UserId"]
    }

    tags = {
        Name        = "dynamodb-table-1"
        Environment = "production"
    }
}

resource "aws_appautoscaling_target" "dynamodb_table_read_target" {
    max_capacity       = 100
    min_capacity       = 5
    resource_id        = "table/${aws_dynamodb_table.basic-dynamodb-table.name}"
    scalable_dimension = "dynamodb:table:ReadCapacityUnits"
    service_namespace  = "dynamodb"
}

resource "aws_appautoscaling_policy" "dynamodb_table_read_policy" {
    name               = "DynamoDBReadCapacityUtilization:${aws_appautoscaling_target.dynamodb_table_read_target.resource_id}"
    policy_type        = "TargetTrackingScaling"
    resource_id        = aws_appautoscaling_target.dynamodb_table_read_target.resource_id
    scalable_dimension = aws_appautoscaling_target.dynamodb_table_read_target.scalable_dimension
    service_namespace  = aws_appautoscaling_target.dynamodb_table_read_target.service_namespace

    target_tracking_scaling_policy_configuration {
        predefined_metric_specification {
            predefined_metric_type = "DynamoDBReadCapacityUtilization"
        }

        target_value = 70
    }
}
