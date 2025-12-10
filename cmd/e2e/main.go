package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8080/api"

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Balance struct {
	UserID int64 `json:"user_id"`
	Amount int64 `json:"amount"`
}

func main() {
	fmt.Println("🚀 Starting End-to-End Test")

	suffix := time.Now().UnixNano()
	userA := register(fmt.Sprintf("alice_%d", suffix), fmt.Sprintf("alice_%d@example.com", suffix), "password")
	fmt.Printf("✅ Registered User A: ID=%d\n", userA.ID)

	// 2. Register User B
	userB := register(fmt.Sprintf("bob_%d", suffix), fmt.Sprintf("bob_%d@example.com", suffix), "password")
	fmt.Printf("✅ Registered User B: ID=%d\n", userB.ID)

	// 3. Deposit to User A
	fmt.Println("\n💳 Depositing 1000 to Alice...")
	createTransaction(nil, &userA.ID, 1000, "deposit")

	// 4. Verify Balance A (Wait for Worker)
	waitForBalance(userA.ID, 1000)
	fmt.Println("✅ Alice's balance confirmed: 1000")

	// 5. Transfer from A to B
	fmt.Println("\n💸 Transferring 500 from Alice to Bob...")
	createTransaction(&userA.ID, &userB.ID, 500, "transfer")

	// 6. Verify Balances
	waitForBalance(userA.ID, 500)
	fmt.Println("✅ Alice's balance confirmed: 500")

	waitForBalance(userB.ID, 500)
	fmt.Println("✅ Bob's balance confirmed: 500")

	fmt.Println("\n🎉 All tests passed! The worker is successfully processing transactions.")
}

func register(username, email, password string) User {
	payload := map[string]string{
		"username": username,
		"email":    email,
		"password": password,
	}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(baseURL+"/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("Register failed: %s", string(b)))
	}

	var user User
	json.NewDecoder(resp.Body).Decode(&user)
	return user
}

func createTransaction(from, to *int64, amount int64, typeStr string) {
	payload := map[string]interface{}{
		"amount": amount,
		"type":   typeStr,
	}
	if from != nil {
		payload["from_user_id"] = *from
	}
	if to != nil {
		payload["to_user_id"] = *to
	}

	body, _ := json.Marshal(payload)
	resp, err := http.Post(baseURL+"/transactions", "application/json", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		b, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("Transaction failed: %s", string(b)))
	}
}

func waitForBalance(userID int64, expectedAmount int64) {
	fmt.Printf("⏳ Waiting for balance %d for User %d...", expectedAmount, userID)
	for i := 0; i < 50; i++ { // Try for 5 seconds
		resp, err := http.Get(fmt.Sprintf("%s/balance?user_id=%d", baseURL, userID))
		if err == nil && resp.StatusCode == http.StatusOK {
			var bal Balance
			json.NewDecoder(resp.Body).Decode(&bal)
			resp.Body.Close()

			if bal.Amount == expectedAmount {
				fmt.Println(" Done.")
				return
			}
		}
		time.Sleep(100 * time.Millisecond)
		fmt.Print(".")
	}
	panic("Timeout waiting for balance update")
}
