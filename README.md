# YouTube Premium Payment Management API

API for managing YouTube Premium payments, including user authentication, payment management, and WhatsApp notifications.

## Features

- User authentication with JWT
- User and payment management
- WhatsApp notifications using Twilio
- Automatic payment status updates
- RESTful API with Echo framework

## Requirements

- Go 1.19 or higher
- MongoDB
- Twilio account (for WhatsApp notifications)
- Environment variables configured (see configuration section)

## Project Structure

```
.
├── auth/           # Authentication and JWT token handling
├── config/         # Configuration management
├── db/            # MongoDB connection and operations
├── handlers/      # HTTP request handlers
├── middleware/    # HTTP middleware (auth, logging, etc.)
├── models/        # Data models and structures
├── notifications/ # WhatsApp notification service
├── repository/    # Data access layer
├── routes/        # API route definitions
├── scheduler/     # Scheduled tasks and cron jobs
├── server/        # Server configuration and setup
├── .env           # Environment variables
├── go.mod         # Go module definition
├── go.sum         # Go module checksums
├── main.go        # Application entry point
└── README.md      # Project documentation
```

## Configuration

1. Clone the repository:
```bash
git clone [repository-url]
cd youtube-premium
```

2. Install dependencies:
```bash
go mod download
```

3. Configure environment variables:
```bash
# MongoDB
export MONGODB_URI="mongodb://localhost:27017"
export MONGODB_DB="youtube_premium"

# JWT
export JWT_SECRET="your-jwt-secret"

# Twilio
export TWILIO_ACCOUNT_SID="your-account-sid"
export TWILIO_AUTH_TOKEN="your-auth-token"
export TWILIO_PHONE_NUMBER="your-twilio-number"
```

## API Documentation

### Authentication

#### Register User
```http
POST /register
Content-Type: application/json

Request Body:
{
    "username": "string",
    "email": "string",
    "password": "string",
    "name": "string",
    "cell_phone": "string",
    "date_to_pay": "string"
}

Response: 201 Created
{
    "id": "string",
    "username": "string",
    "email": "string",
    "name": "string",
    "cell_phone": "string",
    "date_to_pay": "string",
    "created_at": "string"
}
```

#### Login
```http
POST /login
Content-Type: application/json

Request Body:
{
    "email": "string",
    "password": "string"
}

Response: 200 OK
{
    "access_token": "string",
    "refresh_token": "string",
    "token_type": "Bearer",
    "expires_in": 3600
}
```

#### Refresh Token
```http
POST /refresh
Content-Type: application/json

Request Body:
{
    "refresh_token": "string"
}

Response: 200 OK
{
    "access_token": "string",
    "token_type": "Bearer",
    "expires_in": 3600
}
```

### Users

#### Get All Users
```http
GET /api/users
Authorization: Bearer {token}

Response: 200 OK
[
    {
        "id": "string",
        "username": "string",
        "email": "string",
        "name": "string",
        "cell_phone": "string",
        "date_to_pay": "string",
        "created_at": "string",
        "updated_at": "string"
    }
]
```

#### Get User by ID
```http
GET /api/users/{id}
Authorization: Bearer {token}

Response: 200 OK
{
    "id": "string",
    "username": "string",
    "email": "string",
    "name": "string",
    "cell_phone": "string",
    "date_to_pay": "string",
    "created_at": "string",
    "updated_at": "string"
}
```

#### Update User
```http
PUT /api/users/{id}
Authorization: Bearer {token}
Content-Type: application/json

Request Body:
{
    "name": "string",
    "cell_phone": "string",
    "date_to_pay": "string"
}

Response: 200 OK
{
    "id": "string",
    "username": "string",
    "email": "string",
    "name": "string",
    "cell_phone": "string",
    "date_to_pay": "string",
    "updated_at": "string"
}
```

#### Delete User
```http
DELETE /api/users/{id}
Authorization: Bearer {token}

Response: 204 No Content
```

### Payments

#### Get All Payments
```http
GET /api/payments
Authorization: Bearer {token}

Response: 200 OK
[
    {
        "id": "string",
        "user_id": "string",
        "amount": "number",
        "payment_date": "string",
        "status": "string",
        "created_at": "string"
    }
]
```

#### Get User's Payments
```http
GET /api/{userId}/payments
Authorization: Bearer {token}

Response: 200 OK
[
    {
        "id": "string",
        "user_id": "string",
        "amount": "number",
        "payment_date": "string",
        "status": "string",
        "created_at": "string"
    }
]
```

#### Create Payment
```http
POST /api/{userId}/payments
Authorization: Bearer {token}
Content-Type: application/json

Request Body:
{
    "amount": "number",
    "payment_date": "string"
}

Response: 201 Created
{
    "id": "string",
    "user_id": "string",
    "amount": "number",
    "payment_date": "string",
    "status": "string",
    "created_at": "string"
}
```

## Scheduled Tasks

- Daily payment verification at 17:00
- Payment status updates on the 13th and 25th of each month

## Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 