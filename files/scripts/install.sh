#!/bin/bash

export DEBIAN_FRONTEND=noninteractive
export NEEDRESTART_MODE=l

# Just in case for some weird reason, these dirs don't exist.
mkdir -p /usr/share/keyrings/
mkdir -p /etc/apt/sources.list.d/

curl https://lambdalabs-guest-agent.s3.us-west-2.amazonaws.com/keys/lambda-guest-agent.pub.gpg | gpg --dearmor > /usr/share/keyrings/lambda-guest-agent-archive-keyring.pgp
echo 'deb [signed-by=/usr/share/keyrings/lambda-guest-agent-archive-keyring.pgp] https://lambdalabs-guest-agent.s3.us-west-2.amazonaws.com stable main' > /etc/apt/sources.list.d/lambda-guest-agent.list
apt-get -y update
apt-get install lambda-guest-agent
