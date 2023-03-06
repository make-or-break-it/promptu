# Setting up Promptu for the first time

Promptu's code base doesn't just contain its application code, but also its infrastructure code. Therefore, to build and run Promptu in its entirety, some manual, one-off tasks need to be completed for the first time before it can be deployed somwhere. The list of steps to follow might seem overwhelming, but remember you only need to do this once at the start of deploying Promptu! _Build it and forget it!_ If you want to learn more about _why_ we're building these resources, head over to the [infrastructure.md](./infrastructure.md) doc!

Before starting, it's recommended that you run all of these steps [in VS Code and inside the devcontainer provided with the repo](https://code.visualstudio.com/docs/devcontainers/containers#_quick-start-open-an-existing-folder-in-a-container). This will ensure you have all the CLI tools needed to run the below steps. Just make sure you have the [prerequisites to run devcontainers installed](https://code.visualstudio.com/docs/devcontainers/containers#_installation).

If you've forked Promptu and want to get its end-to-end workflow running, here's what you need to do:

1. **Setup a MongoDB Atlas account (free)**
   1. Create an account in [MongoDB Atlas](https://www.mongodb.com/atlas/database)
   2. In _Access Manager_, create an API key for Terraform (this can have any description - example: `Promptu Terraform`) with the permissions: `Organization Project Creator`. Safely record the public and private keys - these keys will be used by Terraform to manage your MongoDB account, so make sure to keep it safe for now! You'll use this later.
   3. From the same page where you got the API key, log organisation ID the MongoDB organisation - this in the URL, and is the long random string in between the URL path (e.g. `https://cloud.mongodb.com/v2#/org/<long-string>/access/users`)
2. **Setup fly.io (free/low cost)**
   1. Create an account in [fly.io](https://fly.io/)
   2. Create an organisation for Promptu (this can have any name - example: `promptu`)
   3. Create an access token for Terraform Cloud (this can have any name - example: `promptu-terraform-cloud`) - this key will be used by Terraform to manage your fly.io account, so make sure to keep it safe for now! You'll use this later.
   4. Create an access token for your Github Workflow (this can have any name - example: `promptu-github-workflow`) - this key will be used by Terraform to manage your fly.io account, so make sure to keep it safe for now! You'll use this later.
   5. You will need to attach a card to your account. You will only be charged for the CPU/memory used and the up time of the applications, which you can destroy at any time with Terraform (**tip:** use a virtual credit card, like ones provided by Revolut, to make it easier to track payments and control funds for the account). This setup described by this doc should not cost you more than $6 a month at the time of writing.
3. **Setup Terraform Cloud (free)**

   1. Create an account in [Terraform Cloud](https://cloud.hashicorp.com/products/terraform)
   2. Create a new workspace for Promptu (choose the _Version control workflow_)
   3. In _Configure settings_, go to _Advanced options_ and set your _Terraform Working Directory_ to `infrastructure`
   4. Create the following environment variables in your `promptu` Terraform Cloud workspace (settings in brackets):
      - (**Sensitive**, Environment Variable) `MONGODB_ATLAS_PUBLIC_KEY` (value secured from step `1.2.`)
      - (**Sensitive**, Environment Variable) `MONGODB_ATLAS_PRIVATE_KEY` (value secured from step `1.2.`)
      - (**Sensitive**, Environment Variable) `FLY_API_TOKEN` (value secured from step `2.3.`)
      - (Terraform Variable) `promptu_mongodb_org_id` (the organisation ID retrieved in step `1.3.`)
      - (Terraform Variable) `promptu_fly_io_org_name` (the organisation name in fly.io created in step `2.2.`)
      - (Terraform Variable) `promptu_fly_io_name_suffix` (application name suffix to make your version of Promptu app globally unique - the name will become `promptu-${promptu_fly_io_name_suffix}` in fly.io)
   5. After all your config variables have been supplied, you can create an initial run to create the base resources necessary to run your deployments - you can do this through the Terraform UI by going to _Runs > Actions > Start new run_. Once you're happy with the plan, you can apply!
   6. Once your plan has been successfully applied, you need to ensure you have the right connection details for your application to be deployed successfully for the first time - in your MongoDB Atlas `promptu` project, go to _Security > Database Access_ and edit the `promptu-feeder` and `promptu-db-updater` users. Edit each users' password and autogenerate a secure password - copy the value and keep it safe!
   7. In your MongoDB Atlas `promptu` project, go to _Deployment > Database_ and click on _Connect_ on your `promptu-db`. Select _Connect your application_ and copy the host name in the connection string. So if you're connection string looks like `mongodb+srv://<username>:<password>@promptu-db.p4jpncm.mongodb.net/?retryWrites=true&w=majority`, you only need `promptu-db.p4jpncm.mongodb.net`

4. **Scale Kafka and Zookeeper memory**
   1. The default application memory in fly.io is 256MB, but our Kafka and Zookeeper needs more memory to run. To ensure it has plenty of head room, we've choosen to give them 2048MB (or 2GB) and 1024MB (or 1GB) respectively. Before you run your Github workflows to deploy your applications for the first time, you have to run the following command to scale its memory:
   ```sh
   # Select 'yes' when the prompt asks you to continue despite the app name being different
   flyctl scale memory -a promptu-kafka-eds 2048
   flyctl scale memory -a promptu-zookeeper-eds 1024
   ```
5. **Prepare Github Workflow**

   1. In your forked Github repo, go to _Settings > Security > Secrets > Actions_ and create the following repository secrets:
      - (**Sensitive**) `FLY_API_TOKEN` (value secured from step `2.4.`)
   2. Update your `fly.toml` files to include the suffix you provided in step `3.4` (if you chose `paper_mache` as your suffix, then your app name will be `promptu-paper_mache` for the `ui` component and `promptu-api-paper_mache` for the `api` component)
   3. Use the `flyctl` CLI tool to create a `PROMPTU_MONGODB_URL` secret with your MongoDB connection URL from `3.8.` with the following script (**NOTE:** make sure you're running this from a safe environment or from within a script, as your secret will be preserved within your shell history if run directly in your shell environment):

   ```sh
   flyctl secrets set -a promptu-feeder-api-eds --detach PROMPTU_MONGODB_URL="mongodb+srv://promptu-feeder:<password from section 3.6>@<hostname from section 3.8>"
   flyctl secrets set -a promptu-db-updater-eds --detach PROMPTU_MONGODB_URL="mongodb+srv://promptu-db-updater:<password from section 3.6>@<hostname from section 3.8>"
   ```

6. **Raise your first PR and merge it into `main` to build your infrastructure and deploy your apps** - now that all the scaffolding has been set up, it's time to add some bricks! Merging your first PR will deploy your application to fly.io for the first time! All subsequent PRs will not only deploy the latest changes to fly.io, but it will also update your infrastructure through Terraform Cloud. But we're not done yet - we still need to secure access to our DB!
7. **Secure your promptu-api to MongoDB Atlas**

   1. Right now, your DB can be accessed by anyone in the world! We need to restrict this so only our fly.io can communicate with it. Find out the public IP address for the applications `feeder-api` and `db-updater` by using `flyctl ssh issue` to issue an SSH key for your fly.io app (entering `/home/vscode/.ssh/promptu-feeder-api-fly-io` and `/home/vscode/.ssh/promptu-db-updater-fly-io` respectively as your path to store the keys if you're in a devcontainer - otherwise, you have to supply the absolute path of your `~/.ssh` directory, followed by any prefix name you want e.g. `promptu-api-fly-io`) and then run `flyctl ssh console` from within the `apps/api` directory
   2. Once inside, run the following commands in order to install `dig` and to find out your app's public IP address - record this IP address (**NOTE:** this step needs to be done because the IP address in the fly.io UI is private and not recognised by Mongo. So we have to get the correct IP using this step):

   ```sh
   apt update
   apt install -y dnsutils
   dig +short txt ch whoami.cloudflare @1.0.0.1
   ```

   3. Modify the IP address you retrieved from step `6.2.` so that the last octet is 0 and it has a 24 bit subnet mask. So for example, if your IP address is `1.2.3.4`, then it should look like `1.2.3.0/24`
   4. Save the values from `6.3.` into Terraform Cloud as a Terraform variable called `promptu_feeder_api_cidr_range` and `promptu_db_updater_cidr_range`, then perform a new run - this should whitelist only your app's IP address to MongoDB Atlas

8. **Accessing the Promptu app** - you're all set! You should now be able to connect to your application from through its fly.io domain name! 🚀
   1. Go to your app on the fly.io dashboard for the `ui` component: `https://fly.io/apps/promptu-<your-suffix>`, where `<your-suffix>` is from step `4.2.`
   2. Find your app's URL via _Overview > Application Information > Hostname_ - it will likely be `https://promptu-<your-suffix>.fly.dev/`

Rest assured, most of this setup only ever happens once when you create your application. From here on out, all of your changes will be automatically applied on merge based on the settings in your `.github/workflows` config.
