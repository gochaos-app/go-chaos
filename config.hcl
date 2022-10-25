app = "TestingApp"
description = "this is a test" 

job "aws" "ec2_autoscaling" {
    region = "us-west-2"
    config {
        tag = "env:prod"
        chaos = "update"
        count = 0
    }
}

#notification "email" {
#    from = "rodriguez.esparza.ramon@gmail.com"
#    emails = ["rodriguez.esparza.ramon@gmail.com"]
#    body = "I'm executing go-chaos, everything is fine"
#}