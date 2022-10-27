app = "inventory app"
description = "Chaos Experiment for inventory app, check how many customers are able to get products in cart" 

job "kubernetes" "pod" {
    namespace = "inventory-app"
    config {
        tag = "app:inventory-prod"  ## search for label app = inventory-prod in inventory-app namespace
        chaos = "terminate"         ## terminate pods
        count = 2                   ## We have 4 pods for this service
    }
}

# reboot instance for vault where secrets and parameters are stored, while these pods start again

job "aws" "ec2" {
    region = "us-west-2"
    config {
        tag = "hashi-vault"
        chaos = "reboot"
        count = 1
    }
}

# elevate the deployment replicas up to 40 with no secrets store
# what does the cluster does? 
# the rest of the pods are fine?
# do the pods initialize with no parameters from vault? if so, why? 

job "kubernetes" "deployment" {
    namespace = "inventory-app"
    config {
        tag = "app:inventory-prod" ## search for label app = inventory-prod in inventory-app namespace
        chaos = "update"             ## Do not destroy lambdas, put concurrency to 0    
        count = 40                   ## Update number of pods up to 40
    }
}




