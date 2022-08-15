#!/bin/bash

set -e 

INSTALL_ROOT=${1:-/}

APP_NAME_SHORT=gateway

__get_cleanup_script_directory_by_os_release() {
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

CLEANUP_SCRIPT_DIRECTORY="$(__get_cleanup_script_directory_by_os_release)"
CLEANUP_SCRIPT_FILENAME="cleanup-${APP_NAME_SHORT}.sh"

CLEANUP_SCRIPT_FILEPATH="${CLEANUP_SCRIPT_DIRECTORY}/${CLEANUP_SCRIPT_FILENAME}"

{
    echo "🟩 Running ${CLEANUP_SCRIPT_FILENAME}..."
    $SHELL "${CLEANUP_SCRIPT_FILEPATH}" "${INSTALL_ROOT}"
} || {
    echo "🟥 ${CLEANUP_SCRIPT_FILENAME} failed."
    exit 1
}

echo "✅ ${CLEANUP_SCRIPT_FILENAME} finished."
