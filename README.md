# SanjanaCMS (Odoo Lite)

SanjanaCMS is a lightweight, production-ready Content Management System (CMS) built to replicate the core website-building experience of Odoo, without the complexity of an ERP. It focuses on speed, simplicity, and ease of deployment.

---

## 🛠️ Technology Stack

This project is built using a modern, efficient, and minimal technology stack:

-   **Backend:** [Go (Golang)](https://go.dev/) - Chosen for its performance, type safety, and ability to compile into a single static binary.
-   **Web Framework:** [Echo (v4)](https://echo.labstack.com/) - A high-performance, minimalist Go web framework.
-   **Database:** [SQLite](https://www.sqlite.org/) - A serverless, zero-configuration database. We use a pure Go driver (`modernc.org/sqlite`) to maintain a zero-dependency environment.
-   **Templating:** [Go `html/template`](https://pkg.go.dev/html/template) - Native Go templates for secure, server-side rendering.
-   **Frontend:** Vanilla CSS & JavaScript - No heavy frameworks or build steps required for the frontend, ensuring instant load times.
-   **Containerization:** [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/) - For seamless deployment and environment consistency.

---

## 🏗️ Architecture & How It Was Built

SanjanaCMS follows a clean, modular architecture:

### 1. **Server-Side Rendering (SSR)**
The CMS uses Go's native templating engine to render pages on the server. This ensures excellent SEO, fast initial load times, and a simpler development flow. Themes and layouts are structured to be easily customizable.

### 2. **Modular Handlers**
The logic is split into three main areas:
-   **Public Handlers:** Serve the live website and handle dynamic page routing.
-   **Admin Handlers:** Manage the backend dashboard, authentication, and CRUD operations for pages, menus, and settings.
-   **Auth Handlers:** Manage secure login sessions and initial setup.

### 3. **Smart Routing**
A catch-all route `/*` in Go Echo processes requests. It checks the database for a matching URL slug. if a page exists, it renders it using the dynamic page builder blocks; otherwise, it serves a custom 404 page.

### 4. **Persistent Storage**
Even though it runs in Docker, all data is persistent.
-   **Database:** Stored in `/data/cms.db` within the container, mapped to a Docker volume.
-   **Media:** Uploaded files are stored in `static/uploads`, also mapped to the host system.

---

## 🚀 Getting Started

### Prerequisites
-   Docker installed on your system.
-   (Optional) Go 1.25+ if you wish to run/build locally without Docker.

### 1. Permission Fix (Linux/WSL)
If you encounter a "permission denied" error when running Docker, run:
```bash
sudo usermod -aG docker $USER
newgrp docker
```

### 2. Build and Start
Run the following command to build the image and start the container in the background:
```bash
docker compose up -d --build
```

### 3. Access the Site
-   **Live Website:** [http://localhost:2512](http://localhost:2512)
-   **Admin Dashboard:** [http://localhost:2512/admin/login](http://localhost:2512/admin/login)

---

## 🔐 Configuration & Defaults

### Initial Setup
The first time you run the CMS, navigate to `/admin/setup` (or try to login) to create the primary administrator account.

**Default Admin (if pre-seeded):**
-   **Email:** `hi@RAKTIMranjit.in`
-   **Password:** `hi@raktim`

### Environment Variables
You can customize the deployment in `docker-compose.yml`:
-   `DB_Path`: Path to the SQLite database file.
-   `PORT`: Internal port the Go server listens on (default: 8080).

---

## 📋 Viewing Logs

If you encounter issues or want to see the server activity, use these commands:

-   **Follow logs in real-time:**
    ```bash
    docker compose logs -f
    ```
-   **View last 100 lines:**
    ```bash
    docker compose logs --tail=100
    ```
-   **Show logs for the CMS container only:**
    ```bash
    docker compose logs cms
    ```

---

## 📦 Features at a Glance

-   **Dynamic Page Builder:** Create pages using Drag-and-Drop style blocks (Text, Images, Buttons, Custom HTML).
-   **Menu Manager:** Drag-and-drop to reorder navigation items and create submenus.
-   **Media Library:** Built-in upload manager for images and PDFs (up to 100MB).
-   **SEO Suite:** Full control over Meta Tags, Open Graph data, Auto-generated `sitemap.xml`, and `robots.txt`.
-   **Branding:** Centralized settings for Logo, Favicon, and Brand Colors.

---

## 🤝 Contributing
Feel free to fork this project and submit pull requests. For major changes, please open an issue first to discuss what you would like to change.

---
© 2026 SanjanaCMS - Built with ❤️ for speed and simplicity.
