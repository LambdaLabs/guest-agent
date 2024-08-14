terraform { 
    backend "remote" { 
    
    organization = "lambdacloud" 

    workspaces { 
        name = "aws-prod-guest-agent" 
    } 
    } 
}

import {
    to = aws_s3_bucket.guest-agent-bucket-prod
    id = "lambdalabs-guest-agent"
}
resource "aws_s3_bucket" "guest-agent-bucket-prod" {
    bucket = "lambdalabs-guest-agent"
}