#!/bin/bash

systemctl enable {{ .GuestAgentServiceFile }}
systemctl enable {{ .GuestAgentUpdaterTimerFile }}
systemctl daemon-reload
systemctl restart {{ .GuestAgentServiceName }}
