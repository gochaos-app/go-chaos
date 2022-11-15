app = "TestingApp"
description = "this is a test" 


#job "gcp" "vm" {
#    project = "go-chaos"
#    region = "us-central1-a"
#    config {
#        tag = "test:value"
#        chaos = "terminate"
#        count = 0
#    }
#}

job "aws" "ec2_autoscaling" {
    region = "us-east-1"
    config {
        tag = "test:value" 
        chaos = "terminate"    
        count = 1   
    }
}

