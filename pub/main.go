package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/nats-io/stan.go"

	"github.com/yalagtyarzh/L0/models"
	"github.com/yalagtyarzh/L0/pub/internal/config"
)

var dir = "./pub/orders"
var validExt = ".json"

// main simulates sending data to channel via api
func main() {
	cmd := exec.Command("/bin/sh", "nats.sh")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("can't run cmd")
	}

	fmt.Print("xd")

	cfg := config.GetConfig()
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalln(err)
	}
	sc, _ := stan.Connect(cfg.STAN.Cluster, cfg.STAN.Client)

	defer sc.Close()

	for _, file := range files {
		time.Sleep(time.Second * time.Duration(cfg.STAN.Delay))
		fileName := file.Name()
		if !checkExt(fileName, validExt) {
			log.Println("Invalid file extension of file", fileName)
			continue
		}

		f, err := os.Open(fmt.Sprintf("%s/%s", dir, fileName))
		if err != nil {
			fmt.Printf("Couldn't open file: %s\n", err)
			continue
		}
		defer f.Close()

		data, err := io.ReadAll(f)
		if err != nil {
			log.Printf("Couldn't read file: %s\n", err)
			continue
		}

		var order models.Order
		err = json.Unmarshal(data, &order)
		if err != nil {
			log.Printf("Invalid json: %s\n", err)
			continue
		}

		jsonData, err := json.Marshal(order)
		if err != nil {
			log.Printf("Couldn't marshal: %s\n", err)
			continue
		}

		err = sc.Publish("foo", jsonData)

		if err != nil {
			log.Printf("error to send data: %s", err)
			continue
		}

		log.Println("sent")
	}

	sub, _ := sc.Subscribe(
		"foo", func(msg *stan.Msg) {
			fmt.Printf("Received a message: %s\n", string(msg.Data))
		},
	)
	sub.Unsubscribe()
}

// checkExt checks for valid file extension
func checkExt(filename string, ext string) bool {
	fileExt := filepath.Ext(filename)
	if fileExt != ext {
		return false
	}

	return true
}
