terraform {
    cloud {
        organisation = "Promptu-mir-demo"

        workspaces {
          name = "promptu"
        }
    }

    required_providers {
        fly = {
            source = "fly-apps/fly"
            version = "0.0.21"
        }

        mongodbatlas = {
            source = "mongodb/mongodbatlas"
            version = "1.8.0"
        }
    }
}