#!/bin/bash

set -e

# functions
__is_version_gt() {
    test "$(echo "$@" | tr " " "\n" | sort -V | head -n 1)" != "$1"
}

__is_migration_needed() {
    local version1
    local version2

    version1="${1}"
    version2="${2}"

    if [ "${version1}" = "${version2}" ]; then
        return 1
    fi

    if [ "CURRENT_VERSION_NOT_FOUND" = "${version1}" ]; then
        return 1
    fi

    if [ "LEGACY_WITHOUT_VERSION" = "${version1}" ]; then
        return 0
    fi

    __is_version_gt "${version2}" "${version1}"
}

BUILD_PATH=$(dirname "${BASH_SOURCE[0]}")/../../..
SOURCE_ROOT=${BUILD_PATH}/sysroot

APP_NAME="casaos-gateway"
APP_NAME_LEGACY="casaos"

# check if migration is needed
SOURCE_BIN_PATH=${SOURCE_ROOT}/usr/bin
SOURCE_BIN_FILE=${SOURCE_BIN_PATH}/${APP_NAME}

CURRENT_BIN_PATH=/usr/bin
CURRENT_BIN_PATH_LEGACY=/usr/local/bin
CURRENT_BIN_FILE=${CURRENT_BIN_PATH}/${APP_NAME}
CURRENT_BIN_FILE_LEGACY=$(realpath -e ${CURRENT_BIN_PATH}/${APP_NAME_LEGACY} || realpath -e ${CURRENT_BIN_PATH_LEGACY}/${APP_NAME_LEGACY} || which ${APP_NAME_LEGACY} || echo CURRENT_BIN_FILE_LEGACY_NOT_FOUND)

SOURCE_VERSION="$(${SOURCE_BIN_FILE} -v)"
CURRENT_VERSION="$(${CURRENT_BIN_FILE} -v || ${CURRENT_BIN_FILE_LEGACY} -v || (stat "${CURRENT_BIN_FILE_LEGACY}" > /dev/null && echo LEGACY_WITHOUT_VERSION) || echo CURRENT_VERSION_NOT_FOUND)"

echo "CURRENT_VERSION: ${CURRENT_VERSION}"
echo "SOURCE_VERSION: ${SOURCE_VERSION}"

NEED_MIGRATION=$(__is_migration_needed "${CURRENT_VERSION}" "${SOURCE_VERSION}" && echo "true" || echo "false")

if [ "${NEED_MIGRATION}" = "false" ]; then
    echo "✅ Migration is not needed."
    exit 0
fi

echo TODO: migrate "${CURRENT_VERSION}" to "${SOURCE_VERSION}"
