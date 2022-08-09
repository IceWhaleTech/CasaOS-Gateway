#!/bin/bash

set -e

APP_NAME="casaos-gateway"

echo "Stopping service..."
systemctl stop --force --no-ask-password "${APP_NAME}.service"

echo "Disabling service..."
systemctl disable --force --no-ask-password "${APP_NAME}.service"
