#!/bin/bash
set -e

TABLE_NAME="entityTable"

echo "Creating DynamoDB table: $TABLE_NAME"
if aws dynamodb create-table \
  --table-name $TABLE_NAME \
  --attribute-definitions \
    AttributeName=id,AttributeType=S \
  --key-schema \
    AttributeName=id,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST; then
  echo "Table creation initiated. Waiting for table to become active..."
else
  echo "Error creating table. Checking if table already exists..."
  aws dynamodb describe-table --table-name $TABLE_NAME > /dev/null || { echo "Failed to create or find table"; exit 1; }
fi

# Wait for the table to be created
aws dynamodb wait table-exists --table-name $TABLE_NAME
echo "Table is now active."

# Function for putting items with error handling
put_item() {
  local id=$1
  local name=$2
  local refId=$3
  local path=$4
  
  echo "Inserting item: $name ($id)"
  
  local refid_attr
  if [ "$refId" == "null" ]; then
    refid_attr='"refId": {"NULL": true}'
  else
    refid_attr="\"refId\": {\"S\": \"$refId\"}"
  fi
  
  if aws dynamodb put-item \
    --table-name $TABLE_NAME \
    --item "{
      \"id\": {\"S\": \"$id\"},
      \"name\": {\"S\": \"$name\"},
      $refid_attr,
      \"type\": {\"S\": \"user\"},
      \"path\": {\"S\": \"$path\"}
    }"; then
    echo "Successfully inserted $name"
  else
    echo "Failed to insert $name"
    return 1
  fi
}

# Insert all items
echo "Inserting hierarchy items..."

# Level 1
# put_item "user001" "CEO" null "/user001" || exit 1

# Level 2
put_item "user010" "VP Sales" "user001" "/user001/user010" || exit 1
put_item "user020" "VP Tech" "user001" "/user001/user020" || exit 1

# Level 3
put_item "user101" "Director NA" "user010" "/user001/user010/user101" || exit 1
put_item "user102" "Director EMEA" "user010" "/user001/user010/user102" || exit 1
put_item "user103" "CTO" "user020" "/user001/user020/user103" || exit 1
put_item "user104" "VP Product" "user020" "/user001/user020/user104" || exit 1

# Level 4
put_item "user201" "Manager East" "user101" "/user001/user010/user101/user201" || exit 1
put_item "user202" "Manager West" "user101" "/user001/user010/user101/user202" || exit 1
put_item "user203" "Architect" "user103" "/user001/user020/user103/user203" || exit 1
put_item "user204" "Security Lead" "user103" "/user001/user020/user103/user204" || exit 1

# Level 5
put_item "user301" "Team Lead NY" "user201" "/user001/user010/user101/user201/user301" || exit 1
put_item "user302" "Team Lead DC" "user201" "/user001/user010/user101/user201/user302" || exit 1

# Level 6
put_item "user401" "Sales Rep 1" "user301" "/user001/user010/user101/user201/user301/user401" || exit 1
put_item "user402" "Sales Rep 2" "user301" "/user001/user010/user101/user201/user301/user402" || exit 1

echo "All items inserted successfully."
echo "Setup complete!"
