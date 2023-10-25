app = "TestingApp"
description = "this is a test" 


job "aws" "ec2" {
    region = "us-west-2"
    config {
        tag = "Name:hola" 
        chaos = "terminate"    
        count = 4
    }
}

job "script" "python3:script.py" {
    region = "us-west-2"
    namespace = "default"
    project = "this is a test"
    config {
        tag = "hello" 
        chaos = "terminate"    
        count = 3
    }
}
