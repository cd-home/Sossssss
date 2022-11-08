package main

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
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
	// 交换机持久化
	err = ch.ExchangeDeclare("logs-fanout", "fanout", true, false, false, false, nil)
	if err != nil {
		log.Println(err)
		return
	}
	for i := 0; i < 1000; i++ {
		t := time.Now().Format(time.RFC3339)
		err = ch.PublishWithContext(context.Background(), "logs-fanout", "", false, false, amqp.Publishing{
			ContentType: "",
			Body:        []byte(t),
			// 消息持久化
			DeliveryMode: amqp.Persistent,
		})
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(t)
		time.Sleep(time.Second * 2)
	}
}
