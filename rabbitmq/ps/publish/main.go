package main

import (
	"github.com/streadway/amqp"
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
	err = ch.ExchangeDeclare("logs", "fanout", true, false, false, false, nil)
	if err != nil {
		log.Println(err)
		return
	}
	err = ch.Publish("logs", "", false, false, amqp.Publishing{
		ContentType: "",
		Body:        []byte("Hello"),
	})
	if err != nil {
		log.Println(err)
		return
	}
}
