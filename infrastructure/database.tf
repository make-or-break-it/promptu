resource "fly_app" "promptu_db" {
  name = "promptu_db"
  org = "promptu"
}

resource "fly_volume" "promptu_db_volume" {
  name   = "promptu_db_volume"
  app    = fly_app.promptu_db.name
  size   = 10
  region = "cdg"
}