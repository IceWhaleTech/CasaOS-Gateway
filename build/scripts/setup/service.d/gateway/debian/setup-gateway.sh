#!/bin/bash

set -e

## base variables
BUILD_PATH=${1:?missing build path}
SOURCE_ROOT=${BUILD_PATH}/sysroot

APP_NAME="casaos-gateway"
APP_NAME_SHORT="gateway"

LEGACY_APP_NAME="casaos"

# functions
__is_version_gt() {
    test "$(echo "$@" | tr " " "\n" | sort -V | head -n 1)" != "$1"
}

__is_migration_needed() {
    local version1
    local version2
    
    version1="${1}"
    version2="${2}"

    if [ "TARGET_VERSION_NOT_FOUND" = "${version2}" ]; then
        false
    fi

    if [ "LEGACY_WITHOUT_VERSION" = "${version2}" ]; then
        true
    fi

    __is_version_gt "${version1}" "${version2}"
}

# main
SOURCE_BIN_PATH=${SOURCE_ROOT}/usr/bin
SOURCE_BIN_FILE=${SOURCE_BIN_PATH}/${APP_NAME}

TARGET_BIN_PATH=/usr/bin
TARGET_BIN_FILE=${TARGET_BIN_PATH}/${APP_NAME}
TARGET_BIN_FILE_LEGACY=${TARGET_BIN_PATH}/${LEGACY_APP_NAME}

CONF_PATH=/etc/casaos
CONF_FILE=${CONF_PATH}/${APP_NAME_SHORT}.ini
CONF_FILE_SAMPLE=${CONF_PATH}/${APP_NAME_SHORT}.ini.sample

if [ ! -f "${CONF_FILE}" ]; then \
    echo "Initializing config file..."
    cp -v "${CONF_FILE_SAMPLE}" "${CONF_FILE}"; \
fi

SOURCE_VERSION="$($SOURCE_BIN_FILE -v)"
TARGET_VERSION="$($TARGET_BIN_FILE -v || $TARGET_BIN_FILE_LEGACY -v || $LEGACY_APP_NAME -v || (which $LEGACY_APP_NAME > /dev/null && echo LEGACY_WITHOUT_VERSION) || echo TARGET_VERSION_NOT_FOUND)"

if __is_migration_needed "${SOURCE_VERSION}" "${TARGET_VERSION}"; then \
    echo "Old version of CasaOS found. Migrating..."


fi

echo "Enabling service..."
systemctl enable --force --no-ask-password "${APP_NAME}.service"

echo "Starting service..."
systemctl start --force --no-ask-password "${APP_NAME}.service"
