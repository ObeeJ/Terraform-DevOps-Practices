# Terraform DevOps Practices

This repository demonstrates best practices for using Terraform to manage AWS infrastructure with a focus on DevOps principles such as modularity, remote state management, and automation.

## Features

- **Modular Terraform Code**: Infrastructure is split into reusable modules (e.g., VPC, EC2).
- **Remote State Management**: Uses AWS S3 for storing the Terraform state file and DynamoDB for state locking to enable team collaboration.
- **Production Environment Example**: Reference implementation under `environments/prod/` showing how to deploy a VPC and EC2 instance.
- **GitHub Actions Integration**: Ready for CI/CD pipelines and automated workflow integrations (see `.github/workflows/terraform.yml`).

## Repository Structure

```
.  
├── environments/  
│   └── prod/  
│       ├── main.tf  
│       ├── variables.tf  
│       └── outputs.tf  
├── modules/  
│   ├── vpc/  
│   └── ec2/  
└── .github/  
    └── workflows/  
        └── terraform.yml  
```

- **environments/prod/**: Example production environment using the shared modules.
- **modules/**: Contains reusable infrastructure modules (vpc, ec2).
- **.github/workflows/**: Contains GitHub Actions workflow for Terraform automation.

## Getting Started

### Prerequisites

- Terraform v1.0 or newer
- AWS account and credentials
- AWS S3 bucket and DynamoDB table for remote state

### Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/ObeeJ/Terraform-DevOps-Practices.git  
   cd Terraform-DevOps-Practices/environments/prod  
   ```

2. **Initialize Terraform**:
   ```bash
   terraform init  
   ```

3. **Review and customize variables (optional)**:
   Edit `variables.tf` or provide overrides via CLI or a `.tfvars` file.

4. **Plan and apply**:
   ```bash
   terraform plan  
   terraform apply  
   ```

## Variables

See `environments/prod/variables.tf` for configurable options such as:

- `aws_region`: AWS region to deploy resources
- `vpc_name`: Name of the VPC
- `vpc_cidr`: CIDR block for the VPC
- `instance_type`: EC2 instance type

## Outputs

After applying, Terraform will output:

- `vpc_id`: The ID of the created VPC
- `ec2_public_ip`: The public IP address of the deployed EC2 instance

## Best Practices Demonstrated

- **Infrastructure as Code**: All AWS resources are managed via code.
- **State Locking**: Prevents concurrent modifications.
- **Modularization**: Encourages reuse and cleaner code.
- **CI/CD Ready**: Example GitHub Actions workflow for automation.

## Contributing

Contributions are welcome! Please open issues or submit pull requests for improvements. 
