provider "aws" {
  region = "us-east-1"
}

terraform {
  required_providers {
    aws = "5.94.1"
  }
}

resource "aws_api_gateway_rest_api" "foo" {
    name        = "foo"
    description = "This is foo API"
}

resource "aws_api_gateway_resource" "foo" {
    rest_api_id = aws_api_gateway_rest_api.foo.id
    parent_id   = aws_api_gateway_rest_api.foo.root_resource_id
    path_part   = "foo"
}

resource "aws_api_gateway_method" "foo" {
    rest_api_id   = aws_api_gateway_rest_api.foo.id
    resource_id   = aws_api_gateway_resource.foo.id
    http_method   = "GET"
    authorization = "NONE"
}

resource "aws_api_gateway_method_response" "response_200" {
    rest_api_id = aws_api_gateway_rest_api.foo.id
    resource_id = aws_api_gateway_resource.foo.id
    http_method = aws_api_gateway_method.foo.http_method
    status_code = "200"
}

resource "aws_api_gateway_integration" "foo" {
    http_method = aws_api_gateway_method.foo.http_method
    resource_id = aws_api_gateway_resource.foo.id
    rest_api_id = aws_api_gateway_rest_api.foo.id
    type        = "MOCK"
}

resource "aws_api_gateway_integration_response" "foo" {
    rest_api_id = aws_api_gateway_rest_api.foo.id
    resource_id = aws_api_gateway_resource.foo.id
    http_method = aws_api_gateway_method.foo.http_method
    status_code = aws_api_gateway_method_response.response_200.status_code

    # Transforms the backend JSON response to XML
    response_templates = {
        "application/xml" = <<EOF
#set($inputRoot = $input.path('$'))
<?xml version="1.0" encoding="UTF-8"?>
<message>
    $inputRoot.body
</message>
EOF
    }
}

resource "aws_api_gateway_rest_api" "bar" {
    name        = "bar"
    description = "This is bar API"
    body = jsonencode({
        openapi = "3.0.1"
        info = {
            title   = "example"
            version = "1.0"
        }
        paths = {
            "/path1" = {
                get = {
                    x-amazon-apigateway-integration = {
                        httpMethod           = "GET"
                        payloadFormatVersion = "1.0"
                        type                 = "MOCK"
                        responses = {
                            "2\\d{2}" : {
                                "statusCode" : "200",
                                "responseTemplates" : {
                                    "application/json" : "#set ($root=$input.path('$')) { \"stage\": \"$root.name\", \"user-id\": \"$root.key\" }",
                                    "application/xml" : "#set ($root=$input.path('$')) <stage>$root.name</stage> "
                                }
                            },
                        }
                    }
                }
            }
        }
    })
}

resource "aws_api_gateway_rest_api" "baz" {
    name        = "baz"
    description = "This is baz API"
    body = jsonencode({
        swagger = "2.0"
        info = {
            title   = "test"
            version = "2017-04-20T04:08:08Z"
        }
        schemes = ["https"]
        paths = {
            "/test" = {
                get = {
                    responses = {
                        "200" = {
                            description = "OK"
                        }
                    }
                    x-amazon-apigateway-integration = {
                        httpMethod = "GET"
                        type       = "HTTP"
                        responses = {
                            default = {
                                statusCode = 200
                            }
                        }
                        uri = "https://aws.amazon.com/"
                    }
                }
            }
        }
    })
}
