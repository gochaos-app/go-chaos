app = "TestingApp"
description = "this is a test" 


job "kubernetes" "pod" {
    namespace = "nginx-dev"
    config {
        tag = "app:nginx" 
        chaos = "terminate"    
        count = 1
    }
}

#job "script" "python3:script.py" {
#    region = "us-west-2"
#    namespace = "default"
#    project = "this is a test"
#    config {
#        tag = "hello" 
#        chaos = "terminate"    
#        count = 3
#    }
#}
