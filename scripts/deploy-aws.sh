#!/bin/bash
set -e

echo "===== AWS ECS Deployment Script ====="
echo ""

# カラー設定
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 変数
AWS_REGION="ap-northeast-1"
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
TERRAFORM_DIR="terraform"

# Terraformディレクトリに移動
cd "$(dirname "$0")/.."

echo -e "${GREEN}Step 1: Terraform Apply${NC}"
echo "========================================"
cd $TERRAFORM_DIR
terraform init
terraform plan
read -p "Do you want to apply these changes? (yes/no): " APPLY_CONFIRM
if [ "$APPLY_CONFIRM" != "yes" ]; then
  echo -e "${RED}Deployment cancelled.${NC}"
  exit 1
fi
terraform apply -auto-approve

# 出力値を取得
FRONTEND_ECR_URL=$(terraform output -raw frontend_ecr_repository_url)
BACKEND_ECR_URL=$(terraform output -raw backend_ecr_repository_url)
ECS_CLUSTER=$(terraform output -raw ecs_cluster_name)
ALB_URL=$(terraform output -raw alb_url)

echo ""
echo -e "${GREEN}Step 2: Docker Build & Push${NC}"
echo "========================================"
cd ..

# ECRログイン
echo "Logging in to ECR..."
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com

# Frontend Build & Push
echo ""
echo "Building Frontend..."
cd frontend
docker build -t $FRONTEND_ECR_URL:latest -f Dockerfile.prod .
echo "Pushing Frontend to ECR..."
docker push $FRONTEND_ECR_URL:latest
cd ..

# Backend Build & Push
echo ""
echo "Building Backend..."
cd backend
docker build -t $BACKEND_ECR_URL:latest -f Dockerfile.prod .
echo "Pushing Backend to ECR..."
docker push $BACKEND_ECR_URL:latest
cd ..

echo ""
echo -e "${GREEN}Step 3: Update ECS Services${NC}"
echo "========================================"

# ECSサービスを強制的に新しいデプロイメントにする
echo "Updating Frontend service..."
aws ecs update-service --cluster $ECS_CLUSTER --service album-battler-frontend --force-new-deployment --region $AWS_REGION > /dev/null

echo "Updating Backend service..."
aws ecs update-service --cluster $ECS_CLUSTER --service album-battler-backend --force-new-deployment --region $AWS_REGION > /dev/null

echo ""
echo -e "${GREEN}Deployment Complete!${NC}"
echo "========================================"
echo "Application URL: $ALB_URL"
echo ""
echo "Note: It may take a few minutes for the services to start."
echo "Check status with: aws ecs describe-services --cluster $ECS_CLUSTER --services album-battler-frontend album-battler-backend --region $AWS_REGION"
