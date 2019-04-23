#**************************************************************************
# Basic project administration settings.                                  *
#**************************************************************************
variable "database_performance" {
  description = "Knobs to use for dynamoDB performance"
  type = "map"
}

variable "deployment_admin" {
  description = "Administrative information about the deployment"
  type = "map"
}

variable "project_admin" {
  description = "Administration information about the project."
  type = "map",
  default = {
    region              = "us-east-1"
    # You must change this to the region in which you are working.
    profile             = "adaptive-dev"
    # You must change this to the name of the AWS profile you wish to use for this deployment
  }
}