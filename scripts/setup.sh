#!/bin/bash
set -e

# バイナリに実行権限を付与
chmod +x /opt/my-vocabulary-book/my-vocabulary-book_backend

# systemd unit ファイルを配置 (初回のみ)
if [ ! -f /etc/systemd/system/my-vocabulary-book.service ]; then
    cat > /etc/systemd/system/my-vocabulary-book.service <<'EOF'
[Unit]
Description=my-vocabulary-book backend
After=network.target

[Service]
ExecStart=/opt/my-vocabulary-book/my-vocabulary-book_backend
Restart=always
EnvironmentFile=/etc/my-vocabulary-book/env
StandardOutput=append:/var/log/app/my-vocabulary-book.log
StandardError=append:/var/log/app/my-vocabulary-book.log

[Install]
WantedBy=multi-user.target
EOF
    systemctl daemon-reload
    systemctl enable my-vocabulary-book
fi

# ログディレクトリを確保
mkdir -p /var/log/app
touch /var/log/app/my-vocabulary-book.log

exit 0
