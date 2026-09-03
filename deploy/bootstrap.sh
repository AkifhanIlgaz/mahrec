#!/usr/bin/env bash
# One-time hardening + setup for a fresh Hetzner VPS (Debian/Ubuntu).
# Run as root: ssh root@<server-ip> 'bash -s' < deploy/bootstrap.sh
#
# After this finishes:
#   - log in as the `mahrec` sudo user from now on (root SSH login is disabled)
#   - clone the repo into /opt/mahrec and run deploy/deploy.sh to build+start the app
set -euo pipefail

DOMAIN="mahrec.net.tr"
APP_USER="mahrec"
APP_DIR="/opt/mahrec"

if [[ $EUID -ne 0 ]]; then
	echo "Run this as root." >&2
	exit 1
fi

echo "==> Updating base system"
apt-get update -y
apt-get upgrade -y

echo "==> Installing base packages"
apt-get install -y ufw fail2ban unattended-upgrades curl git build-essential

echo "==> Enabling unattended security upgrades"
dpkg-reconfigure -f noninteractive unattended-upgrades

echo "==> Creating app user ($APP_USER)"
if ! id -u "$APP_USER" >/dev/null 2>&1; then
	useradd --create-home --shell /usr/sbin/nologin "$APP_USER"
fi
mkdir -p "$APP_DIR/data"
chown -R "$APP_USER:$APP_USER" "$APP_DIR"

echo "==> Locking down SSH (key-only, no root login)"
SSHD_CONFIG=/etc/ssh/sshd_config
sed -i \
	-e 's/^#\?PermitRootLogin.*/PermitRootLogin no/' \
	-e 's/^#\?PasswordAuthentication.*/PasswordAuthentication no/' \
	-e 's/^#\?ChallengeResponseAuthentication.*/ChallengeResponseAuthentication no/' \
	"$SSHD_CONFIG"
systemctl restart ssh

echo "==> Configuring ufw (deny by default, allow SSH/HTTP/HTTPS)"
ufw default deny incoming
ufw default allow outgoing
ufw allow 22/tcp
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable

echo "==> Enabling fail2ban for SSH"
systemctl enable --now fail2ban

echo "==> Installing Go (for building the app)"
GO_VERSION="1.26.3"
if ! command -v go >/dev/null 2>&1; then
	ARCH="$(dpkg --print-architecture)"
	curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-${ARCH}.tar.gz" -o /tmp/go.tar.gz
	tar -C /usr/local -xzf /tmp/go.tar.gz
	rm /tmp/go.tar.gz
	ln -sf /usr/local/go/bin/go /usr/local/bin/go
	ln -sf /usr/local/go/bin/gofmt /usr/local/bin/gofmt
fi

echo "==> Installing templ and Node (for asset generation)"
env PATH="$PATH:/usr/local/go/bin" GOBIN=/usr/local/bin \
	go install github.com/a-h/templ/cmd/templ@latest
if ! command -v node >/dev/null 2>&1; then
	curl -fsSL https://deb.nodesource.com/setup_lts.x | bash -
	apt-get install -y nodejs
fi

echo "==> Installing Caddy"
if ! command -v caddy >/dev/null 2>&1; then
	apt-get install -y debian-keyring debian-archive-keyring apt-transport-https
	curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' \
		| gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
	curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' \
		| tee /etc/apt/sources.list.d/caddy-stable.list
	apt-get update -y
	apt-get install -y caddy
fi
cat >/etc/caddy/Caddyfile <<EOF
${DOMAIN} {
	reverse_proxy 127.0.0.1:3000
	encode gzip
}
EOF
systemctl reload caddy || systemctl restart caddy

echo "==> Installing systemd unit for the app"
cat >/etc/systemd/system/mahrec.service <<EOF
[Unit]
Description=mahrec goal-tracking app
After=network.target

[Service]
Type=simple
User=${APP_USER}
WorkingDirectory=${APP_DIR}
Environment=ADDR=127.0.0.1:3000
Environment=DB_PATH=${APP_DIR}/data/mahrec.db
ExecStart=${APP_DIR}/bin/server
Restart=on-failure
RestartSec=2

# Basic hardening
NoNewPrivileges=true
ProtectSystem=strict
ReadWritePaths=${APP_DIR}/data
PrivateTmp=true

[Install]
WantedBy=multi-user.target
EOF
systemctl daemon-reload
systemctl enable mahrec

cat <<EOF

==> Bootstrap done.

Next steps:
  1. Point DNS for ${DOMAIN} at this server's IP.
  2. Add your SSH public key to /home/${APP_USER}/.ssh/authorized_keys
     (root login is now disabled).
  3. As ${APP_USER} (or via sudo), clone the repo into ${APP_DIR}:
       git clone <repo-url> ${APP_DIR}
       cd ${APP_DIR} && bash deploy/deploy.sh
  4. Caddy will fetch a Let's Encrypt cert for ${DOMAIN} automatically
     once DNS points here and port 443 is reachable.
EOF
