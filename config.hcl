app = "TestingApp"
description = "this is a test" 

variables {
    var_file = "testarc"
}

#job "kubernetes" "daemonSet" {
#    namespace = "kube-system"
#    config {
#        tag = "k8s-app:fluentd-logging"
#        chaos = "terminate"
#        count = 0
#    }
#}

#notification "email" {
#    from = "rodriguez.esparza.ramon@gmail.com"
#    emails = ["rodriguez.esparza.ramon@gmail.com"]
#    body = "I'm executing go-chaos, everything is fine"
#}