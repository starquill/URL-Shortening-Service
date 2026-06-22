#!/bin/bash
set -e

# AWS Deployment Script for URL Shortening Service
# This script helps automate the deployment to AWS ECS

echo "🚀 URL Shortening Service - AWS Deployment"
echo "==========================================="
echo ""

# Check if AWS CLI is installed
if ! command -v aws &> /dev/null; then
    echo "❌ AWS CLI not found. Please install it first:"
    echo "   brew install awscli"
    exit 1
fi

# Check if AWS credentials are configured
if ! aws sts get-caller-identity &> /dev/null; then
    echo "❌ AWS credentials not configured. Please run:"
    echo "   aws configure"
    exit 1
fi

echo "✅ AWS CLI configured"
echo ""

# Get AWS account ID and region
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
AWS_REGION=${AWS_REGION:-us-east-1}
ECR_REPOSITORY="${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/url-shortener"

echo "📋 Configuration:"
echo "   AWS Account: ${AWS_ACCOUNT_ID}"
echo "   Region: ${AWS_REGION}"
echo "   ECR Repository: ${ECR_REPOSITORY}"
echo ""

# Function to create ECR repository if it doesn't exist
create_ecr_repo() {
    echo "🔍 Checking ECR repository..."
    if ! aws ecr describe-repositories --repository-names url-shortener --region ${AWS_REGION} &> /dev/null; then
        echo "📦 Creating ECR repository..."
        aws ecr create-repository \
            --repository-name url-shortener \
            --region ${AWS_REGION} \
            --image-scanning-configuration scanOnPush=true
        echo "✅ ECR repository created"
    else
        echo "✅ ECR repository already exists"
    fi
    echo ""
}

# Function to build and push Docker image
build_and_push() {
    echo "🔨 Building Docker image..."
    docker build -t url-shortener:production .

    echo "🏷️  Tagging image..."
    docker tag url-shortener:production ${ECR_REPOSITORY}:latest
    docker tag url-shortener:production ${ECR_REPOSITORY}:$(git rev-parse --short HEAD)

    echo "🔐 Logging into ECR..."
    aws ecr get-login-password --region ${AWS_REGION} | \
        docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com

    echo "⬆️  Pushing image to ECR..."
    docker push ${ECR_REPOSITORY}:latest
    docker push ${ECR_REPOSITORY}:$(git rev-parse --short HEAD)

    echo "✅ Image pushed successfully"
    echo ""
}

# Function to run database migrations
run_migrations() {
    echo "🗃️  Database migrations"
    echo "⚠️  Manual step required:"
    echo "   1. Connect to RDS instance"
    echo "   2. Run migrations using: golang-migrate"
    echo "   3. Or exec into ECS task and run migrations"
    echo ""
}

# Main menu
echo "What would you like to do?"
echo ""
echo "1. Create ECR repository"
echo "2. Build and push Docker image"
echo "3. Full deployment (ECR + Build + Push)"
echo "4. Show deployment status"
echo "5. Exit"
echo ""
read -p "Enter your choice (1-5): " choice

case $choice in
    1)
        create_ecr_repo
        ;;
    2)
        build_and_push
        ;;
    3)
        create_ecr_repo
        build_and_push
        echo "✅ Deployment complete!"
        echo ""
        echo "Next steps:"
        echo "1. Create RDS PostgreSQL instance"
        echo "2. Create ElastiCache Redis cluster"
        echo "3. Create ECS cluster and service"
        echo "4. Configure Application Load Balancer"
        echo ""
        echo "See AWS_DEPLOYMENT.md for detailed instructions"
        ;;
    4)
        echo "📊 Checking deployment status..."
        echo ""

        echo "ECR Images:"
        aws ecr describe-images \
            --repository-name url-shortener \
            --region ${AWS_REGION} \
            --query 'imageDetails[*].[imageTags[0],imagePushedAt,imageSizeInBytes]' \
            --output table 2>/dev/null || echo "No images found"
        echo ""

        echo "ECS Services:"
        aws ecs list-services \
            --cluster url-shortener-cluster \
            --region ${AWS_REGION} \
            --output table 2>/dev/null || echo "No services found"
        echo ""

        echo "RDS Instances:"
        aws rds describe-db-instances \
            --db-instance-identifier url-shortener-db \
            --query 'DBInstances[*].[DBInstanceIdentifier,DBInstanceStatus,Endpoint.Address]' \
            --output table 2>/dev/null || echo "No RDS instance found"
        echo ""

        echo "ElastiCache Clusters:"
        aws elasticache describe-cache-clusters \
            --cache-cluster-id url-shortener-redis \
            --query 'CacheClusters[*].[CacheClusterId,CacheClusterStatus]' \
            --output table 2>/dev/null || echo "No cache cluster found"
        ;;
    5)
        echo "👋 Goodbye!"
        exit 0
        ;;
    *)
        echo "❌ Invalid choice"
        exit 1
        ;;
esac
