package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type User struct {
	Greeting    string   `json:"greeting"`
	Ipaddresses []string `json:"Ipaddresses"`
}

type Message struct {
	Greeting    string   `json:"greeting"`
	Customer    string   `json:"customer"`
	Ipaddresses []string `json:"Ipaddresses"`
	Orders      []string `json:"orders"`
}

func handleFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	value := r.URL.Path[1:]

	if value != "" {
		greetings := "Hello Roshan Fernando - Linode"
		var ip_addresses []string

		addrs, err := net.InterfaceAddrs()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				ip_addresses = append(ip_addresses, ipnet.IP.String())
			}
		}

		topic := value
		partition := 0
		address := "172.105.117.236:9092"
		var kafka_messages []string

		conn, err := kafka.DialLeader(context.Background(), "tcp", address, topic, partition)
		if err != nil {
			log.Fatal("failed to dial leader:", err)
		}

		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

		defer batch.Close()

		defer conn.Close()

		b := make([]byte, 10e3) // 10KB max per message
		for {
			n, err := batch.Read(b)
			if err != nil {
				break
			}
			r := strings.NewReplacer("\t", "", "\n", "", "\\", "")
			newStr := r.Replace(string(b[:n]))

			kafka_messages = append(kafka_messages, newStr)
		}

		message := []Message{
			{Greeting: greetings, Ipaddresses: ip_addresses, Customer: value, Orders: kafka_messages},
		}

		jsonData, err := json.Marshal(message)
		if err != nil {
			log.Println("Failed to marshal JSON data:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		w.Write(jsonData)

	} else {
		greetings := "Hello Roshan Fernando - Linode"
		var ip_addresses []string

		addrs, err := net.InterfaceAddrs()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				ip_addresses = append(ip_addresses, ipnet.IP.String())
			}
		}

		users := []User{
			{Greeting: greetings, Ipaddresses: ip_addresses},
		}

		jsonData, err := json.Marshal(users)
		if err != nil {
			log.Println("Failed to marshal JSON data:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Set the Content-Type header to application/json
		w.Header().Set("Content-Type", "application/json")

		// Write the JSON data to the HTTP response
		w.Write(jsonData)
	}

}
func main() {
	// REST service

	// Establish a service
	var handler http.ServeMux
	handler.HandleFunc("/", handleFunc)

	server := http.Server{
		Addr:         ":8080", // localhost:8080 (host:port)
		Handler:      &handler,
		ReadTimeout:  0,
		WriteTimeout: 0,
	}

	server.ListenAndServe()
}