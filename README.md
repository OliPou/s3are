# S3 File Upload Service

This project provides a secure and efficient service for uploading files to Amazon S3 with authentication and database integration.

The S3 File Upload Service is designed to handle file uploads to Amazon S3 buckets while incorporating user authentication and maintaining upload records in a database. It offers a robust API for file management, leveraging Go's concurrency features for optimal performance.

Key features include:
- Secure file uploads to Amazon S3
- User authentication for upload requests
- Database integration for tracking upload metadata
- Efficient handling of large file uploads
- Configurable S3 client for flexible deployment

## Repository Structure

- `auth/`: Contains authentication-related code
  - `consumer.go`: Implements authentication logic
- `internal/`: Internal packages
  - `common/`: Shared utilities
    - `json.go`: JSON handling functions
    - `validator.go`: Input validation functions
  - `database/`: Database interaction
    - `db.go`: Database connection and query execution
    - `models.go`: Database models
    - `uploadFile.sql.go`: SQL queries for file uploads
- `main.go`: Entry point of the application
- `middleware/`: HTTP middleware
  - `auth.go`: Authentication middleware
- `s3client/`: Amazon S3 client implementation
  - `main.go`: S3 client setup and operations
- `s3UploadFile/`: File upload handling
  - `handlers.go`: HTTP handlers for file uploads
  - `models.go`: Data models for file uploads
- `sql/`: SQL-related files
  - `queries/`: SQL query files
    - `uploadFile.sql`: SQL queries for file operations
  - `schema/`: Database schema
    - `0001_uploads.sql`: Initial schema for uploads table
- `sqlc.yaml`: Configuration file for sqlc code generation

## Usage Instructions

### Installation

1. Ensure you have Go 1.16 or later installed.
2. Clone the repository:
   ```
   git clone <repository-url>
   cd <repository-directory>
   ```
3. Install dependencies:
   ```
   go mod download
   ```

### Configuration

1. Set up your AWS credentials in `~/.aws/credentials` or use environment variables.
2. Configure the database connection in `internal/database/db.go`.
3. Adjust S3 bucket settings in `s3client/main.go`.

### Running the Service

1. Build the application:
   ```
   go build -o s3uploader
   ```
2. Run the service:
   ```
   ./s3uploader
   ```

### API Usage

#### Upload a File

```http
POST /upload
Authorization: Bearer <token>
Content-Type: multipart/form-data

file: <file-data>
```

Response:
```json
{
  "file_id": "123e4567-e89b-12d3-a456-426614174000",
  "file_name": "example.jpg",
  "file_size": 1024000,
  "upload_time": "2023-04-01T12:00:00Z"
}
```

### Testing

Run the test suite:

```
go test ./...
```

### Troubleshooting

- If you encounter S3 access issues, verify your AWS credentials and bucket permissions.
- For database connection problems, check the connection string and ensure the database server is running.
- Enable debug logging by setting the `DEBUG` environment variable to `true`.

## Data Flow

1. Client sends an authenticated file upload request.
2. Authentication middleware validates the user token.
3. Upload handler receives the file and metadata.
4. File is streamed to S3 using the S3 client.
5. Upload metadata is stored in the database.
6. Response with file details is sent back to the client.

```
Client -> Auth Middleware -> Upload Handler -> S3 Client -> S3 Bucket
                                            -> Database
```

## Infrastructure

The project uses the following AWS resources:

- S3 Bucket: Stores uploaded files
- RDS PostgreSQL: Stores upload metadata and user information

Note: Specific infrastructure details are not provided in the given context. For production deployment, consider using infrastructure-as-code tools like Terraform or AWS CloudFormation to manage these resources.