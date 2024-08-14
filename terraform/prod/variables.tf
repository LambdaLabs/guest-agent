variable "prefix" {
  type        = string
  default     = "lambdalabs"
  description = "A prefix for all resource names for namespacing purposes"
}

variable "github_org" {
  type        = string
  default     = "lambdal"
  description = "Github organization"
}

variable "github_repo" {
  type        = string
  default     = "guest-agent"
  description = "Name of the Github repo for the above orga that is allowed to push to out private ECR"
}
