package main

import (
	"context"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"time"
)

func main() {
	queue := "job_queue"
	addr := "amqp://rabbitmq:rabbitmq@localhost:5672"
	c := New(queue, addr)
	for i := 0; i < 100; i++ {
		t := time.Now().Format(time.RFC3339)
		select {
		case <-time.After(time.Second * 2):
			// Push 失败的原因可能是什么?
			err := c.Push([]byte(t))
			// TODO 是否需要重新初始化
			if err != nil {
				log.Println("Push Failed")
			} else {
				log.Println("Push Succeeded")
			}
		}
	}
}

const (
	reconnectDelay = time.Second * 5
	reInitDelay    = time.Second * 3
	resendDelay    = time.Second * 2
)

type Client struct {
	addr            string
	queue           string
	logger          *log.Logger
	conn            *amqp.Connection
	channel         *amqp.Channel
	done            chan bool
	notifyConnClose chan *amqp.Error
	notifyChanClose chan *amqp.Error
	notifyConfirm   chan amqp.Confirmation
	isReady         bool
}

func New(queue, addr string) *Client {
	client := &Client{
		addr:   addr,
		queue:  queue,
		logger: log.New(os.Stdout, "", log.LstdFlags),
		done:   make(chan bool),
	}
	go client.Reconnect(addr)
	return client
}

func (c *Client) Reconnect(addr string) {
	for {
		c.isReady = false
		c.logger.Println("Attempting to connect")
		conn, err := c.connect(addr)
		// 重连
		if err != nil {
			c.logger.Println("Failed to connect")
			select {
			case <-c.done:
				return
			case <-time.After(reconnectDelay):
			}
			continue
		}
		if done := c.reInit(conn); done {
			break
		}
	}
}

func (c *Client) connect(addr string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		return nil, err
	}
	// new conn
	c.conn = conn
	c.notifyConnClose = make(chan *amqp.Error, 1)
	c.conn.NotifyClose(c.notifyConnClose)
	c.logger.Println("Connected!")
	return conn, nil
}

func (c *Client) reInit(conn *amqp.Connection) bool {
	for {
		c.isReady = false
		err := c.initChannel(conn)
		// 初始化 channel 与 queue 声明
		if err != nil {
			c.logger.Println("Failed Init channel. Retrying...")
			select {
			case <-c.done:
				return true
			case <-time.After(reInitDelay):
			}
			continue
		}
		select {
		case <-c.done:
			return true
		case <-c.notifyConnClose:
			c.logger.Println("Connection closed. Reconnecting...")
			return false
		case <-c.notifyChanClose:
			c.logger.Println("Channel closed. Re-running init...")
			return false
		}
	}
}

func (c *Client) initChannel(conn *amqp.Connection) error {
	// 初始化 channel
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	// 同步等待Push响应
	err = ch.Confirm(false)
	if err != nil {
		return err
	}
	// queue 声明
	_, err = ch.QueueDeclare(c.queue, true, false, false, false, nil)
	if err != nil {
		return err
	}
	// new channel
	c.channel = ch
	c.notifyChanClose = make(chan *amqp.Error, 1)
	c.notifyConfirm = make(chan amqp.Confirmation, 1)
	c.channel.NotifyClose(c.notifyChanClose)
	c.channel.NotifyPublish(c.notifyConfirm)

	c.isReady = true
	c.logger.Println("Setup!")

	return nil
}

func (c *Client) Push(data []byte) error {
	if !c.isReady {
		return errors.New("client Not Ready")
	}
	for {
		// 多次Push是否有问题
		err := c.UnSafePush(data)
		if err != nil {
			c.logger.Println("Push Failed, Retrying...")
			select {
			case <-c.done:
				return errors.New("server shutdown")
			case <-time.After(resendDelay):
			}
			continue
		}
		select {
		case confirm := <-c.notifyConfirm:
			if confirm.Ack {
				c.logger.Println("Push confirmed")
				return nil
			}
		case <-time.After(resendDelay):
		}
		// 没有收到MQ的确认, 进行重发
		c.logger.Println("Push do not confirmed. Retrying...")
	}
}

func (c *Client) UnSafePush(data []byte) error {
	if !c.isReady {
		return errors.New("client Not Ready")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return c.channel.PublishWithContext(
		ctx,
		"",
		c.queue,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: amqp.Persistent,
			Body:         data,
		},
	)
}

func (c *Client) Consume() (<-chan amqp.Delivery, error) {
	if !c.isReady {
		return nil, errors.New("client Not Ready")
	}
	if err := c.channel.Qos(1, 0, false); err != nil {
		return nil, err
	}
	return c.channel.Consume(c.queue, "", false, false, false, false, nil)
}

func (c *Client) Close() error {
	if !c.isReady {
		return errors.New("already Close")
	}
	close(c.done)
	if err := c.channel.Close(); err != nil {
		return err
	}
	if err := c.conn.Close(); err != nil {
		return err
	}
	c.isReady = false
	return nil
}
