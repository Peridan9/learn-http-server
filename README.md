
# Chirpy 🐦  
A Lightweight Twitter-Style API Server Built with Go

## Overview 📖  
**Chirpy** is a production-style HTTP server developed in Go that simulates a simple social media backend, inspired by Twitter. It supports user registration, posting short messages ("chirps"), and secure API key management. This project was created as part of a guided course to learn web server fundamentals in Go, including routing, data storage, and secure authentication.

---

## ✨ Features
- Create and retrieve short messages ("chirps")
- User registration and login with hashed password storage
- Role-based user upgrade with API key authentication
- PostgreSQL integration with type-safe SQL via `sqlc`
- Secure token-based authentication and basic authorization
- Webhook endpoint for external service integrations
- Fully RESTful API with JSON responses

---

## 🛠️ Technologies Used
- **Go** – Main backend language
- **PostgreSQL** – Relational database for persistent storage
- **sqlc** – Generates type-safe Go code from SQL queries
- **Goose** – Manages database migrations
- **bcrypt** – Secure password hashing
- **net/http** – Standard Go HTTP server package
- **JSON, Status Codes & Headers** – For RESTful communication

---

## Prerequisites 🛠️  
Make sure the following are installed:

- Go (version 1.18 or higher recommended)  
- PostgreSQL database  
- `goose` (for running database migrations)  
- `sqlc` (for generating Go code from SQL)

---

## Installation & Setup 🚀

1. **Clone the project**
   ```bash
   git clone https://github.com/Peridan9/learn-http-server.git
   cd learn-http-server
   ```

2. **Set up your environment variables**
   Create a `.env` or export the following in your shell:
   ```
   DB_URL=postgres://user:password@localhost:5432/chirpy?sslmode=disable
   ```

3. **Run database migrations**
   ```bash
   goose -dir migrations postgres "$DB_URL" up
   ```

4. **Generate type-safe DB code with sqlc**
   ```bash
   sqlc generate
   ```

5. **Run the server**
   ```bash
   go run main.go
   ```

---

## API Overview 📜

### 🔐 Authentication
- **Register:** `POST /api/users`
- **Login:** `POST /api/login`
- **Upgrade User:** `PUT /api/users/upgrade`

### 🐦 Chirps
- **Create Chirp:** `POST /api/chirps`
- **List Chirps:** `GET /api/chirps?sort=asc|desc&author_id=`
- **Get Chirp by ID:** `GET /api/chirps/{id}`

### 🧪 Webhooks
- **Validate Webhook:** `POST /api/webhooks`  
  (Only accepts calls from trusted external sources)

---

## Development Notes 💡
This project avoids frameworks to help build a deeper understanding of:
- Routing using `http.ServeMux`
- Middleware logic
- JSON serialization
- Working with raw SQL
- Server-side validation and error handling

---

## Deployment 📦  
To compile a standalone binary:
```bash
go build -o chirpy
./chirpy
```
