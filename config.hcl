app = "TestingApp"
description = "this is a test" 


job "aws" "ec2" {
    region = "us-east-1"
    config {
        tag = "Name:test"
        chaos = "terminate"
        count = 1
    }
}

job "aws" "ec2_autoscaling" {
    region = "us-west-2"
    config {
        tag = "env:prod"
        chaos = "addto"
        count = 6
    }
}

job "aws" "s3" {
    config {
        tag = "PREFIX:app"
        count = 0
        chaos = "terminate"
    }
}

job "aws" "lambda" {
    config {
        tag = "tag:example"
        count = 0
        chaos = "terminate"
    }
}


job "kubernetes" "pod" {
    namespace = "default"
    config {
        tag = "app:nginx"
        count = 0
        chaos = "terminateAll"
    }
}

job "kubernetes" "deployment" {
    namespace = "default"
    config {
        tag = "app:nginx"
        count = 0
        chaos = "terminate"
    }
}

