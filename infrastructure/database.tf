resource "fly_app" "promptu-db" {
  name = "promptu-db"
  org = "promptu"
}

resource "fly_volume" "promptu-db-volume" {
  name   = "promptu_db_volume" # Volumes do not allow dashes, unlike apps
  app    = fly_app.promptu-db.name
  size   = 10
  region = "cdg"
}