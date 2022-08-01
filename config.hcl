app = "TestingApp"
description = "this is a test" 
    
job "aws" "ec2" {
    region = "us-west-1"
    config "chaos" {
        tags = ["env:prod"]
        chaos = "STOP"
        count = 4
    }
}

job "aws" "s3" {
    region = "us-west-1"
    config "chaos" {
        tags = ["SUFFIX:test"]
        count = 2
        chaos = "TERMINATE"
    }
}

job "aws" "lambda" {
    region = "us-east-1"
    config "chaos" {
        tags = ["name:lambda_name"]
        count = 1
        chaos = "TERMINATE"
    }
}

job "kubernetes" "deployment" {
    
    config "chaos" {
        tags = ["name:dev-app`"]
        count = 1
        chaos = "TERMINATE"
    }
}