#!/bin/bash

set -e

readonly CASA_EXEC=casaos-gateway
readonly CASA_SERVICE=casaos-gateway.service

CASA_SERVICE_PATH=$(systemctl show ${CASA_SERVICE} --no-pager  --property FragmentPath | cut -d'=' -sf2)
readonly CASA_SERVICE_PATH

CASA_CONF=$( grep -i ConditionFileNotEmpty "${CASA_SERVICE_PATH}" | cut -d'=' -sf2)
if [[ -z "${CASA_CONF}" ]]; then
    CASA_CONF=/etc/casaos/gateway.ini
fi

readonly aCOLOUR=(
    '\e[38;5;154m' # green  	| Lines, bullets and separators
    '\e[1m'        # Bold white	| Main descriptions
    '\e[90m'       # Grey		| Credits
    '\e[91m'       # Red		| Update notifications Alert
    '\e[33m'       # Yellow		| Emphasis
)

Show() {
    # OK
    if (($1 == 0)); then
        echo -e "${aCOLOUR[2]}[$COLOUR_RESET${aCOLOUR[0]}  OK  $COLOUR_RESET${aCOLOUR[2]}]$COLOUR_RESET $2"
    # FAILED
    elif (($1 == 1)); then
        echo -e "${aCOLOUR[2]}[$COLOUR_RESET${aCOLOUR[3]}FAILED$COLOUR_RESET${aCOLOUR[2]}]$COLOUR_RESET $2"
    # INFO
    elif (($1 == 2)); then
        echo -e "${aCOLOUR[2]}[$COLOUR_RESET${aCOLOUR[0]} INFO $COLOUR_RESET${aCOLOUR[2]}]$COLOUR_RESET $2"
    # NOTICE
    elif (($1 == 3)); then
        echo -e "${aCOLOUR[2]}[$COLOUR_RESET${aCOLOUR[4]}NOTICE$COLOUR_RESET${aCOLOUR[2]}]$COLOUR_RESET $2"
    fi
}

Warn() {
    echo -e "${aCOLOUR[3]}$1$COLOUR_RESET"
}

trap 'onCtrlC' INT
onCtrlC() {
    echo -e "${COLOUR_RESET}"
    exit 1
}

if [[ ! -x "$(command -v ${CASA_EXEC})" ]]; then
    Show 2 "${CASA_EXEC} is not detected, exit the script."
    exit 1
fi

Show 2 "Stopping ${CASA_SERVICE}..."
systemctl disable --now "${CASA_SERVICE}" || Show 3 "Failed to disable ${CASA_SERVICE}"

rm -rvf "$(which ${CASA_EXEC})" || Show 3 "Failed to remove ${CASA_EXEC}"
rm -rvf "${CASA_CONF}" || Show 3 "Failed to remove ${CASA_CONF}"

rm -rvf /var/run/casaos/gateway.pid
rm -rvf /var/run/casaos/management.url
rm -rvf /var/run/casaos/routes.json
rm -rvf /var/run/casaos/static.url
rm -rvf /var/lib/casaos/www
