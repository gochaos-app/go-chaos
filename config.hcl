app = "TestingApp"
description = "this is a test" 

job "do" "droplet" {
    region = "us-east-1"
    config {
        tag = "Name:superAwesomeApp"
        chaos = "terminate"
        count = 3
    }
}
