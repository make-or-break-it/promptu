resource "fly_app" "promptu" {
  name = var.promptu_fly_io_name_suffix == "" ? "promptu" : "promptu-${var.promptu_fly_io_name_suffix}"
  org  = var.promptu_fly_io_org_name
}

resource "fly_app" "promptu-api" {
  name = var.promptu_fly_io_name_suffix == "" ? "promptu-api" : "promptu-api-${var.promptu_fly_io_name_suffix}"
  org  = var.promptu_fly_io_org_name
}

resource "fly_machine" "example" {
  app    = "test-app"
  region = "iad"
  image  = "ubuntu:latest"
  env = {
    testEnvVar = "testValue"
  }
}
