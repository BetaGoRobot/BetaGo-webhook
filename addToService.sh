set -e
pwd=`pwd`

cat > /lib/systemd/system/webhook.service << EOF
[Unit]
Description=go-webhook

[Service]
Type=simple
ExecStart=$pwd/webhook
Restart=on-failure
RestartSec=10s

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl start webhook
systemctl enable webhook