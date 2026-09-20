# Live Server Installation & Deployment Guide

This guide provides step-by-step instructions for deploying the **Admin Dashboard** application on a live production Linux server (Ubuntu 22.04 / 24.04 LTS or Debian 12).

---

## Table of Contents

1. [System Requirements](#1-system-requirements)
2. [Server Preparation & Dependencies](#2-server-preparation--dependencies)
3. [MySQL Database Setup](#3-mysql-database-setup)
4. [Application Setup & Directory Permissions](#4-application-setup--directory-permissions)
5. [Environment Configuration (`.env`)](#5-environment-configuration-env)
6. [Database Migration & Seeding](#6-database-migration--seeding)
7. [Building Production Binaries](#7-building-production-binaries)
8. [Systemd Service Setup (Auto-Restart & Daemon)](#8-systemd-service-setup-auto-restart--daemon)
   - [Web Application Service](#web-application-service)
   - [Queue Worker Service](#queue-worker-service)
9. [Nginx Reverse Proxy & SSL Configuration](#9-nginx-reverse-proxy--ssl-configuration)
10. [Firewall & Security Hardening](#10-firewall--security-hardening)
11. [Monitoring, Logs & Maintenance](#11-monitoring-logs--maintenance)
12. [Zero-Downtime Update / Redeployment](#12-zero-downtime-update--redeployment)
13. [Troubleshooting Guide](#13-troubleshooting-guide)

---

## 1. System Requirements

| Component | Minimum | Recommended |
|---|---|---|
| **Operating System** | Ubuntu 20.04 / 22.04 / 24.04 LTS | Ubuntu 24.04 LTS |
| **CPU** | 1 vCPU | 2+ vCPUs |
| **RAM** | 1 GB | 2 GB+ |
| **Disk** | 20 GB SSD | 40 GB+ SSD |
| **Go Runtime** | Go 1.22+ | Go 1.22 or 1.23 |
| **Database** | MySQL 8.0+ or MariaDB 10.6+ | MySQL 8.0+ |
| **Web Server** | Nginx 1.18+ | Nginx 1.24+ |

---

## 2. Server Preparation & Dependencies

Update packages and install core build tools, MySQL, Nginx, and Git:

```bash
sudo apt update && sudo apt upgrade -y
sudo apt install -y curl wget git build-essential ufw nginx certbot python3-certbot-nginx mysql-server
```

### Install Go (Golang)

Download and install the official Go binary:

```bash
# Download Go 1.23.0 (or check https://go.dev/dl/ for latest)
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz

# Remove previous installation and extract
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz

# Add Go to system PATH
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify Go installation
go version
```

---

## 3. MySQL Database Setup

### Secure MySQL Installation

```bash
sudo mysql_secure_installation
```

Follow the prompts to set a strong root password and remove test databases.

### Create Database and User

Log in to MySQL:

```bash
sudo mysql -u root -p
```

Execute the following SQL commands (replace `your_secure_db_password` with a strong password):

```sql
-- Create database with utf8mb4 charset for full Unicode/emoji support
CREATE DATABASE admin_dashboard CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Create application user
CREATE USER 'admin_user'@'localhost' IDENTIFIED BY 'your_secure_db_password';

-- Grant full privileges on the database
GRANT ALL PRIVILEGES ON admin_dashboard.* TO 'admin_user'@'localhost';

-- Apply privilege changes
FLUSH PRIVILEGES;

-- Exit MySQL
EXIT;
```

---

## 4. Application Setup & Directory Permissions

### Create Application Directory

We will deploy to `/var/www/admin-dashboard`:

```bash
sudo mkdir -p /var/www/admin-dashboard
sudo chown -R $USER:$USER /var/www/admin-dashboard
cd /var/www/admin-dashboard
```

### Transfer / Clone Application Files

Clone your repository or upload your files:

```bash
git clone <YOUR_REPOSITORY_URL> .
```

### Create Runtime Storage Directories

Ensure the required storage folders exist and have appropriate write permissions:

```bash
mkdir -p storage/uploads storage/backup storage/cache storage/sessions bin

# Set write permissions for web user (www-data)
sudo chown -R www-data:www-data storage assets
sudo chmod -R 775 storage
```

---

## 5. Environment Configuration (`.env`)

Copy the template configuration file:

```bash
cp .env.example .env
chmod 600 .env
nano .env
```

### Configure Required Environment Variables:

```env
# Server Configuration
APP_NAME=Admin-Dashboard
APP_KEY=your_random_32_character_secret_key_here
APP_IP=127.0.0.1
APP_DOMAIN=dashboard.yourdomain.com
APP_PORT=3000
APP_BASE_URL=https://dashboard.yourdomain.com

# Assets Configuration
ASSET_IP=127.0.0.1
ASSET_DOMAIN=dashboard.yourdomain.com
ASSET_PORT=3000
ASSET_BASE_URL=https://dashboard.yourdomain.com

# Database Connection
DB_CONNECTION=mysql
DB_RECORD_NOT_FOUND_MESSAGE="record not found"
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=admin_dashboard
DB_USERNAME=admin_user
DB_PASSWORD=your_secure_db_password

# Web Security & Cookie Settings
WEB_COOKIE_NAME=_admin_session
WEB_COOKIE_DOMAIN=dashboard.yourdomain.com
WEB_COOKIE_PATH=/
WEB_COOKIE_MAX_AGE=3600
WEB_COOKIE_SECURE=true
WEB_COOKIE_HTTP_ONLY=true
WEB_COOKIE_SAME_SITE=strict
WEB_CSRF_TOKEN_LOOKUP=header:X-CSRF-Token,form:_csrf
WEB_CSRF_TOKEN_LENGTH=32
WEB_SECURE_HSTS_MAX_AGE=63072000
WEB_RATE_LIMIT=60

# File Storage
FILE_SYSTEM=local
FILE_SYSTEM_LOCATION=./storage/uploads
FILE_BACKUP_LOCATION=./storage/backup

# Session Configuration (database or file)
SESSION_DRIVER=database
SESSION_PATH=storage/sessions
SESSION_TABLE=sessions
SESSION_LIFETIME=120

# Cache Configuration
CACHE_DRIVER=file
CACHE_PATH=storage/cache
CACHE_TIME=60

# Queue Worker Configuration
QUEUE_WORKER_INTERVAL=5
QUEUE_WORKER_BATCH_SIZE=50
QUEUE_WORKER_MAX_TRIES=5
```

Save and exit (`Ctrl+O`, `Enter`, `Ctrl+X`).

---

## 6. Database Migration & Seeding

Run the database migration to create all tables (including users, settings, sessions, audit trails, mail & SMS histories):

```bash
cd /var/www/admin-dashboard
go run ./cmd/internal/databases/migration
```

You should see:
```text
Database migrated successfully
```

---

## 7. Building Production Binaries

Compile optimized, stripped binaries for both the Web Application and the Background Queue Worker:

```bash
cd /var/www/admin-dashboard

# Build Web Server
go build -ldflags="-s -w" -o bin/admin-dashboard ./cmd/app/web

# Build Standalone Queue Worker
go build -ldflags="-s -w" -o bin/admin-worker ./cmd/app/worker

# Verify binaries exist and are executable
chmod +x bin/admin-dashboard bin/admin-worker
ls -lh bin/
```

---

## 8. Systemd Service Setup (Auto-Restart & Daemon)

We will configure systemd to automatically start both services upon server boot and restart them on crashes.

### Web Application Service

Create `/etc/systemd/system/admin-dashboard.service`:

```bash
sudo nano /etc/systemd/system/admin-dashboard.service
```

Paste the following configuration:

```ini
[Unit]
Description=Admin Dashboard Web Application
After=network.target mysql.service
Wants=mysql.service

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/var/www/admin-dashboard
ExecStart=/var/www/admin-dashboard/bin/admin-dashboard
Restart=always
RestartSec=5s
LimitNOFILE=65535
EnvironmentFile=/var/www/admin-dashboard/.env

# Security sandbox options
ProtectSystem=full
ProtectHome=true
NoNewPrivileges=true

[Install]
WantedBy=multi-user.target
```

### Queue Worker Service

Create `/etc/systemd/system/admin-worker.service`:

```bash
sudo nano /etc/systemd/system/admin-worker.service
```

Paste the following configuration:

```ini
[Unit]
Description=Admin Dashboard Pending Email & SMS Queue Worker
After=network.target mysql.service
Wants=mysql.service

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/var/www/admin-dashboard
ExecStart=/var/www/admin-dashboard/bin/admin-worker
Restart=always
RestartSec=5s
LimitNOFILE=65535
EnvironmentFile=/var/www/admin-dashboard/.env

[Install]
WantedBy=multi-user.target
```

### Enable and Start Services

```bash
# Reload systemd daemon
sudo systemctl daemon-reload

# Enable services on system boot
sudo systemctl enable admin-dashboard
sudo systemctl enable admin-worker

# Start both services
sudo systemctl start admin-dashboard
sudo systemctl start admin-worker

# Check service status
sudo systemctl status admin-dashboard --no-pager
sudo systemctl status admin-worker --no-pager
```

---

## 9. Nginx Reverse Proxy & SSL Configuration

### Create Nginx Server Block

Create `/etc/nginx/sites-available/admin-dashboard`:

```bash
sudo nano /etc/nginx/sites-available/admin-dashboard
```

Paste the following Nginx configuration (replace `dashboard.yourdomain.com` with your real domain):

```nginx
server {
    listen 80;
    server_name dashboard.yourdomain.com;

    # Client body limit for file uploads (e.g. 50MB)
    client_max_body_size 50M;

    # Static assets caching
    location /assets/ {
        alias /var/www/admin-dashboard/assets/;
        expires 30d;
        add_header Cache-Control "public, no-transform";
        access_log off;
    }

    # Uploaded files
    location /storage/uploads/ {
        alias /var/www/admin-dashboard/storage/uploads/;
        expires 30d;
        add_header Cache-Control "public, no-transform";
        access_log off;
    }

    # Proxy to Go application
    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_http_version 1.1;

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket support
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";

        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
}
```

### Enable Site and Test Nginx

```bash
# Enable site
sudo ln -sf /etc/nginx/sites-available/admin-dashboard /etc/nginx/sites-enabled/

# Remove default site if present
sudo rm -f /etc/nginx/sites-enabled/default

# Test Nginx syntax
sudo nginx -t

# Reload Nginx
sudo systemctl reload nginx
```

### Obtain Free Let's Encrypt SSL Certificate

Ensure your domain's DNS `A` record points to your server IP, then run:

```bash
sudo certbot --nginx -d dashboard.yourdomain.com
```

Select option `2` to redirect HTTP traffic to HTTPS automatically.

Certbot will automatically update your Nginx configuration and set up automatic renewal cron jobs. Test renewal:

```bash
sudo certbot renew --dry-run
```

---

## 10. Firewall & Security Hardening

Enable the UFW firewall to only permit SSH, HTTP, and HTTPS:

```bash
# Allow SSH
sudo ufw allow OpenSSH

# Allow Web traffic
sudo ufw allow 'Nginx Full'

# Enable firewall
sudo ufw --force enable

# Check status
sudo ufw status verbose
```

### Protect Sensitive Files

```bash
# Ensure .env is only readable by root and www-data
sudo chown root:www-data /var/www/admin-dashboard/.env
sudo chmod 640 /var/www/admin-dashboard/.env
```

---

## 11. Monitoring, Logs & Maintenance

### View Real-Time Application Logs

```bash
# Web application logs
sudo journalctl -u admin-dashboard -f

# Queue worker logs (shows pending email/SMS dispatch in real time)
sudo journalctl -u admin-worker -f
```

### Common Service Commands

| Action | Web Service | Queue Worker |
|---|---|---|
| **Restart** | `sudo systemctl restart admin-dashboard` | `sudo systemctl restart admin-worker` |
| **Stop** | `sudo systemctl stop admin-dashboard` | `sudo systemctl stop admin-worker` |
| **Status** | `sudo systemctl status admin-dashboard` | `sudo systemctl status admin-worker` |

---

## 12. Zero-Downtime Update / Redeployment

When deploying new updates or git commits in the future, use this deploy script:

Save as `/var/www/admin-dashboard/deploy.sh`:

```bash
#!/bin/bash
set -e

echo "=== Starting Deployment ==="
cd /var/www/admin-dashboard

# 1. Pull latest code
git pull origin main

# 2. Build new binaries
go build -ldflags="-s -w" -o bin/admin-dashboard.new ./cmd/app/web
go build -ldflags="-s -w" -o bin/admin-worker.new ./cmd/app/worker

# 3. Swap binaries atomically
mv bin/admin-dashboard.new bin/admin-dashboard
mv bin/admin-worker.new bin/admin-worker
chmod +x bin/admin-dashboard bin/admin-worker

# 4. Run any pending database migrations
go run ./cmd/internal/databases/migration

# 5. Restart services
sudo systemctl restart admin-dashboard
sudo systemctl restart admin-worker

echo "=== Deployment Completed Successfully ==="
```

Make it executable:
```bash
chmod +x /var/www/admin-dashboard/deploy.sh
```

---

## 13. Troubleshooting Guide

### 1. Database Connection Error (`connection refused` or `access denied`)
- Verify MySQL is running: `sudo systemctl status mysql`
- Verify `.env` credentials match your MySQL user:
  ```bash
  mysql -u admin_user -p -h 127.0.0.1 admin_dashboard
  ```
- Check MySQL bind address in `/etc/mysql/mysql.conf.d/mysqld.cnf` (`bind-address = 127.0.0.1`).

### 2. Port 3000 Already in Use
- Check what process is using port 3000:
  ```bash
  sudo lsof -i :3000
  ```

### 3. File Upload / Permission Errors
- Ensure `storage/uploads` and `storage/backup` are owned by `www-data`:
  ```bash
  sudo chown -R www-data:www-data /var/www/admin-dashboard/storage
  sudo chmod -R 775 /var/www/admin-dashboard/storage
  ```

### 4. Nginx 502 Bad Gateway
- Ensure the Go service is active:
  ```bash
  sudo systemctl status admin-dashboard
  ```
- Verify the app is listening on `127.0.0.1:3000`:
  ```bash
  curl -I http://127.0.0.1:3000/
  ```
