package rabbitmq

import (
	"errors"
	"strings"
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var Channel *amqp.Channel
var exchanges string
var queues string
var hasMQ bool = false

func init() {
	configFilePath := "config.ini"
	cfg, _ := goconfig.LoadConfigFile(configFilePath)
	mqcfg, _ := cfg.GetSection("rabbitmq")
	path := strings.Join([]string{"amqp://", mqcfg["user"], ":", mqcfg["password"], "@", mqcfg["host"], ":", mqcfg["port"], "/oj"}, "")
	fmt.Println(path)
	err := SetupRMQ(path)
	fmt.Println(err)
}

type Reader interface {
	Read(msg *string) (err error)
}

// 初始化 参数格式：amqp://用户名:密码@地址:端口号/host
func SetupRMQ(rmqAddr string) (err error) {
	if Channel == nil {
		conn, err = amqp.Dial(rmqAddr)
		if err != nil {
			return err
		}
		Channel, err = conn.Channel()
		if err != nil {
			return err
		}
		hasMQ = true
	}
	return nil
}

// 是否已经初始化
func HasMQ() bool {
	return hasMQ
}

// 测试连接是否正常
func Ping() (err error) {
	if !hasMQ || Channel == nil {
		return errors.New("RabbitMQ is not initialize")
	}
	err = Channel.ExchangeDeclare("ping.ping", "direct", false, true, false, true, nil)
	if err != nil {
		return err
	}
	msgContent := "ping.ping"
	err = Channel.Publish("ping.ping", "ping.ping", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:		[]byte(msgContent),
	})
	if err != nil {
		return err
	}
	err = Channel.ExchangeDelete("ping.ping", false, false)
	return err
}

// 发布消息
func Publish(exchange, routingKey string, msg []byte) (err error) {
	if exchanges == "" || !strings.Contains(exchanges, exchange) {
		err = Channel.ExchangeDeclare(exchange, "direct", true, false, false, true, nil)
		if err != nil {
			return err
		}
		exchanges += "  " + exchange + "  "
	}
	err = Channel.Publish(exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        msg,
	})
	return nil
}

// 监听接收到的消息
func Receive(exchange, queue, routingKey string, reader func(msg *[]byte)) (err error) {
	if exchanges == "" || !strings.Contains(exchanges, exchange) {
		err = Channel.ExchangeDeclare(exchange, "direct", true, false, false, true, nil)
		if err != nil {
			return err
		}
		exchanges += "  " + exchange + "  "
	}
	if queues == "" || !strings.Contains(queues, queue) {
		_, err = Channel.QueueDeclare(queue, true, false, false, true, nil)
		if err != nil {
			return err
		}
		err = Channel.QueueBind(queue, routingKey, exchange, true, nil)
		if err != nil {
			return err
		}
		queues += "  " + queue + "  "
	}
	msgs, err := Channel.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		//fmt.Println(*msgs)
		for d := range msgs {
			// s := bytesToString(&(d.Body))
			
			reader(&(d.Body))
			
		}
	}()
	return nil
}

// 关闭连接
func Close() {
	_ = Channel.Close()
	_ = conn.Close()
	hasMQ = false
}

// func bytesToString(b *[]byte) *string {
// 	s := bytes.NewBuffer(*b)
// 	r := s.String()
// 	return &r
// }
