resource "fly_app" "promptu" {
  name = "promptu"
  org = "promptu"
}

resource "fly_app" "promptu_api" {
  name = "promptu_api"
  org = "promptu"
}