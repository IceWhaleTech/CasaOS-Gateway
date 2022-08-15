#!/bin/bash

set -e

## base variables
APP_NAME="casaos-gateway"
APP_NAME_SHORT="gateway"

# copy config files
CONF_PATH=/etc/casaos
CONF_FILE=${CONF_PATH}/${APP_NAME_SHORT}.ini
CONF_FILE_SAMPLE=${CONF_PATH}/${APP_NAME_SHORT}.ini.sample

if [ ! -f "${CONF_FILE}" ]; then \
    echo "Initializing config file..."
    cp -v "${CONF_FILE_SAMPLE}" "${CONF_FILE}"; \
fi

# enable and start service
systemctl daemon-reload

echo "Enabling service..."
systemctl enable --force --no-ask-password "${APP_NAME}.service"

echo "Starting service..."
systemctl start --force --no-ask-password "${APP_NAME}.service"
