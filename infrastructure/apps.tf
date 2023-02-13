resource "fly_app" "promptu" {
  name = var.promptu_fly_io_name_suffix == "" ? "promptu" : "promptu-${var.promptu_fly_io_name_suffix}"
  org  = var.promptu_fly_io_org_name
}

resource "fly_app" "promptu-feeder-api" {
  name = var.promptu_fly_io_name_suffix == "" ? "promptu-feeder-api" : "promptu-feeder-api-${var.promptu_fly_io_name_suffix}"
  org  = var.promptu_fly_io_org_name
}

resource "fly_app" "promptu-post-api" {
  name = var.promptu_fly_io_name_suffix == "" ? "promptu-post-api" : "promptu-post-api-${var.promptu_fly_io_name_suffix}"
  org  = var.promptu_fly_io_org_name
}

resource "fly_app" "promptu-db-updater" {
  name = var.promptu_fly_io_name_suffix == "" ? "promptu-db-updater" : "promptu-db-updater-${var.promptu_fly_io_name_suffix}"
  org  = var.promptu_fly_io_org_name
}

resource "fly_app" "promptu-kafka" {
  name = var.promptu_fly_io_name_suffix == "" ? "promptu-kafka" : "promptu-kafka-${var.promptu_fly_io_name_suffix}"
  org  = var.promptu_fly_io_org_name
}

resource "fly_app" "promptu-zookeeper" {
  name = var.promptu_fly_io_name_suffix == "" ? "promptu-zookeeper" : "promptu-zookeeper-${var.promptu_fly_io_name_suffix}"
  org  = var.promptu_fly_io_org_name
}
