app = "TestingApp"
description = "this is a test" 

job "gcp" "vm" {
    project = "go-chaos"
    region = "us-central1-a"
    config {
        tag = "gcp:hellovm"
        chaos = "terminate"
        count = 1
    }
}

job "aws" "ec2" {
    region = "us-east-1"
    config {
        tag = "aws:hello" 
        chaos = "terminate"    
        count = 1  
    }
}

job "do" "droplet" {
    config {
        tag = "hello" 
        chaos = "terminate"    
        count = 1  
    }
}

job "kubernetes" "deployment" {
    namespace = "default"
    config {
        tag = "app:nginx" 
        chaos = "terminate"    
        count = 1   
    }
}

#notification "gmail" {
#    from = "chaos-email@gmail.com" #email notification only works with gmail, set up GMAIL_APP_TOKEN
#    emails = ["customers-dev@gmail.com", "customer-qa@gmail.com"] # distribution lists to dev and qa team
#    body = "chaos experiment running, to check dev teams and get latency in qa"
#}