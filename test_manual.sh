#!/bin/bash
set -e

echo "🚀 Starting Manual CURL Test"
echo "-----------------------------------"

echo "1️⃣  Registering Alice..."
ALICE_RESP=$(curl -s -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username": "alice_curl", "email": "alice_curl@example.com", "password": "password"}')
echo "Response: $ALICE_RESP"

# Extract ID using python (simple fallback if jq is missing)
ALICE_ID=$(echo $ALICE_RESP | python3 -c "import sys, json; print(json.load(sys.stdin)['id'])")
echo "✅ Alice ID: $ALICE_ID"

echo "-----------------------------------"

echo "2️⃣  Registering Bob..."
BOB_RESP=$(curl -s -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username": "bob_curl", "email": "bob_curl@example.com", "password": "password"}')
echo "Response: $BOB_RESP"

BOB_ID=$(echo $BOB_RESP | python3 -c "import sys, json; print(json.load(sys.stdin)['id'])")
echo "✅ Bob ID: $BOB_ID"

echo "-----------------------------------"

echo "3️⃣  Depositing 1000 to Alice (ID: $ALICE_ID)..."
curl -s -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d "{\"to_user_id\": $ALICE_ID, \"amount\": 1000, \"type\": \"deposit\"}"
echo -e "\n✅ Request Sent. Waiting 2 seconds for worker..."
sleep 2

echo "-----------------------------------"

echo "4️⃣  Checking Alice Balance..."
curl -s "http://localhost:8080/api/balance?user_id=$ALICE_ID"
echo -e "\n✅ Check Done."

echo "-----------------------------------"

echo "5️⃣  Transferring 300 from Alice ($ALICE_ID) to Bob ($BOB_ID)..."
curl -s -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d "{\"from_user_id\": $ALICE_ID, \"to_user_id\": $BOB_ID, \"amount\": 300, \"type\": \"transfer\"}"
echo -e "\n✅ Request Sent. Waiting 2 seconds for worker..."
sleep 2

echo "-----------------------------------"

echo "6️⃣  Checking Balances..."
echo "Alice:"
curl -s "http://localhost:8080/api/balance?user_id=$ALICE_ID"
echo -e "\nBob:"
curl -s "http://localhost:8080/api/balance?user_id=$BOB_ID"
echo ""

echo "-----------------------------------"
echo "🎉 Manual Test Complete"
