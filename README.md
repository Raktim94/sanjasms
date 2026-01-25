# SanjanaCMS (Odoo Lite)

A detailed, production-ready, lightweight CMS designed to replicate the core website building experience of Odoo without the ERP bloat. Built with **Go** (Echo), **SQLite**, and **Docker**.

## 🚀 Features

-   **Zero Dependency Deployment:** Runs entirely in a single Docker container.
-   **Fast Performance:** Written in Go with SQLite.
-   **Page Builder:** Create pages with custom HTML, Meta Tags, and SEO settings.
-   **Menu Manager:** Drag-and-drop or sequence-based menu ordering.
-   **Media Library:** Upload images (up to 100MB) and easily copy URLs for use in pages.
-   **Modern Admin UI:** Clean, "Apple-esque" design for a premium feel.
-   **Robust Logging:** Detailed logs for debugging production issues.

## 🛠 Prerequisites

-   **Docker** and **Docker Compose** installed on your machine or server.

## 📦 How to Run

### 1. Clone the Repository
```bash
git clone https://github.com/Raktim94/sanjasms.git
cd sanjasms
```

### 2. Start the Server
Run the following command to build and start the container in the background:
```bash
docker compose up -d --build
```
*Note: The first build might take a minute to compile the Go application.*

### 3. Access the Application
-   **Public Website:** [http://localhost:8080](http://localhost:8080)
-   **Admin Panel:** [http://localhost:8080/admin/login](http://localhost:8080/admin/login)

### 4. Default Login Credentials
-   **Email:** `hi@RAKTIMranjit.in`
-   **Password:** `hi@raktim`

> **Security Tip:** Change these credentials immediately after logging in via the database or future profile settings.

## 📝 Usage Guide

### Creating a New Page
1.  Go to **Pages** in the sidebar.
2.  Click **"New Page"**.
3.  Enter a **Slug** (e.g., `services` for `yourdomain.com/services`).
4.  Enter a **Title** and **HTML Content**. 
    *   *Tip: Use the Media Library to upload images first, copy their URL, and paste it into `<img>` tags here.*
5.  Set **Is Published** to `Active`.
6.  Click **Save**.

### Managing Menus
1.  Go to **Menus**.
2.  Add a new menu, linking it to your page slug (e.g., `/services`).
3.  Set the sequence number to order them (1, 2, 3...).

### Uploading Images
1.  Go to **Media**.
2.  Click **Upload File**.
3.  Once uploaded, click **Copy URL**.
4.  Paste this URL into your page content.

## 🔧 Troubleshooting

### "Internal Server Error" (500)
-   **Check Logs:** Run `docker compose logs cms` to see the error detail.
-   **Database:** Ensure `data/cms.db` is writable. If you see "permission denied", run `sudo chown -R $USER:$USER data/`.
-   **Templates:** Ensure you haven't deleted core templates (`templates/admin/layout.html`, `templates/public/layout.html`).

### "Connection Refused"
-   Ensure the container is running: `docker compose ps`
-   Ensure port `8080` is not used by another application.

## 📁 Development

To run locally without Docker (requires Go installed):

```bash
# Install dependencies
go mod tidy

# Run server
go run main.go
```

---
© 2026 SanjanaCMS
