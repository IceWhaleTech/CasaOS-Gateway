#!/bin/bash

set -e

INSTALL_ROOT=${1:-/}

APP_NAME="casaos-gateway"
APP_NAME_SHORT="gateway"

CONF_PATH=${INSTALL_ROOT}/etc/casaos
CONF_FILE=${CONF_PATH}/${APP_NAME_SHORT}.ini

echo -n "checking if ${CONF_FILE} exists..."
cat "${CONF_FILE}" > /dev/null
echo "OK"

echo -n "checking if ${APP_NAME}.service is running..."
systemctl status "${APP_NAME}.service" --no-pager > /dev/null
echo "OK"

# TODO - check if gateway port is responding 200
