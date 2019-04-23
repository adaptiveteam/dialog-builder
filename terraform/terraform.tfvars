#******************************************************************************
# These are tuning variables for DynamoDB database performance and caching.   *
#******************************************************************************
database_performance = {
  dialog_read_capacity       = 5
  dialog_write_capacity      = 5
}

project_admin = {
  #########################################################################
  # Change this to the name of the customer or development environment.   #
  # This key is used for names for many resources in AWS. Consequently,   #
  # hyphens will be replaced where possible with underscore and strings   #
  # will be converted to lower case.                                      #
  # The name must between 1 and 20 characters                             #
  #########################################################################
  infrastructure_id = "ctcreel-adaptive"
}

#**************************************************************************
# Deployment specific information.                                        *
#**************************************************************************
deployment_admin = {
  archive_path = "src/main/resources/data/deployment"
}

dialog_ = {
  archive_path = "src/main/resources/data/deployment"
}
