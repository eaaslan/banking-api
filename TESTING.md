# Testing Guide

This guide explains how to manually test the API endpoints and verify the background worker system.

## 1. Manual Testing with CURL

You can use `curl` to interact with the API directly from your terminal.

### Step 1: Register Users
Create two users to test transactions between them.

**User A (Alice):**
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username": "alice_test", "email": "alice_test@example.com", "password": "password123"}'
```
*Take note of the `id` returned in the response (e.g., `1`).*

**User B (Bob):**
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username": "bob_test", "email": "bob_test@example.com", "password": "password123"}'
```
*Take note of the `id` returned (e.g., `2`).*

### Step 2: Deposit Funds (Triggering the Worker)
Send a "deposit" transaction for Alice. This request is sent to the `Worker Pool`.

```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{"to_user_id": 1, "amount": 1000, "type": "deposit"}'
```
**Response:** `202 Accepted`
The server accepted the request but processed it in the background.

### Step 3: Verify Processing
Check the logs to see the worker picking up the task:

```bash
docker-compose logs -f app
```
You should see a log like:
`INFO msg="Worker processed transaction" worker_id=1 tx_id=...`

### Step 4: Check Balance
Verify that the worker successfully updated the balance.

```bash
curl "http://localhost:8080/api/balance?user_id=1"
```
**Response:** `{"user_id":1,"amount":1000,...}`

---

## 2. Automated E2E Test

We have provided a Go program that runs a full scenario automatically:
1. Registers two valid users.
2. Deposits money to one user.
3. Transfers money between them.
4. Waits for the worker to process the changes.
5. Verifies the final balances.

**Run the test:**
```bash
go run cmd/e2e/main.go
```

## 3. How the Worker Works (Verification)

To verify the asynchronous nature (that the API returns *before* the work is done):

1. **Observe the "Accepted" status**: The transaction endpoint returns `202 Accepted`, not `200 OK`. This is the HTTP standard for "I heard you, I'll do it later."
2. **Read the Logs**: The application logs show when a worker picks up a job from the queue.
3. **Simulate Load**: If you send many requests rapidly, you will see the workers processing them one by one (or 5 at a time, since we have 5 workers).
