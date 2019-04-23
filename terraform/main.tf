#******************************************************************************
# WARNING!!! DO NOT MODIFY THIS FILE DIRECTLY UNLESS YOU KNOW EXACTLY WHAT    *
# YOU ARE DOING!!! Instead, modify the terraform.tfvars file. That file       *
# contain the settings you need to tune the infrastructure.                   *
#******************************************************************************

#******************************************************************************
#                                 Account set up                              *
#******************************************************************************

data "aws_caller_identity" "current" {}

provider "aws" {
  version = ">= 2.5.0"
  region = "${var.project_admin["region"]}"
  profile = "${var.project_admin["profile"]}"
}

#******************************************************************************
#                   Setting up local variables & dependencies                 *
#******************************************************************************
locals {
  resource_name = "${lower(replace(var.project_admin["infrastructure_id"], "-","_"))}"
  api_token = "8b6a8d3e79398ac217de9031fadb0a1eb6b8734b"
}