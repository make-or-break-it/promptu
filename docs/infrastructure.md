# Infrastructure

This doc outlines how we host and run Promptu, including detailing the tooling used.

## Prerequisites

Before we can run anything for Promptu, the following had to be manually made:

* A [new Terraform workspace](https://app.terraform.io/app/sanyia/workspaces/new) in Terraform Cloud (used to host our Terraform state)
* Create a [fly.io organisation](https://fly.io/dashboard/promptu/billing)
* Create an organisation for MongoDB Atlas (and public/private keys) and add to Terraform as environment variabls
* fly.io access tokens for:
    * Github Workflow
    * Terraform Cloud 
* [Manually generated fly.toml](https://fly.io/docs/reference/configuration/) for each application - these will be used to inform Fly how to deploy our applications (this would normally be automatically generated with `fly launch`, but we're using Terraform to create our app instead, so we have to default to manual TOML file generation)
* The infrastructure needs to be created first before the Github workflow/deployments can run
* Run the Terraform apply
* [Add secrets](https://github.com/fly-apps/postgres-ha#set-secrets) to the pg app

## Application runtime

## Databases

TBD - setup instructions [here](https://github.com/fly-apps/postgres-ha)

## Infrastructure as Code

Terraform will be used to maintain Promptu's infrastructure. This is because Terraform is cloud agnostic and has the most active community (so it's likely to find support and documentation for most resources in Terraform).

Terraform Cloud is used to host our Terraform state.