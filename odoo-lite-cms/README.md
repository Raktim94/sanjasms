# Odoo Lite CMS

A lightweight, self-hosted CMS built with Go and SQLite. Docker-based and designed for speed and simplicity.

## Features
- **Admin Dashboard**: Manage your site content securely.
- **Page Management**: Create unlimited pages with custom HTML, clean URLs, and SEO metadata.
- **Menu System**: Dynamic, reorderable navigation menus.
- **File Uploads**: Admin media manager for images and files (Local storage).
- **Site Settings**: Customize title, logo, header, and footer.
- **No Dependencies**: Runs as a single binary or container. No external database required.

## Quick Start (Docker)

1.  Clone the repository:
    ```bash
    git clone <repository-url>
    cd odoo-lite-cms
    ```

2.  Start the service:
    ```bash
    docker compose up -d
    ```

3.  Access the admin panel:
    Open `http://localhost:8080/admin/setup` in your browser.
    Create your first admin account.

4.  Login and start building!

## Development

Prerequisites: Go 1.22+.

1.  Install dependencies:
    ```bash
    go mod tidy
    ```

2.  Run locally:
    ```bash
    go run main.go
    ```

3.  Access at `http://localhost:8080`.

## Architecture
- **Language**: Go
- **Framework**: Echo
- **Database**: SQLite (Pure Go via `glebarez/sqlite`)
- **ORM**: GORM
- **Frontend**: Server-side rendered HTML templates (Go `html/template`).

## License
MIT
