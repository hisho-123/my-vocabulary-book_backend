#!/bin/bash
set -e

# バイナリに実行権限を付与
chmod +x /opt/my-vocabulary-book/my-vocabulary-book_backend

# env ディレクトリを確保し、SSM Parameter Store から env ファイルを取得
mkdir -p /etc/my-vocabulary-book
aws ssm get-parameter \
  --name "/mvb/backend/env" \
  --with-decryption \
  --query Parameter.Value \
  --output text \
  --region ap-northeast-1 \
  > /etc/my-vocabulary-book/env

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

# DB マイグレーション
yum install -y mariadb105 -q
set -a
source /etc/my-vocabulary-book/env
set +a
mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" < /opt/my-vocabulary-book/schema.sql

exit 0
