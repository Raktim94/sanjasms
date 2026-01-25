# sanjanacms (Odoo Lite CMS)

A lightweight, production-ready, self-hosted CMS designed to replicate the core website building features of Odoo. Built with Go (Echo), SQLite, and Docker.

## 🚀 Quick Start (Docker)

1.  **Start the server:**
    ```bash
    docker compose up -d --build
    ```

2.  **Access the Admin Panel:**
    Open [http://localhost:8080/admin/login](http://localhost:8080/admin/login)

    **Default Credentials:**
    -   **Email:** `hi@RAKTIMranjit.in`
    -   **Password:** `hi@raktim`

3.  **View the Website:**
    Open [http://localhost:8080](http://localhost:8080)

## ✨ Features

-   **Lightweight & Fast:** Built on Go and SQLite.
-   **Docker Ready:** Simple deployment with `docker-compose`.
-   **Page Builder:** Create static and landing pages with custom HTML.
-   **Menu Management:** Drag-and-drop menu builder (nested menus supported).
-   **Media Library:** Upload images and files.
-   **SEO Tools:** Manage Meta Titles, Descriptions, OG Tags, and more.
-   **Premium UI:** "Apple-like" aesthetic for the admin interface.

## 🛠 Tech Stack

-   **Backend:** Golang (Echo Framework)
-   **Database:** SQLite (with GORM)
-   **Frontend:** Server-Side Rendered HTML (Go Templates) + Vanilla CSS (Premium Design)
-   **Deployment:** Docker (Alpine based)

## 📁 Project Structure

-   `/handlers`: HTTP request handlers
-   `/models`: Database models (GORM)
-   `/static`: CSS, JS, and user uploads
-   `/templates`: HTML templates
-   `main.go`: Application entry point

---
Odoo Lite © 2026 Odoo Lite CMS
