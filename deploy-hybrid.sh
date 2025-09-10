#!/bin/bash

echo "🚀 Deploying Hybrid Architecture: Parent EC2 + Nitro Enclave"

# Stop any running containers/enclaves
echo "🛑 Stopping existing services..."
docker stop $(docker ps -q) 2>/dev/null || true
sudo nitro-cli terminate-enclave --all 2>/dev/null || true

# Build and run enclave
echo "🔒 Building and starting Nitro Enclave..."
docker build -f Dockerfile.enclave -t safe-code-enclave .
nitro-cli build-enclave --docker-uri safe-code-enclave --output-file safe-code-enclave.eif
sudo nitro-cli run-enclave --cpu-count 1 --memory 512 --eif-path safe-code-enclave.eif --debug-mode &

# Wait for enclave to start
sleep 5

# Build and run parent application
echo "🌐 Building and starting Parent EC2 application..."
docker build -f Dockerfile.parent -t safe-code-parent .
docker run -d -p 8080:8080 -e TELEGRAM_BOT_TOKEN="$TELEGRAM_BOT_TOKEN" safe-code-parent

echo "✅ Hybrid deployment complete!"
echo ""
echo "🔍 Services running:"
echo "  • Nitro Enclave: Cryptographic attestation"
echo "  • Parent EC2: Telegram bot + HTTP proxy"
echo ""
echo "🧪 Test endpoints:"
echo "  curl http://localhost:8080/health"
echo "  curl http://localhost:8080/version"
echo "  curl http://localhost:8080/attestation"
echo ""
echo "🤖 Test Telegram bot:"
echo "  Send /ping, /version, /attestation to your bot"