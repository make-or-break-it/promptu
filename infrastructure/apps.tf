resource "fly_app" "promptu" {
  name = var.promptu_fly_io_name_suffix == "" ? "promptu" : "promptu-${var.promptu_fly_io_name_suffix}"
  org = var.promptu_fly_io_org_name
}
