#!/bin/bash

rm -rf /etc/lambda/guest-agent/systemd/
systemctl stop telegraf
systemctl disable telegraf
echo "ran preinstall"