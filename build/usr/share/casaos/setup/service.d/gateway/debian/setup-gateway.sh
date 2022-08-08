#!/bin/bash

set -e

INSTALL_ROOT=${1:-/}

APP_NAME="casaos-gateway"
APP_NAME_SHORT="gateway"

CONF_PATH=${INSTALL_ROOT}/etc/casaos
CONF_FILE=${CONF_PATH}/${APP_NAME_SHORT}.ini
CONF_FILE_SAMPLE=${CONF_PATH}/${APP_NAME_SHORT}.ini.sample

if [ ! -f "${CONF_FILE}" ]; then \
    echo "Initializing config file..."
    cp -v "${CONF_FILE_SAMPLE}" "${CONF_FILE}"; \
fi

echo "Enabling service..."
systemctl enable --now --force --no-ask-password "${APP_NAME}.service"
