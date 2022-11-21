resource "mongodbatlas_cluster" "promptu-db" {
  project_id    = var.promptu_mongo_db_project_id

  name          = "promptu-db"

  provider_name               = "TENANT"
  backing_provider_name       = "AWS"
  provider_instance_size_name = "M0" # free tier DB

  provider_region_name = "EU_WEST_1"
}

resource "mongodbatlas_database_user" "prompt" {
  project_id    = var.promptu_mongo_db_project_id

  username = "promptu"
  password = "thispasswordisnotreal" # can be changed without affecting the resource

  roles {
   role_name = "readWriteAnyDatabase"
   database_name = "admin"
  }
}