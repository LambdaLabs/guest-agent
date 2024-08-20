#!/bin/bash

systemctl disable {{ .GuestAgentServiceFile }}
systemctl disable {{ .GuestAgentUpdaterServiceFile }}
systemctl disable {{ .GuestAgentUpdaterTimerFile }}
systemctl stop {{ .GuestAgentServiceFile }}
systemctl stop {{ .GuestAgentUpdaterServiceFile }}
systemctl stop {{ .GuestAgentUpdaterTimerFile }}
