app = "TestingApp"
description = "this is a test" 
    
job "aws" "ec2" {
    region = "us-east-1"
    config "chaos" {
        tag = "Name:test"
        chaos = "terminate"
        count = 1
    }
}

job "aws" "lambda" {
    config "chaos" {
        tag = "tag:example"
        count = 1
        chaos = "stop"
    }
}


job "kubernetes" "deployment" {
    
    config "chaos" {
        tag = "name:dev-app"
        count = 1
        chaos = "TERMINATE"
    }
}