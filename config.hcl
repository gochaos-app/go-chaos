app = "TestingApp"
description = "this is a test" 
    
job "aws" "ec2" {
    region = "us-east-1"
    config "chaos" {
        tag = "Name:test"
        chaos = "terminate"
        count = 0
    }
}

job "aws" "ec2_autoscaling" {
    region = "us-west-2"
    config "chaos" {
        tag = "env:prod"
        chaos = "addto"
        count = 6
    }
}

job "aws" "s3" {
    config "chaos" {
        tag = "PREFIX:app"
        count = 0
        chaos = "terminate"
    }
}

job "aws" "lambda" {
    config "chaos" {
        tag = "tag:example"
        count = 0
        chaos = "terminate"
    }
}


job "kubernetes" "pod" {
    namespace = "default"
    config "chaos" {
        tag = "app:nginx"
        count = 0
        chaos = "terminateAll"
    }
}

job "kubernetes" "deployment" {
    namespace = "default"
    config "chaos" {
        tag = "app:nginx"
        count = 0
        chaos = "terminate"
    }
}

