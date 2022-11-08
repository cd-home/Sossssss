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
	// 首先投递失败存在如下两个情况
	// 1. 交换器无法根据自身的类型和路由键匹配到队列
	// 2. 交换器将消息路由到队列时, 发现队列上并不存在任何消费者

	// 监听投递失败的场景
	// 生产者投递消息时指定 mandatory 或者 immediate 参数设为 true
	// RabbitMQ 会把无法投递的消息通过 Basic.Return 命令将消息返回给生产者
	notifyReturn := ch.NotifyReturn(make(chan amqp.Return, 1))
	go func() {
		for {
			body, ok := <-notifyReturn
			if ok {
				fmt.Println("no push", string(body.Body))
			}
		}
	}()

	for i := 0; i < 1000; i++ {
		t := time.Now().Format(time.RFC3339)
		err = ch.PublishWithContext(
			context.Background(),
			"logs-fanout", "",
			true,
			false,
			amqp.Publishing{
				ContentType: "",
				Body:        []byte(t),
				// 消息持久化
				DeliveryMode: amqp.Persistent,
				// 可以设置消息的过期时间, 优先级
			})
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(t)
		time.Sleep(time.Second * 2)
	}
}
