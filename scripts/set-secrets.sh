#/usr/bin/env bash

PROMPTU_DB_HOST='promptu-db.az3ptb6.mongodb.net'

# Being lazy here and statically defining the path of the API code.
# This means we always have to run the script from where we found it.
# We can update this code later to make it runnable from anywhere.
cd ../apps/api

flyctl secrets set --detach PROMPTU_MONGODB_URL="mongodb+srv://${PROMPTU_DB_USERNAME}:${PROMPTU_DB_PASSWORD}@${PROMPTU_DB_HOST}"