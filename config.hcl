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

job "aws" "s3" {
    config "chaos" {
        tag = "PREFIX:test"
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
        tag = "role:myrole"
        count = 1
        chaos = "terminateAll"
    }
}

