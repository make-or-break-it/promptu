resource "fly_app" "kafka" {
  name = "kafka"
  org = var.promptu_fly_io_org_name
}

// TODO we need to raise the memory for the kafka machine to 2GB - ideally from terraform

resource "fly_app" "zookeeper" {
  name = "zookeeper"
  org = var.promptu_fly_io_org_name
}

resource "fly_volume" "kafka" {
  name       = "kafka_data"
  app        = "kafka"
  size       = 1
  region     = "lhr"
  depends_on = [fly_app.kafka]
}

resource "fly_volume" "zookeeper" {
  name       = "zookeeper_data"
  app        = "zookeeper"
  size       = 1
  region     = "lhr"
  depends_on = [fly_app.zookeeper]
}

resource "fly_ip" "kafkaIP" {
  app  = "kafka"
  type = "v6"
}

resource "fly_ip" "zookeeperIP" {
  app  = "zookeeper"
  type = "v6"
}