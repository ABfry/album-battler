# ========================================
# Terraform 管理コマンド
# ========================================

# 環境変数
AWS_REGION ?= ap-northeast-1
TF_VAR_FILE ?= terraform.tfvars

# Terraform初期化
init:
	@echo "🔧 Terraformを初期化しています..."
	terraform init
	@echo "✅ 初期化完了"

# Terraform設定検証
validate:
	@echo "🔍 Terraform設定を検証しています..."
	terraform validate
	@echo "✅ 検証完了"

# フォーマット
fmt:
	@echo "✨ Terraformファイルをフォーマットしています..."
	terraform fmt -recursive
	@echo "✅ フォーマット完了"

# 実行計画の確認
plan:
	@echo "📋 Terraform実行計画を作成しています..."
	terraform plan -var-file=$(TF_VAR_FILE)
	@echo "✅ 計画作成完了"

# インフラ構築
apply:
	@echo "🚀 Terraformを適用しています..."
	terraform apply -var-file=$(TF_VAR_FILE)
	@echo "✅ 適用完了"

# インフラ構築（自動承認）
apply-auto:
	@echo "🚀 Terraformを自動適用しています..."
	terraform apply -var-file=$(TF_VAR_FILE) -auto-approve
	@echo "✅ 適用完了"

# インフラ削除
destroy:
	@echo "🗑️  インフラを削除しています..."
	terraform destroy -var-file=$(TF_VAR_FILE)
	@echo "✅ 削除完了"

# インフラ削除（自動承認）
destroy-auto:
	@echo "🗑️  インフラを自動削除しています..."
	terraform destroy -var-file=$(TF_VAR_FILE) -auto-approve
	@echo "✅ 削除完了"

# 状態確認
show:
	@echo "📊 現在のインフラ状態を表示しています..."
	terraform show

# 出力値の表示
output:
	@echo "📤 Terraform出力値:"
	@terraform output

# ALB URLの取得
get-alb-url:
	@echo "🌐 ALB URL:"
	@terraform output -raw alb_url 2>/dev/null || echo "ALBがまだ作成されていません"

# ECRリポジトリURLの取得
get-ecr-urls:
	@echo "📦 ECR リポジトリURL:"
	@echo "Frontend: $$(terraform output -raw frontend_ecr_repository_url 2>/dev/null || echo 'N/A')"
	@echo "Backend:  $$(terraform output -raw backend_ecr_repository_url 2>/dev/null || echo 'N/A')"

# 状態のリフレッシュ
refresh:
	@echo "🔄 Terraform状態をリフレッシュしています..."
	terraform refresh -var-file=$(TF_VAR_FILE)
	@echo "✅ リフレッシュ完了"

# tfvarsファイルの作成
setup-tfvars:
	@if [ ! -f $(TF_VAR_FILE) ]; then \
		echo "📝 terraform.tfvars を作成しています..."; \
		cp terraform.tfvars.example $(TF_VAR_FILE); \
		echo "✅ terraform.tfvars を作成しました"; \
		echo "⚠️  terraform.tfvars を編集して、必要な値を設定してください"; \
		echo "   - db_password"; \
		echo "   - gemini_api_key"; \
	else \
		echo "✅ terraform.tfvars は既に存在します"; \
	fi

# クリーンアップ
clean:
	@echo "🧹 一時ファイルをクリーンアップしています..."
	rm -rf .terraform
	rm -f .terraform.lock.hcl
	rm -f terraform.tfstate.backup
	@echo "✅ クリーンアップ完了"

# ========================================
# AWS ECS デプロイコマンド
# ========================================

# 環境変数
AWS_ACCOUNT_ID ?= $(shell aws sts get-caller-identity --query Account --output text 2>/dev/null || echo "")
ECR_REGISTRY = $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
FRONTEND_ECR = $(ECR_REGISTRY)/album-battler/frontend
BACKEND_ECR = $(ECR_REGISTRY)/album-battler/backend

# ECRログイン
ecr-login:
	@echo "🔐 ECRにログインしています..."
	@aws ecr get-login-password --region $(AWS_REGION) | \
		docker login --username AWS --password-stdin $(ECR_REGISTRY)
	@echo "✅ ECRログイン完了"

