# AWS Nitro Enclaves Setup Guide

## Prerequisites

1. **AWS Account** with appropriate permissions
2. **Terraform** installed ([install docs](https://developer.hashicorp.com/terraform/install))
3. **AWS CLI** configured ([install docs](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html))
4. **EC2 Key Pair** created

## Supported Instance Types

Nitro Enclaves require specific instance types:
- **m5.large** (minimum, 2 vCPU, 8 GB RAM)
- **m5.xlarge** (4 vCPU, 16 GB RAM)
- **m5.2xlarge** (8 vCPU, 32 GB RAM)
- **c5.large** (2 vCPU, 4 GB RAM)
- **c5.xlarge** (4 vCPU, 8 GB RAM)

## Setup Steps

### 1. Configure Terraform Variables

```bash
cd terraform
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your values
```

### 2. Deploy Infrastructure

```bash
terraform init
terraform plan
terraform apply
```

### 3. Build Nitro Image

```bash
# Build Docker image for Nitro
docker build -f Dockerfile.nitro -t safe-code-nitro .

# Convert to Nitro Enclave Image Format (EIF)
nitro-cli build-enclave --docker-uri safe-code-nitro --output-file safe-code.eif
```

### 4. Deploy to Nitro Enclave

```bash
# SSH to EC2 instance
ssh -i your-key.pem ec2-user@<instance-ip>

# Copy EIF file to instance
scp -i your-key.pem safe-code.eif ec2-user@<instance-ip>:~/

# Run enclave
sudo nitro-cli run-enclave \
  --cpu-count 2 \
  --memory 1024 \
  --eif-path safe-code.eif \
  --debug-mode
```

## Verification

### Check Enclave Status
```bash
nitro-cli describe-enclaves
```

### Get Attestation Document
```bash
curl http://localhost:8080/attestation
```

### Verify with Script
```bash
python3 verify.py http://<instance-ip>:8080
```

## Security Benefits

- **Hardware-level isolation**
- **Cryptographic attestation**
- **No root access to enclave**
- **Encrypted memory**
- **Secure boot process**

## Costs

- **m5.large**: ~$0.096/hour
- **Data transfer**: Standard AWS rates
- **EBS storage**: ~$0.10/GB/month

## Troubleshooting

### Enclave Won't Start
```bash
# Check logs
sudo journalctl -u nitro-enclaves-allocator
dmesg | grep nitro
```

### Memory Issues
```bash
# Allocate more memory to enclaves
echo 'memory=1024' | sudo tee /etc/nitro_enclaves/allocator.yaml
sudo systemctl restart nitro-enclaves-allocator
```

### Network Issues
```bash
# Check if vsock is working
ls -la /dev/vsock
```