#!/bin/bash
set -e

# サービスが active であることを確認
if ! systemctl is-active --quiet my-vocabulary-book; then
    echo "ERROR: my-vocabulary-book service is not running"
    exit 1
fi

# ヘルスチェックエンドポイントへの疎通確認 (最大30秒待機)
for i in $(seq 1 6); do
    if curl -sf http://localhost:8080/health > /dev/null 2>&1; then
        echo "Health check passed"
        exit 0
    fi
    echo "Waiting for service to be ready... (attempt $i/6)"
    sleep 5
done

echo "ERROR: Health check failed after 30 seconds"
exit 1
