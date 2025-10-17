# Image Dock Server

A robust, production-ready image upload and management server built with Go, designed for handling image uploads, storage, and retrieval with cloud storage integration.

## 🚀 About the Project

Image Dock Server is a lightweight, high-performance HTTP server that provides secure image upload, storage, and retrieval capabilities. Built with Go for maximum performance and reliability, it offers a simple REST API for managing images with the following key features:

- **Secure Image Upload**: Handle multipart form uploads with file validation
- **Cloud Storage Integration**: Compatible with AWS S3, Cloudflare R2, and other S3-compatible storage services
- **Database Integration**: PostgreSQL support for metadata storage and tracking
- **CORS Support**: Ready for web application integration
- **RESTful API**: Simple endpoints for upload and listing operations

## 📋 Uses

This server is ideal for:

- **Web Applications**: Backend service for image uploads in web apps
- **Mobile Apps**: API endpoint for mobile application image handling
- **Content Management Systems**: Image storage backend for CMS platforms
- **Social Media Platforms**: Profile picture and media upload handling
- **E-commerce Platforms**: Product image management
- **Portfolio Websites**: Image gallery management

## 🛠️ Pre-requirements

Before setting up the Image Dock Server, ensure you have the following:

### System Requirements
- **Go 1.25.2 or higher** - [Download Go](https://golang.org/dl/)
- **PostgreSQL Database** - Version 12 or higher
- **Cloud Storage Account** - AWS S3, Cloudflare R2, or compatible S3 storage

### Cloud Storage Setup

#### Option 1: Cloudflare R2 (Recommended for free tier)
1. Create a Cloudflare account at [cloudflare.com](https://cloudflare.com)
2. Navigate to R2 storage and create a new bucket
3. Create an API token with R2 storage permissions
4. Note down your Account ID, Access Key ID, and Secret Access Key

#### Option 2: AWS S3
1. Create an AWS account at [aws.amazon.com](https://aws.amazon.com)
2. Create an S3 bucket in your preferred region
3. Create an IAM user with S3 permissions
4. Generate access keys for the IAM user

### Database Setup
1. Set up a PostgreSQL database (local or cloud-hosted)
2. Note down your database connection details:
   - Host/Endpoint
   - Port (default: 5432)
   - Database name
   - Username
   - Password

## 🔧 Environment Configuration

Create a `.env` file in the project root using the provided `.env.example` as a template:

```bash
# Copy the example file
cp .env.example .env
```

### Required Environment Variables

```env
# Server Configuration
PORT=8081

# AWS/Cloudflare R2 Configuration
AWS_REGION=auto
S3_BUCKET=your-bucket-name
S3_ENDPOINT=https://your-account-id.r2.cloudflarestorage.com
S3_ACCESS_KEY_ID=your-access-key-id
S3_SECRET_ACCESS_KEY=your-secret-access-key

# Upload Configuration
UPLOAD_DIR=uploads/

# Database Configuration
DB_USER=your-db-username
DB_PASSWORD=your-db-password
DATABASE_URL=postgres://username:password@host:port/database?sslmode=require

# Public Access Configuration
PUBLIC_URL_BASE=https://your-public-domain.r2.dev
```

### Environment Variables Explained

| Variable | Description | Example |
|----------|-------------|---------|
| `PORT` | Server port number | `8081` |
| `AWS_REGION` | AWS region or 'auto' for R2 | `auto` |
| `S3_BUCKET` | Your storage bucket name | `my-images-bucket` |
| `S3_ENDPOINT` | Storage service endpoint | `https://abc123.r2.cloudflarestorage.com` |
| `S3_ACCESS_KEY_ID` | Your access key ID | `abc123def456` |
| `S3_SECRET_ACCESS_KEY` | Your secret access key | `secret-key-here` |
| `UPLOAD_DIR` | Directory prefix for uploads | `uploads/` |
| `DB_USER` | Database username | `myuser` |
| `DB_PASSWORD` | Database password | `mypassword` |
| `DATABASE_URL` | Full database connection string | `postgres://user:pass@host:5432/db` |
| `PUBLIC_URL_BASE` | Public URL for accessing files | `https://cdn.example.com` |

## 🏗️ Build Instructions

### Step 1: Install Dependencies
```bash
# Download and install required Go modules
go mod download

# Verify installation
go mod verify
```

### Step 2: Build the Application
```bash
# Build for current platform
go build -o image-dock-server .

# Build for Linux (for Docker deployment)
GOOS=linux GOARCH=amd64 go build -o image-dock-server-linux .
```

### Step 3: Verify Build
```bash
# Check if binary was created
ls -la image-dock-server*

# Test run (will fail without proper .env setup)
./image-dock-server
```

## 🚢 Deployment Steps

### Local Deployment

#### Step 1: Set Up Environment
```bash
# 1. Configure your .env file with actual values
# 2. Ensure PostgreSQL is running and accessible
# 3. Verify your cloud storage credentials
```

#### Step 2: Start the Server
```bash
# Run in foreground
./image-dock-server

# Or run in background
nohup ./image-dock-server > server.log 2>&1 &
```

#### Step 3: Verify Deployment
```bash
# Check if server is running
curl http://localhost:8081

# Test upload endpoint
curl -X POST http://localhost:8081/upload -F "image=@test.jpg"

# Test list endpoint
curl http://localhost:8081/images
```

### Docker Deployment

#### Step 1: Create Dockerfile
```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o image-dock-server .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/image-dock-server .
COPY --from=builder /app/.env.example .env

EXPOSE 8081
CMD ["./image-dock-server"]
```

#### Step 2: Build Docker Image
```bash
# Build the image
docker build -t image-dock-server:latest .

# Verify image was created
docker images image-dock-server
```

#### Step 3: Run with Docker Compose
```yaml
# docker-compose.yml
version: '3.8'
services:
  image-server:
    image: image-dock-server:latest
    ports:
      - "8081:8081"
    environment:
      - PORT=8081
      - AWS_REGION=${AWS_REGION}
      - S3_BUCKET=${S3_BUCKET}
      - S3_ENDPOINT=${S3_ENDPOINT}
      - S3_ACCESS_KEY_ID=${S3_ACCESS_KEY_ID}
      - S3_SECRET_ACCESS_KEY=${S3_SECRET_ACCESS_KEY}
      - UPLOAD_DIR=${UPLOAD_DIR}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DATABASE_URL=${DATABASE_URL}
      - PUBLIC_URL_BASE=${PUBLIC_URL_BASE}
    restart: unless-stopped
```

```bash
# Start with Docker Compose
docker-compose up -d

# Check logs
docker-compose logs -f image-server
```

### Cloud Deployment (Heroku/AWS/Google Cloud)

#### Step 1: Prepare for Cloud Deployment
```bash
# Create a Procfile for Heroku
echo "web: image-dock-server" > Procfile

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o image-dock-server .
```

#### Step 2: Set Environment Variables in Cloud Platform

**For Heroku:**
```bash
# Set environment variables
heroku config:set PORT=8081
heroku config:set AWS_REGION=auto
heroku config:set S3_BUCKET=your-bucket-name
heroku config:set S3_ENDPOINT=https://your-account-id.r2.cloudflarestorage.com
heroku config:set S3_ACCESS_KEY_ID=your-access-key-id
heroku config:set S3_SECRET_ACCESS_KEY=your-secret-access-key
heroku config:set UPLOAD_DIR=uploads/
heroku config:set DATABASE_URL=your-database-url
heroku config:set PUBLIC_URL_BASE=https://your-public-domain.r2.dev
```

**For AWS/Google Cloud:**
Set the same environment variables in your cloud platform's environment configuration.

#### Step 3: Deploy to Cloud Platform
```bash
# For Heroku
git add .
git commit -m "Deploy to production"
git push heroku main

# For other platforms, follow their specific deployment guides
```

## 🔐 Secret Key Management

### Security Best Practices

#### 1. Environment Variables (Recommended)
Store all sensitive information in environment variables, never in code:
```bash
# Good - Environment variables
S3_SECRET_ACCESS_KEY=your-secret-key
DB_PASSWORD=your-db-password

# Bad - Hardcoded in code
secretKey := "hardcoded-secret"
```

#### 2. Use Secret Management Services

**AWS Secrets Manager:**
```bash
# Create a secret
aws secretsmanager create-secret \
  --name "image-dock-server/secrets" \
  --secret-string '{"db_password":"your-password","s3_secret":"your-s3-secret"}'

# Retrieve in your application
secret := awsClient.GetSecretValue("image-dock-server/secrets")
```

**Google Secret Manager:**
```bash
# Create secret
gcloud secrets create image-dock-secrets --data-file=secrets.json

# Access in application
secret := gcpClient.AccessSecret("image-dock-secrets")
```

#### 3. Docker Secrets (for Docker deployments)
```bash
# Create secret files
echo "your-s3-secret-key" | docker secret create s3_secret -

# Use in docker-compose.yml
secrets:
  - s3_secret
environment:
  S3_SECRET_ACCESS_KEY_FILE: /run/secrets/s3_secret
```

### Production Security Checklist

- [ ] Use strong, unique passwords for all services
- [ ] Enable SSL/TLS encryption (HTTPS)
- [ ] Set up proper firewall rules
- [ ] Use secret management services
- [ ] Rotate access keys regularly
- [ ] Monitor access logs
- [ ] Set up alerts for security events

## 📡 API Endpoints

### Upload Image
```bash
curl -X POST http://localhost:8081/upload \
  -F "image=@your-image.jpg"
```

**Response:**
```json
{
  "message": "Upload successful",
  "url": "https://your-domain.com/uploads/your-image.jpg"
}
```

### List Images
```bash
curl http://localhost:8081/images
```

**Response:**
```json
[
  "https://your-domain.com/uploads/image1.jpg",
  "https://your-domain.com/uploads/image2.png"
]
```

## 🧪 Testing

### Manual Testing
```bash
# 1. Start the server
./image-dock-server

# 2. Upload a test image
curl -X POST http://localhost:8081/upload \
  -F "image=@test-image.jpg"

# 3. List uploaded images
curl http://localhost:8081/images

# 4. Verify image accessibility
curl -I https://your-public-domain/uploads/test-image.jpg
```

### Automated Testing
```bash
# Create test script
#!/bin/bash
# Test upload
echo "Testing image upload..."
curl -X POST http://localhost:8081/upload \
  -F "image=@test.jpg" | jq '.'

# Test listing
echo "Testing image listing..."
curl http://localhost:8081/images | jq '.'
```

## 🔍 Troubleshooting

### Common Issues

**1. Server won't start**
```bash
# Check if port is in use
netstat -tulpn | grep :8081

# Check environment variables
./image-dock-server 2>&1 | head -20
```

**2. Upload fails**
```bash
# Check S3 credentials
aws s3 ls s3://your-bucket --endpoint-url=your-endpoint

# Check file permissions
ls -la uploads/
```

**3. Database connection fails**
```bash
# Test database connection
psql "your-database-url" -c "SELECT 1;"

# Check database logs
tail -f /var/log/postgresql/postgresql.log
```

### Logs and Monitoring

```bash
# View application logs
tail -f server.log

# Monitor server performance
htop

# Check system resources
df -h  # Disk usage
free -h # Memory usage
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

If you encounter any issues or need help:

1. Check the troubleshooting section above
2. Review the logs for error messages
3. Ensure all environment variables are correctly set
4. Verify your cloud storage and database configurations

For additional support, please open an issue in the repository.

---

**Happy uploading! 🚀**
