# Album Battler - Terraform Infrastructure

## 構成内容

### リソース

- **S3 バケット**: 画像ストレージ用（Public Read 許可）
- **VPC**: 新規 VPC、Internet Gateway、Public Subnets（マルチ AZ）

### ディレクトリ構造

```
terraform/
├── README.md                    # このファイル
├── main.tf                      # メイン設定
├── variables.tf                 # 変数定義
├── outputs.tf                   # 出力値
├── terraform.tfvars.example     # 変数の設定例
└── modules/
    └── s3_storage/              # S3ストレージモジュール
        ├── main.tf
        ├── variables.tf
        └── outputs.tf
```

## 前提条件

### 必要なツール

- [Terraform](https://www.terraform.io/downloads) >= 1.0
- AWS CLI（設定済み）

### AWS 認証情報の設定

```bash
aws configure
```

または、環境変数で設定：

```bash
export AWS_ACCESS_KEY_ID="your-access-key"
export AWS_SECRET_ACCESS_KEY="your-secret-key"
export AWS_DEFAULT_REGION="ap-northeast-1"
```

## 初回セットアップ

### 1. 変数ファイルの作成

```bash
cd terraform
cp terraform.tfvars.example terraform.tfvars
```

`terraform.tfvars` を編集して、プロジェクトに合わせた値を設定してください。

### 2. Terraform の初期化

```bash
terraform init
```

### 3. 実行計画の確認

```bash
terraform plan
```

作成されるリソースを確認します：

- VPC、Internet Gateway、Public Subnets
- S3 バケット（album-battler-images）
- バケットポリシー（Public Read 許可）
- CORS 設定

### 4. インフラストラクチャの適用

```bash
terraform apply
```

`yes` を入力して実行します。

### 5. 出力値の確認

```bash
terraform output
```

以下の情報が表示されます：

- `s3_bucket_name`: S3 バケット名
- `s3_bucket_url`: S3 バケットの URL
- `s3_bucket_region`: リージョン
- `vpc_id`: VPC ID
- `public_subnet_ids`: パブリックサブネット ID

## 日常運用

### リソースの状態確認

```bash
terraform show
```

### 変更の適用

設定ファイルを編集後：

```bash
terraform plan   # 変更内容を確認
terraform apply  # 変更を適用
```

### リソースの削除

```bash
terraform destroy
```

## 使用方法

### Go アプリケーションからのアクセス

環境変数の設定（`backend/.env`）：

```env
AWS_REGION=ap-northeast-1
S3_BUCKET_NAME=album-battler-images
```
