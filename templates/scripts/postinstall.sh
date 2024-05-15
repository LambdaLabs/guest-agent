#!/bin/bash

systemctl enable {{ .GuestAgentServiceFile }}
systemctl enable {{ .GuestAgentUpdaterServiceFile }}
systemctl enable {{ .GuestAgentUpdaterTimerFile }}
systemctl daemon-reload
systemctl restart {{ .GuestAgentServiceName }}
systemctl restart {{ .GuestAgentUpdaterTimerName }}
