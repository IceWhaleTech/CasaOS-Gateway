#!/bin/bash

APP_NAME_SHORT=gateway

__get_setup_script_directory_by_os_release() {
	pushd "$(dirname "${BASH_SOURCE[0]}")/../service.d/${APP_NAME_SHORT}" >/dev/null

	{
		# shellcheck source=/dev/null
		{
			source /etc/os-release
			{
				pushd "${ID}"/"${VERSION_CODENAME}" >/dev/null
			} || {
				pushd "${ID}" >/dev/null
			} || {
				pushd "${ID_LIKE}" >/dev/null
			} || {
				echo "Unsupported OS: ${ID} ${VERSION_CODENAME} (${ID_LIKE})"
				exit 1
			}

			pwd

			popd >/dev/null

		} || {
			echo "Unsupported OS: unknown"
			exit 1
		}

	}

	popd >/dev/null
}

set -e

SETUP_SCRIPT_DIRECTORY="$(__get_setup_script_directory_by_os_release)"
SETUP_SCRIPT_FILENAME="setup-${APP_NAME_SHORT}.sh"

SETUP_SCRIPT_FILEPATH="${SETUP_SCRIPT_DIRECTORY}/${SETUP_SCRIPT_FILENAME}"

echo "✅ Running ${SETUP_SCRIPT_FILEPATH}..."

$SHELL "${SETUP_SCRIPT_FILEPATH}" | sed  "s/^/[${SETUP_SCRIPT_FILENAME}] /"
