// Created by zhbinary on 2023/8/19.
package main

import (
	"bytes"
	"encoding/json"
	"fbchatbot/pkg"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	// Load .env
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Failed to load .env:", err)
		return
	}

	router := mux.NewRouter()

	router.HandleFunc("/", HomePageHandler).Methods("GET")
	router.HandleFunc("/webhook", WebhookVerificationHandler).Methods("GET")
	router.HandleFunc("/webhook", WebhookEventHandler).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Your app is listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func WebhookVerificationHandler(w http.ResponseWriter, r *http.Request) {
	verifyToken := os.Getenv("VERIFY_TOKEN")
	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
	challenge := r.URL.Query().Get("hub.challenge")

	if mode == "subscribe" && token == verifyToken {
		fmt.Println("WEBHOOK_VERIFIED")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(challenge))
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func WebhookEventHandler(w http.ResponseWriter, r *http.Request) {
	// ready body in the request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// unmarshal []byte data into message
	var event *pkg.WebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("failed to unmarshal webhook: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if event.Object != "page" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	result := true
	for _, entry := range event.Entry {
		// Gets the body of the webhook event
		Messaging := entry.Messaging[0]
		log.Println(Messaging)

		// Get the sender PSID
		senderPsid := Messaging.Sender.ID
		log.Printf("Sender PSID:  %v", senderPsid)

		if Messaging.Message != nil {
			response := pkg.HandleMessage(senderPsid, Messaging.Message)
			callSendAPI(senderPsid, response)
		} else if Messaging.FeedBack != nil {
			result = pkg.HandleFeedback(senderPsid, Messaging.FeedBack)
		}
	}

	if result {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("EVENT_RECEIVED"))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func callSendAPI(senderPsid string, response interface{}) {
	pageAccessToken := os.Getenv("PAGE_ACCESS_TOKEN")
	requestBody := map[string]interface{}{
		"recipient": map[string]interface{}{
			"id": senderPsid,
		},
		"message": response,
	}

	bodyJson, err := json.Marshal(requestBody)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("MESSAGE:", string(bodyJson))

	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://graph.facebook.com/v2.6/me/messages", bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Println("Failed to create HTTP request:", err)
		return
	}

	q := req.URL.Query()
	q.Add("access_token", pageAccessToken)
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send message:", err)
		pkg.ResendReply(senderPsid, string(bodyJson))
		return
	}
	defer resp.Body.Close()

	log.Println("Message sent!")
}
