terraform {
  backend "s3" {
    bucket         = "my-prod-terraform-state-bucket"
    key            = "prod/terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "terraform-locks"
    encrypt        = true
  }
}

provider "aws" {
  region = var.aws_region
}

module "vpc" {
  source = "../../modules/vpc"
  vpc_name = var.vpc_name
  cidr_block = var.vpc_cidr
}

module "ec2" {
  source = "../../modules/ec2"
  instance_type = var.instance_type
  vpc_id = module.vpc.vpc_id
}

output "vpc_id" {
  value = module.vpc.vpc_id
}

output "ec2_public_ip" {
  value = module.ec2.public_ip
}
