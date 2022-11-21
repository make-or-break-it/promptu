terraform {
  cloud {
    organization = "sanyia"

    workspaces {
      name = ["promptu"]
    }
  }

  required_providers {
    fly = {
      source = "fly-apps/fly"
      version = "~> 0.0.20"
    }

    mongodbatlas = {
      source = "mongodb/mongodbatlas"
      version = "~> 1.6"
    }
  }
}
