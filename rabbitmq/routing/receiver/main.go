package main

import (
	"flag"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func main() {
	name := flag.String("c", "c", "consumer name")
	queue := flag.String("q", "q", "queue name")
	key := flag.String("k", "k", "routing key")
	flag.Parse()
	conn, err := amqp.Dial("amqp://rabbitmq:rabbitmq@10.211.55.18:5672")
	if err != nil {
		log.Println(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
	}
	err = ch.ExchangeDeclare("logs_direct", "direct", true, false, false, false, nil)
	if err != nil {
		log.Println(err)
		return
	}
	q, err := ch.QueueDeclare(*queue, false, false, true, false, nil)
	if err != nil {
		log.Println(err)
		return
	}
	err = ch.QueueBind(q.Name, *key, "logs_direct", false, nil)
	if err != nil {
		log.Println(err)
		return
	}
	msgs, err := ch.Consume(q.Name, *name, true, false, false, false, nil)
	if err != nil {
		log.Println(err)
		return
	}
	for msg := range msgs {
		log.Printf(" [x] %s", msg.Body)
	}
}
