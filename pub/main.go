package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/nats-io/stan.go"

	"github.com/yalagtyarzh/L0/models"
	"github.com/yalagtyarzh/L0/pub/internal/config"
)

func main() {
	cfg := config.GetConfig()
	var orders []models.Order
	file, err := os.Open("./pub/internal/config/orders.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(data, &orders)
	if err != nil {
		log.Fatalln(err)
	}

	sc, _ := stan.Connect(cfg.STAN.Cluster, cfg.STAN.Client)
	defer sc.Close()

	for _, order := range orders {
		jsonData, err := json.Marshal(order)
		if err != nil {
			log.Fatalln(err)
		}

		err = sc.Publish("foo", jsonData)
		log.Println("sent")

		if err != nil {
			log.Printf("error to send data: %s", err)
		}
		time.Sleep(time.Second * time.Duration(cfg.STAN.Delay))
	}
}
