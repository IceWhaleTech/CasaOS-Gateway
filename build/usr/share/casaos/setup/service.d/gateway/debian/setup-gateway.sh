#!/bin/bash

set -e

INSTALL_ROOT=${1:-/}

APP_NAME="casaos-gateway"
APP_NAME_SHORT="gateway"

BIN_PATH=${INSTALL_ROOT}/usr/bin
BIN_FILE=${BIN_PATH}/${APP_NAME}

CONF_PATH=${INSTALL_ROOT}/etc/casaos
CONF_FILE=${CONF_PATH}/${APP_NAME_SHORT}.ini
CONF_FILE_SAMPLE=${CONF_PATH}/${APP_NAME_SHORT}.ini.sample

VERSION="$($BIN_FILE -v)"
VERSION_PATH=${INSTALL_ROOT}/var/lib/casaos
VERSION_FILE=${VERSION_PATH}/${APP_NAME_SHORT}-version

echo "Writing version number '${VERSION}' to ${VERSION_FILE}..."
echo "${VERSION}" > "${VERSION_FILE}"

if [ ! -f "${CONF_FILE}" ]; then \
    echo "Initializing config file..."
    cp -v "${CONF_FILE_SAMPLE}" "${CONF_FILE}"; \
fi

echo "Enabling service..."
systemctl enable --force --no-ask-password "${APP_NAME}.service"

echo "Starting service..."
systemctl start --force --no-ask-password "${APP_NAME}.service"