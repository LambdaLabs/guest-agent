#!/bin/bash


if [ "$#" -ne 2 ]; then
    echo "Usage: test_script.sh DEB_PATH EXPECTED_VERSION"
    exit 1
fi
DEB_PATH="$1"; shift
EXPECTED_VERSION="$1"; shift

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

function error {
    msg="$1"
    echo "ERROR: $msg"
}

function warning {
    msg="$1"
    echo "WARNING: $msg"
}

function exit_with_error {
    error "$1"
    exit 1
}

function info {
    echo "INFO: $1"
}

function success {
    echo "SUCCESS: $1"
}

dpkg -r lambda-guest-agent
if [ $? -ne 0 ]; then
    exit 1
fi
dpkg -i "$DEB_PATH"
if [ $? -ne 0 ]; then
    exit 1
fi

INSTALLED_VERSION="$(dpkg -s lambda-guest-agent |& grep Version | awk '{print $2}' )"

if [ "$INSTALLED_VERSION" != "$EXPECTED_VERSION" ]; then
    exit_with_error "Installed version does not match expected version: $INSTALLED_VERSION != $EXPECTED_VERSION"
fi
success "Installed version matches expected value: $INSTALLED_VERSION == $EXPECTED_VERSION"

systemctl status lambda-guest-agent.service | grep 'active (running)'
rt=$?
if [ $rt -ne 0 ]; then
    exit_with_error "lambda-guest-agent.service is not running"
fi
success "lambda-guest-agent.service is running"

systemctl status lambda-guest-agent-updater.timer |& egrep -E 'active \(waiting\)|active \(running\)'
rt=$?
if [ $rt -ne 0 ]; then
    exit_with_error "lambda-guest-agent.timer is not running"
fi
success "lambda-guest-agent-updater.timer is active"
