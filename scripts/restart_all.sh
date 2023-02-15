#/usr/bin/env bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
APP_SUFFIX="${1:?App suffix required -- e.g. ./restart_all.sh <app_suffix>}"

cd "${SCRIPT_DIR}/../apps/zookeeper"
flyctl apps restart "promptu-zookeeper-${APP_SUFFIX}"
sleep 10s

cd "${SCRIPT_DIR}/../apps/kafka"
flyctl apps restart "promptu-kafka-${APP_SUFFIX}"
sleep 10s

cd "${SCRIPT_DIR}/../apps/post-api"
flyctl apps restart "promptu-post-api-${APP_SUFFIX}"

cd "${SCRIPT_DIR}/../apps/feeder-api"
flyctl apps restart "promptu-feeder-api-${APP_SUFFIX}"

cd "${SCRIPT_DIR}/../apps/db-updater"
flyctl apps restart "promptu-db-updater-${APP_SUFFIX}"