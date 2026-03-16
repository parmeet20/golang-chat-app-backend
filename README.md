# 🚀 Golang Chat App Backend

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.25.6-00ADD8?style=for-the-badge&logo=go">
  <img src="https://img.shields.io/badge/Chi_Router-v5-000000?style=for-the-badge&logo=go">
  <img src="https://img.shields.io/badge/MongoDB-Atlas-47A248?style=for-the-badge&logo=mongodb">
  <img src="https://img.shields.io/badge/WebSocket-Gorilla-0088CC?style=for-the-badge">
  <img src="https://img.shields.io/badge/JWT-golang--jwt-000000?style=for-the-badge">
  <img src="https://img.shields.io/badge/Docker-Supported-2496ED?style=for-the-badge&logo=docker">
</p>

<p align="center">
<strong>High-performance real-time chat backend built with Go, Chi, MongoDB & WebSockets</strong>
</p>

<p align="center">
<a href="#-features">✨ Features</a> •
<a href="#-tech-stack">🛠 Tech Stack</a> •
<a href="#-architecture">🏗 Architecture</a> •
<a href="#-getting-started">🚀 Getting Started</a> •
<a href="#-project-structure">📁 Project Structure</a> •
<a href="#-api-endpoints">📚 API</a> •
<a href="#-environment-variables">⚙️ Environment</a> •
<a href="#-deployment">🚢 Deployment</a>
</p>

---

# 🎯 Overview

A **production-ready chat backend** providing:

- REST APIs for authentication and chat history  
- WebSocket server for **real-time messaging**  
- JWT-based authentication  
- MongoDB persistence  
- Docker-ready deployment  

Designed to integrate with the **Chatly Frontend**.

---

# ✨ Features

### 🔐 Authentication
- JWT authentication
- bcrypt password hashing
- input validation
- CORS middleware

### 💬 Messaging
- WebSocket real-time messaging
- MongoDB message persistence
- paginated chat history
- scalable goroutine-based connections

### 🧑‍💻 Developer Experience
- Clean architecture
- Chi router middleware
- Docker support
- Go modules dependency management
- structured logging

---

# 🛠 Tech Stack

**Core**

- Go 1.25
- Chi Router v5
- MongoDB Atlas
- Gorilla WebSocket
- JWT (golang-jwt)

**Libraries**

- `go-chi/chi`
- `gorilla/websocket`
- `golang-jwt/jwt`
- `mongo-driver`
- `go-playground/validator`
- `golang.org/x/crypto`

---

# 🏗 Architecture
<p> API Flow</p>
<img width="1920" height="828" alt="API Flow" src="https://github.com/user-attachments/assets/b316ab65-3e2d-4b61-9bab-09bf0549a6ee" />
<p> WebSocket Flow</p>
<img width="1882" height="828" alt="WebSocket Flow" src="https://github.com/user-attachments/assets/de7846c3-ddae-4ddb-ac5c-4b72d5d1f843" />
<p> Contribution Flow</p>

<img width="1920" height="289" alt="Contribution Flow" src="https://github.com/user-attachments/assets/ef2d25cb-b560-4ffc-800c-37cf63c320a5" />
Overall Architecture
```
Client (Web / Mobile)
        │
        ▼
   REST API (Chi)
        │
 Auth / Chat Handlers
        │
 Services / Business Logic
        │
 MongoDB Database
        │
 WebSocket Hub
```

Key principles:

- Clean architecture (`cmd → internal → pkg`)
- Dependency injection
- goroutines + channels for concurrency
- context propagation

---

# 🚀 Getting Started

### Prerequisites

- Go ≥ 1.25
- MongoDB (Local or Atlas)

### Clone Repository

```bash
git clone https://github.com/parmeet20/golang-chat-app-backend.git
cd golang-chat-app-backend
```

### Setup Environment

```bash
cp .env.sample .env
```

### Run Server

```bash
go run cmd/server/main.go
```

Server runs on:

```
http://localhost:8080
```

---

# 📁 Project Structure

```
golang-chat-app-backend

├── cmd
│   └── server
│       └── main.go

├── internal
│   ├── auth
│   ├── user
│   ├── room
│   ├── message
│   ├── websocket
│   └── config

├── pkg
│   └── utils

├── .env.sample
└── go.mod
```

---

# 📚 API Endpoints

### Authentication

```
POST /auth/register
POST /auth/login
POST /auth/refresh
```

### Chat

```
GET /rooms
GET /messages/{roomId}
POST /messages
```

### WebSocket

```
ws://localhost:8080/api/v1/ws/{roomId}
```

Example message payload:

```json
{
  "type": "message",
  "payload": {
    "content": "Hello world"
  }
}
```

---

# ⚙️ Environment Variables

Example `.env`:

```env
PORT=8080
MONGO_URL=mongodb://localhost:27017
DATABASE_NAME=chatly

JWT_SECRET=super-secret-key
JWT_TOKEN_EXPIRATION_TIME=15m

CORS_ALLOWED_ORIGINS=*
```

⚠️ Never commit `.env` files.

---

# 📜 Commands

```bash
go run cmd/server/main.go     # run server
go build -o server cmd/server/main.go
go test ./...                 # run tests
```

---

# 🚢 Deployment

### Docker

Build image

```bash
docker build -t chatly-backend .
```

Run container

```bash
docker run -p 8080:8080 --env-file .env chatly-backend
```

Supports deployment to:

- AWS
- GCP
- Kubernetes
- Render / Railway / Fly.io

---

# 🤝 Contributing

1. Fork repository
2. Create branch

```
feat/new-feature
fix/bug-name
docs/update
```

3. Run tests

```bash
go test ./...
```

4. Open Pull Request

---

<p align="center">
Built with ❤️ using Go by <strong>Parmeet Singh</strong>
</p>
