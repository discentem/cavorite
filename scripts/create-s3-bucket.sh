#!/bin/bash
# Create or verify an S3 bucket on rustfs for cavorite testing, using AWS CLI

set -e

BUCKET_NAME="cavorite-test"
S3_ENDPOINT="http://localhost:9000"
S3_REGION="us-east-1"
ACCESS_KEY="${AWS_ACCESS_KEY_ID:-blah}"
SECRET_KEY="${AWS_SECRET_ACCESS_KEY:-blah}"

while [[ $# -gt 0 ]]; do
    case $1 in
        --bucket)
            BUCKET_NAME="$2"
            shift 2
            ;;
        --endpoint)
            S3_ENDPOINT="$2"
            shift 2
            ;;
        --region)
            S3_REGION="$2"
            shift 2
            ;;
        --access-key)
            ACCESS_KEY="$2"
            shift 2
            ;;
        --secret-key)
            SECRET_KEY="$2"
            shift 2
            ;;
        --help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Creates an S3 bucket for cavorite testing using AWS CLI"
            echo ""
            echo "Options:"
            echo "  --bucket NAME          S3 bucket name (default: cavorite-test)"
            echo "  --endpoint URL         S3 endpoint URL (default: http://localhost:9000)"
            echo "  --region REGION        AWS region (default: us-east-1)"
            echo "  --access-key KEY       S3 access key (default: AWS_ACCESS_KEY_ID env var or 'blah')"
            echo "  --secret-key KEY       S3 secret key (default: AWS_SECRET_ACCESS_KEY env var or 'blah')"
            echo "  --help                 Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

if ! command -v aws &> /dev/null; then
    echo "❌ Error: AWS CLI not found. Please install it using scripts/install-awscli.sh"
    exit 1
fi

echo "Checking if bucket '$BUCKET_NAME' exists..."
if AWS_ACCESS_KEY_ID="$ACCESS_KEY" \
   AWS_SECRET_ACCESS_KEY="$SECRET_KEY" \
   aws s3 ls "s3://$BUCKET_NAME" \
   --endpoint-url "$S3_ENDPOINT" \
   --region "$S3_REGION" &>/dev/null 2>&1; then
    echo "✓ Bucket '$BUCKET_NAME' already exists"
    exit 0
fi

echo "Creating S3 bucket: $BUCKET_NAME"
AWS_ACCESS_KEY_ID="$ACCESS_KEY" \
AWS_SECRET_ACCESS_KEY="$SECRET_KEY" \
aws s3 mb "s3://$BUCKET_NAME" \
    --endpoint-url "$S3_ENDPOINT" \
    --region "$S3_REGION" || {
    echo "❌ Failed to create bucket"
    exit 1
}

echo "✓ S3 bucket '$BUCKET_NAME' initialized"
