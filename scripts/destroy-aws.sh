#!/bin/bash
set -e

echo "===== AWS ECS Teardown Script ====="
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
PROJECT_NAME="album-battler"

# Terraformディレクトリに移動
cd "$(dirname "$0")/.."

echo -e "${YELLOW}WARNING: This will destroy all AWS resources!${NC}"
echo "This includes:"
echo "  - ECS Cluster and Services"
echo "  - RDS Database (all data will be lost!)"
echo "  - Application Load Balancer"
echo "  - CloudWatch Logs"
echo ""
read -p "Are you sure you want to continue? (yes/no): " DESTROY_CONFIRM
if [ "$DESTROY_CONFIRM" != "yes" ]; then
  echo -e "${GREEN}Teardown cancelled.${NC}"
  exit 0
fi

echo ""
echo -e "${YELLOW}Step 1: Terraform Destroy${NC}"
echo "========================================"
echo "Destroying ECS Services, RDS, ALB, and other resources..."
cd $TERRAFORM_DIR
terraform destroy -auto-approve

echo ""
echo -e "${YELLOW}Step 2: Delete ECR Images${NC}"
echo "========================================"
echo "ECS services are now stopped. Safe to delete ECR images."
cd ..

# ECR リポジトリの画像を削除
echo "Deleting Frontend ECR images..."
aws ecr batch-delete-image \
  --repository-name $PROJECT_NAME/frontend \
  --image-ids "$(aws ecr list-images --repository-name $PROJECT_NAME/frontend --query 'imageIds[*]' --output json)" \
  --region $AWS_REGION 2>/dev/null || echo "No images to delete or repository doesn't exist"

echo "Deleting Backend ECR images..."
aws ecr batch-delete-image \
  --repository-name $PROJECT_NAME/backend \
  --image-ids "$(aws ecr list-images --repository-name $PROJECT_NAME/backend --query 'imageIds[*]' --output json)" \
  --region $AWS_REGION 2>/dev/null || echo "No images to delete or repository doesn't exist"

echo ""
echo -e "${GREEN}Teardown Complete!${NC}"
echo "All AWS resources have been destroyed."
echo ""
echo "Note: S3 bucket (album-battler-images) is NOT destroyed by this script."
echo "To delete it manually, run:"
echo "  aws s3 rb s3://album-battler-images --force"
