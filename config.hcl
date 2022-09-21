app = "TestingApp"
description = "this is a test" 

job "do" "droplet" {
    config {
        tag = "app"
        chaos = "terminate"
        count = 3
    }
}
