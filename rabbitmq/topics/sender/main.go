package main

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func main() {
	conn, err := amqp.Dial("amqp://rabbitmq:rabbitmq@10.211.55.18:5672")
	if err != nil {
		log.Println(err)
		return
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
		return
	}
	err = ch.ExchangeDeclare("logs_topic", "topic", true, false, false, false, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// . #
	// *.orange.*
	// .*.rabbit
	// lazy.#
	err = ch.PublishWithContext(context.Background(), "logs_topic", "lazy.#", false, false, amqp.Publishing{
		ContentType: "",
		Body:        []byte("Hello"),
	})
	if err != nil {
		log.Println(err)
		return
	}
}
