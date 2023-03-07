# Promptu

Promptu Make IT is an anonymous social media that allows people to post answers to daily questions once at a given time.

## Repo Structure

```sh
.
├── .github 
│   └── workflows   # Contains Github Worklows
├── .devcontainer   # Devcontainer config - single devcontainer for entire repo
├── apps            # Contains apps deployed to fly.io - should only contain directories
│   └── <component> # Component directory - should contain at least fly.toml and Dockerfile
├── docs            # Central docs for entire repo
├── infrastructure  # Terraform config for entire repo
└── scripts         # Helpful utilities and scripts for maintaining projects and repo
```

## Contributing

### To an app in Promptu

This repo is structured as a monorepo: each directory in the `apps` directory of this repo should map to a component of Promptu and should be as self contained as possible. Each directory should be deployed as its own service.

As such, each component will have its own contributing guide, so to contribute, please first decide which component you need to contribute to, and navigate to that component's directory.

### To the infrastructure

Promptu uses Terraform and Terraform Cloud to manage its state. All infrastructure resources are maintained in the `./infrastructure` directory. Infrastructure changes are only locked to the admins of the Terraform Cloud workspace. If you're an admin, raise a PR with your infrastructure changes and a plan will be automatically be produced via the Terraform/Github integration.

## Workflow (for contributors/deep dive speakers)

Because this app was built for educational purposes, we've created a `main` branch for each deep dive that relies on the app as a demo. This is to facilitate independent experimentation for each talk.

* `main` - any changes that should be applied to all talks will be applied to this branch. All deep dive branches should regularly rebase off of this branch to stay up to date. At the end of the Make or Break It event, all changes from all deep dives will be coalesced to this branch.
* `eds/main` - `main` branch belonging to the EDS deep dive (Event Driven Systems)
* `mir/main` - `main` branch belonging to the MIR deep dive (Managing Infrastructure Resources)