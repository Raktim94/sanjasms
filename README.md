# SanjanaCMS (Odoo Lite)

A detailed, production-ready, lightweight CMS designed to replicate the core website building experience of Odoo without the ERP bloat. Built with **Go** (Echo), **SQLite**, and **Docker**.

## 🚀 Features

### 🔐 Login & Access
-   Create an admin account using email and password.
-   First account becomes the main administrator.
-   Turn off public sign-ups for private access.
-   Clean Admin Login button on the website.
-   Central dashboard to manage everything after login.

### 🧭 Admin Dashboard
Control your entire website from one place:
-   **Pages:** Create unlimited normal or landing pages.
-   **Menus:** Manage main and footer menus with drag-and-drop ordering.
-   **Media:** Upload images (PNG, JPG, WebP) and PDFs (up to 100MB).
-   **SEO:** Edit meta titles, descriptions, and Open Graph tags.
-   **Settings:** Customize logo, favicon, title, and brand colors.
-   **Users:** Manage additional admin users.

### 📄 Page Management
-   Create unlimited pages.
-   **Page Builder:** Build pages using simple blocks (Text, Image, Button, Custom HTML).
-   **Embeds:** Easily embed YouTube, Maps, or Forms using the HTML block.
-   **Fast Loading:** All pages are server-side rendered for instant speed.

### 🧭 Menu Builder
-   Create multiple menus (Main, Footer, etc.).
-   Drag items to reorder.
-   Create nested dropdowns.
-   Enable/Disable items without deleting.

### 🎨 Branding
-   Upload website logo and favicon.
-   Set primary brand color (affects buttons and links).
-   Edit footer text.
-   Updates appear instantly—no rebuild required.

### 🔍 SEO Manager
-   **Meta Tags:** Edit Title, Description, Keywords for every page.
-   **Social Sharing:** Control Open Graph image and title.
-   **Indexing:** Control search engine visibility (index/noindex).
-   **Automated:** Sitemap (`/sitemap.xml`) and Robots rules generated automatically.

### ⚡ Technical
-   **Zero Dependency:** Runs entirely in a single Docker container.
-   **Port:** Default `8080` (Configurable).

---

## 🔒 Security
-   Admin area is protected and private.
-   File uploads are safe and controlled.
-   Login sessions are secure.

## 🧩 What This CMS Is NOT
-   No ERP
-   No Billing / Inventory
-   No Ads
-   No AI content generation

---

## 📦 How to Run

### 1. Clone & Start
```bash
git clone https://github.com/Raktim94/sanjasms.git
cd sanjasms
docker compose up -d --build
```

### 2. Access
-   **Website:** `http://localhost:2512`
-   **Admin:** `http://localhost:2512/admin/login`

### 3. Default Credentials
-   **Email:** `hi@RAKTIMranjit.in`
-   **Password:** `hi@raktim`

> **Note:** If you see any errors after updating, run `docker compose down --volumes` to clear old cache.

---
© 2026 SanjanaCMS
