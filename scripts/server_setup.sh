#!/bin/bash
set -e

APP_USER=wrzapi
APP_DIR=/opt/wrzapi
# 替换为你的 GitHub Actions 私钥对应的公钥
SSH_PUBLIC_KEY="ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCziozWZC8R8zMGY6W/e6fxcnynoP+OxyCHAG4vSb5SlYmk1/Tj5jqcriTE3hkYEhm/Z6tg1CmTW4FUCpkd9bf3C3GXAlqZ0LxBag7EeRt9lscIQ+hZXc9EisgJU1lF+0iH1VKJ8W9GJInqzxBhUA1H4naEQNZQ3FAveuBTlESZh6+ft2UOyqwyd3KN/kBGEnNDW76cJcwLKGGJpRAWSFXDwuRVx38lZbKc6kIswZDwRqMxfgGAwVMWciUCzuCurSAqYHNHx1hWYoAqorv1CsK0rKFD6sGsWj/dnq+wtO0dw8VQROqIeL5jcCJfmNuTU6weMFbqXRyErRtxafFCE5xvh6/BaMYxQfNMwzXOor1prxjdHT9ARuSMlyCpGfEuG/0P3W/zmqm0WsyDPQ2wM6viNG1CYGkMQCvZkU1iyAyDtyZ0BydYP9G2yOwL4cYb2z/hYV2Mg6enekpVORHBoAPRsu8aMmZUaiLGNrYlWoP00kW7OSZ3yq+IgeOq7oabpBk= home@DESKTOP-PGKGGJH"

# 1. 容错：用户已存在则跳过（正常提示）
sudo useradd -r -s /usr/sbin/nologin "$APP_USER" || true

# 2. 创建应用目录（仅 releases）
sudo mkdir -p "$APP_DIR"/releases
sudo chown -R "$APP_USER":"$APP_USER" "$APP_DIR"

# ========== 核心修复：彻底重置 wrzapi 家目录权限 ==========
# 先删除错误权限的目录（如果存在）
sudo rm -rf /home/"$APP_USER"
# 重新创建家目录，直接设置所有权为 wrzapi（root 权限执行）
sudo mkdir -p /home/"$APP_USER"
sudo chown -R "$APP_USER":"$APP_USER" /home/"$APP_USER"
sudo chmod 755 /home/"$APP_USER"  # 基础权限：用户可读写执行，其他只读

# 3. 配置 .ssh 目录（此时 wrzapi 有足够权限）
sudo -u "$APP_USER" mkdir -p /home/"$APP_USER"/.ssh
sudo -u "$APP_USER" touch /home/"$APP_USER"/.ssh/authorized_keys
# 避免重复添加公钥
if ! grep -qxF "$SSH_PUBLIC_KEY" /home/"$APP_USER"/.ssh/authorized_keys; then
    sudo -u "$APP_USER" echo "$SSH_PUBLIC_KEY" >> /home/"$APP_USER"/.ssh/authorized_keys
fi
# 强制设置 SSH 严格权限（必须是 700/600）
sudo chmod 700 /home/"$APP_USER"/.ssh
sudo chmod 600 /home/"$APP_USER"/.ssh/authorized_keys
sudo chown -R "$APP_USER":"$APP_USER" /home/"$APP_USER"/.ssh

# 4. 配置 sudo 免密
SUDO_CONFIG="$APP_USER ALL=(ALL) NOPASSWD: /usr/bin/systemctl daemon-reload, /usr/bin/systemctl restart wrzapi, /usr/bin/systemctl status wrzapi"
if ! grep -qxF "$SUDO_CONFIG" /etc/sudoers.d/wrzapi; then
    echo "$SUDO_CONFIG" | sudo tee /etc/sudoers.d/wrzapi >/dev/null
    sudo chmod 0440 /etc/sudoers.d/wrzapi
fi

# 5. 创建服务文件
sudo tee /etc/systemd/system/wrzapi.service >/dev/null <<'EOF'
[Unit]
Description=WRZ API Service
After=network.target

[Service]
Type=simple
User=wrzapi
WorkingDirectory=/opt/wrzapi
Environment=PORT=12088
Environment=SERVER_URL=https://api.notelook.me
ExecStart=/opt/wrzapi/current
Restart=on-failure
RestartSec=3

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable wrzapi

echo "✅ 服务器初始化完成！所有权限配置正常。"
