app = "TestingApp"
description = "this is a test" 
    
job "aws" "ec2" {
    region = "us-east-1"
    config "chaos" {
        tag = "Name:test"
        chaos = "stop"
        count = 2
    }
}

job "aws" "s3" {
    region = "us-west-1"
    config "chaos" {
        tag = "SUFFIX:test"
        count = 2
        chaos = "TERMINATE"
    }
}


job "kubernetes" "deployment" {
    
    config "chaos" {
        tag = "name:dev-app"
        count = 1
        chaos = "TERMINATE"
    }
}