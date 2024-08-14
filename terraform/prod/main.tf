terraform { 
    backend "remote" { 
    
    organization = "lambdacloud" 

    workspaces { 
        name = "aws-prod-guest-agent" 
    } 
    } 
}

import {
    to = aws_s3_bucket.guest_agent_bucket_prod
    id = "lambdalabs-guest-agent"
}
resource "aws_s3_bucket" "guest_agent_bucket_prod" {
    bucket = "lambdalabs-guest-agent"
}