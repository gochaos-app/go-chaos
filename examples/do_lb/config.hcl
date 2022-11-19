app = "Customers app"
description = "Chaos Experiment for the customers application " 

job "do" "droplet" {
    config {
        tag = "customers-plugin-qa"   ## search for tag: customer-plugin-qa
        chaos = "reboot"              ## reboot droplets    
        count = 5                     ## We have 10 droplets in total, half will be rebooted
    }
}

job "do" "load_balancer" {
    config {
        tag = "customers-plugin-dev"  ## search for tag: customers-plugin-dev
        chaos = "removeDroplets"      ## Remove droplets from dev load balancer
        count = 2                     ## 2 droplets our of 3 will be removed  
    }
}

notifications "gmail" {
    from = "chaos-email@gmail.com" #email notification only works with gmail, set up GMAIL_APP_TOKEN
    emails = ["customers-dev@gmail.com", "customer-qa@gmail.com"] # distribution lists to dev and qa team
    body = "chaos experiment running, to check dev teams and get latency in qa"
}