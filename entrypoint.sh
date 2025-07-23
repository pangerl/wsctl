#!/usr/bin/env sh

# 设置权限
chmod +x /app/wsctl

exec /app/wsctl "$@"