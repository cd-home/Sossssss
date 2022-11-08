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

	// 持久化
	err = ch.ExchangeDeclare("logs", "direct", true, false, false, false, nil)
	if err != nil {
		log.Println(err)
		return
	}

	q, err := ch.QueueDeclare("hello", true, false, false, false, nil)
	if err != nil {
		log.Println(err)
	}
	msgs, err := ch.Consume(
		q.Name,
		"",
		false, // autoAck = false, 采用手动确认方式, 自动模式可能存在消息丢失情况
		false,
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
		log.Printf("Recevied Message: %s\n", msg.Body)
		// 处理成功
		msg.Ack(false)

		// 处理失败
		//msg.Nack(true, true)
	}
}
