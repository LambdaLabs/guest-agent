#!/bin/bash

systemctl disable {{ .GuestAgentServiceName }}
systemctl disable {{ .GuestAgentUpdaterServiceName }}
systemctl disable {{ .GuestAgentUpdaterTimerName }}
systemctl stop {{ .GuestAgentServiceName }}
systemctl stop {{ .GuestAgentUpdaterTimerName }}
# Note: we can't stop the updater service itself because it needs
# to run to completion during the course of the package installation.
# If we stop it in this script, the apt install command triggered by
# lambda-guest-agent-updater.service will fail because this script would
# be prematurely halting it before completion. It's equivalent to SIGTERM-ing
# "apt install" halfway through its execution.
