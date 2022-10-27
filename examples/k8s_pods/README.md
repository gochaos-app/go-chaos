
# Chaos experiment #1

This chaos experiment for aws and k8s

* Job 1: terminates 2 pods with label app:inventory-prod
* Job 2: reboots 1 ec2 instance with tag hashi-vault in region us-west-2
* Job 3: updates to 40 the replicas of the inventory deployment

