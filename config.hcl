app = "TestingApp"
description = "this is a test" 

job "linux" "load_balancer" {
    config {
        tag = "nyc3-load-balancer-01"
        chaos = "3emoveRules"
        count = 1
    }
}

job "aws" "ec2" {
    region = "us-east-1"
    config {
        tag = "Name:test"
        chaos = "terminate"
        count = 0
    }
}

job "aws" "ec2_autoscaling" {
    region = "us-west-2"
    config {
        tag = "env:prod"
        chaos = "addto"
        count = 0
    }
}

job "awz" "s3" {
    config {
        tag = "PREFIX:app"
        count = 0
        chaos = "terminate"
    }
}

job "aws" "lambda" {
    config {
        tag = "tag:example"
        count = 1
        chaos = "minate"
    }
}


job "kubernetes" "pod" {
    namespace = "default"
    config {
        tag = "app:nginx"
        count = 0
        chaos = "terminateall"
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