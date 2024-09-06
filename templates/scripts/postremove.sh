#!/bin/bash

systemctl disable {{ .GuestAgentServiceName }}
systemctl disable {{ .GuestAgentUpdaterServiceName }}
systemctl disable {{ .GuestAgentUpdaterTimerName }}
systemctl stop {{ .GuestAgentServiceName }}
systemctl stop {{ .GuestAgentUpdaterServiceName }}
systemctl stop {{ .GuestAgentUpdaterTimerName }}
