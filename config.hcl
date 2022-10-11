app = "TestingApp"
description = "this is a test" 

job "aws" "ec2" {
    region = "us-west-2"
    config {
        tag = "Name:My-awesome-app"
        chaos = "terminate"
        count = 2
    }
}

job "aws" "s3" {
    config {
        tag = "PREFIX:awesome-app-prod"
        chaos = "terminate"
        count = 2
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

job "aws" "s3" {
    config {
        tag = "PREFIX:awesome-app-prod"
        chaos = "terminate"
        count = 2
    }
}

notification "email" {
    from = "rodriguez.esparza.ramon@gmail.com"
    emails = ["ramon.esparza@wizeline.com", 
            "feribelles@gmail.com"
            ]
    body  = "this is another test perrosinpatas" 
}
