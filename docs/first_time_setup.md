# Setting up Promptu for the first time

Promptu's code base doesn't just contain its application code, but also its infrastructure code. Therefore, to build and run Promptu in its entirety, some manual, one-off tasks need to be completed for the first time before it can be deployed somwhere. The list of steps to follow might seem overwhelming, but remember you only need to do this once at the start of deploying Promptu! _Build it and forget it!_

Before starting, it's recommended that you run all of these steps [in VS Code and inside the devcontainer provided with the repo](https://code.visualstudio.com/docs/devcontainers/containers#_quick-start-open-an-existing-folder-in-a-container). This will ensure you have all the CLI tools needed to run the below steps. Just make sure you have the [prerequisites to run devcontainers installed](https://code.visualstudio.com/docs/devcontainers/containers#_installation).

If you've forked Promptu and want to get its end-to-end workflow running, here's what you need to do:

1. **Setup a MongoDB Atlas account (free)**
    1. Create an account in [MongoDB Atlas](https://www.mongodb.com/atlas/database)
    2. In _Access Manager_, create an API key for Terraform (this can have any description - suggestion: `Promptu Terraform`) with the permissions: `Organization Member, Organization Project Creator`. Safely record the public and private keys - these keys will be used by Terraform to manage your MongoDB account, so make sure to keep it safe for now! You'll use this later.
2. **Setup fly.io (free/low cost)**
    1. Create an account in [fly.io](https://fly.io/)
    2. Create an organisation for Promptu (this can have any name - suggestion: `promptu`)
    3. Create an access token for Terraform Cloud (this can have any name - suggestion: `promptu-terraform-cloud`) - this key will be used by Terraform to manage your fly.io account, so make sure to keep it safe for now! You'll use this later.
    4. Create an access token for your Github Workflow (this can have any name - suggestion: `promptu-github-workflow`) - this key will be used by Terraform to manage your fly.io account, so make sure to keep it safe for now! You'll use this later.
    5. _(OPTIONAL)_ If more than 2 applications are needed to be deployed, you will need to attach a card to your account. You will only be charged for the CPU used and the up time of the applications, which you can destroy at any time with Terraform (**tip:** use a virtual credit card, like ones provided by Revolut, to make it easier to track payments and control funds for the account)    
3. **Setup Terraform Cloud (free)**
    1. Create an account in [Terraform Cloud](https://cloud.hashicorp.com/products/terraform)
    2. Create a new workspace for Promptu (choose the _Version control workflow_)
    3. In _Configure settings_, go to _Advanced options_ and set your _Terraform Working Directory_ to `infrastructure`
    4. Create the following environment variables in your `promptu` Terraform Cloud workspace:
        * (**Sensitive**, Environment Variable) `MONGODB_ATLAS_PUBLIC_KEY` (value secured from step `1.2.`)
        * (**Sensitive**, Environment Variable) `MONGODB_ATLAS_PRIVATE_KEY` (value secured from step `1.2.`)
        * (**Sensitive**, Environment Variable) `FLY_API_TOKEN` (value secured from step `2.3.`)
        * (**Sensitive**, Terraform Variable) `promptu_mongodb_fake_init_password` (can be any string - this is only needed to create the MongoDB Atlas cluster for the first time)
        * (Terraform Variable) `promptu_mongodb_org_id` (the organisation ID of the account created in step `1.1.`)
4. **Prepare Github Workflow**
    1. In your forked Github repo, go to _Settings > Security > Secrets > Actions_ and create the following repository secrets:
        * (**Sensitive**) `FLY_API_TOKEN` (value secured from step `2.4.`)
5. **Raise your first PR and merge it into `main` to build your infrastructure and deploy your apps** - now that all the scaffolding has been set up, it's time to dress it up with some bricks! Merging your first PR will create the infrastructure in MongoDB and fly.io, while also deploying the latest state of Promptu to fly.io. But you're not finished yet! You still need to give your backend access to the DBs!
6. **Connect promptu-api to MongoDB Atlas**
    1. Find out your `promptu-api` public IP address by using `flyctl ssh issue` to issue an SSH key for your fly.io aap (entering `~/.ssh/promptu-api-fly-io` as your path to store the keys) and then `flyctl ssh console` from with the `apps/api` directory
    2. Once inside, run the following commands in order to install `dig` and to find out your app's public IP address - record this IP address
    ```sh
    apt update
    apt install dnsutils
    dig +short txt ch whoami.cloudflare @1.0.0.1
    ```
    3. Modify the IP address you retrieved from step `6.2.` so that the last octet is 0 and it has a 24 bit subnet mask. So for example, if your IP address is `1.2.3.4`, then it should look like `1.2.3.0/24`
    4. Save the value from `6.3.` into Terraform Cloud as a Terraform variable called `prompt_api_cidr_range`, then perform a new run - this should whitelist your app's IP address to MongoDB Atlas 
    5. In your MongoDB Atlas `promptu` project, go to _Security > Database Access_ and edit the `promptu` user. Edit the `promptu` user's password and autogenerate a secure password - copy the value and keep it safe!
    6. In your MongoDB Atlas `promptu` project, go to _Deployment > Database_ and click on _Connect_ on your `promptu-db`. Select _Connect your application_ and copy the connection string up until (but not including!) the last `/`. So if you're connection string looks like `mongodb+srv://promptu:<password>@promptu-db.p4jpncm.mongodb.net/?retryWrites=true&w=majority`, you only need `mongodb+srv://promptu:<password>@promptu-db.p4jpncm.mongodb.net`.
    7. Replace `<password>` in the connection string from `6.6.` with the value secured from `6.5.` - keep this value secure!
    8. Use the `flyctl` CLI tool to create a `PROMPTU_MONGODB_URL` secret with your MongoDB connection URL from `6.7.` with the following script (**note:** make sure you're running this from a safe environment or from within a script, as your secret will be preserved within your shell history if run directly in your shell environment): 
    ```sh
    flyctl secrets set --detach PROMPTU_MONGODB_URL="<secure value from 6.7.>"
    ```

And you're all set! You should now be able to connect to your application from through its fly.io domain name! 🚀
