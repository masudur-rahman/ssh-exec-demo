package main

import (
	"log"
	"os"
	"time"

	"github.com/hpcloud/tail"
	"github.com/nats-io/nats.go"

	"github.com/tamalsaha/ssh-exec-demo/internal/util"
)

func publishFile(source string, subject string, nc *nats.Conn) error {
	t, err := tail.TailFile(source, tail.Config{Follow: true})
	if err != nil {
		return err
	}

	log.Printf("Publishing lines from %s to %s", source, subject)

	for line := range t.Lines {
		nc.Publish(subject, []byte(line.Text))
	}

	return nil
}

func main() {
	source := os.Getenv("SHIPPER_FILE")
	if source == "" {
		log.Fatalf("Please set a file to publish using SHIPPER_FILE\n")
	}

	subject := os.Getenv("SHIPPER_SUBJECT")
	if subject == "" {
		log.Fatalf("Please set a NATS subject to publish to using SHIPPER_SUBJECT\n")
	}

	//partition, name, err := util.Partition()
	//if err != nil {
	//	log.Fatal(err.Error())
	//}

	addr := "this-is-nats.appscode.ninja:4222"
	nc, err := util.NewConnection(addr, "")
	if err != nil {
		log.Fatalf("Could not connect to NATS: %s\n", err)
	}

	for {
		// fmt.Sprintf("%s.p%d.%s", subject, partition, name)
		err = publishFile(source, subject, nc)
		if err != nil {
			log.Printf("Could not publish file: %s", err)
		}

		time.Sleep(500 * time.Millisecond)
	}
}
