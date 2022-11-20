resource "fly_app" "promptu-db" {
  name = "promptu-db"
  org = "promptu"
}

resource "fly_volume" "promptu-db-volume" {
  name   = "promptu-db-volume"
  app    = "promptu-db"
  size   = 10
  region = "cdg"
}