terraform {
  cloud {
    organization = "make-or-break-it"

    workspaces {
      name = "promptu-mir"
    }
  }

  required_providers {
    fly = {
      source  = "fly-apps/fly"
      version = "~> 0.0.20"
    }

    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.6"
    }
  }
}
