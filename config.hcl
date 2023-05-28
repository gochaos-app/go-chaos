app = "TestingApp"
description = "this is a test" 
function = "random"


hypothesis {
    name = "this is a test"
    description = "My hypothesis is that the latency will not go down"
    workers = 10
    url = "http://localhost:8080"
    report = "testing3.txt"

}

job "do" "droplet" {
    config {
        tag = "hello" 
        chaos = "terminate"    
        count = 3
    }
}
