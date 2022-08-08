app = "TestingApp"
description = "this is a test" 
    
job "aws" "ec2" {
    region = "us-east-1"
    config "chaos" {
        tag = "Name:test"
        chaos = "terminate"
        count = 2
    }
}

job "aws" "s3" {
    config "chaos" {
        tag = "test:blog"
        count = 20
        chaos = "terminate"
    }
}


job "kubernetes" "deployment" {
    
    config "chaos" {
        tag = "name:dev-app"
        count = 1
        chaos = "TERMINATE"
    }
}