#!/bin/bash

systemctl daemon-reload
systemctl restart lambda-guest-agent.service
