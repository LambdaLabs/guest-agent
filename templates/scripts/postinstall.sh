#!/bin/bash

systemctl enable {{ .GuestAgentServiceFile }}
systemctl daemon-reload
systemctl restart {{ .GuestAgentServiceName }}
