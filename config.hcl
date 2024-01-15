app = "TestingApp"
description = "this is a test" 


job "script" "python3:script.py" {
    region = "us-west-2"
    namespace = "default"
    config {
        tag = "hello" 
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
