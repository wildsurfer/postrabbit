package main

import (
	"crypto/tls"
	"log"
	"time"

	pq "github.com/lib/pq"
	"github.com/streadway/amqp"
)

func errorReporter(ev pq.ListenerEventType, err error) {
	if err != nil {
		log.Print(err)
	}
}

func run(config Config) {
	listener := pq.NewListener(config.PostgresURL, 10*time.Second, time.Minute, errorReporter)
	err := listener.Listen("urlwork")
	if err != nil {
		log.Fatal(err)
	}

	rabbitchannel := make(chan string, 100)

	go func() {
		cfg := new(tls.Config)
		cfg.InsecureSkipVerify = true
		conn, err := amqp.DialTLS(config.RabbitMQURL, cfg)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		ch, err := conn.Channel()
		if err != nil {
			log.Fatal(err)
		}
		defer ch.Close()

		for {
			payload := <-rabbitchannel
			log.Println(payload)
			err := ch.Publish("urlwork", "todo", false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(payload),
			})
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	for {
		select {
		case notification := <-listener.Notify:
			rabbitchannel <- notification.Extra
		case <-time.After(90 * time.Second):
			go func() {
				err := listener.Ping()
				if err != nil {
					log.Fatal(err)
				}
			}()
		}
	}
}
