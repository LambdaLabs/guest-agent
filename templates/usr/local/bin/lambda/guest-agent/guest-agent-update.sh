#!/bin/bash

apt update -o Dir::Etc::sourcelist="sources.list.d/lambda-guest-agent.list" -o Dir::Etc::sourceparts="-" -o APT::Get::List-Cleanup="0"
old_version="$(apt info lambda-guest-agent |& grep 'Version' | awk '{print $2}')"
apt install lambda-guest-agent
new_version="$(apt info lambda-guest-agent |& grep 'Version' | awk '{print $2}')"
if [ "$old_version" -ne "$new_version"]; then   
    echo "upgraded $old_version to $new_version, restarting service".
    systemctl restart {{ .GuestAgentServiceName }}
else
    echo "no new version installed"
fi
