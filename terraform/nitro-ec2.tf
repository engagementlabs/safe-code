# AWS Nitro Enclaves EC2 Configuration

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

# Variables
variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "instance_type" {
  description = "EC2 instance type (must support Nitro Enclaves)"
  type        = string
  default     = "m6g.large"  # Graviton minimum for Nitro Enclaves
}

variable "key_name" {
  description = "EC2 Key Pair name"
  type        = string
}

variable "telegram_bot_token" {
  description = "Telegram Bot Token"
  type        = string
  sensitive   = true
}

# Data sources
data "aws_ami" "amazon_linux" {
  most_recent = true
  owners      = ["amazon"]
  
  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-*-arm64-gp2"]
  }
  
  filter {
    name   = "architecture"
    values = ["arm64"]
  }
}

# Security Group
resource "aws_security_group" "nitro_sg" {
  name_prefix = "safe-code-nitro-"
  description = "Security group for Safe-Code Nitro Enclave"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "safe-code-nitro-sg"
  }
}

# IAM Role for EC2
resource "aws_iam_role" "nitro_role" {
  name = "safe-code-nitro-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      }
    ]
  })
}

# Instance Profile
resource "aws_iam_instance_profile" "nitro_profile" {
  name = "safe-code-nitro-profile"
  role = aws_iam_role.nitro_role.name
}

# EC2 Instance with Nitro Enclaves
resource "aws_instance" "nitro_instance" {
  ami                    = data.aws_ami.amazon_linux.id
  instance_type          = var.instance_type
  key_name              = var.key_name
  vpc_security_group_ids = [aws_security_group.nitro_sg.id]
  iam_instance_profile   = aws_iam_instance_profile.nitro_profile.name
  
  # Enable Nitro Enclaves
  enclave_options {
    enabled = true
  }

  user_data = base64encode(<<-EOF
    #!/bin/bash
    yum update -y
    yum install -y docker aws-nitro-enclaves-cli
    systemctl start docker
    systemctl enable docker
    systemctl start nitro-enclaves-allocator
    systemctl enable nitro-enclaves-allocator
    
    # Set environment variable
    echo "TELEGRAM_BOT_TOKEN=${var.telegram_bot_token}" >> /etc/environment
  EOF
  )

  tags = {
    Name = "safe-code-nitro-enclave"
  }
}

# Outputs
output "instance_ip" {
  value = aws_instance.nitro_instance.public_ip
}

output "instance_dns" {
  value = aws_instance.nitro_instance.public_dns
}