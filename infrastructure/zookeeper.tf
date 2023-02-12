resource "fly_volume" "zookeeper" {
  name       = "zookeeper_data"
  app        = "zookeeper"
  size       = 1
  region     = "lhr"
  depends_on = [fly_app.zookeeper]
}

resource "fly_ip" "zookeeper-ip" {
  app  = "zookeeper"
  type = "v6"
}
