#!/bin/bash

apt update -o Dir::Etc::sourcelist="sources.list.d/lambda-guest-agent.list" -o Dir::Etc::sourceparts="-" -o APT::Get::List-Cleanup="0"
apt install lambda-guest-agent
