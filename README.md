# Admin Dashboard

A modern, high-performance Go-based Admin Dashboard web application built with Echo framework, GORM, MySQL, and Vanilla CSS/JS.

## Features

- **System & Business Settings**: Full dynamic configuration of business profile, timezone, currency, and general settings.
- **Email Gateway Management**: Live configuration and credentials encryption (AES-GCM 256-bit) for **SendGrid**, **SMTP**, and **AWS SES** drivers with test email dispatch and audit logging.
- **SMS Gateway Management**: SMS gateway configuration with SSLCommerz integration and audit trail.
- **Continuous Background Queue Worker**: Resilient queue worker process that continuously polls and dispatches pending emails and SMS with retry tracking and audit logging.
- **Role-Based Access Control (RBAC)**: Comprehensive permission group, category, and right action management.
- **Audit Trails & Activity Logs**: Real-time logging of user actions and administrative activities across modules.
- **Database Backup & Restore**: Automated database backup generator with support for local and AWS S3 storage.
- **Authentication & Security**: CSRF protection, secure cookie sessions, rate limiting, and password reset workflows.

---

## Quick Start (Development)

### 1. Prerequisites
- Go 1.22+ installed
- MySQL 8.0+ running

### 2. Setup Environment
```bash
cp .env.example .env
# Edit .env to match your local database and credentials
```

### 3. Run Database Migration & Seeds
```bash
go run ./cmd/internal/databases/migration
```

### 4. Run Web Application
```bash
go run ./cmd/app/web
```
Open your browser at [http://localhost:3000](http://localhost:3000).

### 5. Run Standalone Queue Worker (Optional)
```bash
go run ./cmd/app/worker
```

---

## Production Deployment

For complete, step-by-step instructions on deploying to an Ubuntu/Debian production server with Nginx, SSL, and Systemd services, see:

👉 [**Live Server Installation Guide (INSTALLATION.md)**](INSTALLATION.md)
