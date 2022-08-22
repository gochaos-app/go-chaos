# chaosctl

chaos engineering as code. chaosctl is an app that lets you inject chaos in your infrastructure
with declarative code and run chaos experiments in your cloud infrastructure. Although chaosctl works mainly with the ***server-shutdown perturbation model***, it does have some extra capabilities such as data deletion (objects and data) and increase or decrease of resources.

***If you want to perform chaos experiments as code in your infrastructure, but do not have monitoring, IaC, automated CI/CD or an easy way to recover, then sorry, CHAOS ENGINEERING IS NOT FOR YOU***

chaosctl read chaos experiments in HCL format.

Create a directory in you computer

```
mkdir chaosExperiment && cd chaosExperiment
```

Create a new file with the following content and name it `config.hcl`

```
app = "TestingApp"
description = "this is a test" 
    
job "aws" "ec2" {
    region = "us-east-1"
    config "chaos" {
        tag = "Name:test"
        chaos = "terminate"
        count = 5
    }
}
```
Once the file saved, execute the file with command `chaosctl d config.hcl`

Several jobs are possible with chaosctl, jobs will execute from top to bottom. 
```
job "aws" "ec2" {
    region = "us-east-1"
    config "chaos" {
        tag = "Name:first"
        chaos = "terminate"
        count = 5
    }
}

job "aws" "s3" {
    region = "us-east-1"
    config "chaos" {
        tag = "PREFIX:second"
        chaos = "terminate"
        count = 30
    }
}

```
# Limit blast radius

chaosctl limits itself with the use of config options on each job: 
* region:  will not destroy or delete anything on other cloud regions.
* namespace: limits the blast radius to a single namespace (for K8s only).
* tag:     single tag, that will find the specified resources and kill those. 
* count:   Option to limit the number of resources that chaosctl will destroy. 

# What it can do? 
chaosctl has a number of predifined chaos actions on several resources (AWS and K8s)

## AWS

* Compute Autoscaling
```
terminate
update
addto
```

* Compute EC2:
```
terminate
reboot
stop
```

* Compute Lambda
```
terminate
stop
```

* Storage S3:
```
terminate
delete_content
```

## K8s

* Pods:
```
terminate
terminateAll
```

* Deployments:
```
terminate
```

chaosctl can also execute a single script at the beginning of the file
```
script {
    executor   = "bash"
    source     = "destroy.sh"
}
```



## Roadmap and news

chaosctl has been in development for quite sometime, there are some important modifications over the first versions: 
* uses HCL instead of JSON, YAML or TOML. 
* Go 1.18.5, instead of 1.14 or 1.16.
* aws-sdk-go-v2 instead of v1. 
* Entirely new code base. 
* Same interface to all functions.
* Remove conditions in favor of a simple hashmap for events. 

### Rules: 
* For any new functionality or bug fix create a PR to main branch.
* Only tags are made for production.
* Working and compiling code over perfect code. (we are not making an operating system or flying a plane here) 
* No random function (is difficult to manage randomness for small sets).
* No Windows version (it will require some serious rewriting).
