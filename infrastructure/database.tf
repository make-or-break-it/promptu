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
    project_id = mongodbatlas_project.promptu.id

    name = local.promptu_mongodb_name

    provider_name               = "TENANT"
    backing_provider_name       = "AWS"
    provider_instance_size_name = "M0" # free tier DB

    provider_region_name = "EU_WEST_1"
}

resource "mongodbatlas_database_user" "promptu" {
    project_id = mongodbatlas_project.promptu.id

    username = "promptu"
    password = var.promptu_mongodb_fake_init_password

    auth_database_name = "admin"

    roles {
        role_name = "readWrite"
        database_name = local.promptu_mongodb_name
    }

    scopes {
        name = local.promptu_mongodb_name
        type = "CLUSTER"
    }

    lifecycle {
        ignore_changes = [
          password,
        ]
    }
}

resource "mongodbatlas_project_ip_access_list" "promptu-api" {
  project_id    = mongodbatlas_project.promptu.id

  cidr_block = var.promptu_api_cidr_range
  comment    = "IP address range for fly.io app ${fly_app.promptu-api.id}"
}