#!/bin/sh
set -e

# 检查目录是否存在
binDirPath="/betago-webhook"
if [ ! -d "$binDirPath" ]; then
    mkdir $binDirPath
fi

# 如果文件存在，则删除目标文件
binPath="/betago-webhook/webhook"
if [ -f "$binPath" ]; then
    rm $binPath
fi

configDirPath="/config/betago-webhook"
if [ ! -d "$configDirPath" ]; then
    mkdir $configDirPath
fi

configPath="/config/betago-webhook/env"
if [ ! -f "$configPath" ]; then
    printf "ConfigPath /config/betago-webhook/env not found, please configure it with username&password\n configure it like DOCKER_USERNAME_TENCENT=xxx\nDOCKER_PASSWORD_TENCENT=xxx"
    exit 1
fi

go build -o webhook ./*.go
cp webhook /betago-webhook/webhook

cat > /lib/systemd/system/betago-webhook.service << EOF
[Unit]
Description=betago-webhook

[Service]
Type=simple
ExecStart=$binPath
EnvironmentFile=$configPath
Restart=on-failure
RestartSec=10s

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl restart betago-webhook
systemctl enable betago-webhook