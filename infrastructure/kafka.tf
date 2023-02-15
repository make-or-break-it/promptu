resource "fly_volume" "kafka" {
  name   = "kafka_data"
  app    = fly_app.promptu-kafka.name
  size   = 1
  region = "lhr"
}

resource "fly_ip" "kafka-ip" {
  app  = fly_app.promptu-kafka.name
  type = "v4"
}
