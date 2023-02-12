resource "fly_app" "promptu" {
  name = var.promptu_fly_io_name_suffix == "" ? "promptu" : "promptu-${var.promptu_fly_io_name_suffix}"
  org  = var.promptu_fly_io_org_name
}

resource "fly_app" "promptu-feeder-api" {
  name = "promptu-feeder-api"
  org  = var.promptu_fly_io_org_name
}

resource "fly_app" "promptu-post-api" {
  name = "promptu-post-api"
  org  = var.promptu_fly_io_org_name
}

resource "fly_app" "promptu-db-updater" {
  name = "promptu-db-updater"
  org  = var.promptu_fly_io_org_name
}
