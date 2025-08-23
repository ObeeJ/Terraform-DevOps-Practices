#!/bin/bash

# CarbonAPI Deployment Script
# This script builds and deploys the CarbonAPI to AWS

set -e

echo "🌱 Starting CarbonAPI deployment..."

# Check prerequisites
command -v go >/dev/null 2>&1 || { echo "❌ Go is required but not installed. Aborting." >&2; exit 1; }
command -v terraform >/dev/null 2>&1 || { echo "❌ Terraform is required but not installed. Aborting." >&2; exit 1; }
command -v aws >/dev/null 2>&1 || { echo "❌ AWS CLI is required but not installed. Aborting." >&2; exit 1; }

# Set environment (default to production)
ENVIRONMENT=${1:-production}
echo "📦 Deploying to environment: $ENVIRONMENT"

# Build the Lambda function
echo "🔨 Building Lambda function..."
cd api
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o main main.go database.go carbon_service.go
zip carbonapi.zip main
echo "✅ Lambda function built successfully"
cd ..

# Deploy infrastructure with Terraform
echo "🏗️ Deploying infrastructure with Terraform..."
cd terraform

# Initialize Terraform
terraform init

# Plan the deployment
echo "📋 Planning infrastructure changes..."
terraform plan -var="environment=$ENVIRONMENT" -out=tfplan

# Apply the changes
echo "🚀 Applying infrastructure changes..."
terraform apply tfplan

# Get the API URL
API_URL=$(terraform output -raw api_gateway_url)
echo "✅ Infrastructure deployed successfully!"

# Wait for Lambda to be ready
echo "⏳ Waiting for Lambda function to be ready..."
sleep 30

# Test the deployment
echo "🧪 Testing the deployment..."
if curl -f "$API_URL/health" >/dev/null 2>&1; then
    echo "✅ Health check passed!"
else
    echo "❌ Health check failed!"
    exit 1
fi

# Test the carbon calculation endpoint
echo "🧮 Testing carbon calculation..."
RESPONSE=$(curl -s -X POST "$API_URL/api/v1/calculate" \
    -H "Content-Type: application/json" \
    -d '{"activity":"shipping","weight":500,"from":"NYC","to":"London","transport":"air"}')

if echo "$RESPONSE" | grep -q "carbon_footprint"; then
    echo "✅ Carbon calculation test passed!"
else
    echo "❌ Carbon calculation test failed!"
    echo "Response: $RESPONSE"
    exit 1
fi

cd ..

echo ""
echo "🎉 CarbonAPI deployment completed successfully!"
echo "📊 API URL: $API_URL"
echo "🔗 Health Check: $API_URL/health"
echo "📚 Documentation: $API_URL/api/v1/docs"
echo ""
echo "🚀 Test commands:"
echo "curl $API_URL/health"
echo "curl $API_URL/api/v1/docs"
echo "curl -X POST $API_URL/api/v1/calculate -H 'Content-Type: application/json' -d '{\"activity\":\"shipping\",\"weight\":500,\"from\":\"NYC\",\"to\":\"London\",\"transport\":\"air\"}'"
echo ""
echo "💚 CarbonAPI is now live and ready to calculate carbon footprints!"
