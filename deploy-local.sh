#!/bin/bash

set -e

echo "🚀 Safe-Code Local Deploy Script"
echo "================================="

# Check if we're in git repo
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Not in a git repository"
    exit 1
fi

# Check for uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo "⚠️  You have uncommitted changes. Commit them first:"
    git status --porcelain
    echo ""
    read -p "🤔 Commit changes automatically? (y/N): " auto_commit
    if [[ $auto_commit =~ ^[Yy]$ ]]; then
        git add .
        echo "📝 Enter commit message:"
        read -r commit_message
        git commit -m "$commit_message"
    else
        echo "❌ Please commit your changes first"
        exit 1
    fi
fi

# Get commit info
COMMIT_SHA=$(git rev-parse HEAD)
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

echo "📋 Deploy Information:"
echo "   Commit SHA: $COMMIT_SHA"
echo "   Build Time: $BUILD_TIME"
echo ""

# Push to GitHub
echo "📤 Pushing to GitHub..."
git push origin main

# Check if EC2 info is configured
if [ -z "$EC2_HOST" ] || [ -z "$EC2_KEY_PATH" ]; then
    echo "⚙️  Configure EC2 connection:"
    echo "   export EC2_HOST=52.58.172.224"
    echo "   export EC2_KEY_PATH=./safe-code-ec2-key-pair.pem"
    echo ""
    
    # Set defaults if not configured
    EC2_HOST=${EC2_HOST:-"52.58.172.224"}
    EC2_KEY_PATH=${EC2_KEY_PATH:-"./safe-code-ec2-key-pair.pem"}
fi

# Check if key file exists
if [ ! -f "$EC2_KEY_PATH" ]; then
    echo "❌ EC2 key file not found: $EC2_KEY_PATH"
    exit 1
fi

# Check if TELEGRAM_BOT_TOKEN is set
if [ -z "$TELEGRAM_BOT_TOKEN" ]; then
    echo "❌ TELEGRAM_BOT_TOKEN not set"
    echo "   export TELEGRAM_BOT_TOKEN=your_token_here"
    exit 1
fi

echo "🔗 Connecting to EC2: $EC2_HOST"

# Copy files to EC2
echo "📁 Copying files to EC2..."
scp -i "$EC2_KEY_PATH" parent-main.go Dockerfile.parent deploy-hybrid.sh verify-real.py ec2-user@$EC2_HOST:~/

# Deploy on EC2
echo "🚀 Deploying on EC2..."
ssh -i "$EC2_KEY_PATH" ec2-user@$EC2_HOST << EOF
export TELEGRAM_BOT_TOKEN="$TELEGRAM_BOT_TOKEN"
export COMMIT_SHA="$COMMIT_SHA"
export BUILD_TIME="$BUILD_TIME"

echo "🛑 Stopping existing services..."
docker stop \$(docker ps -q) 2>/dev/null || true
sudo nitro-cli terminate-enclave --all 2>/dev/null || true

echo "🔨 Building with real commit info..."
docker build -f Dockerfile.parent \\
  --build-arg VERSION=parent-v1.0 \\
  --build-arg COMMIT_SHA=\$COMMIT_SHA \\
  --build-arg BUILD_TIME=\$BUILD_TIME \\
  -t safe-code-parent .

echo "🚀 Starting application..."
docker run -d -p 8080:8080 \\
  -e TELEGRAM_BOT_TOKEN="\$TELEGRAM_BOT_TOKEN" \\
  safe-code-parent

echo ""
echo "✅ Deployment completed successfully!"
echo "📋 Details:"
echo "   Commit: \$COMMIT_SHA"
echo "   Build Time: \$BUILD_TIME"
echo "   URL: http://$EC2_HOST:8080"
echo ""
echo "🧪 Test commands:"
echo "   curl http://$EC2_HOST:8080/health"
echo "   python3 verify-real.py http://$EC2_HOST:8080"
EOF

echo ""
echo "🎉 Local deployment completed!"
echo "🔍 Test the deployment:"
echo "   python3 verify-real.py http://$EC2_HOST:8080"