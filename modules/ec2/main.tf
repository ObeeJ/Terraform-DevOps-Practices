variable "instance_type" {
  description = "EC2 instance type"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID to launch the instance in"
  type        = string
}

resource "aws_subnet" "main" {
  vpc_id                  = var.vpc_id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = "us-east-1a"
  map_public_ip_on_launch = true
}

resource "aws_instance" "this" {
  ami           = "ami-0c94855ba95c71c99" # Amazon Linux 2 AMI for us-east-1
  instance_type = var.instance_type
  subnet_id     = aws_subnet.main.id
  tags = {
    Name = "prod-ec2"
  }
}

output "public_ip" {
  value = aws_instance.this.public_ip
}
