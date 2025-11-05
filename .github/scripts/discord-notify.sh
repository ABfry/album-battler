#!/bin/bash

WEBHOOK_URL="$1"
TITLE="$2"
DESCRIPTION="$3"
COLOR="$4"
FIELD1_NAME="$5"
FIELD1_VALUE="$6"
FIELD2_NAME="$7"
FIELD2_VALUE="$8"

TIMESTAMP=$(date -u +%Y-%m-%dT%H:%M:%S.000Z)

curl -H "Content-Type: application/json" \
  -X POST \
  -d "{
    \"embeds\": [{
      \"title\": \"$TITLE\",
      \"description\": \"$DESCRIPTION\",
      \"color\": $COLOR,
      \"fields\": [
        {
          \"name\": \"$FIELD1_NAME\",
          \"value\": \"$FIELD1_VALUE\",
          \"inline\": true
        },
        {
          \"name\": \"$FIELD2_NAME\",
          \"value\": \"$FIELD2_VALUE\",
          \"inline\": true
        }
      ],
      \"footer\": {
        \"text\": \"GitHub Actions\"
      },
      \"timestamp\": \"$TIMESTAMP\"
    }]
  }" \
  "$WEBHOOK_URL"