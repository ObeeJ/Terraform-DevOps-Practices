# Production environment configuration with placeholders
environment     = "prod"
aws_region      = "us-east-1"
owner          = "YOUR_NAME_HERE"

# Network configuration
vpc_cidr = "10.1.0.0/16"

# Lambda configuration - Higher resources for production
lambda_memory_size = 256
lambda_timeout     = 60
log_level         = "INFO"

# API Gateway configuration - Higher limits for production
api_throttle_rate_limit  = 1000
api_throttle_burst_limit = 2000

# Monitoring configuration
alert_email = "your-prod-alerts@example.com"
