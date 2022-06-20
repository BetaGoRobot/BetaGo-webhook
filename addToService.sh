#!/bin/sh
set -e

# 检查目录是否存在
binDirPath="/betago-webhook"
if [ ! -d "$binDirPath" ]; then
    mkdir $binDirPath
fi

# 如果文件存在，则删除目标文件
binPath="/betago-webhook/betago-hook.sh"
if [ -f "$binPath" ]; then
    rm $binPath
fi

chmod +x betago-hook.sh
cp betago-hook.sh /betago-webhook/betago-hook.sh

cat > /lib/systemd/system/betago-webhook.service << EOF
[Unit]
Description=betago-webhook

[Service]
Type=simple
ExecStart=sh binPath
Restart=on-failure
RestartSec=10s

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl restart betago-webhook
systemctl enable betago-webhook