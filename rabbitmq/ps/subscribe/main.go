package main

import (
	"flag"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func main() {
	name := flag.String("c", "c1", "consumer name")
	queue := flag.String("q", "q1", "queue name")
	flag.Parse()
	conn, err := amqp.Dial("amqp://rabbitmq:rabbitmq@10.211.55.18:5672")
	if err != nil {
		log.Println(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
	}
	// 绑定交换机
	err = ch.ExchangeDeclare("logs-fanout", "fanout", true, false, false, false, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// durable=true持久化, exclusive=false连接断开不删除queue
	q, err := ch.QueueDeclare(*queue, true, false, false, false, nil)
	if err != nil {
		log.Println(err)
		return
	}
	err = ch.QueueBind(q.Name, "", "logs-fanout", false, nil)
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
		fmt.Printf(" [x] %s\n", msg.Body)
	}
}
