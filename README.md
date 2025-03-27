# Real-Time Polling Application

A real-time polling application built with Go, featuring room-based polls, user authentication, and live voting.

## Features

- User Authentication (Google OAuth & Email/Password)
- Room Creation and Management
- Invite Link Generation
- Real-time Polling
- Guest Access
- Live Results
- Time-based Scoring System

## Tech Stack

- Backend: Go
- WebSocket: gorilla/websocket
- Database: PostgreSQL
- Cache: Redis
- Authentication: Google OAuth2, JWT
- Frontend: React with TypeScript

## Project Structure

```
.
├── server/             # Go backend server
│   ├── cmd/           # Application entry points
│   ├── internal/      # Private application code
│   ├── pkg/           # Public library code
│   └── config/        # Configuration files
├── interface/         # Frontend React application
└── docker/           # Docker configuration files
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 15 or higher
- Redis 7 or higher
- Node.js 18 or higher

### Installation

1. Clone the repository
2. Set up environment variables (see `.env.example`)
3. Run database migrations
4. Start the server
5. Start the frontend application

## Development

### Backend

```bash
cd server
go mod download
go run cmd/server/main.go
```

### Frontend

```bash
cd interface
npm install
npm run dev
```

## Environment Variables

Create a `.env` file in the server directory with the following variables:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=polling_app
REDIS_URL=redis://localhost:6379
JWT_SECRET=your_jwt_secret
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
```

## License

MIT
