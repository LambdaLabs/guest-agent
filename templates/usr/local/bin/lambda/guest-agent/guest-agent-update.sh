#!/bin/bash

/usr/bin/apt update -o Dir::Etc::sourcelist="sources.list.d/lambda-guest-agent.list" -o Dir::Etc::sourceparts="-" -o APT::Get::List-Cleanup="0"
/usr/bin/apt install {{ .GuestAgentServiceName }}
