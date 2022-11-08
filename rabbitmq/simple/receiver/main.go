package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func main() {
	conn, err := amqp.Dial("amqp://rabbitmq:rabbitmq@10.211.55.18:5672")
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
		return
	}
	defer ch.Close()
	q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	if err != nil {
		log.Println(err)
	}
	// 消费消息
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,  // autoAck 是否自动的确认
		false, // exclusive 是否只允许一个消费者
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Wait.....")
	for msg := range msgs {
		log.Printf("Recevied Message: %s", msg.Body)
	}
}
