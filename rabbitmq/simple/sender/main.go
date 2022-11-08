package main

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func main() {
	// 获取连接
	conn, err := amqp.Dial("amqp://rabbitmq:rabbitmq@10.211.55.18:5672")
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// 获取连接中的channel
	ch, err := conn.Channel()
	defer ch.Close()

	// 声明队列
	q, err := ch.QueueDeclare(
		"hello", // 队列名称
		false,   // durable 是否持久化, 通常给true
		false,   // autoDelete 是否自动删除, 当没有消费链接的时候是否自动删除, 通常给false
		false,   // exclusive 是否独占, 不允许其他的操作 declare, bind, consume, purge or delete a queue with the same name.
		false,   // noWait 为true, 即是假定已经声明在服务端, 不允许修改, 否则跑出错误
		nil,
	)
	if err != nil {
		log.Println(err)
	}
	// 简单模式下, 使用默认的交换机
	body := "Hello World" + time.Now().Format(time.RFC3339)
	err = ch.PublishWithContext(
		context.Background(),
		"",     // exchange  发送到默认交换机
		q.Name, // key 路由到queue
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		log.Println(err)
	}
}
