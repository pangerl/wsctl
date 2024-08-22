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

func produceMessage(topic, nameserver string) bool {
	p, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{nameserver}),
	)
	if err != nil {
		log.Println("创建生产者失败:", err)
		return false
	}
	defer func(p rocketmq.Producer) {
		log.Println("开始关闭生产者", err)
		err := p.Shutdown()
		if err != nil {
			log.Println("关闭生产者失败:", err)
		}
	}(p)

	err = p.Start()
	if err != nil {
		log.Println("启动生产者失败:", err)
		return false
	}

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

func isConsumptionComplete() bool {
	// 根据您的实际业务逻辑来判断消费是否完成
	// 例如，检查已消费的消息数量是否达到预期，或者是否收到特定的结束标志消息等
	return true
}

func consumeMessage(topic, groupName, nameserver string) bool {
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{nameserver}),
		consumer.WithGroupName(groupName),
	)
	if err != nil {
		log.Println("创建消费者失败:", err)
		return false
	}
	defer func(c rocketmq.PushConsumer) {
		log.Println("开始关闭消费者", err)
		err := c.Shutdown()
		if err != nil {
			log.Println("关闭消费者失败:", err)
		}
	}(c)

	err = c.Subscribe(topic, consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			log.Printf("收到消息: %s\n", string(msg.Body))
		}
		if isConsumptionComplete() {
			err := c.Shutdown()
			if err != nil {
				return 0, err
			}
			return 0, nil
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		log.Println("订阅失败:", err)
		return false
	}

	err = c.Start()
	if err != nil {
		log.Println("启动消费者失败:", err)
		return false
	}

	time.Sleep(5 * time.Minute)
	return true
}

func ProbeRocketMq(nameserver string) {
	topic := "demo_topic"
	groupName := "consumer_group_demo"
	mqStatus := ""
	if produceMessage(topic, nameserver) {
		if consumeMessage(topic, groupName, nameserver) {
			mqStatus = "正常"
		}

	}
	mqStatus = "异常"
	log.Println("RocketMq 集群状态:", mqStatus)
}
