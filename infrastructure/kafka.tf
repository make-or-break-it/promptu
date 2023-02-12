resource "fly_machine" "kafka" {
  app      = fly_app.kafka.name
  region   = "lhr"
  name     = "kafka"
  image    = "wurstmeister/kafka:latest"
  memorymb = 2046

  env = {
    KAFKA_ZOOKEEPER_CONNECT              = "zookeeper.fly.dev:2181"
    KAFKA_ADVERTISED_LISTENERS           = "INSIDE://:9092,OUTSIDE://${fly_app.kafka.appurl}:9094"
    KAFKA_LISTENERS                      = "INSIDE://:9092,OUTSIDE://:9094"
    KAFKA_LISTENER_SECURITY_PROTOCOL_MAP = "INSIDE:PLAINTEXT,OUTSIDE:PLAINTEXT"
    KAFKA_INTER_BROKER_LISTENER_NAME     = "INSIDE"
    KAFKA_HEAP_OPTS                      = "-Xmx1024M -Xms1024M"
    KAFKA_CREATE_TOPICS                  = "posts:10:1:compact"
  }

  services = [
    {
      ports = [
        {
          port = 9094
        }
      ]
      "protocol" : "tcp",
      "internal_port" : 9092
    }
  ]

  mounts = [{
    path   = "/data",
    volume = "${fly_volume.kafka.name}"
  }]
}

resource "fly_volume" "kafka" {
  name       = "kafka_data"
  app        = "kafka"
  size       = 1
  region     = "lhr"
  depends_on = [fly_app.kafka]
}

resource "fly_ip" "kafka-ip" {
  app  = "kafka"
  type = "v6"
}
