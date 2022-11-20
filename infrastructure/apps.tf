resource "fly_app" "promptu" {
  name = "promptu"
  org = "promptu"
}

resource "fly_app" "promptu-api" {
  name = "promptu-api"
  org = "promptu"
}