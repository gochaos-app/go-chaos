app = "Login"
description = "Chaos Experiment for login app" 

job "aws" "ec2" {
    region = "us-east-1"
    config {
        tag = "Name:login-app-prod" ## search for tag: Name = login-app-prod in Virginia region
        chaos = "stop"              ## Stop instances    
        count = 3                   ## We have 5 instances in total, stop only 3
    }
}

job "aws" "lambda" {
    region = "us-east-1"
    config {
        tag = "Name:resolution-scale-prod" ## search for tag: Name = resolution-scale-prod in Virginia region
        chaos = "stop"              ## Do not destroy lambdas, put concurrency to 0    
        count = 1                   ## There's only one function in region us-east-1
    }
}

script {
    executor    = "bash"
    source      = "notification.sh" 
}