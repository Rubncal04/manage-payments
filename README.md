# YouTube Premium Payment Management API

API for managing YouTube Premium payments, including user authentication, client management, payment processing, and automated notifications.

## Features

- User authentication with JWT
- Client management per user
- Payment processing with state machine
- Automated WhatsApp notifications for payment reminders
- Scheduled payment status updates
- RESTful API with Echo framework
- MongoDB for data persistence

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
├── notifications/ # WhatsApp notification service using Twilio
├── repository/    # Data access layer
├── routes/        # API route definitions
├── scheduler/     # Scheduled tasks for payment reminders and status updates
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
git clone [git@github.com:Rubncal04/manage-payments.git]
cd manage-payments
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
export TWILIO_FROM_WHATSAPP="your-twilio-number"
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
    "name": "string"
}

Response: 201 Created
{
    "id": "string",
    "username": "string",
    "email": "string",
    "name": "string",
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
    "expires_in": 3600,
    "user": {
        "id": "string",
        "username": "string",
        "email": "string",
        "name": "string"
    }
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

### Clients

#### Get All Clients
```http
GET /api/clients
Authorization: Bearer {token}

Response: 200 OK
[
    {
        "id": "string",
        "user_id": "string",
        "name": "string",
        "cell_phone": "string",
        "day_to_pay": "string",
        "status": "string",
        "last_payment_date": "string",
        "created_at": "string",
        "updated_at": "string"
    }
]
```

#### Get Client by ID
```http
GET /api/clients/{id}
Authorization: Bearer {token}

Response: 200 OK
{
    "id": "string",
    "user_id": "string",
    "name": "string",
    "cell_phone": "string",
    "day_to_pay": "string",
    "status": "string",
    "last_payment_date": "string",
    "created_at": "string",
    "updated_at": "string"
}
```

#### Create Client
```http
POST /api/clients
Authorization: Bearer {token}
Content-Type: application/json

Request Body:
{
    "name": "string",
    "cell_phone": "string",
    "day_to_pay": "string"
}

Response: 201 Created
{
    "id": "string",
    "user_id": "string",
    "name": "string",
    "cell_phone": "string",
    "day_to_pay": "string",
    "status": "active",
    "last_payment_date": null,
    "created_at": "string",
    "updated_at": "string"
}
```

#### Update Client
```http
PUT /api/clients/{id}
Authorization: Bearer {token}
Content-Type: application/json

Request Body:
{
    "name": "string",
    "cell_phone": "string",
    "day_to_pay": "string"
}

Response: 200 OK
{
    "id": "string",
    "user_id": "string",
    "name": "string",
    "cell_phone": "string",
    "day_to_pay": "string",
    "status": "string",
    "last_payment_date": "string",
    "updated_at": "string"
}
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
        "client_id": "string",
        "amount": "number",
        "payment_date": "string",
        "status": "string",
        "error": "string",
        "created_at": "string",
        "updated_at": "string"
    }
]
```

#### Get Client's Payments
```http
GET /api/clients/{clientId}/payments
Authorization: Bearer {token}

Response: 200 OK
[
    {
        "id": "string",
        "client_id": "string",
        "amount": "number",
        "payment_date": "string",
        "status": "string",
        "error": "string",
        "created_at": "string",
        "updated_at": "string"
    }
]
```

#### Create Payment
```http
POST /api/clients/{clientId}/payments
Authorization: Bearer {token}
Content-Type: application/json

Request Body:
{
    "amount": "number"
}

Response: 201 Created
{
    "id": "string",
    "client_id": "string",
    "amount": "number",
    "payment_date": "string",
    "status": "processing",
    "created_at": "string",
    "updated_at": "string"
}
```

## Scheduled Tasks

The system includes automated tasks for payment management:

1. **Daily Payment Verification** (16:00 every day)
   - Checks for pending payments
   - Sends WhatsApp reminders to clients who haven't paid
   - Uses Twilio for WhatsApp notifications

2. **Monthly Payment Status Updates**
   - Runs on the 13th of each month for clients with payment dates 15-20
   - Runs on the 25th of each month for clients with payment dates 28-30
   - Updates payment statuses and sends notifications

## Payment States

The payment system implements a state machine with the following states:

1. **Processing** (`processing`)
   - Initial state when a payment is created
   - Payment is being processed

2. **Completed** (`completed`)
   - Final state when payment is successful
   - Updates client's last payment date

3. **Rejected** (`rejected`)
   - Final state when payment fails
   - Includes error message explaining the failure
   - Does not update client's last payment date

## Security

- All routes except `/register`, `/login`, and `/refresh` require authentication
- Users can only access their own clients and payments
- JWT tokens are used for authentication
- Passwords are hashed before storage

## Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 