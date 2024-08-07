#!/bin/bash

systemctl enable /etc/lambda/guest-agent/systemd/lambda-guest-agent.service
systemctl daemon-reload
systemctl restart lambda-guest-agent.service
