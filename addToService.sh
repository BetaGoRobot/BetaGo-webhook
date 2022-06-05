set -e
pwd=`pwd`

go build -o webhook *.go

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
systemctl restart webhook
systemctl enable webhook