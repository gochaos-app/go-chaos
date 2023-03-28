app = "TestingApp"
description = "this is a test" 
#function = "random"
job "module" "vm" {
    project = "go-chaos"
    region = "us-central1-a"
    config {
        tag = "gcp:hellovm"
        chaos = "terminate"
        count = 1
    }

}
