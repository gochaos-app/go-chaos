app = "TestingApp"
description = "this is a test" 



job "kubernetes" "deployment" {
    namespace = "default"
    config {
        tag = "app:nginx"
        count = 15
        chaos = "update"
    }
}

#notification "email" {
#    from = "rodriguez.esparza.ramon@gmail.com"
#    emails = ["rodriguez.esparza.ramon@gmail.com"]
#    body = "I'm executing go-chaos, everything is fine"
#}