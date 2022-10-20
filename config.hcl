app = "TestingApp"
description = "this is a test" 



job "kubernetes" "node" {
    namespace = "default"
    config {
        tag = "app:nginx"
        count = 0
        chaos = "terminateall"
    }
}

