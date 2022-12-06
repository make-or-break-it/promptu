resource "mongodbatlas_project" "promptu" {
  name = "promptu"
  org_id = var.promptu_mongodb_org_id
  
  lifecycle {
    ignore_changes = [
      api_keys,
    ]
  }
}

resource "mongodbatlas_cluster" "promptu-db" {
  project_id    = mongodbatlas_project.promptu.id

  name          = local.promptu_mongodb_name

  provider_name               = "TENANT"
  backing_provider_name       = "AWS"
  provider_instance_size_name = "M0" # free tier DB

  provider_region_name = "EU_WEST_1"
}

resource "mongodbatlas_database_user" "promptu" {
  project_id    = mongodbatlas_project.promptu.id

  # In order to create a DB, we have to set a fictional first time password then create
  # a real password in the MongoDB Atlas UI. Updating the password in the UI will not
  # cause configuration drift. 
  password = var.promptu_mongodb_fake_init_password
  username = "promptu"

  auth_database_name = "admin"

  roles {
   role_name     = "readWrite"
   database_name = local.promptu_mongodb_name
  }

  scopes {
    name = local.promptu_mongodb_name
    type = "CLUSTER"
  }
}

resource "mongodbatlas_project_ip_access_list" "promptu-api" {
  project_id    = mongodbatlas_project.promptu.id

  ip_address = data.fly_app.prompt-api.ipaddresses[0]
  comment    = "IP address for fly.io app ${fly_app.promptu-api.id}"
}