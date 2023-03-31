app = "TestingApp"
description = "this is a test" 
#function = "random"
#job "module" "vm" {
#    project = "go-chaos"
#    region = "us-central1-a"
#    config {
#        tag = "modules tags"
#        chaos = "terminate"
#        count = 1
#    }

#}

notification "slack" {
    //from = "chaos-email@gmail.com" #email notification only works with gmail, set up GMAIL_APP_TOKEN
    to = ["C04N5HS91MY","C045XNKCY1H","C04N5HS91MY"] # distribution lists to dev and qa team
    body = "chaos experiment running, to check dev teams and get latency in qa"
}