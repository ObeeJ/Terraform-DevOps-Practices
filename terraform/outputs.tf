# API Outputs
output "api_gateway_url" {
  description = "The URL of the API Gateway"
  value       = aws_apigatewayv2_stage.main.invoke_url
}

output "api_gateway_id" {
  description = "The ID of the API Gateway"
  value       = aws_apigatewayv2_api.main.id
}

# Database Outputs
output "database_endpoint" {
  description = "The RDS instance endpoint"
  value       = aws_db_instance.main.endpoint
  sensitive   = true
}

output "database_port" {
  description = "The RDS instance port"
  value       = aws_db_instance.main.port
}

output "database_name" {
  description = "The database name"
  value       = aws_db_instance.main.db_name
}

# Cache Outputs
output "redis_endpoint" {
  description = "The ElastiCache Redis endpoint"
  value       = aws_elasticache_replication_group.main.primary_endpoint_address
  sensitive   = true
}

output "redis_port" {
  description = "The ElastiCache Redis port"
  value       = aws_elasticache_replication_group.main.port
}

# Lambda Outputs
output "lambda_function_name" {
  description = "The name of the Lambda function"
  value       = aws_lambda_function.carbonapi.function_name
}

output "lambda_function_arn" {
  description = "The ARN of the Lambda function"
  value       = aws_lambda_function.carbonapi.arn
}

# Network Outputs
output "vpc_id" {
  description = "The ID of the VPC"
  value       = aws_vpc.main.id
}

output "vpc_cidr_block" {
  description = "The CIDR block of the VPC"
  value       = aws_vpc.main.cidr_block
}

output "public_subnet_ids" {
  description = "The IDs of the public subnets"
  value       = aws_subnet.public[*].id
}

output "private_subnet_ids" {
  description = "The IDs of the private subnets"
  value       = aws_subnet.private[*].id
}

# Security Outputs
output "database_secret_arn" {
  description = "The ARN of the database password secret"
  value       = aws_secretsmanager_secret.db_password.arn
}

output "redis_secret_arn" {
  description = "The ARN of the Redis auth token secret"
  value       = aws_secretsmanager_secret.redis_auth.arn
}

# Storage Outputs
output "s3_bucket_name" {
  description = "The name of the S3 bucket for app storage"
  value       = aws_s3_bucket.app_storage.bucket
}

# Monitoring Outputs
output "cloudwatch_dashboard_url" {
  description = "URL to the CloudWatch dashboard"
  value       = "https://${var.aws_region}.console.aws.amazon.com/cloudwatch/home?region=${var.aws_region}#dashboards:name=${aws_cloudwatch_dashboard.main.dashboard_name}"
}

output "lambda_logs_url" {
  description = "URL to Lambda function logs"
  value       = "https://${var.aws_region}.console.aws.amazon.com/cloudwatch/home?region=${var.aws_region}#logsV2:log-groups/log-group/${replace(aws_cloudwatch_log_group.lambda.name, "/", "$252F")}"
}

# Environment Info
output "environment" {
  description = "The deployment environment"
  value       = var.environment
}

output "project_name" {
  description = "The project name"
  value       = var.project_name
}

output "aws_region" {
  description = "The AWS region"
  value       = var.aws_region
}

# Quick Start Commands
output "api_test_commands" {
  description = "Commands to test the API"
  value = {
    health_check = "curl ${aws_apigatewayv2_stage.main.invoke_url}/health"
    documentation = "curl ${aws_apigatewayv2_stage.main.invoke_url}/api/v1/docs"
    calculate_shipping = "curl -X POST ${aws_apigatewayv2_stage.main.invoke_url}/api/v1/calculate -H 'Content-Type: application/json' -d '{\"activity\":\"shipping\",\"weight\":500,\"from\":\"NYC\",\"to\":\"London\",\"transport\":\"air\"}'"
    get_activities = "curl ${aws_apigatewayv2_stage.main.invoke_url}/api/v1/activities"
  }
}
