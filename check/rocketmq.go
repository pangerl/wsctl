// Package check @Author lanpang
// @Date 2024/8/21 下午3:58:00
// @Desc
package check

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"log"
	"time"
)

func produceMessage(p rocketmq.Producer, topic string) bool {
	message := &primitive.Message{
		Topic: topic,
		Body:  []byte("这是一条测试消息"),
	}
	res, err := p.SendSync(context.Background(), message)
	if err != nil {
		log.Println("发送消息失败:", err)
		return false
	}
	log.Printf("发送消息成功: %s\n", res.MsgID)
	return true
}

func consumeMessage(c rocketmq.PushConsumer, topic string) bool {
	// 定义消费处理函数
	consumeFunc := func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			log.Printf("Received message: %s\n", string(msg.Body))
			//consumer.Ack(msg) // 确保正确导入并使用 c 来调用 Ack
		}
		return consumer.ConsumeSuccess, nil
	}

	// 执行订阅
	err := c.Subscribe(topic, consumer.MessageSelector{}, consumeFunc)
	if err != nil {
		log.Println("订阅失败:", err)
		return false
	}

	// 启动消费
	err = c.Start()
	if err != nil {
		log.Println("启动消费者失败:", err)
		return false
	}

	// 为了避免程序一直运行，设置一个超时等待
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			log.Println("Consume timeout")
			return true
		default:
			// 可以添加一些其他处理逻辑或日志
		}
	}
	//return true
}

func ProbeRocketMq(nameserver string) bool {
	topic := "demo_topic"
	groupName := "consumer_group_demo"
	// 创建生产者
	p, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{nameserver}),
		producer.WithRetry(2),
	)
	if err != nil {
		log.Println("创建生产者失败:", err)
		return false
	}
	defer func() {
		log.Println("开始关闭生产者")
		err := p.Shutdown()
		if err != nil {
			log.Println("关闭生产者失败:", err)
		}
	}()
	// 启动生产者
	err = p.Start()
	if err != nil {
		log.Println("启动生产者失败:", err)
		return false
	}
	// 创建消费者
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{nameserver}),
		consumer.WithGroupName(groupName),
	)
	if err != nil {
		log.Println("创建消费者失败:", err)
		return false
	}
	defer func() {
		log.Println("开始关闭消费者")
		err := c.Shutdown()
		if err != nil {
			log.Println("关闭消费者失败:", err)
		}
	}()
	if produceMessage(p, topic) {
		if consumeMessage(c, topic) {
			return true
		}
	}
	return false
}
