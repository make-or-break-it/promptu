resource "fly_app" "promptu" {
  name = "promptu"
  org = "promptu"
}

resource "fly_app" "promptu-api" {
  name = "promptu-api"
  org = "promptu"
}

data "fly_app" "prompt-api" {
  name = fly_app.promptu-api.id
}