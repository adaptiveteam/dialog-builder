#******************************************************************************
# This is for all dialog in the bot.                                          *
#******************************************************************************
resource "aws_dynamodb_table" "dialog" {
  name = "${local.resource_name}_dialog"
  read_capacity = "${var.database_performance["dialog_read_capacity"]}"
  write_capacity = "${var.database_performance["dialog_write_capacity"]}"
  hash_key = "dialog_id"

  attribute = [
    {
      name = "dialog_id"
      type = "S"
    },
    {
      name = "context"
      type = "S"
    },
    {
      name = "subject"
      type = "S"
    }
  ]

  server_side_encryption {
    enabled = false
  }

  tags {
    Customer = "${local.resource_name}"
  }

  global_secondary_index {
    name               = "context-subject-index"
    hash_key           = "context"
    range_key          = "subject"
    write_capacity     = "${var.database_performance["dialog_read_capacity"]}"
    read_capacity      = "${var.database_performance["dialog_read_capacity"]}"
    projection_type    = "ALL"
  }
}

#******************************************************************************
# This is for application aliases to dialog context in the bot.               *
#******************************************************************************
resource "aws_dynamodb_table" "context-alias" {
  name = "${local.resource_name}_dialog_alias"
  read_capacity = "${var.database_performance["dialog_read_capacity"]}"
  write_capacity = "${var.database_performance["dialog_write_capacity"]}"
  hash_key = "application_alias"

  attribute = [
    {
      name = "application_alias"
      type = "S"
    }
  ]

  server_side_encryption {
    enabled = false
  }

  tags {
    Customer = "${local.resource_name}"
  }
}