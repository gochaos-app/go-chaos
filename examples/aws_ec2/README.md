
# Chaos experiment #1

This chaos experiment consist of two jobs in aws ec2

* Job 1: stops 3 instances with tag Name:login-app-prod in region us-east-1
* Job 2: stops 1 lambda function with tag name:resolution-scale-prod in region us-east-1
* Script: notification.sh, sends a message to a slack channel.

