package main

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

// 消息投递结果通知
var cf chan amqp.Confirmation

func main() {
	conn, err := amqp.Dial("amqp://rabbitmq:rabbitmq@10.211.55.18:5672")
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	defer ch.Close()

	// 第一阶段: 发布确认机制confirm, 持久化到硬盘通知 [异步的]
	ch.Confirm(false)
	cf = ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	// 持久化
	err = ch.ExchangeDeclare("logs", "direct", true, false, false, false, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// 队列持久化是指: 重启恢复队列中的数据
	q, err := ch.QueueDeclare("hello", true, false, false, false, nil)

	if err != nil {
		log.Println(err)
	}

	go func() {
		// 消息确认投递
		for {
			if confirmed := <-cf; confirmed.Ack {
				// 那些消息投递成功，那些失败了
				// map 中去判断一下
				//fmt.Println(confirmed)
			}
		}
	}()
	// 存起来？ map: seq: message
	// seq := ch.GetNextPublishSeqNo()

	for i := 0; i < 1000; i++ {
		body := "Hello World" + time.Now().Format(time.RFC3339)
		// 第二阶段：mandatory returnBack
		err = ch.PublishWithContext(
			context.Background(),
			"",
			q.Name,
			false, // mandatory 需要在看下
			false, // immediate 需要在看下, 废弃
			amqp.Publishing{
				// 消息持久化到磁盘，配合上前面的队列持久化，即使服务挂掉, 重启也会恢复队列中的数据
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(body),
			},
		)
		fmt.Println(body)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(time.Second * 2)
	}
}
