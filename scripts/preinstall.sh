#!/bin/bash

rm -rf /etc/systemd/system/lambda-guest-agent.service.d/
systemctl stop telegraf
systemctl disable telegraf
echo "ran preinstall"