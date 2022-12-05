locals {
  promptu_mongodb_name = "promptu-db"
}

variable "promptu_mongo_db_project_id" {
    type = string
    default = "5d3a13b979358e125fc1fd48"
    description = "The MongoDB project ID for Promptu"
}

variable "promptu_fake_mongodb_init_password" {
    type = string
    description = <<EOF
(Source from environment variable only) The initial password used to 
create the Promptu MongoDB Atlas. In order to create a DB, we have to set 
a fictional first time password then create a real password in the MongoDB 
Atlas UI. Updating the password in the UI will not cause configuration drift 
in Terraform.
EOF
}