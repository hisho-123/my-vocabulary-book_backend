#!/bin/bash
set -e

# サービスが active であることを確認
if ! systemctl is-active --quiet my-vocabulary-book; then
    echo "ERROR: my-vocabulary-book service is not running"
    exit 1
fi

# ヘルスチェックエンドポイントへの疎通確認 (最大30秒待機)
for i in $(seq 1 6); do
    STATUS=$(curl -s -o /dev/null -w "%{http_code}" \
        --max-time 5 http://localhost:80/api/home 2>/dev/null)
    if [ -n "$STATUS" ] && [ "$STATUS" != "000" ]; then
        echo "Health check passed (HTTP $STATUS)"
        exit 0
    fi
    echo "Waiting for service to be ready... (attempt $i/6)"
    sleep 5
done

echo "ERROR: Health check failed after 30 seconds"
exit 1
