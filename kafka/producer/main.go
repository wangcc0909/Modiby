package main

import (
	"fmt"
	"github.com/Shopify/sarama"
)

func main() {
	config := sarama.NewConfig()
	//
	config.Producer.RequiredAcks = sarama.WaitForAll
	//
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	//
	producer, err := sarama.NewAsyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		fmt.Printf("producer create producer error :%s\n", err.Error())
		return
	}
	defer producer.Close()

	//构建发送的消息
	msg := &sarama.ProducerMessage{
		//Topic: "test", //包含了消息的主题
		Partition: int32(10),
		Key:       sarama.StringEncoder("key"),
	}
	var value string
	var msgType string
	for {
		_, err := fmt.Scanf("%s", &value)
		if err != nil {
			break
		}
		fmt.Scanf("%s", &msgType)
		fmt.Println("msgType = ", msgType, " ,value = ", value)
		msg.Topic = msgType
		//将字符串转换未字节数组
		msg.Value = sarama.ByteEncoder(value)
		//SendMessage: 该方法是生产者生产给定的消息
		//生产成功的时候返回该消息的分区和所在的偏移量
		//生产失败的时候返回error
		//partition,offset,err := producer.SendMessage(msg)
		producer.Input() <- msg
		select {
		case suc := <-producer.Successes():
			fmt.Printf("partition = %d, offset = %d\n", suc.Partition, suc.Offset)
		case fail := <-producer.Errors():
			fmt.Printf("error = %s\n", fail.Err.Error())
		}
	}
}
