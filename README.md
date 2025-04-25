# Zolaris Backend

A Go web application for managing IoT devices and sensor data with a clean architecture approach.

## Project Structure

```
├── api/
│   └── handlers/         # HTTP request handlers
├── internal/
│   ├── config/          # Application configuration
│   ├── db/              # Database initialization
│   ├── middleware/      # HTTP middleware components
│   ├── models/          # Data models and DTOs
│   ├── repositories/    # Data access layer
│   ├── services/        # Business logic
│   ├── transport/       # HTTP response formatting
│   └── utils/           # Utility functions
├── Dockerfile           # Docker configuration
├── docker-compose.yml   # Docker Compose configuration
├── go.mod              # Go module definition
└── main.go             # Application entry point
```

## Architecture

This application follows a clean architecture approach with the following layers:

1. **Handlers** (Presentation Layer):
   - Handle HTTP requests and responses
   - Parse and validate input
   - Call appropriate service methods
   - Format responses

2. **Services** (Business Logic Layer):
   - Implement business rules and workflows
   - Coordinate between multiple repositories
   - Handle higher-level operations

3. **Repositories** (Data Access Layer):
   - Interface with DynamoDB and other AWS services
   - Handle CRUD operations
   - Deal with low-level data concerns

4. **Models** (Domain Layer):
   - Define domain entities and data structures
   - Contain validation rules

## Prerequisites

- Go 1.24+
- AWS Credentials (for DynamoDB access)
- DynamoDB tables:
  - `devices` - for device information
  - `machine_data` - for sensor data
  - `users` - for user information

## Configuration

The application uses environment variables for configuration:

| Environment Variable | Description | Default |
|----------------------|-------------|--------|
| `PORT` | The port on which the server listens | 8080 |
| `DEVICE_TABLE_NAME` | DynamoDB table for devices | devices |
| `DATA_TABLE_NAME` | DynamoDB table for sensor data | machine_data |
| `USER_TABLE_NAME` | DynamoDB table for users | users |
| `AWS_REGION` | AWS region | us-east-1 |
| `IOT_POLICY_NAME` | Name of the IoT policy | DefaultIoTPolicy |
| `AWS_ACCESS_KEY_ID` | AWS access key | - |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key | - |

## Running the Application

### Using Go directly

```bash
# Build and run the application
go build -o zolaris-backend
./zolaris-backend
```

### Using Docker

```bash
# Build and run the Docker container
docker-compose up --build
```

## API Endpoints

### Add Device

```
POST /device/add
```

Request Body:

```json
{
  "deviceId": "myDevice123",
  "deviceName": "Living Room Sensor"
}
```

### Attach IoT Policy

```
POST /device/attach-policy
```

Request Body:

```json
{
  "identityId": "us-east-1:12345678-1234-1234-1234-123456789012"
}
```

### Get Device Sensor Data

```
POST /device/sensor-data
```

Request Body:

```json
{
  "deviceMacId": "00:11:22:33:44:55",
  "timestamp": "1684160445500",
  "dateMode": "daily"
}
```

### List User Devices

```
GET /user/devices
```

Header:

```
X-User-ID: user123
```

### Health Check

```
GET /health
```

## Authentication

Authentication is handled using the `X-User-ID` header. In a production environment, this should be replaced with proper JWT or OAuth2 authentication.

## Error Handling

Error responses follow this format:

```json
{
  "status": false,
  "message": "Error description"
}
```

## Successful Responses

Successful responses follow this format:

```json
{
  "status": true,
  "data": [response data]
}
```

or

```json
{
  "status": true,
  "message": "Success message"
}
```

## Development

### Adding New Endpoints

To add a new endpoint:

1. Create a new handler in the `api/handlers` directory
2. Add models in the `internal/models` directory if needed
3. Add repository methods in the `internal/repositories` directory
4. Add business logic in the `internal/services` directory
5. Register the endpoint in `main.go`

### Testing

Run tests with:

```bash
go test ./...
```

