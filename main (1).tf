# Keeping the provider and IAM role as they are
provider "aws" {
  region = "us-east-1"
}

resource "aws_iam_role" "lambda_exec_role" {
  name = "lambda-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          Service = "lambda.amazonaws.com"
        },
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_lambda_function" "lambda_one" {
  function_name = "lambda-one"
  filename      = "lambda_one_payload.zip"
  handler       = "index.handler"
  runtime       = "nodejs18.x"
  role          = aws_iam_role.lambda_exec_role.arn

  source_code_hash = filebase64sha256("lambda_one_payload.zip")
}

# Adding lambda-two since Jamie needs it for downstream processing
resource "aws_lambda_function" "lambda_two" {
  function_name = "lambda-two"
  filename      = "lambda_two_payload.zip" # Using the stub ZIP Jamie mentioned
  handler       = "index.handler" # Same as lambda-one, should work with the stub
  runtime       = "nodejs18.x"    # Matching lambda-one's runtime
  role          = aws_iam_role.lambda_exec_role.arn

  source_code_hash = filebase64sha256("lambda_two_payload.zip")
}
# Note to self: Copied lambda-one’s setup and adjusted for lambda-two. Checked Terraform docs at https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function to get this right.

# Adding logging permissions for both Lambdas
resource "aws_iam_role_policy_attachment" "lambda_logging" {
  role       = aws_iam_role.lambda_exec_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}
# Note to self: Found this policy ARN in the AWS Lambda docs at https://docs.aws.amazon.com/lambda/latest/dg/lambda-intro-execution-role.html so they can log to CloudWatch.

# Setting up S3 to trigger both Lambdas when files upload
resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = "my-beta-bucket" # Need to ask Jamie for the real bucket name later

  lambda_function {
    lambda_function_arn = aws_lambda_function.lambda_one.arn
    events             = ["s3:ObjectCreated:*"] # Should catch all new uploads
  }

  lambda_function {
    lambda_function_arn = aws_lambda_function.lambda_two.arn
    events             = ["s3:ObjectCreated:*"] # Same for lambda-two
  }

  depends_on = [
    aws_lambda_permission.allow_s3_lambda_one,
    aws_lambda_permission.allow_s3_lambda_two
  ]
}
# Note to self: Followed the Terraform S3 notification guide at https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_notification to set this up.

# Allowing S3 to call lambda-one
resource "aws_lambda_permission" "allow_s3_lambda_one" {
  statement_id  = "AllowS3InvokeLambdaOne"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda_one.function_name
  principal     = "s3.amazonaws.com"
  source_arn    = "arn:aws:s3:::my-beta-bucket" # Placeholder, need real ARN
}

# Allowing S3 to call lambda-two
resource "aws_lambda_permission" "allow_s3_lambda_two" {
  statement_id  = "AllowS3InvokeLambdaTwo"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda_two.function_name
  principal     = "s3.amazonaws.com"
  source_arn    = "arn:aws:s3:::my-beta-bucket" # Placeholder, need real ARN
}
# Note to self: Got this permission setup from the AWS Lambda permissions page at https://docs.aws.amazon.com/lambda/latest/dg/access-control-resource-based.html.

# TODO: Rename 'lambda_one' to something cooler before the team sees it