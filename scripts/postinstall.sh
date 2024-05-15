#!/bin/bash

systemctl daemon-reload
systemctl restart lambda-metrics-stack.service
