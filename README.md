# 💬 Comment API
> A modern social platform backend - Think Twitter/X, but built with Go!

[![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![JWT](https://img.shields.io/badge/JWT-000000?style=for-the-badge&logo=JSON%20web%20tokens&logoColor=white)](https://jwt.io/)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge)](LICENSE)

## 🚀 Overview

**Comment API** is a REST API service built with Go and Vue that powers a social media platform similar to Twitter/X. Users can register, create posts with file attachments, engage with content through comments, and administrators can moderate the entire ecosystem.

### ✨ Key Features

- 🔐 **Secure Authentication** - JWT-based auth with bcrypt password hashing
- 📝 **Social Posting** - Create, edit, and delete posts with file attachments
- 💬 **Interactive Comments** - Full CRUD operations on comments with nested discussions
- 👥 **Role-Based Access** - User and Admin roles with different permissions
- 📎 **File Management** - Upload and manage attachments (images, PDFs, documents)
- 📊 **Admin Dashboard** - Comprehensive moderation tools and analytics
- 🔍 **Advanced Filtering** - Search and filter by user, date, and content
- 📄 **Pagination** - Efficient data loading with page-based navigation

## 🏗️ Architecture

```
comment-api/
├── 📁 .github/                 # GitHub workflows
├── 📁 cmd/
│   └── app/
│       └── main.go             # Application entry point
├── 📁 db/                      # Database configurations
├── 📁 frontend/                # Frontend written in Vue
├── 📁 internal/
│   ├── 📁 repo/                # Repository layer (Data Access)
│   │   ├── comment/
│   │   ├── post/
│   │   ├── token/
│   │   └── user/
│   ├── 📁 rest/                # REST API handlers
│   │   ├── 📁 handler/         # HTTP request handlers
│   │   │   ├── comment/
│   │   │   ├── health/
│   │   │   ├── post/
│   │   │   ├── token/
│   │   │   └── user/
│   │   └── 📁 middleware/      # HTTP middleware
│   │       └── middleware.go
│   └── 📁 service/             # Business logic layer
│       ├── comment/
│       ├── post/
│       ├── token/
│       └── user/
├── 📁 pkg/                     # Shared packages
│   ├── 📁 db/                  # Database utilities
│   ├── 📁 errs/                # Error definitions
│   ├── 📁 logger/              # Logging utilities
│   ├── 📁 types/               # Common types
│   └── 📁 utils/               # Utility functions
├── 📄 .env                     # Environment variables
├── 📄 .env.example             # Environment template
├── 📄 .gitignore               # Git ignore rules
├── 📄 .golangci.yml            # Linter configuration
├── 📄 api.log                  # Application logs
├── 📄 coverage.out             # Test coverage report
├── 📄 coverage.out.mods        # Coverage modifications
├── 📄 docker-compose.yml       # Docker compose config
├── 📄 go.mod                   # Go module dependencies
├── 📄 go.sum                   # Go module checksums
└── 📄 README.md                # Project documentation
```

## 🛠️ Tech Stack

| Component | Technology |
|-----------|------------|
| **Language** | Go 1.24.2 |
| **Framework** | Gin
| **Database** | PostgreSQL |
| **DB Driver** | pgx with raw SQL |
| **Authentication** | JWT + bcrypt |
| **File Storage** | [Upload Care](https://uploadcare.com/) |

## 🚦 Getting Started

### Prerequisites

- Go 1.20 or higher
- PostgreSQL 12+
- Git
- Docker
- Node.js LTS (for frontend development)
- Upload Care account (for file storage)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/hc-b666/twitter-api.git
   cd twitter-api
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your upload care credentials
   ```

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Start the server**
   ```bash
   go run cmd/main.go
   ```

The server will start on `http://localhost:9999` 🎉

## 📡 API Endpoints

> **Base URL:** `/api/v1`

### 🔧 Public Endpoints
```http
GET  /api/v1/health                # Health check
POST /api/v1/create-admin          # Create admin user (setup)
```

### 🔐 Authentication (`/auth`)
```http
POST /api/v1/auth/register         # User registration
POST /api/v1/auth/login            # User login
POST /api/v1/auth/refresh          # Refresh JWT token
```

### 👨‍💼 User Management (`/user`) 🔒
```http
GET /api/v1/user/:userID           # Get user by ID
GET /api/v1/user/profile           # Get current user profile
```

### 📝 Posts Management (`/posts`) 🔒
```http
GET  /api/v1/posts                 # Get all posts
GET  /api/v1/posts/u/:userID       # Get posts by specific user
GET  /api/v1/posts/:postID         # Get specific post
POST /api/v1/posts                 # Create new post
PUT  /api/v1/posts/:postID         # Update existing post
POST /api/v1/posts/:postID         # Soft delete post
```

### 💬 Comments System (`/comments`) 🔒
```http
GET  /api/v1/comments/:postID      # Get all comments for a post
GET  /api/v1/comments/u/:userID    # Get comments by specific user
POST /api/v1/comments/:postID      # Create new comment on post
PUT  /api/v1/comments/:commentID   # Update existing comment
POST /api/v1/comments/delete/:commentID  # Soft delete comment
```

### 🛡️ Admin Panel (`/admin`) 🔒👑
```http
GET    /api/v1/admin/users         # Get all users
GET    /api/v1/admin/comments      # Get all comments
DELETE /api/v1/admin/comment/:commentID  # Hard delete comment
DELETE /api/v1/admin/post/:postID  # Hard delete post
```

## 🔑 Authentication

All protected endpoints require JWT token in the Authorization header:

```http
Authorization: Bearer <your_jwt_token>
```

### User Roles

- **👤 User**: Can create, edit, and delete own content
- **👨‍💼 Admin**: Full access to all content and moderation tools

## 📋 Request/Response Examples

### Register New User
```bash
curl -X POST http://localhost:9999/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com", 
    "password": "securepass123"
  }'
```

### Environment Variables
```env
UPLOADCARE_PUBLIC_KEY=your_public_key
UPLOADCARE_SECRET_KEY=your_secret_key
```

## 🤝 Contributing

We welcome contributions! Here's how you can help:

1. 🍴 Fork the repository
2. 🌿 Create a feature branch (`git checkout -b feature/amazing-feature`)
3. 💾 Commit your changes (`git commit -m 'Add amazing feature'`)
4. 📤 Push to the branch (`git push origin feature/amazing-feature`)
5. 🔄 Open a Pull Request

## 👥 Contributors

This project was developed by a dedicated team of 2 contributors:

- **[Muhammadbobur](https://github.com/hc-b666)** - Full-Stack Developer & Architecture
- **[Nazokat](https://github.com/NazokatSabitova)** - Backend Development & Testing

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📞 Support

Having issues? We're here to help!

- 🐛 **Bug Reports**: [Create an issue](https://github.com/hc-b666/twitter-api/issues)
- 💡 **Feature Requests**: [Submit a feature request](https://github.com/hc-b666/twitter-api/issues)
- 📧 **Contact**: bobur0218programmer@gmail.com

---

<div align="center">

**Built with ❤️ using Vue and Go**

⭐ **Star this repo if you found it helpful!** ⭐

</div>