# Frontendイメージのビルド
build-frontend:
	@echo "🏗️  Frontendイメージをビルドしています..."
	@cd .. && docker build --platform linux/amd64 -t $(FRONTEND_ECR):latest \
		-f frontend/Dockerfile.prod \
		./frontend
	@echo "✅ Frontendビルド完了"

# Backendイメージのビルド
build-backend:
	@echo "🏗️  Backendイメージをビルドしています..."
	@cd .. && docker build --platform linux/amd64 -t $(BACKEND_ECR):latest \
		-f backend/Dockerfile.prod \
		./backend
	@echo "✅ Backendビルド完了"

# すべてのイメージをビルド
build-all: build-frontend build-backend
	@echo "✅ すべてのイメージのビルドが完了しました"

# ビルド & プッシュ（すべて）
deploy-all: build-all ecr-login
	@echo "📤 すべてのイメージをプッシュしています..."
	@docker push $(FRONTEND_ECR):latest
	@docker push $(BACKEND_ECR):latest
	@echo "🚀 すべてのデプロイが完了しました"

# ECSサービスを強制更新
ecs-update-all:
	@echo "🔄 ECSサービスを更新しています..."
	@aws ecs update-service \
		--cluster album-battler-cluster \
		--service album-battler-frontend \
		--force-new-deployment \
		--region $(AWS_REGION) > /dev/null 2>&1 && \
	echo "  ✅ Frontendサービス更新完了" || echo "  ⚠️ Frontendサービス更新スキップ"
	@aws ecs update-service \
		--cluster album-battler-cluster \
		--service album-battler-backend \
		--force-new-deployment \
		--region $(AWS_REGION) > /dev/null 2>&1 && \
	echo "  ✅ Backendサービス更新完了" || echo "  ⚠️ Backendサービス更新スキップ"

# フルデプロイ（ビルド → プッシュ → ECS更新）
full-deploy: deploy-all ecs-update-all
	@echo "🎉 フルデプロイが完了しました！"

# セットアップ（初回のみ）
setup: setup-tfvars init
	@echo "✅ Terraformセットアップが完了しました"
	@echo ""
	@echo "次のステップ:"
	@echo "  1. terraform.tfvars を編集して必要な値を設定"
	@echo "  2. make plan でプランを確認"
	@echo "  3. make apply でインフラを構築"
	@echo "  4. make deploy-all で初回イメージをプッシュ"

# ヘルプ
help:
	@echo "========================================="
	@echo "Terraform & ECS 管理コマンド"
	@echo "========================================="
	@echo ""
	@echo "セットアップ:"
	@echo "  make setup         - 初回セットアップ（tfvars作成 + init）"
	@echo "  make setup-tfvars  - terraform.tfvars ファイルを作成"
	@echo "  make init          - Terraform初期化"
	@echo ""
	@echo "開発:"
	@echo "  make fmt           - コードフォーマット"
	@echo "  make validate      - 設定検証"
	@echo "  make plan          - 実行計画の確認"
	@echo ""
	@echo "インフラ:"
	@echo "  make apply         - インフラ構築（確認あり）"
	@echo "  make apply-auto    - インフラ構築（自動承認）"
	@echo "  make destroy       - インフラ削除（確認あり）"
	@echo ""
	@echo "ECSデプロイ:"
	@echo "  make build-all     - すべてのイメージをビルド"
	@echo "  make deploy-all    - すべてビルド&プッシュ"
	@echo "  make full-deploy   - ビルド&プッシュ&ECS更新"
	@echo ""
	@echo "情報取得:"
	@echo "  make show          - 現在の状態を表示"
	@echo "  make output        - 出力値を表示"
	@echo "  make get-alb-url   - ALB URLを取得"
	@echo "  make get-ecr-urls  - ECRリポジトリURLを取得"
	@echo ""
	@echo "メンテナンス:"
	@echo "  make refresh       - 状態のリフレッシュ"
	@echo "  make clean         - 一時ファイルの削除"
	@echo ""

.PHONY: init validate fmt plan apply apply-auto destroy destroy-auto \
		show output get-alb-url get-ecr-urls refresh setup-tfvars clean setup \
		ecr-login build-frontend build-backend build-all deploy-all \
		ecs-update-all full-deploy help
