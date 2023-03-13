package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Message struct {
	Ebay_order_id int    `json:"ebay_order_id"`
	Product_id    int    `json:"product_id"`
	Name          string `json:"name"`
	City          string `json:"city"`
}

func main() {

	// to produce messages
	topic := "customer-001"
	partition := 0
	address := "172.105.117.236:9092" // kafka broker address - Linode

	conn, err := kafka.DialLeader(context.Background(), "tcp", address, topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	message := Message{
		Ebay_order_id: 100,
		Product_id:    1,
		Name:          "Ray Ban Sunglasses",
		City:          "Helsinki",
	}
	data, err := json.Marshal(message)

	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte(data)},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

}
