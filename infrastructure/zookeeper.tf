resource "fly_volume" "zookeeper" {
  name   = "zookeeper_data"
  app    = fly_app.promptu-zookeeper.name
  size   = 1
  region = "lhr"
}

resource "fly_ip" "zookeeper-ip" {
  app  = fly_app.promptu-zookeeper.name
  type = "v6"
}
