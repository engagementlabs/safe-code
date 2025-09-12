#!/bin/bash

echo "🚀 Safe-Code Setup Script"
echo "========================="

# Check if running from correct directory
if [[ ! -f "go.mod" ]]; then
    echo "❌ Please run this script from the safe-code root directory"
    exit 1
fi

echo "📋 This script will install:"
echo "   • Terraform"
echo "   • AWS CLI"
echo ""

read -p "🤔 Continue with installation? (y/N): " confirm
if [[ ! $confirm =~ ^[Yy]$ ]]; then
    echo "❌ Installation cancelled"
    exit 0
fi

echo ""
echo "🔧 Installing dependencies..."

# Install Terraform
echo "1️⃣  Installing Terraform..."
./bin/install-terraform.sh

echo ""

# Install AWS CLI
echo "2️⃣  Installing AWS CLI..."
./bin/install-aws-cli.sh

echo ""
echo "✅ Setup completed!"
echo ""
echo "🔧 Next steps:"
echo "   1. Configure AWS: aws configure"
echo "   2. Copy terraform vars: cp terraform/terraform.tfvars.example terraform/terraform.tfvars"
echo "   3. Edit terraform/terraform.tfvars with your values"
echo "   4. Deploy: cd terraform && terraform init && terraform apply"