#!/bin/bash
set -e

# サービスが存在する場合のみ停止 (初回デプロイ時はスキップ)
if systemctl is-active --quiet my-vocabulary-book; then
    systemctl stop my-vocabulary-book
fi

exit 0
