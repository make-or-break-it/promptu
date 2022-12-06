# Infrastructure

This doc outlines how we host and run Promptu, including detailing the tooling used.

The most up to date documentation we have on our infrastructure is in our [infrastructure as code directory](../infrastructure/), so please take a look there if you'd like more detail on how Promptu's infrastructure is built.

## Prerequisites

Before we can run anything for Promptu, there are some first-time manual setup that needed to be performed. Checkout the [setup guide](./first_time_setup.md) if you've forked the repo and want Promptu to run on your own infrastructure.

## Application runtime

[fly.io](fly.io) is used to host and run our applications. We chose fly.io for this project because, at the time of writing, they have more hobby friendly pricing and do not collect bills that are less than $5 a month.

## Configuration management

We use [fly.io's secrets management](https://fly.io/docs/reference/secrets/#setting-secrets) and [the fly.toml to maintain environment variables](https://fly.io/docs/reference/configuration/) for Promptu and configure it (a la [12 factor app methodology](https://12factor.net/)).

## State management and databases

We use [MongoDB Atlas](https://www.mongodb.com/atlas/database) to store Promptu feed data because it has a hobby friendly pricing tier and has a mature Terraform provider.

## Infrastructure as Code

Terraform will be used to maintain Promptu's infrastructure. This is because Terraform is cloud agnostic and has the most active community (so it's likely to find support and documentation for most resources in Terraform).

Terraform Cloud is used to host our Terraform state.