#!/usr/bin/env bash
# Run on the VPS, from /opt/mahrec, to update and restart the app.
set -euo pipefail

cd /opt/mahrec

git pull --ff-only
make build

sudo systemctl restart mahrec
sudo systemctl status --no-pager mahrec
