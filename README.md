# Zip.ly - URL Shortener

## 📝 Description  
**Zip.ly** is a modern URL shortener that converts lengthy URLs into concise, shareable links. With a focus on clean architecture and user experience, Zip.ly offers easy link management and built-in analytics.

## ✨ Features  
- **URL Shortening:** Quickly generate short, user-friendly links for seamless sharing.  
- **Link Expiration:** Set custom expiry dates for time-sensitive campaigns or resources.  
- **Custom Aliases:** Create branded links for better recognition and engagement.  
- **Intuitive Interface:** Simple and mobile-optimized for effortless navigation.

## 📱 Application
![Home](https://github.com/rakheshkrishna2005/Zip.ly/blob/main/img/1.png)
![Features](https://github.com/rakheshkrishna2005/Zip.ly/blob/main/img/2.png)
![Contact](https://github.com/rakheshkrishna2005/Zip.ly/blob/main/img/3.png)

## 🛠️ Tech Stack

### Backend
- **Go** - Core programming language
- **Gorilla Mux** - HTTP router and URL matcher
- **PostgreSQL** - Database for storing URL mappings and analytics
- **sqlx** - Extensions to Go's database/sql package
- **Golang-migrate** - Database migration tool

### Frontend
- **HTML5/CSS3** - Structure and styling
- **JavaScript** - Client-side functionality

### Architecture
- Clean architecture with separation of concerns:
  - **Handlers**: Process HTTP requests and responses
  - **Services**: Implement business logic
  - **Repositories**: Handle data persistence
  - **Models**: Define data structures

## 📋 API Endpoints

### URL Operations
- `POST /api/v1/urls` - Create a new short URL
- `GET /api/v1/urls/:id` - Get URL details by ID
- `PUT /api/v1/urls/:id` - Update URL (custom alias)
- `DELETE /api/v1/urls/:id` - Delete URL
- `GET /:shortCode` - Redirect to original URL

### System Operations
- `GET /health` - System health check
- `GET /` - Web interface

## 🚀 Getting Started

### Prerequisites
- Go 1.16+
- PostgreSQL 12+
- Make (optional, for using Makefile commands)

### Installation

1. Clone the repository
```bash
git clone https://github.com/yourusername/zip-ly.git
cd zip-ly
```

2. Set up environment variables (copy from example)
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Initialize the database
```bash
make init-db
```

4. Run database migrations
```bash
make migrate-up
```

5. Build and run the application
```bash
make build
make run
```

6. Access the application
```
http://localhost:8080
```

## 📦 Project Structure

```
Zip.ly/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── api/                  # API layer
│   ├── config/               # Application configuration
│   ├── models/               # Data models
│   ├── repository/           # Database operations
│   ├── service/              # Business logic
│   └── utils/                # Helper utilities
├── migrations/               # Database migrations
├── web/
│   ├── templates/            # HTML templates
│   └── static/               # Static assets (CSS, JS)
├── .env.example              # Example environment variables
├── go.mod                    # Go module file
├── Makefile                  # Common commands
└── README.md                 # Project documentation
```

## 🧪 Running Tests

```bash
make test
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
